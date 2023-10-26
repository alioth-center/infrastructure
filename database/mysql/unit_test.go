package mysql

import (
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"testing"
)

func TestMysqlDb(t *testing.T) {
	// you must create a database named "test" and a table named "test_table" in your mysql server,
	// or you can change the config to your own database, but you must change the table name in the test code.
	cfg := Config{
		Server:    "192.168.1.140",
		Port:      3306,
		Username:  "public",
		Password:  "123456",
		Database:  "payment",
		Charset:   "utf8mb4",
		Location:  "Local",
		ParseTime: true,
		Debug:     true,
	}
	mysql, e := NewMysqlDb(cfg)
	if e != nil {
		t.Fatal(e)
	}

	ctx := trace.NewContextWithTraceID()
	result := []map[string]any{}
	qe := mysql.QueryRawWithCtx(ctx, &result, "select * from payment.merchant_data limit 10")
	if qe != nil {
		t.Fatal(qe)
	}

	for _, v := range result {
		logger.Info(logger.NewFields(ctx).WithMessage("query result").WithData(v))
	}
}
