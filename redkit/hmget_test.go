package redkit

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHGetAll(t *testing.T) {
	ctx := context.Background()

	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
	})
	uc.HMSet(ctx, "test", "foo", `{"id":1,"name":"foo"}`, "bar", `{"id":2,"name":"bar"}`, "hello", `{"id":3,"name":"hello"}`)
	defer uc.Del(ctx, "test")

	type Demo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	ret, err := HGetAll[Demo](ctx, uc, "test")
	assert.Nil(t, err)
	t.Log(ret)
}

func TestHMGetMap(t *testing.T) {
	ctx := context.Background()

	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
	})
	uc.HMSet(ctx, "test", "foo", `{"id":1,"name":"foo"}`, "bar", `{"id":2,"name":"bar"}`, "hello", `{"id":3,"name":"hello"}`)
	defer uc.Del(ctx, "test")

	type Demo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	ret, err := HMGetMap[Demo](ctx, uc, "test", []string{"foo", "bar", "hello", "none"})
	assert.Nil(t, err)
	t.Log(ret)
}

func TestHMGetStringMap(t *testing.T) {
	ctx := context.Background()

	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
	})
	uc.HMSet(ctx, "test", "foo", `{"id":1,"name":"foo"}`, "bar", `{"id":2,"name":"bar"}`, "hello", `{"id":3,"name":"hello"}`)
	defer uc.Del(ctx, "test")

	ret, err := HMGetStrMap(ctx, uc, "test", []string{"foo", "bar", "hello", "none"})
	assert.Nil(t, err)
	t.Log(ret)
}
