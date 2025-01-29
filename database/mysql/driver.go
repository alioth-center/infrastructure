package mysql

import (
	"context"
	"fmt"

	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/utils/shortcut"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DriverName = "mysql"

// DefaultDriver returns a new instance of the mysqlDriver.
// Use database.NewDatabaseConnection with this driver to create a new mysql connection.
func DefaultDriver() database.Driver {
	return mysqlDriver{}
}

type mysqlDriver struct{}

func (m mysqlDriver) DriverName() string {
	return DriverName
}

func (m mysqlDriver) BuildDataSource(dsn database.DSN) string {
	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		dsn.Username, dsn.Password,
		dsn.Host, dsn.Port, dsn.Database,
		dsn.Charset, shortcut.Ternary(dsn.ParseTime, "True", "False"), dsn.Location,
	)
	dataSource += shortcut.Ternary(dsn.Timeout > 0, fmt.Sprintf("&timeout=%ds", dsn.Timeout), "")

	return dataSource
}

func (m mysqlDriver) Connect(_ context.Context, option database.Options) (database.Database, error) {
	if option.Logger == nil {
		option.Logger = logger.Default()
	}

	db, connectErr := gorm.Open(mysql.Open(option.DataSource), &gorm.Config{
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
