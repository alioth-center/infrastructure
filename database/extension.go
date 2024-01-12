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
}

type Extension[extended Extended] interface {
	Database
	InitializeExtension(base Database) extended
}

type Extended interface {
	Database
	ExtensionName() string
}

type BaseExtensionImplement[extended Extended] struct {
	Database
	Ext extended
}

func (b *BaseExtensionImplement[extended]) InitializeExtension(base ExtMethods) extended {
	panic("extension must be initialized")
}

type BaseExtendedImplement struct {
	Database
}

func (b *BaseExtendedImplement) ExtensionName() string {
	panic("extension name must be implemented")
}
