package redkit

import (
	"github.com/noble-gase/ne/helper"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

// Discard 丢弃数据，不缓存
const Discard = helper.NilError("caches: discarded")

var scriptHSet = redis.NewScript(`
redis.call('HSET', KEYS[1], ARGV[1], ARGV[2])
if redis.call('TTL', KEYS[1]) == -1 then
    redis.call('EXPIRE', KEYS[1], ARGV[3])
end
`)

var scriptLock = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`)
