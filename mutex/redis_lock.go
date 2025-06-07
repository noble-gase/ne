package mutex

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrClientNil = errors.New("redis client is nil (forgotten initialize?)")

// ErrLockNil 未获取到锁
var ErrLockNil = errors.New("redlock: lock not acquired")

const script = `
if redis.call('get', KEYS[1]) == ARGV[1] then
	return redis.call('del', KEYS[1])
else
	return 0
end
`

// Mutex 分布式锁
type Mutex interface {
	// Lock 获取锁；未获取到会返回`ErrLockNil`
	Lock(ctx context.Context) error
	// TryLock 尝试获取锁；未获取到会返回`ErrLockNil`
	TryLock(ctx context.Context, attempts int, interval time.Duration) error
	// UnLock 释放锁
	UnLock(ctx context.Context) error
}

// redLock 基于「Redis」实现的分布式锁
type redLock struct {
	cli   redis.UniversalClient
	key   string
	ttl   time.Duration
	token string
}

func (l *redLock) Lock(ctx context.Context) error {
	select {
	case <-ctx.Done(): // timeout or canceled
		return ctx.Err()
	default:
	}

	if err := l.lock(ctx); err != nil {
		return err
	}
	if len(l.token) != 0 {
		return nil
	}
	return ErrLockNil
}

func (l *redLock) TryLock(ctx context.Context, attempts int, interval time.Duration) error {
	for i := 0; i < attempts; i++ {
		select {
		case <-ctx.Done(): // timeout or canceled
			return ctx.Err()
		default:
		}

		// attempt to acquire lock
		if err := l.lock(ctx); err != nil {
			return err
		}
		if len(l.token) != 0 {
			return nil
		}
		time.Sleep(interval)
	}
	return ErrLockNil
}

func (l *redLock) UnLock(ctx context.Context) error {
	if len(l.token) == 0 {
		return nil
	}
	if l.cli == nil {
		return ErrClientNil
	}
	return l.cli.Eval(context.WithoutCancel(ctx), script, []string{l.key}, l.token).Err()
}

func (l *redLock) lock(ctx context.Context) error {
	if l.cli == nil {
		return ErrClientNil
	}

	token := uuid.New().String()

	ok, err := l.cli.SetNX(ctx, l.key, token, l.ttl).Result()
	if err != nil {
		// 尝试GET一次：避免因redis网络错误导致误加锁
		v, _err := l.cli.Get(ctx, l.key).Result()
		if _err != nil {
			if errors.Is(_err, redis.Nil) {
				return err
			}
			return fmt.Errorf("SET-NX: %w; GET: %w", err, _err)
		}
		if v == token {
			l.token = token
		}
		return nil
	}
	if ok {
		l.token = token
	}
	return nil
}

// RedLock 基于Redis实现的分布式锁实例
func RedLock(cli redis.UniversalClient, key string, ttl time.Duration) Mutex {
	mutex := &redLock{
		cli: cli,
		key: key,
		ttl: ttl,
	}
	if mutex.ttl == 0 {
		mutex.ttl = time.Second * 10
	}
	return mutex
}
