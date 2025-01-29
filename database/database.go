package database

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/alioth-center/infrastructure/logger"
)

// DSN is the data source definition for the database connection.
type DSN struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Database  string
	Charset   string
	Location  string
	ParseTime bool
	Timeout   int
}

// Options is the configuration for the database connection.
type Options struct {
	DataSource    string
	MaxIdle       int
	MaxOpen       int
	MaxLife       time.Duration
	Logger        logger.Logger
	MigrateModels []any
}

// Database is the interface that wraps the basic database operations.
// The implementation of this interface should be thread-safe.
type Database interface {
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

type Driver interface {
	// DriverName returns the name of the database driver.
	DriverName() string

	// BuildDataSource builds the data source link for the database connection.
	BuildDataSource(dsn DSN) string

	// Connect establishes a connection to the database with the provided context and options.
	Connect(ctx context.Context, option Options) (Database, error)
}

// NewDatabaseConnection creates a new database connection with the provided context, driver, DSN, and options.
// If the data source is not provided in the options, it will be built using the Driver.BuildDataSource method.
//
// Parameters:
//
//	ctx (context.Context): The context for the database connection.
//	driver (Driver): The database driver to use for the connection.
//	dsn (DSN): The data source definition for the database connection.
//	option (Options): The configuration for the database connection.
//
// Returns:
//
//	Database: The database connection instance.
//	error: An error if the operation fails, otherwise nil.
func NewDatabaseConnection(ctx context.Context, driver Driver, dsn DSN, option Options) (Database, error) {
	if option.DataSource == "" {
		option.DataSource = driver.BuildDataSource(dsn)
	}

	return driver.Connect(ctx, option)
}
