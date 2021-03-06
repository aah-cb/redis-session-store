package redis

import "time"

const (
	// DefaultRedisNetwork the redis network option, "tcp"
	DefaultRedisNetwork = "tcp"
	// DefaultRedisAddr the redis address option, "127.0.0.1:6379"
	DefaultRedisAddr = "127.0.0.1:6379"
	// DefaultRedisIdleTimeout the redis idle timeout option, time.Duration(30) * time.Minute
	DefaultRedisIdleTimeout = time.Duration(30) * time.Minute
	// DefaultRedisMaxAgeSeconds the redis storage last parameter (SETEX), 31556926.0 (1 year)
	DefaultRedisMaxAgeSeconds = 31556926.0 //1 year
)

type Config struct {
	// Network "tcp"
	Network string
	// Addr "127.0.0.1:6379"
	Addr string
	// Password string .If no password then no 'AUTH'. Default ""
	Password string
	// If Database is empty "" then no 'SELECT'. Default ""
	Database string
	// MaxIdle 0 no limit
	MaxIdle int
	// MaxActive 0 no limit
	MaxActive int
	// IdleTimeout  time.Duration(5) * time.Minute
	IdleTimeout time.Duration
	// Prefix "myprefix-for-this-website". Default ""
	Prefix string
	// MaxAgeSeconds how much long the redis should keep the session in seconds. Default 31556926.0 (1 year)
	MaxAgeSeconds int
}
