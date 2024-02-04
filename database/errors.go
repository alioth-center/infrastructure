package database

import "github.com/alioth-center/infrastructure/utils/values"

type TemplateFormatError struct {
	Template string `json:"template"`
}

func (err *TemplateFormatError) Error() string {
	return values.BuildStrings("template syntax format error: ", err.Template)
}

func NewTemplateFormatError(template []rune) error {
	if len(template) > 100 {
		return &TemplateFormatError{Template: values.BuildStrings(string(template[:100]), "...")}
	}

	return &TemplateFormatError{Template: string(template[:100])}
}

type ExecuteSqlError struct {
	Sql           string `json:"sql"`
	ErrorOccurred error  `json:"error"`
}

func (err *ExecuteSqlError) Error() string {
	return values.BuildStrings("execute sql error: ", err.Sql, ", error: ", err.ErrorOccurred.Error())
}

func NewExecuteSqlError(sql string, errorOccurred error) error {
	return &ExecuteSqlError{Sql: sql, ErrorOccurred: errorOccurred}
}
