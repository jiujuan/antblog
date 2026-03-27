// Package crypto 提供密码加密与验证，使用 bcrypt 算法。
// 同时提供通用 HMAC-SHA256 签名工具。
package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost bcrypt 默认加密强度（10 适合生产环境）
	DefaultCost = bcrypt.DefaultCost
	// MinCost 最低加密强度（测试环境使用，加快速度）
	MinCost = bcrypt.MinCost
)

// HashPassword 对密码进行 bcrypt 哈希
// cost 建议生产环境用 DefaultCost(10)，测试环境用 MinCost(4)
func HashPassword(password string, cost ...int) (string, error) {
	c := DefaultCost
	if len(cost) > 0 && cost[0] >= MinCost {
		c = cost[0]
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), c)
	if err != nil {
		return "", fmt.Errorf("crypto: hash password: %w", err)
	}
	return string(hashed), nil
}

// CheckPassword 验证密码与哈希是否匹配
// 返回 true 表示密码正确
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ─── HMAC-SHA256 ─────────────────────────────────────────────────────────────

// HmacSHA256 使用 HMAC-SHA256 对数据签名，返回十六进制字符串
func HmacSHA256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 对数据进行 SHA256 哈希，返回十六进制字符串
func SHA256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256Bytes 对字节切片进行 SHA256 哈希，返回十六进制字符串
func SHA256Bytes(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// ─── 随机工具 ────────────────────────────────────────────────────────────────

// RandomHex 生成指定字节长度的随机十六进制字符串（长度 = n*2）
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("crypto: generate random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// MustRandomHex 生成随机十六进制字符串，失败时 panic
func MustRandomHex(n int) string {
	s, err := RandomHex(n)
	if err != nil {
		panic(err)
	}
	return s
}
