package rdb

import (
	"context"
	"github.com/spf13/cast"
	"go-ssip/app/common/tools"
	g "go-ssip/app/service/rpc/msg/global"
)

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

func GetUserStatus(ctx context.Context, u int64) (string, error) {
	//exist, err := g.Rdb.Exists(ctx, cast.ToString(u)).Result()
	//if err != nil {
	//	return "", err
	//}
	//if exist == 1 {
	//	return "", nil
	//}
	res, err := g.Rdb.Get(ctx, cast.ToString(u)).Result()
	return res, err
}
