package eto

import "github.com/alioth-center/infrastructure/database"

func DtoGetOne[po, dto any](db database.Database, model po, condition dto) (result dto, err error) {
	gorm := db.ExtMethods().GetGorm()
	err = gorm.Model(model).First(condition).Scan(result).Error
	return result, err
}
