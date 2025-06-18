package caches

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func HGet[T any](ctx context.Context, cli redis.UniversalClient, key, field string, fn func(ctx context.Context) (T, error), ttl time.Duration) (T, error) {
	var ret T

	if cli == nil {
		return ret, ErrClientNil
	}

	str, err := cli.HGet(ctx, key, field).Result()
	if err == nil {
		_err := json.Unmarshal([]byte(str), &ret)
		return ret, _err
	}
	if !errors.Is(err, redis.Nil) {
		return ret, err
	}

	// 缓存未命中
	sfKey := key + ":" + field
	data, err, _ := sf.Do(sfKey, func() (any, error) {
		// 调用fn获取数据
		data, _err := fn(ctx)
		if _err != nil {
			sf.Forget(sfKey)
			return nil, _err
		}
		// 缓存数据
		if b, _err := json.Marshal(data); _err == nil {
			cli.HSet(ctx, key, field, string(b))
			if ttl > 0 && cli.TTL(ctx, key).Val() == -1 {
				cli.Expire(ctx, key, ttl)
			}
		}
		return data, nil
	})
	if err != nil {
		return ret, err
	}
	return data.(T), nil
}
