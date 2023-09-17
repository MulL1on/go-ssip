package tools

import (
	"context"
	"errors"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-redis/redis/v8"
)

const (
	// 过期时间
	ttl = time.Second * 30
	// 重置过期时间间隔
	resetTTLInterval = ttl / 3
	// 重新获取锁间隔
	tryLockInterval = time.Second
	// 解锁脚本
	unlockScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end`
)

var (
	ErrLockFailed = errors.New("lock failed")
	ErrTimeout    = errors.New("timeout")
)

type Locker struct {
	client          *redis.Client
	script          *redis.Script
	ttl             time.Duration
	tryLockInterval time.Duration
}

func NewLocker(client *redis.Client) *Locker {
	return &Locker{
		client:          client,
		script:          redis.NewScript(unlockScript),
		ttl:             ttl,
		tryLockInterval: tryLockInterval,
	}
}

type Lock struct {
	client          *redis.Client
	script          *redis.Script
	resource        string
	randomValue     string
	watchDog        chan struct{}
	ttl             time.Duration
	tryLockInterval time.Duration
}

func (l *Locker) GetLock(resource string) *Lock {
	// try to get a getLock
	return &Lock{
		client:          l.client,
		script:          l.script,
		resource:        resource,
		randomValue:     gofakeit.UUID(),
		watchDog:        make(chan struct{}),
		ttl:             l.ttl,
		tryLockInterval: l.tryLockInterval,
	}
}

func (l *Lock) TryLock(ctx context.Context) error {
	success, err := l.client.SetNX(ctx, l.resource, l.randomValue, l.ttl).Result()
	if err != nil {
		return err
	}

	if !success {
		return ErrLockFailed
	}

	go l.startWatchDog()
	return nil
}

func (l *Lock) startWatchDog() {
	ticker := time.NewTicker(l.ttl / 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// extent lock expire time
			ctx, cancel := context.WithTimeout(context.Background(), l.ttl/3*2)
			ok, err := l.client.Expire(ctx, l.resource, l.ttl).Result()
			cancel()
			if err != nil || !ok {
				return
			}
		case <-l.watchDog:
			return
		}
	}
}

func (l *Lock) Lock(ctx context.Context) error {
	err := l.TryLock(ctx)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrLockFailed) {
		return err
	}

	// lock failed retry
	ticker := time.NewTicker(l.tryLockInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ErrTimeout
		case <-ticker.C:
			// retry to get lock
			err = l.TryLock(ctx)
			if err == nil {
				return nil
			}
			if !errors.Is(err, ErrLockFailed) {
				return err
			}
		}
	}
}

func (l *Lock) Unlock(ctx context.Context) error {
	err := l.script.Run(ctx, l.client, []string{l.resource}, l.randomValue).Err()
	close(l.watchDog)
	return err
}
