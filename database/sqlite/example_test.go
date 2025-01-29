package sqlite_test

import (
	"context"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/sqlite"
)

func ExampleDefaultDriver() {
	driver := sqlite.DefaultDriver()
	dsn := database.DSN{Host: "test.db"}
	conn, connectErr := database.NewDatabaseConnection(context.Background(), driver, dsn, database.Options{})
	if connectErr != nil {
		// handle error...
		panic(connectErr)
	}

	// do something with conn...
	conn.GetGormCore(context.Background()).Select("ack as column").Scan(&map[string]interface{}{})
}
