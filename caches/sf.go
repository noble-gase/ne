package caches

import (
	"context"

	"github.com/noble-gase/ne/helper"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

// OmitEmpty 不缓存数据
const OmitEmpty = helper.NilError("caches: omitempty")

func Del(ctx context.Context, uc redis.UniversalClient, key string) error {
	sf.Forget(key)
	return uc.Del(ctx, key).Err()
}

func HDel(ctx context.Context, uc redis.UniversalClient, key, field string) error {
	sf.Forget(key + ":" + field)
	return uc.HDel(ctx, key, field).Err()
}
