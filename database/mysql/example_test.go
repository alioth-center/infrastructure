package mysql_test

import (
	"context"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/mysql"
)

func ExampleDefaultDriver() {
	driver := mysql.DefaultDriver()
	dsn := database.DSN{
		Host:      "127.0.0.1",
		Port:      3306,
		Username:  "your_username",
		Password:  "your_password",
		Database:  "your_database",
		Charset:   "uft8mb4",
		Location:  "Local",
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
