package data

import (
	"database/sql"
	"fmt"

	"github.com/go-kratos/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func makeDsn(c *conf.Data_Database) string {
	var auth string

	if c.Username != "" {
		auth = c.Username + ":" + c.Password
	}

	dsn := fmt.Sprintf("%s@%s(%s)/%s", auth, c.Protocol, c.Addr, c.Schema)
	if c.Params != "" {
		dsn = dsn + "?" + c.Params
	}

	//log.Infof("dsn: %s", dsn)
	return dsn
}

func newSqlDB(c *conf.Data, logger log.Logger) *sql.DB {
	dsn := makeDsn(c.Database)

	db, err := otelsql.Open(c.Database.Driver, dsn,
		otelsql.WithDBName(c.Database.Schema),
		otelsql.WithAttributes(semconv.DBSystemKey.String(c.Database.Driver)),
	)
	if err != nil {
		log.NewHelper(logger).Fatalf("failed to initialize sql.DB: %v", err)
	}
	// TODO 连接池
	return db
}

func newGormDB(c *conf.Data, logger log.Logger) *gorm.DB {
	db := newSqlDB(c, logger)

	var err error
	var gdb *gorm.DB
	if c.Database.Driver == "mysql" {
		gdb, err = gorm.Open(mysql.New(mysql.Config{
			Conn: db,
		}), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			TranslateError:         true,
		})
		if err != nil {
			log.NewHelper(logger).Fatalf("failed to initialize gorm.db: %v", err)
		}
	} else {
		err = fmt.Errorf("unsupported driver: %s", c.Database.Driver)
	}

	if err := gdb.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		log.NewHelper(logger).Fatalf("failed to install otel plugin for GORM: %v", err)
	}
	return gdb
}
