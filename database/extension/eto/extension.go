package eto

import "github.com/alioth-center/infrastructure/database"

type Extension struct{}

func (e *Extension) InitializeExtension(base database.Database) Extended {
	return &extended{
		Database: base,
		methods:  base.ExtMethods(),
	}
}
