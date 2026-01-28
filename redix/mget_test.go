package redix

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestMGetMap(t *testing.T) {
	ctx := context.Background()

	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
	})
	uc.MSet(ctx, "foo", `{"id":1,"name":"foo"}`, "bar", `{"id":2,"name":"bar"}`, "hello", `{"id":3,"name":"hello"}`)
	defer uc.Del(ctx, "foo", "bar", "hello")

	type Demo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	ret, err := MGetMap[Demo](ctx, uc, []string{"foo", "bar", "hello", "none"})
	assert.Nil(t, err)
	t.Log(ret)
}

func TestMGetStringMap(t *testing.T) {
	ctx := context.Background()

	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
	})
	uc.MSet(ctx, "foo", "test-foo", "bar", "test-bar", "hello", "test-hello")
	defer uc.Del(ctx, "foo", "bar", "hello")

	ret, err := MGetStringMap(ctx, uc, []string{"foo", "bar", "hello", "none"})
	assert.Nil(t, err)
	t.Log(ret)
}
