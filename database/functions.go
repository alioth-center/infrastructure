package database

import (
	"reflect"

	"github.com/alioth-center/infrastructure/utils/values"
)

func FromSlice(s any) bool {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Kind() == reflect.Slice
}

func EmptySlice(s any) bool {
	v := reflect.ValueOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Kind() == reflect.Slice && v.Len() == 0
}

func FromMap(m any) bool {
	v := reflect.ValueOf(m)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Kind() == reflect.Map
}

func Receivable(v any) bool {
	return v != nil && reflect.ValueOf(v).Kind() == reflect.Ptr
}

func Column(tableName, columnName string) string {
	return values.BuildStrings(tableName, ".", columnName)
}

func ColumnAlias(tableName, columnName, alias string) string {
	return values.BuildStrings(tableName, ".", columnName, " as ", alias)
}
