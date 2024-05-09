package mysql

import (
	"os"
	"testing"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
)

func TestMysqlDb(t *testing.T) {
	if os.Getenv("MYSQL_TEST_ENABLE") != "true" {
		t.Skip("skip mysql test")
	}

	// you must create a database named "test" and a table named "test_table" in your mysql server,
	// or you can change the config to your own database, but you must change the table name in the test code.
	cfg := Config{
		Server:    "localhost",
		Port:      3306,
		Username:  "root",
		Password:  "123456",
		Database:  "test_db",
		Charset:   "utf8mb4",
		Location:  "Local",
		ParseTime: true,
		Debug:     true,
	}
	mysql, e := NewMysqlDb(cfg)
	if e != nil {
		t.Fatal(e)
	}

	ctx := trace.NewContext()
	result := []map[string]any{}
	qe := mysql.QueryRawWithCtx(ctx, &result, "select * from test_db.test_table limit 10")
	if qe != nil {
		t.Fatal(qe)
	}

	for _, v := range result {
		logger.Info(logger.NewFields(ctx).WithMessage("query result").WithData(v))
	}
}
