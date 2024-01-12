package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alioth-center/infrastructure/logger"
	"gorm.io/gorm"
	"time"
)

type BaseDatabaseImplement struct {
	Db      *gorm.DB
	Logger  logger.Logger
	Timeout time.Duration
	randCmd string
	driver  string
}

type BaseExtMethodGroup struct {
	core *BaseDatabaseImplement
}

func (b *BaseExtMethodGroup) DriverName() string {
	return b.core.driver
}

func (b *BaseExtMethodGroup) GetGorm() *gorm.DB {
	return b.core.Db
}

func (b *BaseExtMethodGroup) Exec(cmd func(db *gorm.DB) *gorm.DB) error {
	return b.core.exec(cmd)
}

func (b *BaseExtMethodGroup) ExecCtx(ctx context.Context, cmd func(db *gorm.DB) *gorm.DB) error {
	return b.core.exec(cmd, ctx)
}

func (s *BaseDatabaseImplement) Init(_ Options) error {
	panic("init function must be implemented")
}

func (s *BaseDatabaseImplement) SetRandCommand(command string) {
	if command == "" {
		s.randCmd = "rand()"
	}

	s.randCmd = command
}

func (s *BaseDatabaseImplement) SetDriverName(name string) {
	s.driver = name
}

func (s *BaseDatabaseImplement) ParseLoggerOptions(opts Options) {
	loggerOpts := logger.Options{
		LogLevel:     logger.LevelInfo,
		StdoutWriter: logger.ConsoleWriter(),
		StderrWriter: logger.ConsoleErrorWriter(),
	}
	if opts.DebugLog {
		loggerOpts.LogLevel = logger.LevelDebug
	}
	if opts.Stdout != "" {
		stdout, fwe := logger.FileWriter(opts.Stdout)
		if fwe == nil {
			loggerOpts.StdoutWriter = stdout
		}
	}
	if opts.Stderr != "" {
		stderr, fwe := logger.FileWriter(opts.Stderr)
		if fwe == nil {
			loggerOpts.StderrWriter = stderr
		}
	}

	s.Logger = logger.NewLoggerWithOptions(loggerOpts)
}

func (s *BaseDatabaseImplement) ParseDatabaseOptions(db *sql.DB, opts Options) {
	if opts.MaxIdle > 0 {
		db.SetMaxIdleConns(opts.MaxIdle)
	}
	if opts.MaxOpen > 0 {
		db.SetMaxOpenConns(opts.MaxOpen)
	}
	if opts.MaxLife > 0 {
		db.SetConnMaxLifetime(opts.MaxLife)
	}
	if opts.Timeout > 0 {
		s.Timeout = opts.Timeout
	}
}

func (s *BaseDatabaseImplement) Migrate(models ...any) error {
	s.Logger.Info(logger.NewFields().WithMessage("start migrate BaseDatabaseImplement database"))
	migrateErr := s.Db.AutoMigrate(models...)
	if migrateErr != nil {
		s.Logger.Error(logger.NewFields().WithMessage(migrateErr.Error()).WithData(models))
	}
	return migrateErr
}

func (s *BaseDatabaseImplement) exec(command func(tx *gorm.DB) *gorm.DB, ctx ...context.Context) error {
	var db *gorm.DB
	var trace context.Context
	if len(ctx) != 1 {
		db = s.Db.Session(&gorm.Session{})
		trace = nil
	} else if ctx[0] == nil {
		db = s.Db.Session(&gorm.Session{})
		trace = nil
	} else if s.Timeout > 0 {
		timeout, cancel := context.WithTimeout(ctx[0], s.Timeout)
		trace = timeout
		defer cancel()
		db = s.Db.Session(&gorm.Session{}).WithContext(trace)
	} else {
		trace = ctx[0]
		db = s.Db.Session(&gorm.Session{}).WithContext(trace)
	}

	sqlCommand := s.Db.ToSQL(command)
	s.Logger.Debug(logger.NewFields(trace).WithMessage("sql executed").WithData(sqlCommand))

	if err := db.Transaction(func(tx *gorm.DB) error { return command(tx).Error }); err != nil {
		err = fmt.Errorf("sql execution failed: %w", err)
		s.Logger.Error(logger.NewFields(trace).WithMessage(err.Error()).WithData(sqlCommand))
		return err
	}

	return nil
}

func (s *BaseDatabaseImplement) Has(table, query string, args ...any) (exist bool, err error) {
	var count int64
	err = s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(table).Where(query, args...).Limit(100).Count(&count)
	})
	return count > 0, err
}

func (s *BaseDatabaseImplement) Count(table, query string, args ...any) (count int64, err error) {
	err = s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(table).Where(query, args...).Count(&count)
	})
	return count, err
}

func (s *BaseDatabaseImplement) GetOne(receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.First(receiver, append([]any{query}, args...)...)
	})
}

func (s *BaseDatabaseImplement) GetAll(receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(receiver, append([]any{query}, args...)...)
	})
}

func (s *BaseDatabaseImplement) GetPage(receiver any, offset, limit int, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(offset*limit).Limit(limit).Find(receiver, append([]any{query}, args...)...)
	})
}

func (s *BaseDatabaseImplement) PickOne(receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Order(s.randCmd).Take(receiver, append([]any{query}, args...)...)
	})
}

func (s *BaseDatabaseImplement) PickAll(receiver any, length int, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Order(s.randCmd).Limit(length).Find(receiver, append([]any{query}, args...)...)
	})
}

func (s *BaseDatabaseImplement) InsertOne(data any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(data)
	})
}

func (s *BaseDatabaseImplement) InsertAll(data any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.CreateInBatches(data, 100)
	})
}

func (s *BaseDatabaseImplement) UpdateOne(data any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(data).Where(query, args...).Limit(1).Updates(data)
	})
}

func (s *BaseDatabaseImplement) UpdateAll(data any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(data).Where(query, args...).Updates(data)
	})
}

func (s *BaseDatabaseImplement) DeleteOne(query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Delete(query, args...).Limit(1)
	})
}

func (s *BaseDatabaseImplement) DeleteAll(query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Delete(query, args...)
	})
}

func (s *BaseDatabaseImplement) ExecRaw(sql string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Exec(sql, args...)
	})
}

func (s *BaseDatabaseImplement) QueryRaw(receiver any, sql string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(sql, args...).Scan(receiver)
	})
}

func (s *BaseDatabaseImplement) HasWithCtx(ctx context.Context, table, query string, args ...any) (exist bool, err error) {
	var count int64
	err = s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(table).Where(query, args...).Limit(100).Count(&count)
	}, ctx)
	return count > 0, err
}

func (s *BaseDatabaseImplement) CountWithCtx(ctx context.Context, table, query string, args ...any) (count int64, err error) {
	err = s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.WithContext(ctx).Table(table).Where(query, args...).Count(&count)
	}, ctx)
	return count, err
}

func (s *BaseDatabaseImplement) GetOneWithCtx(ctx context.Context, receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.First(receiver, append([]any{query}, args...)...)
	}, ctx)
}

func (s *BaseDatabaseImplement) GetAllWithCtx(ctx context.Context, receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(receiver, append([]any{query}, args...)...)
	}, ctx)
}

func (s *BaseDatabaseImplement) PickOneWithCtx(ctx context.Context, receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Order(s.randCmd).Take(receiver, append([]any{query}, args...)...)
	}, ctx)
}

func (s *BaseDatabaseImplement) PickAllWithCtx(ctx context.Context, receiver any, length int, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Order(s.randCmd).Limit(length).Find(receiver, append([]any{query}, args...)...)
	}, ctx)
}

func (s *BaseDatabaseImplement) InsertOneWithCtx(ctx context.Context, data any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(data)
	}, ctx)
}

func (s *BaseDatabaseImplement) InsertAllWithCtx(ctx context.Context, data any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.CreateInBatches(data, 100)
	}, ctx)
}

func (s *BaseDatabaseImplement) UpdateOneWithCtx(ctx context.Context, data any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(data).Where(query, args...).Limit(1).Updates(data)
	}, ctx)
}

func (s *BaseDatabaseImplement) UpdateAllWithCtx(ctx context.Context, data any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(data).Where(query, args...).Updates(data)
	}, ctx)
}

func (s *BaseDatabaseImplement) DeleteOneWithCtx(ctx context.Context, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Delete(query, args...).Limit(1)
	}, ctx)
}

func (s *BaseDatabaseImplement) DeleteAllWithCtx(ctx context.Context, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Delete(query, args...)
	}, ctx)
}

func (s *BaseDatabaseImplement) ExecRawWithCtx(ctx context.Context, sql string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Exec(sql, args...)
	}, ctx)
}

func (s *BaseDatabaseImplement) QueryRawWithCtx(ctx context.Context, receiver any, sql string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(sql, args...).Scan(receiver)
	}, ctx)
}

func (s *BaseDatabaseImplement) ExtMethods() ExtMethods {
	return &BaseExtMethodGroup{core: s}
}
