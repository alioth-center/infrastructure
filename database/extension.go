package database

import (
	"context"

	"gorm.io/gorm"
)

type ExtMethods interface {
	DriverName() string
	GetGorm() *gorm.DB
	Exec(cmd func(db *gorm.DB) *gorm.DB) error
	ExecCtx(ctx context.Context, cmd func(db *gorm.DB) *gorm.DB) error
	Transaction(cmd func(tx *gorm.DB) error) error
	TransactionCtx(ctx context.Context, cmd func(tx *gorm.DB) error) error
}

type Extension[extended Extended] interface {
	InitializeExtension(base Database) extended
}

type Extended interface {
	Database
	ExtensionName() string
}
