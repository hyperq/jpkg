package cache

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RC struct {
	rc *redis.Client
	mu sync.Mutex
}

func New(option *redis.Options) (rc *RC, err error) {
	redisclient := redis.NewClient(option)
	_, err = redisclient.Ping(ctx).Result()
	if err != nil {
		return
	}
	rc = &RC{rc: redisclient, mu: sync.Mutex{}}
	return
}

func (r *RC) SET(key string, value interface{}, expiration int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rc.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

func (r *RC) GET(key string) (res string, err error) {
	return r.rc.Get(ctx, key).Result()
}

func (r *RC) DEL(key ...string) error {
	return r.rc.Del(ctx, key...).Err()
}

func (r *RC) KEYS(key string) ([]string, error) {
	return r.rc.Keys(ctx, key).Result()
}

func (r *RC) INCR(key string) error {
	return r.rc.Incr(ctx, key).Err()
}

func (r *RC) DECR(key string) error {
	return r.rc.Decr(ctx, key).Err()
}
