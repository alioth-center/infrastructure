package sqlite

import (
	"github.com/alioth-center/infrastructure/trace"
	"testing"
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
	ctx := trace.NewContextWithTraceID()
	sqlite, e := NewSqliteDb(opt, &table{})
	if e != nil {
		t.Fatal(e)
	}

	t.Log("insert error:", sqlite.InsertOneWithCtx(ctx, &table{Value: "test"}))
	instance := table{}
	t.Log("query error:", sqlite.GetOneWithCtx(ctx, &instance, "value = ?", "test"))
	t.Log("query result:", instance)
	t.Log("pick error:", sqlite.PickOneWithCtx(ctx, &instance, "value = ?", "test"))
}
