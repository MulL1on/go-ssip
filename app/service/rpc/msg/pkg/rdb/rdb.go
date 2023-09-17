package rdb

import (
	"context"
	"github.com/spf13/cast"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/rpc/msg/global"
)

func UpdateMaxSeq(ctx context.Context, u int64) error {
	ok, _ := g.Rdb.HExists(ctx, "seq", cast.ToString(u)).Result()
	if !ok {
		err := g.Rdb.HSet(ctx, "seq", u, 1).Err()
		if err != nil {
			return err
		}
	} else {
		err := g.Rdb.HIncrBy(ctx, "seq", cast.ToString(u), 1).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAndUpdateSeq(ctx context.Context, u int64) (int, error) {
	//get lock
	locker := tools.NewLocker(g.Rdb)
	l := locker.GetLock("u")
	err := l.Lock(ctx)
	if err != nil {
		return 0, err
	}
	ok, _ := g.Rdb.HExists(ctx, "seq", cast.ToString(u)).Result()
	if !ok {
		err = g.Rdb.HSet(ctx, "seq", u, 0).Err()
		if err != nil {
			return 0, err
		}
	}

	res, err := g.Rdb.HIncrBy(ctx, "seq", cast.ToString(u), 1).Result()
	if err != nil {
		return 0, err
	}
	l.Unlock(ctx)

	return cast.ToInt(res), nil
}
