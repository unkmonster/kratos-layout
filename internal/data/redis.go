package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func newRedisClient(conf *conf.Data, logger log.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Network:      conf.Redis.Network,
		Addr:         conf.Redis.Addr,
		Username:     conf.Redis.Username,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.NewHelper(logger).Fatalf("failed to connect to Redis: %v", err)
		return nil
	}

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.NewHelper(logger).Fatalf("failed to instrument Redis tracing:", err)
	}
	return rdb
}
