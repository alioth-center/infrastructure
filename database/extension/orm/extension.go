package orm

import "github.com/alioth-center/infrastructure/database"

type Extension struct {
	database.Database
}

func (e *Extension) InitializeExtension(base database.Database) Extended {
	return &extended{
		Database: base,
		methods:  base.ExtMethods(),
	}
}

func NewExtension(base database.Database) database.Extension[Extended] {
	return &Extension{
		Database: base,
	}
}
