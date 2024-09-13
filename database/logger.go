package database

import (
	"context"
	"github.com/alioth-center/infrastructure/logger"
	glog "gorm.io/gorm/logger"
	"time"
)

type DBLogger struct {
	log logger.Logger
}

func NewDBLogger(log logger.Logger) *DBLogger {
	return &DBLogger{log: log}
}

func (dl *DBLogger) LogMode(level glog.LogLevel) glog.Interface {
	return dl
}

func (dl *DBLogger) Info(ctx context.Context, s string, i ...interface{}) {
	dl.log.Infof(logger.NewFields(ctx), s, i...)
}

func (dl *DBLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	dl.log.Warnf(logger.NewFields(ctx), s, i...)
}

func (dl *DBLogger) Error(ctx context.Context, s string, i ...interface{}) {
	dl.log.Errorf(logger.NewFields(ctx), s, i...)
}

func (dl *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	if err != nil {
		dl.log.Errorf(logger.NewFields(ctx), "trace error: %v, sql: %s, rows affected: %d", err, sql, rows)
		return
	}

	dl.log.Infof(logger.NewFields(ctx), "trace sql: %s, rows affected: %d", sql, rows)
}
