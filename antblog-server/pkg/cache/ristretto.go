package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

// ─── Ristretto 本地缓存实现 ──────────────────────────────────────────────────

// localItem 缓存条目，携带过期时间
type localItem struct {
	Value     []byte
	ExpiresAt time.Time // zero = 永不过期
}

func (i *localItem) expired() bool {
	return !i.ExpiresAt.IsZero() && time.Now().After(i.ExpiresAt)
}

// RistrettoCache 基于 dgraph-io/ristretto 的本地内存缓存
// 实现 ICache 接口（注意：Incr/SetNX 通过 mutex 保证简单原子性）
type RistrettoCache struct {
	cache      *ristretto.Cache
	defaultTTL time.Duration
	keyPrefix  string
	mu         sync.Mutex
	counters   map[string]int64 // 用于 Incr（ristretto 无原生计数器）
}

// NewRistretto 创建本地缓存实例
func NewRistretto(opts ...Option) (*RistrettoCache, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: o.ristrettoNumCounters,
		MaxCost:     o.ristrettoMaxCost,
		BufferItems: o.ristrettoBufferItems,
	})
	if err != nil {
		return nil, fmt.Errorf("cache: create ristretto: %w", err)
	}

	return &RistrettoCache{
		cache:      cache,
		defaultTTL: o.defaultTTL,
		keyPrefix:  o.keyPrefix,
		counters:   make(map[string]int64),
	}, nil
}

func (c *RistrettoCache) k(key string) string {
	return c.keyPrefix + key
}

func (c *RistrettoCache) ttl(d time.Duration) time.Duration {
	if d == 0 {
		return c.defaultTTL
	}
	return d
}

// Set 设置缓存（自动 JSON 序列化）
func (c *RistrettoCache) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: marshal value: %w", err)
	}

	item := &localItem{
		Value:     data,
		ExpiresAt: time.Now().Add(c.ttl(ttl)),
	}

	cost := int64(len(data))
	c.cache.SetWithTTL(c.k(key), item, cost, c.ttl(ttl))
	c.cache.Wait()
	return nil
}

// Get 获取缓存
func (c *RistrettoCache) Get(_ context.Context, key string) ([]byte, error) {
	val, ok := c.cache.Get(c.k(key))
	if !ok {
		return nil, &ErrCacheMiss{Key: key}
	}
	item, ok := val.(*localItem)
	if !ok || item.expired() {
		c.cache.Del(c.k(key))
		return nil, &ErrCacheMiss{Key: key}
	}
	return item.Value, nil
}

// Delete 删除缓存
func (c *RistrettoCache) Delete(_ context.Context, key string) error {
	c.cache.Del(c.k(key))
	return nil
}

// Exists 判断 key 是否存在
func (c *RistrettoCache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.Get(ctx, key)
	if IsMiss(err) {
		return false, nil
	}
	return err == nil, err
}

// SetNX 仅在 key 不存在时设置
func (c *RistrettoCache) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	exists, _ := c.Exists(ctx, key)
	if exists {
		return false, nil
	}
	return true, c.Set(ctx, key, value, ttl)
}

// Expire 重置过期时间（ristretto 无直接 API，通过 get+set 模拟）
func (c *RistrettoCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	item := &localItem{
		Value:     data,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.cache.SetWithTTL(c.k(key), item, int64(len(data)), ttl)
	return nil
}

// Incr key 自增 1（本地计数器，非持久化）
func (c *RistrettoCache) Incr(_ context.Context, key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[c.k(key)]++
	return c.counters[c.k(key)], nil
}

// IncrBy key 自增 delta
func (c *RistrettoCache) IncrBy(_ context.Context, key string, delta int64) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[c.k(key)] += delta
	return c.counters[c.k(key)], nil
}

// Close 释放资源
func (c *RistrettoCache) Close() error {
	c.cache.Close()
	return nil
}

// Ping 本地缓存始终可用
func (c *RistrettoCache) Ping(_ context.Context) error {
	return nil
}
