package orm

import "github.com/alioth-center/infrastructure/database"

type Extension struct{}

func (e *Extension) InitializeExtension(base database.Database) Extended {
	return &extended{
		Database: base,
		methods:  base.ExtMethods(),
	}
}

// NewExtension creates a new extension instance
// example:
//
//	var db database.Database
//	extension := orm.NewExtension().InitializeExtension(db)
//
// then can use extension to execute gorm function
//
//	extension.ExecuteGormFunction(func(db *gorm.DB) *gorm.DB {
//		return db.Model(&User{}).Where("id = ?", 1)
//	})
func NewExtension() database.Extension[Extended] {
	return &Extension{}
}
