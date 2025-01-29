package postgres

import (
	"context"
	"fmt"

	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/utils/shortcut"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DriverName = "postgres"

// DefaultDriver returns a new instance of the postgresDriver.
// Use database.NewDatabaseConnection with this driver to create a new postgres connection.
func DefaultDriver() database.Driver {
	return postgresDriver{}
}

type postgresDriver struct{}

func (p postgresDriver) DriverName() string {
	return DriverName
}

func (p postgresDriver) BuildDataSource(dsn database.DSN) string {
	dataSource := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		dsn.Host, dsn.Username, dsn.Password, dsn.Database, dsn.Port,
		shortcut.Ternary(dsn.Location == "", "Asia/Shanghai", dsn.Location),
	)
	dataSource += shortcut.Ternary(dsn.Timeout > 0, fmt.Sprintf("connect_timeout=%d", dsn.Timeout), "")

	return dataSource
}

func (p postgresDriver) Connect(_ context.Context, option database.Options) (database.Database, error) {
	if option.Logger == nil {
		option.Logger = logger.Default()
	}

	db, connectErr := gorm.Open(postgres.Open(option.DataSource), &gorm.Config{
		Logger: database.NewDBLogger(option.Logger),
	})
	if connectErr != nil {
		return nil, connectErr
	}

	if applyErr := database.ApplyCommonOptions(db, option); applyErr != nil {
		return nil, applyErr
	}

	return database.NewBaseDatabaseImplement(db), nil
}
