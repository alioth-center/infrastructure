package database

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/alioth-center/infrastructure/logger"
)

type Options struct {
	DataSource string
	MaxIdle    int
	MaxOpen    int
	MaxLife    time.Duration
	Timeout    time.Duration
	Logger     logger.Logger
}

// Database is the interface that wraps the basic database operations.
// The implementation of this interface should be thread-safe.
//
// Deprecated: Use DatabaseV2 instead.
type Database interface {
	Init(options Options) error
	Migrate(models ...any) error
	Has(table string, query string, args ...any) (exist bool, err error)
	Count(table string, query string, args ...any) (count int64, err error)
	GetOne(receiver any, query string, args ...any) error
	GetAll(receiver any, query string, args ...any) error
	GetPage(receiver any, offset, limit int, query string, args ...any) error
	PickOne(receiver any, query string, args ...any) error
	PickAll(receiver any, length int, query string, args ...any) error
	InsertOne(data any) error
	InsertAll(data any) error
	UpdateOne(data any, query string, args ...any) error
	UpdateAll(data any, query string, args ...any) error
	DeleteOne(query string, args ...any) error
	DeleteAll(query string, args ...any) error
	ExecRaw(sql string, args ...any) error
	QueryRaw(receiver any, sql string, args ...any) error
	HasWithCtx(ctx context.Context, table string, query string, args ...any) (exist bool, err error)
	CountWithCtx(ctx context.Context, table string, query string, args ...any) (count int64, err error)
	GetOneWithCtx(ctx context.Context, receiver any, query string, args ...any) error
	GetAllWithCtx(ctx context.Context, receiver any, query string, args ...any) error
	PickOneWithCtx(ctx context.Context, receiver any, query string, args ...any) error
	PickAllWithCtx(ctx context.Context, receiver any, length int, query string, args ...any) error
	InsertOneWithCtx(ctx context.Context, data any) error
	InsertAllWithCtx(ctx context.Context, data any) error
	UpdateOneWithCtx(ctx context.Context, data any, query string, args ...any) error
	UpdateAllWithCtx(ctx context.Context, data any, query string, args ...any) error
	DeleteOneWithCtx(ctx context.Context, query string, args ...any) error
	DeleteAllWithCtx(ctx context.Context, query string, args ...any) error
	ExecRawWithCtx(ctx context.Context, sql string, args ...any) error
	QueryRawWithCtx(ctx context.Context, receiver any, sql string, args ...any) error
	ExtMethods() ExtMethods
	SetLogger(logger logger.Logger)
}

// DatabaseV2 is the interface that wraps the basic database operations.
// The implementation of this interface should be thread-safe.
type DatabaseV2 interface {
	// GetGormCore retrieves the core *gorm.DB instance with the provided context.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//
	// Returns:
	//	*gorm.DB: The GORM database instance with the provided context.
	GetGormCore(ctx context.Context) *gorm.DB

	// GetDataBySingleCondition retrieves data from the database based on a single column condition.
	// The result is stored in the receiver.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	receiver (any): The destination where the query result will be stored.
	//	column (string): The column name to apply the condition.
	//	condition (any): The condition value for the specified column.
	//	needFields (...string): Optional fields to select in the query.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	GetDataBySingleCondition(ctx context.Context, receiver any, column string, condition any, needFields ...string) error

	// GetDataByCustomCondition retrieves data from the database based on a custom condition.
	// The result is stored in the receiver.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	receiver (any): The destination where the query result will be stored.
	//	condition (any): The custom condition for the query.
	//	needFields (...string): Optional fields to select in the query.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	GetDataByCustomCondition(ctx context.Context, receiver, condition any, needFields ...string) error

	// ListDataWithPage retrieves a paginated list of data from the database based on the provided filter and ordering.
	// The result is stored in the receiver.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	receiver (any): The destination where the query result will be stored.
	//	filter (any): The filter condition for the query.
	//	order (string): The column name to order by.
	//	desc (bool): Whether to order in descending order.
	//	offset (int): The offset for pagination.
	//	limit (int): The limit for pagination.
	//	needFields (...string): Optional fields to select in the query.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	ListDataWithPage(ctx context.Context, receiver any, filter any, order string, desc bool, offset, limit int, needFields ...string) error

	// CreateSingleDataIfNotExist creates a single data record in the database if it does not already exist.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	data (any): The data record to be created.
	//
	// Returns:
	//	created (bool): True if the data record was created, false if it already exists.
	//	err (error): An error if the operation fails, otherwise nil.
	CreateSingleDataIfNotExist(ctx context.Context, data any) (created bool, err error)

	// CreateDataOnDuplicateKeyUpdate creates a data record in the database or updates it if a duplicate key is found.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	data (any): The data record to be created or updated.
	//	indexKeys ([]string): The index keys to check for duplication.
	//	updateFields ([]string): The fields to update if a duplicate key is found.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	CreateDataOnDuplicateKeyUpdate(ctx context.Context, data any, indexKeys, updateFields []string) error

	// UpdateDataBySingleCondition updates data in the database based on a single column condition.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	updates (any): The data to be updated.
	//	column (string): The column name to apply the condition.
	//	condition (any): The condition value for the specified column.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	UpdateDataBySingleCondition(ctx context.Context, updates any, column string, condition any) error

	// UpdateDataByCustomCondition updates data in the database based on a custom condition.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	updates (any): The data to be updated.
	//	condition (any): The custom condition for the query.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	UpdateDataByCustomCondition(ctx context.Context, updates, condition any) error

	// ExecuteRawSqlTemplateQuery executes a raw SQL template query with the provided context.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	receiver (any): The destination where the query result will be stored.
	//  sql (string): The raw SQL template to execute.
	//	template (RawSqlTemplate): The raw SQL template value to execute.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	ExecuteRawSqlTemplateQuery(ctx context.Context, receiver any, sql string, template RawSqlTemplate) error

	// ExecuteRawSqlTemplate executes a raw SQL template with the provided context.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//  sql (string): The raw SQL template to execute.
	//	template (RawSqlTemplate): The raw SQL template value to execute.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	ExecuteRawSqlTemplate(ctx context.Context, sql string, template RawSqlTemplate) error

	// ExecuteRawSqlQuery executes a raw SQL query with the provided context.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//	receiver (any): The destination where the query result will be stored.
	//  sql (string): The raw SQL to execute.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	ExecuteRawSqlQuery(ctx context.Context, receiver any, sql string) error

	// ExecuteRawSql executes a raw SQL with the provided context.
	//
	// Parameters:
	//	ctx (context.Context): The context for the database operation.
	//  sql (string): The raw SQL to execute.
	//
	// Returns:
	//	error: An error if the operation fails, otherwise nil.
	ExecuteRawSql(ctx context.Context, sql string) error
}
