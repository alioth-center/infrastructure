package values

import (
	"reflect"
	"strconv"
	"strings"
)

func CheckStruct(structure any, tag ...string) string {
	t := checkBodyTag
	if len(tag) > 0 {
		t = tag[0]
	}

	return checkBody("", reflect.ValueOf(structure), t)
}

const (
	checkBodyTag         = "vc"
	checkBodyTagKey      = "key:"
	checkBodyTagRequired = "required"
)

type checkTagConfig struct {
	key      string
	required bool
}

func parseCheckTag(tag string) *checkTagConfig {
	config := &checkTagConfig{}

	if tag == "" {
		return config
	}

	tagValues := strings.Split(tag, ",")
	for _, tagValue := range tagValues {
		switch {
		case tagValue == checkBodyTagRequired:
			config.required = true
		case strings.HasPrefix(tagValue, checkBodyTagKey):
			config.key = strings.TrimPrefix(tagValue, checkBodyTagKey)
		}
	}

	return config
}

func checkBody(prefix string, value reflect.Value, tag string) string {
	if tag == "" {
		tag = checkBodyTag
	}

	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)
			fieldType := value.Type().Field(i)

			cfg := parseCheckTag(fieldType.Tag.Get(tag))
			if cfg.key == "" {
				continue
			}
			fieldName := cfg.key
			if prefix != "" {
				fieldName = BuildStrings(prefix, ".", fieldName)
			}

			for fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			if fieldValue.Kind() == reflect.Struct {
				subResult := checkBody(fieldName, fieldValue, tag)
				if subResult != "" {
					return subResult
				}
			}

			if cfg.required && fieldValue.IsZero() {
				return fieldName
			}
		}
	}

	return ""
}

// ConvertStructToStringMap 函数将结构体转换为map[string]string，为了template使用
//   - val: 结构体，需要指定cpc tag
//   - keyPrefix: key的前缀，选填
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
//	result := ConvertStructToStringMap(user) // map[male:true name:test]
func ConvertStructToStringMap(val any, keyPrefix ...string) map[string]string {
	value := reflect.ValueOf(val)

	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// 提前退出
	if value.Kind() == reflect.Invalid {
		return map[string]string{}
	}
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return map[string]string{}
	}

	prefix, result := "", make(map[string]string, value.NumField())
	if len(keyPrefix) > 0 {
		prefix = keyPrefix[0]
	}

	structToMap(prefix, value, result)
	return result
}

func structToMap(prefix string, value reflect.Value, m map[string]string) {
	// 如果是指针类型，则获取其指向的值
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// 目前只处理结构体类型
	if value.Kind() == reflect.Struct {
		// 遍历结构体的所有字段
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)
			fieldType := value.Type().Field(i)

			// 获取字段的转换配置
			cpc := getConvertParserConfig(fieldType.Tag.Get(convertParserConfigTagName))

			// 如果字段被忽略，则跳过
			if cpc.ignore {
				continue
			}

			// 如果未指定字段的名称，则使用字段名称
			if cpc.key == "" {
				cpc.key = fieldType.Name
			}

			fieldName := cpc.key
			if prefix != "" {
				fieldName = prefix + "." + fieldName
			}

			// 如果字段是指针类型，则获取其指向的值
			for fieldValue.Kind() == reflect.Pointer {
				fieldValue = fieldValue.Elem()
			}

			// 如果字段是结构体，则递归调用structToMap，将其子字段添加到map中
			if fieldValue.Kind() == reflect.Struct {
				structToMap(fieldName, fieldValue, m)
				continue
			}

			// 根据字段的类型进行相应的处理，将字段值转换为字符串
			convertedValue, converted := convertFieldValueToString(fieldValue)
			// 对于不支持的类型，跳过
			if !converted {
				continue
			}

			// 处理omitempty标签
			if cpc.omitEmpty && fieldValue.IsZero() {
				continue
			}

			// 如果转换后的值为空，并且有默认值，则使用默认值
			if fieldValue.IsZero() && cpc.defaultValue != "" {
				convertedValue = cpc.defaultValue
			}

			// 将字段名和值添加到map中
			m[fieldName] = convertedValue
		}
	}
}

// convertFieldValueToString 函数根据字段的类型将其值转换为字符串
func convertFieldValueToString(fieldValue reflect.Value) (string, bool) {
	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(fieldValue.Int(), 10), true
	case reflect.String:
		return fieldValue.String(), true
	case reflect.Bool:
		return strconv.FormatBool(fieldValue.Bool()), true
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(fieldValue.Uint(), 10), true
	default:
		return "", false
	}
}

const (
	convertParserConfigTagName             = "cpc"
	convertParserConfigIgnoreTag           = "-"
	convertParserConfigKeyTagName          = "key:"
	convertParserConfigOmitShortTagName    = "oe"
	convertParserConfigOmitTagName         = "omitempty"
	convertParserConfigDefaultShortTagName = "dft:"
	convertParserConfigDefaultTagName      = "default:"
)

type convertParserConfig struct {
	ignore       bool
	key          string
	omitEmpty    bool
	defaultValue string
}

func getConvertParserConfig(tag string) *convertParserConfig {
	if tag == "" || tag == convertParserConfigIgnoreTag {
		return &convertParserConfig{ignore: true}
	}

	config := &convertParserConfig{}
	tags := strings.Split(tag, ",")
	for _, tag := range tags {
		switch {
		case tag == convertParserConfigOmitShortTagName || tag == convertParserConfigOmitTagName:
			config.omitEmpty = true
		case strings.HasPrefix(tag, convertParserConfigKeyTagName):
			config.key = strings.TrimPrefix(tag, convertParserConfigKeyTagName)
		case strings.HasPrefix(tag, convertParserConfigDefaultShortTagName):
			config.defaultValue = strings.TrimPrefix(tag, convertParserConfigDefaultShortTagName)
		case strings.HasPrefix(tag, convertParserConfigDefaultTagName):
			config.defaultValue = strings.TrimPrefix(tag, convertParserConfigDefaultTagName)
		}
	}

	return config
}
