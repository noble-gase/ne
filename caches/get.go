package caches

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func Get[T any](ctx context.Context, cli redis.UniversalClient, key string, fn func(ctx context.Context) (T, error), ttl time.Duration) (T, error) {
	var ret T

	if cli == nil {
		return ret, ErrClientNil
	}

	str, err := cli.Get(ctx, key).Result()
	if err == nil {
		_err := json.Unmarshal([]byte(str), &ret)
		return ret, _err
	}
	if !errors.Is(err, redis.Nil) {
		return ret, err
	}

	// 缓存未命中
	data, err, _ := sf.Do(key, func() (any, error) {
		// 调用fn获取数据
		data, _err := fn(ctx)
		if _err != nil {
			sf.Forget(key)
			return nil, _err
		}
		// 缓存数据
		if b, _err := json.Marshal(data); _err == nil {
			cli.Set(ctx, key, string(b), ttl)
		}
		return data, nil
	})
	if err != nil {
		return ret, err
	}
	return data.(T), nil
}
