package sqlite

import (
	"fmt"
	"testing"

	"github.com/alioth-center/infrastructure/trace"
)

type table struct {
	ID    int    `gorm:"primary_key;column:id;autoIncrement"`
	Value string `gorm:"type:varchar(100);column:value"`
}

func (t table) TableName() string {
	return "test_table"
}

func TestSqliteDb(t *testing.T) {
	opt := Config{
		Database:      ":memory:",
		Debug:         true,
		TimeoutSecond: 1,
	}
	ctx := trace.NewContext()
	sqlite, e := NewSqliteDb(opt, &table{})
	if e != nil {
		t.Fatal(e)
	}

	t.Log("insert error:", sqlite.InsertOneWithCtx(ctx, &table{Value: "test"}))
	instance := table{}
	t.Log("query error:", sqlite.GetOneWithCtx(ctx, &instance, "value = ?", "test"))
	t.Log("query result:", instance)
	t.Log("pick error:", sqlite.PickOneWithCtx(ctx, &instance, "value = ?", "test"))

	t.Log(sqlite.QueryRaw(&instance, fmt.Sprintf("select * from test_table where value = '%s'", "'; insert into test_table(value) values('test_inject'); --")))

	var values []table
	t.Log(sqlite.GetAll(&values, ""))
	t.Log(values)
}
