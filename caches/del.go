package caches

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Del(ctx context.Context, uc redis.UniversalClient, key string) error {
	sf.Forget(key)
	return uc.Del(ctx, key).Err()
}

func HDel(ctx context.Context, uc redis.UniversalClient, key, field string) error {
	sf.Forget(key + ":" + field)
	return uc.HDel(ctx, key, field).Err()
}
