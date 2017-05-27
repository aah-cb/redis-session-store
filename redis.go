package redis

import (
	"errors"
	"time"

	"fmt"

	config "aahframework.org/config.v0"
	log "aahframework.org/log.v0"
	"aahframework.org/security.v0/session"
	"github.com/garyburd/redigo/redis"
)

func init() {
	_ = session.AddStore("redis", &RedisStore{})
}

type RedisStore struct {
	Config    *Config
	pool      *redis.Pool
	Connected bool
}

func (r *RedisStore) Init(cfg *config.Config) error {
	r.Config = &Config{
		Network:       DefaultRedisNetwork,
		Addr:          DefaultRedisAddr,
		Password:      "",
		Database:      "",
		MaxIdle:       0,
		MaxActive:     0,
		IdleTimeout:   DefaultRedisIdleTimeout,
		Prefix:        "",
		MaxAgeSeconds: DefaultRedisMaxAgeSeconds,
	}

	r.Config.Network = cfg.StringDefault("security.session.store.redis.network", "tcp")
	r.Config.Addr = cfg.StringDefault("security.session.store.redis.addr", "127.0.0.1:6379")
	r.Config.Password = cfg.StringDefault("security.session.store.redis.password", "")
	r.Config.Database = cfg.StringDefault("security.session.store.redis.database", "")
	r.Config.Prefix = cfg.StringDefault("security.session.store.redis.prefix", "")

	r.Config.MaxIdle = cfg.IntDefault("security.session.store.redis.max_idle", 10)
	r.Config.MaxActive = cfg.IntDefault("security.session.store.redis.max_active", 30)

	r.connect()

	_, err := r.pingPong()
 
	if err != nil {
		return errors.New("Redis Connection error on Connect:" + err.Error())
	}
	return nil
}

func (r *RedisStore) Read(id string) string {

	c := r.pool.Get()
	defer c.Close()
	if err := c.Err(); err != nil {
		log.Errorf("session: redis store - read error: %v", err)
		return ""
	}

	redisVal, err := c.Do("GET", r.Config.Prefix+id)

	if err != nil {
		log.Errorf("session: redis store - read error: %v", err)
		return ""
	}
	if redisVal == nil {
		log.Errorf("session: redis store - Key '%s' doesn't", id)
		return ""
	}

	sVal, err := redis.String(redisVal, err)
	if err != nil {
		return ""
	}

	return sVal

}
func (r *RedisStore) Save(id, value string) (err error) {

	c := r.pool.Get()
	defer c.Close()
	if err = c.Err(); err != nil {
		return
	}
	_, err = c.Do("SETEX", r.Config.Prefix+id, r.Config.MaxAgeSeconds, value)

	return

}
func (r *RedisStore) Delete(id string) error {
	c := r.pool.Get()
	defer c.Close()

	if _, err := c.Do("DEL", r.Config.Prefix+id); err != nil {
		return err
	}
	return nil
}
func (r *RedisStore) IsExists(id string) bool {
	c := r.pool.Get()
	defer c.Close()

	if existed, err := redis.Int(c.Do("EXISTS", r.Config.Prefix+id)); err != nil || existed == 0 {
		return false
	}
	return true

}

func (r *RedisStore) Cleanup(m *session.Manager) {
	sessions, err := r.getAll("session")
	if err != nil {
		log.Error(err)
		return
	}
	cnt := 0
	for sid, sess := range sessions {

		if _, err := m.DecodeToSession(sess); err == session.ErrCookieTimestampIsExpired {
			if err := r.Delete(sid); err != nil {
				log.Error(err)
			} else {
				cnt++
			}
		}
	}

	log.Infof("%v expired session redis cleaned up", cnt)
}

func (r *RedisStore) getAll(key string) (map[string]string, error) {
	c := r.pool.Get()
	defer c.Close()
	if err := c.Err(); err != nil {
		return nil, err
	}

	reply, err := c.Do("HGETALL", r.Config.Prefix+key)

	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, errors.New(fmt.Sprintf("Key '%s' doesn't found", key))
	}

	return redis.StringMap(reply, err)

}

func (r *RedisStore) pingPong() (bool, error) {
	c := r.pool.Get()
	defer c.Close()
	msg, err := c.Do("PING")
	if err != nil || msg == nil {
		return false, err
	}
	return (msg == "PONG"), nil
}

//************************redis connect****************
func (r *RedisStore) connect() {
	c := r.Config

	if c.IdleTimeout <= 0 {
		c.IdleTimeout = DefaultRedisIdleTimeout
	}

	if c.Network == "" {
		c.Network = DefaultRedisNetwork
	}

	if c.Addr == "" {
		c.Addr = DefaultRedisAddr
	}

	if c.MaxAgeSeconds <= 0 {
		c.MaxAgeSeconds = DefaultRedisMaxAgeSeconds
	}

	pool := &redis.Pool{IdleTimeout: DefaultRedisIdleTimeout, MaxIdle: c.MaxIdle, MaxActive: c.MaxActive}
	pool.TestOnBorrow = func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	}

	if c.Database != "" {
		pool.Dial = func() (redis.Conn, error) {
			red, err := dial(c.Network, c.Addr, c.Password)
			if err != nil {
				return nil, err
			}
			if _, err = red.Do("SELECT", c.Database); err != nil {
				red.Close()
				return nil, err
			}
			return red, err
		}
	} else {
		pool.Dial = func() (redis.Conn, error) {
			return dial(c.Network, c.Addr, c.Password)
		}
	}
	r.Connected = true
	r.pool = pool

}
func dial(network string, addr string, pass string) (redis.Conn, error) {
	if network == "" {
		network = DefaultRedisNetwork
	}
	if addr == "" {
		addr = DefaultRedisAddr
	}
	log.Print(addr)
	c, err := redis.Dial(network, addr)
	log.Print("init")
	log.Print(err)
	if err != nil {
		return nil, err
	}
	if pass != "" {
		if _, err = c.Do("AUTH", pass); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}
