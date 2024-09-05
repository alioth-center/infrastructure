package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm/clause"

	"github.com/alioth-center/infrastructure/logger"
	"gorm.io/gorm"
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

func (b *BaseExtMethodGroup) Transaction(cmd func(tx *gorm.DB) error) error {
	return b.core.execTransaction(cmd)
}

func (b *BaseExtMethodGroup) TransactionCtx(ctx context.Context, cmd func(tx *gorm.DB) error) error {
	return b.core.execTransaction(cmd, ctx)
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

func (s *BaseDatabaseImplement) logSql(command func(tx *gorm.DB) *gorm.DB) (sql string) {
	return s.Db.ToSQL(command)
}

func (s *BaseDatabaseImplement) exec(command func(tx *gorm.DB) *gorm.DB, ctx ...context.Context) error {
	var trace context.Context
	var db *gorm.DB
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

	if err := db.Transaction(func(tx *gorm.DB) error { return command(tx).Error }); err != nil {
		sqlCommand := s.logSql(command)
		execErr := NewExecuteSqlError(sqlCommand, err)
		s.Logger.Error(logger.NewFields(trace).WithMessage("sql command execute failed").WithData(execErr))
		return execErr
	}

	return nil
}

func (s *BaseDatabaseImplement) execTransaction(command func(tx *gorm.DB) error, ctx ...context.Context) error {
	var trace context.Context
	var db *gorm.DB
	if len(ctx) != 1 {
		db = s.Db.Session(&gorm.Session{})
	} else if ctx[0] == nil {
		db = s.Db.Session(&gorm.Session{})
	} else if s.Timeout > 0 {
		timeout, cancel := context.WithTimeout(ctx[0], s.Timeout)
		trace = timeout
		defer cancel()
		db = s.Db.Session(&gorm.Session{}).WithContext(trace)
	} else {
		trace = ctx[0]
		db = s.Db.Session(&gorm.Session{}).WithContext(trace)
	}

	if err := db.Transaction(command); err != nil {
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

func (s *BaseDatabaseImplement) SetLogger(logger logger.Logger) {
	s.Logger = logger
}

type BaseDatabaseImplementV2 struct {
	Db *gorm.DB
}

func (v2 *BaseDatabaseImplementV2) GetGormCore(ctx context.Context) *gorm.DB {
	return v2.Db.WithContext(ctx)
}

func (v2 *BaseDatabaseImplementV2) GetDataBySingleCondition(ctx context.Context, receiver any, column string, condition any, needFields ...string) error {
	if column == "" || condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	if len(needFields) == 0 {
		needFields = append(needFields, "*")
	}

	return v2.Db.WithContext(ctx).Model(receiver).Where(column, condition).Select(needFields).Scan(receiver).Error
}

func (v2 *BaseDatabaseImplementV2) GetDataByCustomCondition(ctx context.Context, receiver, condition any, needFields ...string) error {
	if condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	if len(needFields) == 0 {
		needFields = append(needFields, "*")
	}

	return v2.Db.WithContext(ctx).Model(receiver).Where(condition).Select(needFields).Scan(receiver).Error
}

func (v2 *BaseDatabaseImplementV2) ListDataWithPage(ctx context.Context, receiver any, filter any, order string, desc bool, offset, limit int, needFields ...string) error {
	if filter == nil || EmptySlice(filter) {
		return ErrInvalidCondition
	}

	if len(needFields) == 0 {
		needFields = append(needFields, "*")
	}

	return v2.Db.WithContext(ctx).Model(receiver).Where(filter).Order(clause.OrderByColumn{
		Column: clause.Column{Name: order}, Desc: desc,
	}).Limit(limit).Offset(offset * limit).Select(needFields).Scan(receiver).Error
}

func (v2 *BaseDatabaseImplementV2) CreateSingleDataIfNotExist(ctx context.Context, data any) (created bool, err error) {
	if data == nil || FromSlice(data) || FromMap(data) {
		return false, ErrInvalidSingleData
	}

	session := v2.Db.WithContext(ctx).Model(data).Clauses(clause.OnConflict{DoNothing: true}).Create(data)
	if session.Error != nil {
		return false, session.Error
	}

	return session.RowsAffected > 0, nil
}

func (v2 *BaseDatabaseImplementV2) CreateDataOnDuplicateKeyUpdate(ctx context.Context, data any, indexKeys, updateFields []string) error {
	if len(indexKeys) == 0 || len(updateFields) == 0 {
		return ErrInvalidCondition
	}

	duplicatedColumns := make([]clause.Column, len(indexKeys))
	for i, key := range indexKeys {
		duplicatedColumns[i] = clause.Column{Name: key}
	}

	return v2.Db.WithContext(ctx).Model(data).Clauses(clause.OnConflict{
		Columns:   duplicatedColumns,
		DoUpdates: clause.AssignmentColumns(updateFields),
	}).Create(data).Error
}

func (v2 *BaseDatabaseImplementV2) UpdateDataBySingleCondition(ctx context.Context, updates any, column string, condition any) error {
	if column == "" || condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	return v2.Db.WithContext(ctx).Model(updates).Where(column, condition).Updates(updates).Error
}

func (v2 *BaseDatabaseImplementV2) UpdateDataByCustomCondition(ctx context.Context, updates, condition any) error {
	if condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	return v2.Db.WithContext(ctx).Model(updates).Where(condition).Updates(updates).Error
}

func (v2 *BaseDatabaseImplementV2) ExecuteRawSqlTemplateQuery(ctx context.Context, receiver any, sql string, template RawSqlTemplate) error {
	return v2.Db.WithContext(ctx).Raw(template.ParseTemplate(sql)).Scan(receiver).Error
}

func (v2 *BaseDatabaseImplementV2) ExecuteRawSqlTemplate(ctx context.Context, sql string, template RawSqlTemplate) error {
	return v2.Db.WithContext(ctx).Exec(template.ParseTemplate(sql)).Error
}

var (
	ErrInvalidCondition  = errors.New("invalid condition")
	ErrInvalidSingleData = errors.New("invalid single data")
)
