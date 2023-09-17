package initialize

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	g "go-ssip/app/service/rpc/msg/global"
	"go.uber.org/zap"
	"time"
)

func InitRdb() {
	cfg := g.ServerConfig.RedisInfo
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		g.Logger.Fatal("connect to redis failed", zap.Error(err))
	}
	g.Rdb = rdb
	g.Logger.Info("initialize redis successfully.")
}
