package database

type RawSqlTemplate interface {
	ParseTemplate(tmpl string) (result string)
}
