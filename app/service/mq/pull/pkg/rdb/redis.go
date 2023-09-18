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
