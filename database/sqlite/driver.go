package sqlite

import (
	"context"

	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DriverName = "sqlite"

// DefaultDriver returns a new instance of the sqliteDriver.
// Use database.NewDatabaseConnection with this driver to create a new sqlite connection.
func DefaultDriver() database.Driver {
	return sqliteDriver{}
}

type sqliteDriver struct{}

func (s sqliteDriver) DriverName() string {
	return DriverName
}

func (s sqliteDriver) BuildDataSource(dsn database.DSN) string {
	return dsn.Host
}

func (s sqliteDriver) Connect(ctx context.Context, option database.Options) (database.Database, error) {
	if option.Logger == nil {
		option.Logger = logger.Default()
	}

	db, connectErr := gorm.Open(sqlite.Open(option.DataSource))
	if connectErr != nil {
		return nil, connectErr
	}

	if applyErr := database.ApplyCommonOptions(db, option); applyErr != nil {
		return nil, applyErr
	}

	return database.NewBaseDatabaseImplement(db), nil
}
