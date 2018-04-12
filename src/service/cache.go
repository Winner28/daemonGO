package service

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// Cache struct
type Cache struct {
	pool      *redis.Pool
	redisHost string
}

var cache *Cache

const redisHost = ":6379"

func init() {
	cache = &Cache{
		pool:      getPool(redisHost),
		redisHost: redisHost,
	}
}

func getPool(server string) (pool *redis.Pool) {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, _ time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}

func (c *Cache) getPool() *redis.Pool {
	return c.pool
}

func (c *Cache) setKeyValue(key, value string) bool {
	conn := c.getPool().Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return false
	}

	return true
}

func (c *Cache) getKeyValue(key string) (string, bool) {
	conn := c.getPool().Get()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return data, false
	}

	return data, true
}
