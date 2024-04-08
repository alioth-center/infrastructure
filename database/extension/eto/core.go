package eto

import "github.com/alioth-center/infrastructure/database"

const ExtensionName = "eto"

type (
	// EntityObject 数据实体对象，用于在代码中方便地传递数据表所对应的结构
	EntityObject[po PersistentObject, dto any] interface {
		EoDefinition(po, dto)
		FromPO(po)
		FromDTO(dto)
		ToPO() po
		ToDTO() dto
	}

	// QueryObject 数据库查询对象，
	QueryObject[dto any] interface {
		QoDefinition(dto)
		FromDTO(dto)
	}

	DataTransferObject[po PersistentObject] interface {
		DtoDefinition(po)
		FromPO(po)
		ToPO() po
	}

	PersistentObject interface {
		PoDefinition()
		TableName() string
	}
)

type Extended interface {
	database.Extended
}

type extended struct {
	database.Database
	methods database.ExtMethods
}

func (e *extended) ExtensionName() string {
	return ExtensionName
}
