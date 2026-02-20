package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/cyc1ones/go-kit/db/transaction"
	gormtx "github.com/cyc1ones/go-kit/db/transaction/gorm"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewTransaction,
	NewDiscovery,
	NewRegistrar,
	NewGreeterRepo,
)

var _ biz.Transaction = (*Data)(nil)

// Data .
type Data struct {
	conf *conf.Data
	db   *gorm.DB
	rdb  *redis.Client
	tx   transaction.Transaction
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	gdb := newGormDB(c, logger)
	rdb := newRedisClient(c, logger)

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if sqlDB, _ := gdb.DB(); sqlDB != nil {
			sqlDB.Close()
		}
		rdb.Close()
	}

	data := &Data{
		conf: c,
		db:   gdb,
		rdb:  rdb,
		tx:   gormtx.New(gdb),
	}

	if err := data.migrate(); err != nil {
		panic(err)
	}
	return data, cleanup, nil
}

func NewTransaction(data *Data) biz.Transaction {
	return data
}

// Exec implements biz.Transaction.
func (d *Data) Exec(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.tx.Exec(ctx, fn)
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	return d.tx.DBFromContext(ctx).(*gorm.DB)
}
