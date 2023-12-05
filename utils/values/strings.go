package values

import (
	"strconv"
	"strings"
)

// BuildStrings 构建字符串，使用strings.Builder
func BuildStrings(parts ...string) string {
	builder := strings.Builder{}
	for _, part := range parts {
		builder.WriteString(part)
	}

	return builder.String()
}

// BuildStringsWithReplacement 构建字符串，使用strings.Builder，同时替换字符串
func BuildStringsWithReplacement(replacement map[string]string, parts ...string) string {
	builder := strings.Builder{}
	for _, part := range parts {
		builder.WriteString(replacement[part])
	}

	result := builder.String()
	for original, replace := range replacement {
		result = strings.ReplaceAll(result, original, replace)
	}

	return result
}

// BuildStringsWithTemplate 构建字符串并替换模板，使用strings.ReplaceAll
func BuildStringsWithTemplate(template string, args ...string) string {
	if args == nil || len(args) == 0 {
		return template
	}

	for i, arg := range args {
		template = strings.ReplaceAll(template, BuildStrings("{", strconv.Itoa(i), "}"), arg)
	}

	return template
}
