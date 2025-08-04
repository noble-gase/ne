package caches

import (
	"context"
	"errors"

	"github.com/noble-gase/ne/helper"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

// OmitEmpty 不缓存数据
const OmitEmpty = helper.NilError("caches: omitempty")

var ErrClientNil = errors.New("redis client is nil (forgotten initialize?)")

func Del(ctx context.Context, cli redis.UniversalClient, key string) error {
	if cli == nil {
		return ErrClientNil
	}

	sf.Forget(key)
	return cli.Del(ctx, key).Err()
}

func HDel(ctx context.Context, cli redis.UniversalClient, key, field string) error {
	if cli == nil {
		return ErrClientNil
	}

	sf.Forget(key + ":" + field)
	return cli.HDel(ctx, key, field).Err()
}
