package redix

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHGet(t *testing.T) {
	ctx := context.Background()

	uc := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
		DB:    0,
	})
	defer uc.Del(ctx, "hello")

	type Demo struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	ret, err := HGet(ctx, uc, "hello", "foo", func(ctx context.Context) (*Demo, error) {
		t.Log(">> callback ")
		return &Demo{
			ID:   1,
			Name: "hello",
		}, nil
	}, time.Minute)
	assert.Nil(t, err)
	t.Logf("%+v", ret)
	t.Log(uc.HGet(ctx, "hello", "foo").String())
}
