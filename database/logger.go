package database

import (
	"context"
	"time"

	"github.com/alioth-center/infrastructure/logger"
	glog "gorm.io/gorm/logger"
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
		logMessage := map[string]any{"sql": sql, "error": err.Error(), "rows": rows}
		dl.log.Error(logger.NewFields(ctx).WithMessage("tracing sql with error").WithData(logMessage))
		return
	}

	logMessage := map[string]any{"sql": sql, "rows": rows}
	dl.log.Info(logger.NewFields(ctx).WithMessage("tracing sql").WithData(logMessage))
}
