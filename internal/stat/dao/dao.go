package dao

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/zqkgo/goim-enhanced/internal/stat/conf"
)

type Dao struct {
	redis *redis.Pool
}

func NewDao(c *conf.Config) *Dao {
	return &Dao{
		redis: newRedis(c.Redis),
	}
}

func newRedis(c *conf.Redis) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     c.Idle,
		MaxActive:   c.Active,
		IdleTimeout: time.Duration(c.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(c.Network, c.Addr,
				redis.DialConnectTimeout(time.Duration(c.DialTimeout)),
				redis.DialReadTimeout(time.Duration(c.ReadTimeout)),
				redis.DialWriteTimeout(time.Duration(c.WriteTimeout)),
				redis.DialPassword(c.Auth),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}
