package utils

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	reSpaces     = regexp.MustCompile(`[\s\p{Zs}]+`)
	reInvalidSlug = regexp.MustCompile(`[^a-z0-9\-]`)
	reMultiDash  = regexp.MustCompile(`-{2,}`)
)

// Slugify 将任意字符串转为 URL 友好的 slug
// 示例: "Hello World 文章" → "hello-world"
func Slugify(s string) string {
	// 1. 转小写
	s = strings.ToLower(s)

	// 2. Unicode 标准化（NFD 分解，去掉组合字符）
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn = 非间距标记（音调符号）
	}), norm.NFC)
	result, _, err := transform.String(t, s)
	if err == nil {
		s = result
	}

	// 3. 将空白替换为 -
	s = reSpaces.ReplaceAllString(s, "-")

	// 4. 去除非 a-z0-9- 字符（中文、特殊符号等）
	s = reInvalidSlug.ReplaceAllString(s, "")

	// 5. 合并多个 -
	s = reMultiDash.ReplaceAllString(s, "-")

	// 6. 去除首尾 -
	s = strings.Trim(s, "-")

	return s
}

// SlugifyWithFallback 生成 slug，若结果为空则使用时间戳作为后备
func SlugifyWithFallback(s string) string {
	slug := Slugify(s)
	if slug == "" {
		return time.Now().Format("20060102150405")
	}
	return slug
}

// SlugifyWithSuffix 生成 slug 并追加数字后缀（解决重复问题）
// suffix=0 时不加后缀，suffix>0 时加 "-n"
func SlugifyWithSuffix(s string, suffix int) string {
	slug := SlugifyWithFallback(s)
	if suffix <= 0 {
		return slug
	}
	return slug + "-" + intToStr(suffix)
}

// IsValidSlug 校验 slug 格式
func IsValidSlug(s string) bool {
	if s == "" || len(s) > 200 {
		return false
	}
	return !reInvalidSlug.MatchString(s) && !strings.HasPrefix(s, "-") && !strings.HasSuffix(s, "-")
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 10)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
