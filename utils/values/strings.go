package values

import (
	"strconv"
	"strings"
)

// BuildStrings 构建字符串，使用strings.Builder
// example:
//
//	BuildStrings("a", "b", "c") -> "abc"
func BuildStrings(parts ...string) string {
	builder := strings.Builder{}
	for _, part := range parts {
		builder.WriteString(part)
	}

	return builder.String()
}

// BuildStringsWithJoin 构建字符串，使用strings.Join
// example:
//
//	BuildStringsWithJoin("/", "a", "b", "c") -> "a/b/c"
func BuildStringsWithJoin(sep string, parts ...string) string {
	if parts == nil || len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, sep)
}

// BuildStringsWithJoinIgnoreEmpty use strings.Join to build string, all empty string in parts will be ignored
// example:
//
//	BuildStringsWithJoinIgnoreEmpty("/", "a", "", "b", "c", "") -> "a/b/c"
func BuildStringsWithJoinIgnoreEmpty(sep string, parts ...string) string {
	realParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			realParts = append(realParts, part)
		}
	}

	return BuildStringsWithJoin(sep, realParts...)
}

// BuildStringsWithReplacement 构建字符串，使用strings.Builder，同时替换字符串
// example:
//
//	BuildStringsWithReplacement(map[string]string{"a": "1", "b": "2", "c": "3"}, "a", "b", "c") -> "123"
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
// example:
//
//	BuildStringsWithTemplate("a{1}b{2}c{3}", "1", "2", "3") -> "a1b2c3"
func BuildStringsWithTemplate(template string, args ...string) string {
	if args == nil || len(args) == 0 {
		return template
	}

	for i, arg := range args {
		template = strings.ReplaceAll(template, BuildStrings("{", strconv.Itoa(i+1), "}"), arg)
	}

	return template
}

// StringToInt 字符串转换为int
// example:
//
//	StringToInt("1", 0) -> 1
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

// StringToUint 字符串转换为uint
// example:
//
//	StringToUint("1", 0) -> 1
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

// StringToFloat64 字符串转换为float64
// example:
//
//	StringToFloat64("1.0", 0) -> 1.0
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

// StringToBool 字符串转换为bool
// example:
//
//	StringToBool("true", false) -> true
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

// StringToStringPtr 字符串转换为*string
// example:
//
//	StringToStringPtr("test") -> "test"
func StringToStringPtr[T ~string](raw string) *T {
	if raw == "" {
		return nil
	}

	result := T(raw)
	return &result
}

// StringToIntPtr 字符串转换为*int
// example:
//
//	StringToIntPtr("1") -> 1
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

// StringToUintPtr 字符串转换为*uint
// example:
//
//	StringToUintPtr("1") -> 1
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

// StringToFloat64Ptr 字符串转换为*float64
// example:
//
//	StringToFloat64Ptr("1.0") -> 1.0
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

// StringToBoolPtr 字符串转换为*bool
// example:
//
//	StringToBoolPtr("true") -> true
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

type StringTemplate struct {
	template  []rune
	arguments map[string]string
}

func (t *StringTemplate) findVariable(start int, endSignals map[rune]bool) (variable string, end int) {
	index, length := start, len(t.template)
	for index < length && !endSignals[t.template[index]] {
		if index-start > 32 {
			// variable name too long
			return "", 0
		}

		index++
	}

	if index == length {
		// did not find the endSignal
		return "", 0
	}

	return string(t.template[start:index]), index
}

func (t *StringTemplate) Parse() string {
	builder, length := strings.Builder{}, len(t.template)
	builder.Grow(length)
	endSignals := map[rune]bool{'$': true, '}': true, ']': true, ':': true}

	for index := 0; index < length; {
		current := t.template[index]
		switch current {
		case '$', '{', '[', ':':
			// cases: $variable_name$, {variable_name}, [variable_name], :variable_name:
			if index+1 < length {
				// only check the next character when it exists
				variable, end := t.findVariable(index+1, endSignals)
				if end == 0 || variable == "" {
					// not found the endSignal, write the current character and continue
					builder.WriteRune(current)
					index++
					continue
				}

				argument, exist := t.arguments[variable]
				// found the endSignal but argument not given, write the current character and continue
				if !exist {
					builder.WriteRune(current)
					index++
					continue
				}

				// found the endSignal and argument, replace the variable with the argument
				builder.WriteString(argument)
				index = end + 1
				continue
			}

			// do not have more than one character, write the current character and continue
			builder.WriteRune(current)
			index++
		default:
			// normal character, write the current character and continue
			builder.WriteRune(current)
			index++
		}
	}

	return builder.String()
}

func NewStringTemplate(template string, arguments map[string]string) *StringTemplate {
	if arguments == nil {
		arguments = map[string]string{}
	}

	return &StringTemplate{template: []rune(template), arguments: arguments}
}
