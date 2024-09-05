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
	if len(parts) == 0 {
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
	if len(args) == 0 {
		return template
	}

	for i, arg := range args {
		template = strings.ReplaceAll(template, BuildStrings("{", strconv.Itoa(i+1), "}"), arg)
	}

	return template
}

// SecretString returns a string with the middle part replaced by paddingChar.
//
// Parameters:
//
//	raw (string): The original string.
//	prefixDisplay (int): The number of characters to display at the beginning of the string.
//	suffixDisplay (int): The number of characters to display at the end of the string.
//	paddingChar (string): The character to replace the middle part of the string.
//
// Returns:
//
//	string: The string with the middle part replaced by paddingChar.
//
// Example:
//
//	SecretString("1234567890", 3, 3, "*") -> "123****890"
//	SecretString("123", 3, 3, "*") -> "123"
func SecretString(raw string, prefixDisplay, suffixDisplay int, paddingChar string) string {
	runes := []rune(raw)
	length := len(runes)
	if length <= prefixDisplay+suffixDisplay {
		return string(runes)
	}

	return BuildStrings(string(runes[:prefixDisplay]), strings.Repeat(paddingChar, length-prefixDisplay-suffixDisplay), string(runes[length-suffixDisplay:]))
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
	template     string
	arguments    map[string]string
	getVariable  func(key string, arguments map[string]string) (exist bool, result string)
	startSignals charSet
	endSignals   charSet

	escapeCache map[string]string
}

type RawSqlTemplate struct {
	*StringTemplate
}

// NewStringTemplate 创建一个字符串模板，使用支持的内容作为参数，支持的内容有：结构体(指针)，map[string]string，*map[string]string
//   - template: 模板字符串
//   - args: 参数，如果是结构体(指针)，需要指定cpc tag，并且目前只支持uint, int, float, string, bool和嵌套结构体类型的成员字段
//
// example:
//
//	type User struct {
//		Name string `cpc:"key:name"`
//		Age  int    `cpc:"key:age,omitempty"`
//		Male bool   `cpc:"key:male,default:true"`
//	}
//	user := &User{
//		Name: "test",
//		Age:  0,
//	}
//	template := "name: ${name}, age: ${age}, male: ${male}"
//	result := NewStringTemplate(template, user).Parse()
//	fmt.Println(result) // name: test, age: ${age}, male: true
func NewStringTemplate(template string, args any) *StringTemplate {
	// 如果args是map[string]string，调用NewStringTemplateWithMap
	if m, ok := args.(map[string]string); ok {
		return NewStringTemplateWithMap(template, m)
	}

	// 如果args是*map[string]string，调用NewStringTemplateWithMap
	if m, ok := args.(*map[string]string); ok {
		return NewStringTemplateWithMap(template, *m)
	}

	// 使用ConvertStructToStringMap将args转换为map[string]string，调用NewStringTemplateWithMap
	return NewStringTemplateWithMap(template, ConvertStructToStringMap(args))
}

// NewStringTemplateWithMap 创建一个字符串模板，使用map[string]string作为参数
//   - template: 模板字符串
//   - arguments: 参数
//
// example:
//
//	template := "name: ${name}, age: ${age}"
//	arguments := map[string]string{
//		"name": "test",
//		"age":  "18",
//	}
//	result := NewStringTemplateWithMap(template, arguments).Parse()
//	fmt.Println(result) // name: test, age: 18
func NewStringTemplateWithMap(template string, arguments map[string]string) *StringTemplate {
	return NewStringTemplateWithCustomSignals(template, arguments, "", "")
}

// NewStringTemplateWithCustomSignals 创建一个字符串模板，使用map[string]string作为参数，并且可以自定义模板开始和结束符号
//   - template: 模板字符串
//   - arguments: 参数
//   - startSignals: 模板开始符号，长度必须和endSignals一样
//   - endSignals: 模板结束符号，长度必须和startSignals一样
//
// example:
//
//	template := "name: $@name@, age: $@age@"
//	arguments := map[string]string{
//		"name": "test",
//		"age":  "18",
//	}
//	result := NewStringTemplateWithCustomSignals(template, arguments, "@", "@").Parse()
//	fmt.Println(result) // name: test, age: 18
func NewStringTemplateWithCustomSignals(template string, arguments map[string]string, startSignals, endSignals string) *StringTemplate {
	// not definition by var charSet, avoid duplicate memory allocate
	startSignalsCharSet := defaultStartSignalsCharSet
	endSignalsCharSet := defaultEndSignalsCharSet

	if len(startSignals) > 0 && len(endSignals) > 0 && len(startSignals) == len(endSignals) {
		startSignalsCharSet, endSignalsCharSet = toArrayCharset(startSignals, endSignals)
	}

	rtn := &StringTemplate{
		template:     template,
		arguments:    arguments,
		startSignals: startSignalsCharSet,
		endSignals:   endSignalsCharSet,
	}
	rtn.getVariable = rtn.getVariableUnescaped
	return rtn
}

// NewRawSqlTemplate 创建一个字符串模板，使用支持的内容作为参数，支持的内容有：结构体(指针)，map[string]string，*map[string]string
//   - template: 模板字符串
//   - args: 参数，如果是结构体(指针)，需要指定cpc tag，并且目前只支持uint, int, float, string, bool和嵌套结构体类型的成员字段
//
// example:
//
//	type User struct {
//		Name string `cpc:"key:name"`
//		Age  int    `cpc:"key:age,default:18"`
//	}
//	user := &User{
//		Name: "'; DROP TABLE users WHERE 1=1 --",
//	}
//	template := "SELECT * FROM user WHERE name = '${name}' AND age = ${age}"
//	result := NewRawSqlTemplate(template, user).Parse()
//
// then
//
//	SELECT * FROM user WHERE name = '\'; DROP TABLE users WHERE 1=1 --' AND age = 18
func NewRawSqlTemplate(template string, arguments any) *RawSqlTemplate {
	// 如果args是map[string]string，调用NewRawSqlTemplateWithMap
	if m, ok := arguments.(map[string]string); ok {
		return NewRawSqlTemplateWithMap(template, m)
	}

	// 如果args是*map[string]string，调用NewRawSqlTemplateWithMap
	if m, ok := arguments.(*map[string]string); ok {
		return NewRawSqlTemplateWithMap(template, *m)
	}

	// 使用ConvertStructToStringMap将args转换为map[string]string，调用NewRawSqlTemplateWithMap
	return NewRawSqlTemplateWithMap(template, ConvertStructToStringMap(arguments))
}

// NewRawSqlTemplateWithMap 创建一个字符串模板，使用map[string]string作为参数，并且模板中的参数会被转义
//   - template: 模板字符串
//   - arguments: 参数
//
// example:
//
//	template := "SELECT * FROM user WHERE name = '${name}'"
//	arguments := map[string]string{
//		"name": "test",
//	}
//	result := NewRawSqlTemplateWithMap(template, arguments).Parse()
//	fmt.Println(result) // SELECT * FROM user WHERE name = 'test'
func NewRawSqlTemplateWithMap(template string, arguments map[string]string) *RawSqlTemplate {
	rtn := &RawSqlTemplate{
		&StringTemplate{
			template:     template,
			arguments:    arguments,
			startSignals: defaultRawStartSignalsCharSet,
			endSignals:   defaultRawEndSignalsCharSet,

			escapeCache: make(map[string]string, len(arguments)),
		},
	}
	rtn.getVariable = rtn.getVariableEscaped
	return rtn
}

func (t *StringTemplate) Parse() string {
	template := t.template
	if len(t.arguments) == 0 || len(t.startSignals) == 0 || len(t.endSignals) == 0 {
		// 没有参数或者没有模板标志，原样输出
		return t.template
	}
	var (
		lastCopyIndex int
		variable      string
		exist         bool
		argument      string
		startSignal   byte

		buffer           = make([]string, 0, 16)
		startSignalIndex = -1
		started          = false
	)

	for i, current := range UnsafeStringToBytes(template) {
		switch {
		// 跳过非ascii字符
		case current > 127:
			continue
		// 模板以$开头，用于标记是否开始
		case !started && current == '$':
			started = true
		// 先判断结束标记，因为开始标记和结束标记可能一样
		case t.endSignals[current] != 0 && started:
			// 开始标记和结束标记一样的结果，将当前位置设置为开始位置
			if t.endSignals[current] == current && startSignalIndex == -1 {
				startSignalIndex = i
				startSignal = current
				continue
			}

			// 找到结束标志，获取参数
			if t.endSignals[current] == startSignal {
				variable = template[startSignalIndex+1 : i]
				if len(variable) > 0 {
					// 找到结束标志但是没有给出参数，原样输出当前字符并继续
					exist, argument = t.getVariable(variable, t.arguments)
					if exist {
						if startSignalIndex-lastCopyIndex > 0 {
							// 从上次复制的位置到开始标记$之间的字符
							buffer = append(buffer, template[lastCopyIndex:startSignalIndex-1])
						}
						// 找到参数，输出参数对应值并继续
						buffer = append(buffer, argument)
						lastCopyIndex = i + 1
					}
				}
				// 重置开始标记
				startSignalIndex = -1
				startSignal = 0
				started = false
			}
		case t.startSignals[current] != 0 && started:
			// 匹配到变量模板头
			startSignalIndex = i
			startSignal = current
		}
	}

	if len(buffer) == 0 {
		return template
	}
	return strings.Join(append(buffer, template[lastCopyIndex:]), "")
}

func (t *StringTemplate) escape(raw string) (result string) {
	var lastCopyIndex int
	buffer := make([]string, 0, 16)
	bytes := UnsafeStringToBytes(raw)
	for i, b := range bytes {
		if b == '\'' || b == '"' || b == '\\' {
			if i-lastCopyIndex > 0 {
				buffer = append(buffer, raw[lastCopyIndex:i])
			}
			buffer = append(buffer, transferSymbol)
			lastCopyIndex = i
		}
	}

	if len(buffer) == 0 {
		return raw
	}

	buffer = append(buffer, raw[lastCopyIndex:])
	return strings.Join(buffer, "")
}

func (t *StringTemplate) getVariableUnescaped(template string, arguments map[string]string) (exist bool, result string) {
	result, exist = arguments[template]
	return exist, result
}

func (t *StringTemplate) getVariableEscaped(template string, arguments map[string]string) (exist bool, result string) {
	result, exist = arguments[template]
	if !exist {
		return false, result
	}
	cache, ok := t.escapeCache[template]
	if ok {
		return true, cache
	}
	cache = t.escape(result)
	t.escapeCache[template] = cache
	return true, cache
}

const (
	defaultStartSignals    = `{[`
	defaultEndSignals      = `}]`
	defaultRawStartSignals = "{"
	defaultRawEndSignals   = "}"
	transferSymbol         = `\`
)

var (
	defaultStartSignalsCharSet, defaultEndSignalsCharSet       = toArrayCharset(defaultStartSignals, defaultEndSignals)
	defaultRawStartSignalsCharSet, defaultRawEndSignalsCharSet = toArrayCharset(defaultRawStartSignals, defaultRawEndSignals)
)

type charSet []byte

func toArrayCharset(start, end string) (startSet, endSet charSet) {
	startSet, endSet = make(charSet, 128), make(charSet, 128)
	if len(start) > 0 && len(end) > 0 && len(start) == len(end) {
		startBytes := UnsafeStringToBytes(start)
		endBytes := UnsafeStringToBytes(end)
		for i := 0; i < len(startBytes); i++ {
			startSet[startBytes[i]] = endBytes[i]
			endSet[endBytes[i]] = startBytes[i]
		}
		return startSet, endSet
	}

	return nil, nil
}
