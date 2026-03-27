package jwt

import "time"

// Option JWT 构建选项（函数选项模式）
type Option func(*options)

type options struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
	signingMethod string // HS256 | HS384 | HS512
}

func defaultOptions() *options {
	return &options{
		secret:        "antblog-secret-change-me",
		accessExpiry:  2 * time.Hour,
		refreshExpiry: 7 * 24 * time.Hour,
		issuer:        "antblog",
		signingMethod: "HS256",
	}
}

// WithSecret 设置签名密钥
func WithSecret(secret string) Option {
	return func(o *options) { o.secret = secret }
}

// WithAccessExpiry 设置 Access Token 有效期
func WithAccessExpiry(d time.Duration) Option {
	return func(o *options) { o.accessExpiry = d }
}

// WithRefreshExpiry 设置 Refresh Token 有效期
func WithRefreshExpiry(d time.Duration) Option {
	return func(o *options) { o.refreshExpiry = d }
}

// WithIssuer 设置 Token 签发方
func WithIssuer(issuer string) Option {
	return func(o *options) { o.issuer = issuer }
}

// WithSigningMethod 设置签名算法（HS256|HS384|HS512）
func WithSigningMethod(method string) Option {
	return func(o *options) { o.signingMethod = method }
}
