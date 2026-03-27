// Package storage 本地文件存储驱动，实现 domain/media.IStorageDriver 接口。
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	domain "antblog/internal/domain/media"
	"antblog/pkg/config"
)

// LocalStorage 本地磁盘存储实现
type LocalStorage struct {
	localPath string // 文件根目录（绝对路径），如 /app/uploads
	baseURL   string // 对外访问 URL 前缀，如 http://localhost:8080/uploads
}

// NewLocalStorage 创建本地存储驱动（fx provider）
func NewLocalStorage(cfg *config.Config) (domain.IStorageDriver, error) {
	absPath, err := filepath.Abs(cfg.Upload.LocalPath)
	if err != nil {
		return nil, fmt.Errorf("storage: resolve local_path: %w", err)
	}

	// 确保根目录存在
	if err = os.MkdirAll(absPath, 0o755); err != nil {
		return nil, fmt.Errorf("storage: mkdir %s: %w", absPath, err)
	}

	baseURL := strings.TrimRight(cfg.Upload.BaseURL, "/")
	return &LocalStorage{localPath: absPath, baseURL: baseURL}, nil
}

// Save 将文件字节写入 localPath/YYYY/MM/<uuid>.<ext>，返回相对路径和访问 URL
func (s *LocalStorage) Save(originalName string, data []byte) (storagePath, url string, err error) {
	ext := strings.ToLower(filepath.Ext(originalName))
	if ext == "" {
		ext = ".bin"
	}

	// 按年月分目录存储（避免单目录文件过多）
	now := time.Now()
	subDir := fmt.Sprintf("%04d/%02d", now.Year(), now.Month())
	absDir := filepath.Join(s.localPath, subDir)
	if err = os.MkdirAll(absDir, 0o755); err != nil {
		return "", "", fmt.Errorf("storage: mkdir %s: %w", absDir, err)
	}

	// 用 UUID 生成不可猜测的文件名，防止枚举
	filename := uuid.NewString() + ext
	absFile := filepath.Join(absDir, filename)

	if err = os.WriteFile(absFile, data, 0o644); err != nil {
		return "", "", fmt.Errorf("storage: write file %s: %w", absFile, err)
	}

	// storagePath 使用 Unix 格式相对路径（数据库存储，跨平台一致）
	relativePath := subDir + "/" + filename
	accessURL := s.baseURL + "/" + relativePath

	return relativePath, accessURL, nil
}

// Delete 按相对路径删除物理文件
func (s *LocalStorage) Delete(storagePath string) error {
	// 安全校验：防止路径穿越
	clean := filepath.Clean(storagePath)
	if strings.Contains(clean, "..") {
		return fmt.Errorf("storage: invalid path %q", storagePath)
	}

	absFile := filepath.Join(s.localPath, clean)

	// 确认文件在根目录内
	if !strings.HasPrefix(absFile, s.localPath) {
		return fmt.Errorf("storage: path %q escapes storage root", storagePath)
	}

	if err := os.Remove(absFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage: remove %s: %w", absFile, err)
	}
	return nil
}
