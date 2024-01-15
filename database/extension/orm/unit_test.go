package orm

import (
	"github.com/alioth-center/infrastructure/database/sqlite"
	"gorm.io/gorm"
	"testing"
)

type table struct {
	ID    int    `gorm:"primary_key;column:id;autoIncrement"`
	Value string `gorm:"type:varchar(100);column:value"`
}

func (t table) TableName() string {
	return "test_table"
}

func TestOrmExtension(t *testing.T) {
	opt := sqlite.Config{
		Database:      ":memory:",
		Debug:         true,
		TimeoutSecond: 1,
	}
	db, e := sqlite.NewSqliteDb(opt, &table{})
	if e != nil {
		t.Fatal(e)
	}

	ext := NewExtension(db)
	extMethods := ext.InitializeExtension(db)
	_ = extMethods.ExecuteGormFunction(func(db *gorm.DB) *gorm.DB {
		return db.Create(&table{Value: "test"})
	})

	var tb table
	_ = extMethods.PickOne(&tb, "value = ?", "test")
	t.Logf("pick result: %v", tb)
}
