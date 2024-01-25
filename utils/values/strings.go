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

func StringToInt[T ~int](raw string, defaultValue T) T {
	if raw == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}

	return T(result)
}

func StringToUint[T ~uint](raw string, defaultValue T) T {
	if raw == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(raw)
	if err != nil {
		return defaultValue
	}

	return T(result)
}

func StringToFloat64[T ~float64](raw string, defaultValue T) T {
	if raw == "" {
		return defaultValue
	}

	result, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return defaultValue
	}

	return T(result)
}

func StringToBool[T ~bool](raw string, defaultValue T) T {
	if raw == "" {
		return defaultValue
	}

	result, err := strconv.ParseBool(raw)
	if err != nil {
		return defaultValue
	}

	return T(result)
}

func StringToStringPtr[T ~string](raw string) *T {
	if raw == "" {
		return nil
	}

	result := T(raw)
	return &result
}

func StringToIntPtr[T ~int](raw string) *T {
	if raw == "" {
		return nil
	}

	result, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}

	resultT := T(result)
	return &resultT
}

func StringToUintPtr[T ~uint](raw string) *T {
	if raw == "" {
		return nil
	}

	result, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}

	resultT := T(result)
	return &resultT
}

func StringToFloat64Ptr[T ~float64](raw string) *T {
	if raw == "" {
		return nil
	}

	result, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil
	}

	resultT := T(result)
	return &resultT
}

func StringToBoolPtr[T ~bool](raw string) *T {
	if raw == "" {
		return nil
	}

	result, err := strconv.ParseBool(raw)
	if err != nil {
		return nil
	}

	resultT := T(result)
	return &resultT
}
