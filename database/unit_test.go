package database

import (
	"context"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"testing"
)

type User struct {
	ID   int `gorm:"primaryKey"`
	Name string
	Age  int
}

func TestBaseDatabaseImplementV2(t *testing.T) {
	// 创建内存中的 SQLite 数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移数据库
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// 实例化 BaseDatabaseImplementV2
	baseDB := &BaseDatabaseImplementV2{Db: db}

	// 测试 CreateSingleDataIfNotExist
	ctx := context.Background()
	user := User{Name: "Alice", Age: 25}
	created, err := baseDB.CreateSingleDataIfNotExist(ctx, &user)
	if err != nil {
		t.Fatalf("failed to create single data if not exist: %v", err)
	}
	if !created {
		t.Fatalf("expected user to be created")
	}

	// 测试 GetDataBySingleCondition
	var retrievedUser User
	err = baseDB.GetDataBySingleCondition(ctx, &retrievedUser, "name", "Alice")
	if err != nil {
		t.Fatalf("failed to get data by single condition: %v", err)
	}
	if retrievedUser.Name != "Alice" {
		t.Fatalf("expected user name to be 'Alice', got '%s'", retrievedUser.Name)
	}

	// 测试 ListDataWithPage
	var users []User
	err = baseDB.ListDataWithPage(ctx, &users, &User{Age: 25}, "id", false, 0, 10)
	if err != nil {
		t.Fatalf("failed to list data with page: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	// 测试 UpdateDataBySingleCondition
	err = baseDB.UpdateDataBySingleCondition(ctx, &User{Age: 26}, "name", "Alice")
	if err != nil {
		t.Fatalf("failed to update data by single condition: %v", err)
	}
	err = baseDB.GetDataBySingleCondition(ctx, &retrievedUser, "name", "Alice")
	if err != nil {
		t.Fatalf("failed to get data by single condition: %v", err)
	}
	if retrievedUser.Age != 26 {
		t.Fatalf("expected user age to be 26, got %d", retrievedUser.Age)
	}
}
