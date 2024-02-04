package database

import (
	"strings"
	"testing"
)

func TestTemplate(t *testing.T) {
	registrationServiceCheckExistInstanceSQL := `select id, connected from {table_prefix}regions where code = '${}'`
	t.Log(NewTemplate(
		registrationServiceCheckExistInstanceSQL,
		map[string]string{
			"table_prefix":     "test_",
			"service_endpoint": "{table_r",
		},
	).Prepare())
}

func BenchmarkTestTemplate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		template := NewTemplate(
			`select * from {table_prefix}database from ${table_name} where id = [idx] and name = '$[name]'`,
			map[string]string{
				"table_prefix": "test_",
				"table_name":   "test_table",
				"idx":          "1",
				"name":         "'; drop table test_table; --",
			},
		)
		template.Prepare()
	}
}
func BenchmarkTestTemplate2(b *testing.B) {
	t := &template{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template := `select * from {table_prefix}database from ${table_name} where id = [idx] and name = '$[name]'`
		args := map[string]string{
			"table_prefix": "test_",
			"table_name":   "test_table",
			"idx":          "1",
			"name":         "'; drop table test_table; --",
		}

		for k, v := range args {
			template = strings.ReplaceAll(template, "{"+k+"}", t.escape(v))
			template = strings.ReplaceAll(template, "${"+k+"}", t.escape(v))
			template = strings.ReplaceAll(template, "["+k+"]", t.escape(v))
			template = strings.ReplaceAll(template, "$["+k+"]", t.escape(v))
		}
	}
}
