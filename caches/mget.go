package caches

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func MGetMap(ctx context.Context, cli redis.UniversalClient, keys []string, omitempty bool) (map[string]string, error) {
	values, err := cli.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	if len(values) != len(keys) {
		return nil, errors.New("the number of keys and values mismatch")
	}

	if omitempty {
		return buildMapOmitEmpty(keys, values), nil
	}
	return buildMapOmitReserve(keys, values), nil
}

func buildMapOmitEmpty(keys []string, values []any) map[string]string {
	m := make(map[string]string, len(keys))
	for i, k := range keys {
		if v := values[i]; v != nil {
			if s, ok := v.(string); ok && len(s) != 0 {
				m[k] = s
			}
		}
	}
	return m
}

func buildMapOmitReserve(keys []string, values []any) map[string]string {
	m := make(map[string]string, len(keys))
	for i, k := range keys {
		if v := values[i]; v != nil {
			s, _ := v.(string)
			m[k] = s
		} else {
			m[k] = ""
		}
	}
	return m
}
