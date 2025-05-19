package rdb

import (
	"github.com/EduardMikhrin/forecaster/internal/data"
	"github.com/pkg/errors"
	"github.com/redis/go-redis"
	"time"
)

type cacheQ struct {
	r       *redis.Client
	timeout time.Duration
}

func (c cacheQ) SetCode(key string, value string) error {
	return errors.Wrap(c.r.Set(key, value, c.timeout).Err(), "failed to set code to cache")

}

func (c cacheQ) GetCode(key string) (string, error) {
	val, err := c.r.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", errors.Wrap(err, "failed to get code")
	}
	return val, nil
}

func (c cacheQ) DelCode(key string) error {
	if err := c.r.Del(key).Err(); err != nil {
		return errors.Wrap(err, "failed to delete code")
	}

	return nil
}

func NewCacheQ(addr, password string, db uint, timeout int) data.CacheQ {
	return &cacheQ{
		r: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       int(db),
		}),
		timeout: time.Duration(timeout) * time.Second,
	}
}
