package sqlite_test

import (
	"context"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/sqlite"
	"testing"
)

func TestSqliteDriver(t *testing.T) {
	t.Run("DriverName", func(t *testing.T) {
		driver := sqlite.DefaultDriver()
		if driver.DriverName() != sqlite.DriverName {
			t.Errorf("DriverName() = %s, want %s", driver.DriverName(), sqlite.DriverName)
		}
	})

	t.Run("BuildDataSource", func(t *testing.T) {
		driver := sqlite.DefaultDriver()
		dsn := database.DSN{
			Host: "test.db",
		}
		dataSource := driver.BuildDataSource(dsn)
		if dataSource != dsn.Host {
			t.Errorf("BuildDataSource() = %s, want %s", dataSource, dsn.Host)
		}
	})

	t.Run("Connect", func(t *testing.T) {
		driver := sqlite.DefaultDriver()
		option := database.Options{
			DataSource: ":memory:",
		}
		db, err := driver.Connect(nil, option)
		if err != nil {
			t.Errorf("Connect() error = %v, want nil", err)
		}
		if db == nil {
			t.Errorf("Connect() db = nil, want not nil")
		}

		result := map[string]any{}
		if db.GetGormCore(context.Background()).Raw(`select "ack" as column;`).Scan(&result).Error != nil {
			t.Errorf("Connect() db.Exec() error = %v, want nil", err)
		}
		if result["column"] != "ack" {
			t.Errorf("Connect() db.Exec() result = %v, want ack", result["column"])
		}
	})
}
