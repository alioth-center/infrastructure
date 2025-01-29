package postgres_test

import (
	"context"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/postgres"
)

func ExampleDefaultDriver() {
	driver := postgres.DefaultDriver()
	dsn := database.DSN{
		Host:      "1217.0.0.1",
		Port:      5432,
		Username:  "your_username",
		Password:  "your_password",
		Database:  "your_database",
		Charset:   "utf8mb4",
		Location:  "Asia/Shanghai",
		ParseTime: true,
		Timeout:   10,
	}
	conn, connectErr := database.NewDatabaseConnection(context.Background(), driver, dsn, database.Options{})
	if connectErr != nil {
		// handle error...
		panic(connectErr)
	}

	// do something with conn...
	conn.GetGormCore(context.Background()).Select("ack as column").Scan(&map[string]interface{}{})
}
