package rdb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
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

func (rm *RedisManager) GetMaxSeq(ctx context.Context, u int64) (int64, error) {
	res := rm.rdb.HGet(ctx, "seq", cast.ToString(u))
	err := res.Err()
	if err != nil {
		return 0, err
	}
	seq, err := res.Int64()
	if err != nil {
		return 0, err
	}
	return seq, nil
}
