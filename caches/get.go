package caches

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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
		if _err := json.Unmarshal([]byte(str), &ret); _err != nil {
			return ret, fmt.Errorf("unmarshal(%s): %w", str, _err)
		}
		return ret, nil
	}
	if !errors.Is(err, redis.Nil) {
		return ret, err
	}

	// 缓存未命中
	data, err, _ := sf.Do(key, func() (any, error) {
		// 调用fn获取数据
		data, _err := fn(ctx)
		if _err != nil {
			if errors.Is(_err, OmitEmpty) {
				return data, nil
			}
			sf.Forget(key)
			return nil, _err
		}

		// 缓存数据
		b, _err := json.Marshal(data)
		if _err != nil {
			slog.ErrorContext(ctx, "[caches:Get] marshal data failed", slog.String("key", key), slog.String("error", _err.Error()))
			return data, nil
		}
		if _err = cli.Set(ctx, key, string(b), ttl).Err(); _err != nil {
			slog.ErrorContext(ctx, "[caches:Get] set data failed", slog.String("key", key), slog.String("value", string(b)), slog.String("error", _err.Error()))
		}
		return data, nil
	})
	if err != nil {
		return ret, err
	}
	return data.(T), nil
}
