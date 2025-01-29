package database

import (
	"context"
	"errors"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type BaseDatabaseImplement struct {
	conn *gorm.DB
}

func (impl *BaseDatabaseImplement) GetGormCore(ctx context.Context) *gorm.DB {
	return impl.conn.WithContext(ctx)
}

func (impl *BaseDatabaseImplement) GetDataBySingleCondition(ctx context.Context, receiver any, column string, condition any, needFields ...string) error {
	if column == "" || condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	if len(needFields) == 0 {
		needFields = append(needFields, "*")
	}

	return impl.GetGormCore(ctx).Model(receiver).Where(column, condition).Select(needFields).Scan(receiver).Error
}

func (impl *BaseDatabaseImplement) GetDataByCustomCondition(ctx context.Context, receiver, condition any, needFields ...string) error {
	if condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	if len(needFields) == 0 {
		needFields = append(needFields, "*")
	}

	return impl.GetGormCore(ctx).Model(receiver).Where(condition).Select(needFields).Scan(receiver).Error
}

func (impl *BaseDatabaseImplement) CreateSingleDataIfNotExist(ctx context.Context, data any) (created bool, err error) {
	if data == nil || FromSlice(data) || FromMap(data) {
		return false, ErrInvalidSingleData
	}

	session := impl.GetGormCore(ctx).Model(data).Clauses(clause.OnConflict{DoNothing: true}).Create(data)
	if session.Error != nil {
		return false, session.Error
	}

	return session.RowsAffected > 0, nil
}

func (impl *BaseDatabaseImplement) CreateDataOnDuplicateKeyUpdate(ctx context.Context, data any, indexKeys, updateFields []string) error {
	if len(indexKeys) == 0 || len(updateFields) == 0 {
		return ErrInvalidCondition
	}

	duplicatedColumns := make([]clause.Column, len(indexKeys))
	for i, key := range indexKeys {
		duplicatedColumns[i] = clause.Column{Name: key}
	}

	return impl.GetGormCore(ctx).Model(data).Clauses(clause.OnConflict{
		Columns:   duplicatedColumns,
		DoUpdates: clause.AssignmentColumns(updateFields),
	}).Create(data).Error
}

func (impl *BaseDatabaseImplement) UpdateDataBySingleCondition(ctx context.Context, updates any, column string, condition any) error {
	if column == "" || condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	return impl.GetGormCore(ctx).Model(updates).Where(column, condition).Updates(updates).Error
}

func (impl *BaseDatabaseImplement) UpdateDataByCustomCondition(ctx context.Context, updates, condition any) error {
	if condition == nil || EmptySlice(condition) {
		return ErrInvalidCondition
	}

	return impl.GetGormCore(ctx).Model(updates).Where(condition).Updates(updates).Error
}

func (impl *BaseDatabaseImplement) ExecuteRawSqlTemplateQuery(ctx context.Context, receiver any, sql string, template RawSqlTemplate) error {
	return impl.GetGormCore(ctx).Raw(template.ParseTemplate(sql)).Scan(receiver).Error
}

func (impl *BaseDatabaseImplement) ExecuteRawSqlTemplate(ctx context.Context, sql string, template RawSqlTemplate) error {
	return impl.GetGormCore(ctx).Exec(template.ParseTemplate(sql)).Error
}

func (impl *BaseDatabaseImplement) ExecuteRawSqlQuery(ctx context.Context, receiver any, sql string) error {
	return impl.GetGormCore(ctx).Raw(sql).Scan(receiver).Error
}

func (impl *BaseDatabaseImplement) ExecuteRawSql(ctx context.Context, sql string) error {
	return impl.GetGormCore(ctx).Exec(sql).Error
}

var (
	ErrInvalidCondition  = errors.New("invalid condition")
	ErrInvalidSingleData = errors.New("invalid single data")
)

func NewBaseDatabaseImplement(conn *gorm.DB) *BaseDatabaseImplement {
	return &BaseDatabaseImplement{conn: conn}
}

func ApplyCommonOptions(db *gorm.DB, opts Options) error {
	rawDatabase, getRawDatabaseErr := db.DB()
	if getRawDatabaseErr != nil {
		return getRawDatabaseErr
	}
	if opts.MaxIdle > 0 {
		rawDatabase.SetMaxIdleConns(opts.MaxIdle)
	}
	if opts.MaxOpen > 0 {
		rawDatabase.SetMaxOpenConns(opts.MaxOpen)
	}
	if opts.MaxLife > 0 {
		rawDatabase.SetConnMaxLifetime(opts.MaxLife)
	}
	if len(opts.MigrateModels) > 0 {
		migrateErr := db.AutoMigrate(opts.MigrateModels...)
		if migrateErr != nil {
			return migrateErr
		}
	}

	return nil
}
