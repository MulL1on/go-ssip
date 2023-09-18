package rdb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"go-ssip/app/common/tools"
)

type RedisManager struct {
	rdb *redis.Client
}

func NewRedisManager(rdb *redis.Client) *RedisManager {
	return &RedisManager{
		rdb: rdb,
	}
}

func (rm *RedisManager) GetUserStatus(ctx context.Context, u int64) (string, error) {
	res, err := rm.rdb.Get(ctx, cast.ToString(u)).Result()
	return res, err
}

func (rm *RedisManager) GetAndUpdateSeq(ctx context.Context, u int64) (int, error) {
	//get lock
	locker := tools.NewLocker(rm.rdb)
	l := locker.GetLock("u")
	err := l.Lock(ctx)
	if err != nil {
		return 0, err
	}
	ok, _ := rm.rdb.HExists(ctx, "seq", cast.ToString(u)).Result()
	if !ok {
		err = rm.rdb.HSet(ctx, "seq", u, 0).Err()
		if err != nil {
			return 0, err
		}
	}

	res, err := rm.rdb.HIncrBy(ctx, "seq", cast.ToString(u), 1).Result()
	if err != nil {
		return 0, err
	}
	l.Unlock(ctx)

	return cast.ToInt(res), nil
}
