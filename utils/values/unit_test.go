package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Run("BuildStrings", func(t *testing.T) {
		args := []string{"a", "b", "c"}
		result := BuildStrings(args...)
		if result != "abc" {
			t.Errorf("BuildStrings() = %v, want %v", result, "abc")
		}
	})

	t.Run("BuildStringsWithJoin", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			args := []string{"a", "b", "c"}
			result := BuildStringsWithJoin("-", args...)
			if result != "a-b-c" {
				t.Errorf("BuildStringsWIthJoin() = %v, want %v", result, "a-b-c")
			}
		})

		t.Run("NilParts", func(t *testing.T) {
			var args []string = nil
			result := BuildStringsWithJoin(".", args...)
			if result != "" {
				t.Errorf("BuildStringsWIthJoin() = %v, want %v", result, "(empty string)")
			}
		})

		t.Run("EmptyParts", func(t *testing.T) {
			args := []string{}
			result := BuildStringsWithJoin(".", args...)
			if result != "" {
				t.Errorf("BuildStringsWIthJoin() = %v, want %v", result, "(empty string)")
			}
		})
	})

	t.Run("BuildStringsWithJoinIgnoreEmpty", func(t *testing.T) {
		args := []string{"a", "", "b", "c", ""}
		result := BuildStringsWithJoinIgnoreEmpty("-", args...)
		if result != "a-b-c" {
			t.Errorf("BuildStringsWithJoinIgnoreEmpty() = %v, want %v", result, "a-b-c")
		}
	})

	t.Run("BuildStringsWithReplacement", func(t *testing.T) {
		args := []string{"a", "b", "c"}
		replacement := map[string]string{"a": "1", "b": "2", "c": "3"}
		result := BuildStringsWithReplacement(replacement, args...)
		if result != "123" {
			t.Errorf("BuildStringsWithReplacement() = %v, want %v", result, "123")
		}
	})

	t.Run("BuildStringsWithTemplate", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			args := []string{"a", "b", "c"}
			template := "{1}-{2}-{3}"
			result := BuildStringsWithTemplate(template, args...)
			if result != "a-b-c" {
				t.Errorf("BuildStringsWithTemplate() = %v, want %v", result, "a-b-c")
			}
		})

		t.Run("NilArgs", func(t *testing.T) {
			var args []string = nil
			template := "{1}-{2}-{3}"
			result := BuildStringsWithTemplate(template, args...)
			if result != "{1}-{2}-{3}" {
				t.Errorf("BuildStringsWithTemplate() = %v, want %v", result, "{1}-{2}-{3}")
			}
		})

		t.Run("EmptyArgs", func(t *testing.T) {
			args := []string{}
			template := "{1}-{2}-{3}"
			result := BuildStringsWithTemplate(template, args...)
			if result != "{1}-{2}-{3}" {
				t.Errorf("BuildStringsWithTemplate() = %v, want %v", result, "{1}-{2}-{3}")
			}
		})
	})

	t.Run("SecretString", func(t *testing.T) {
		testCases := []struct {
			name   string
			secret string
			prefix int
			suffix int
			char   string
			want   string
		}{
			{
				name:   "Normal",
				secret: "1234567890",
				prefix: 3,
				suffix: 3,
				char:   "*",
				want:   "123****890",
			},
			{
				name:   "Empty",
				secret: "",
				prefix: 3,
				suffix: 3,
				char:   "*",
				want:   "",
			},
			{
				name:   "Short",
				secret: "123",
				prefix: 3,
				suffix: 3,
				char:   "*",
				want:   "123",
			},
			{
				name:   "Prefix",
				secret: "123",
				prefix: 1,
				suffix: 0,
				char:   "*",
				want:   "1**",
			},
			{
				name:   "Suffix",
				secret: "123",
				prefix: 0,
				suffix: 1,
				char:   "*",
				want:   "**3",
			},
			{
				name:   "Hans",
				secret: "我是刻晴",
				prefix: 1,
				suffix: 2,
				char:   "*",
				want:   "我*刻晴",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := SecretString(tc.secret, tc.prefix, tc.suffix, tc.char)
				if result != tc.want {
					t.Errorf("SecretString() = %v, want %v", result, tc.want)
				}
			})
		}
	})

	t.Run("StringToInt", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToInt("123", 0)
			if result != 123 {
				t.Errorf("StringToInt() = %v, want %v", result, 123)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			result := StringToInt("", 123)
			if result != 123 {
				t.Errorf("StringToInt() = %v, want %v", result, 123)
			}
		})

		t.Run("Error", func(t *testing.T) {
			result := StringToInt("fuck you", 123)
			if result != 123 {
				t.Errorf("StringToInt() = %v, want %v", result, 123)
			}
		})
	})

	t.Run("StringToUint", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToUint("123", uint(0))
			if result != 123 {
				t.Errorf("StringToUint() = %v, want %v", result, 123)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			result := StringToUint("", uint(123))
			if result != 123 {
				t.Errorf("StringToUint() = %v, want %v", result, 123)
			}
		})

		t.Run("Error", func(t *testing.T) {
			result := StringToUint("fuck you", uint(123))
			if result != 123 {
				t.Errorf("StringToUint() = %v, want %v", result, 123)
			}
		})
	})

	t.Run("StringToFloat64", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToFloat64("123.456", 0.0)
			if result != 123.456 {
				t.Errorf("StringToFloat64() = %v, want %v", result, 123.456)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			result := StringToFloat64("", 123.456)
			if result != 123.456 {
				t.Errorf("StringToFloat64() = %v, want %v", result, 123.456)
			}
		})

		t.Run("Error", func(t *testing.T) {
			result := StringToFloat64("fuck you", 123.456)
			if result != 123.456 {
				t.Errorf("StringToFloat64() = %v, want %v", result, 123.456)
			}
		})
	})

	t.Run("StringToBool", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToBool("true", false)
			if result != true {
				t.Errorf("StringToBool() = %v, want %v", result, true)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			result := StringToBool("", true)
			if result != true {
				t.Errorf("StringToBool() = %v, want %v", result, true)
			}
		})

		t.Run("Error", func(t *testing.T) {
			result := StringToBool("fuck", true)
			if result != true {
				t.Errorf("StringToBool() = %v, want %v", result, true)
			}
		})
	})

	t.Run("StringToStringPtr", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToStringPtr[string]("123")
			if *result != "123" {
				t.Errorf("StringToStringPtr() = %v, want %v", *result, "123")
			}
		})

		t.Run("Empty", func(t *testing.T) {
			result := StringToStringPtr[string]("")
			if result != nil {
				t.Errorf("StringToStringPtr() = %v, want nil", result)
			}
		})
	})

	t.Run("StringToIntPtr", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToIntPtr[int]("123")
			if *result != 123 {
				t.Errorf("StringToIntPtr() = %v, want 123", result)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			val := 1
			result := &val
			result = StringToIntPtr[int]("")
			if result != nil {
				t.Errorf("StringToIntPtr() = %v, want nil", result)
			}
		})

		t.Run("Error", func(t *testing.T) {
			val := 1
			result := &val
			result = StringToIntPtr[int]("fuck you")
			if result != nil {
				t.Errorf("StringToIntPtr() = %v, want nil", result)
			}
		})
	})

	t.Run("StringToUintPtr", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToUintPtr[uint]("123")
			if *result != 123 {
				t.Errorf("StringToUintPtr() = %v, want 123", result)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			val := uint(1)
			result := &val
			result = StringToUintPtr[uint]("")
			if result != nil {
				t.Errorf("StringToUintPtr() = %v, want nil", result)
			}
		})

		t.Run("Error", func(t *testing.T) {
			val := uint(1)
			result := &val
			result = StringToUintPtr[uint]("fuck you")
			if result != nil {
				t.Errorf("StringToUintPtr() = %v, want nil", result)
			}
		})
	})

	t.Run("StringToFloat64Ptr", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToFloat64Ptr[float64]("123.456")
			if *result != 123.456 {
				t.Errorf("StringToFloat64Ptr() = %v, want 123.456", result)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			val := 123.456
			result := &val
			result = StringToFloat64Ptr[float64]("")
			if result != nil {
				t.Errorf("StringToFloat64Ptr() = %v, want nil", result)
			}
		})

		t.Run("Error", func(t *testing.T) {
			val := 123.456
			result := &val
			result = StringToFloat64Ptr[float64]("fuck you")
			if result != nil {
				t.Errorf("StringToFloat64Ptr() = %v, want nil", result)
			}
		})
	})

	t.Run("StringToBoolPtr", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			result := StringToBoolPtr[bool]("true")
			if *result != true {
				t.Errorf("StringToBoolPtr() = %v, want true", result)
			}
		})

		t.Run("Empty", func(t *testing.T) {
			val := true
			result := &val
			result = StringToBoolPtr[bool]("")
			if result != nil {
				t.Errorf("StringToBoolPtr() = %v, want nil", result)
			}
		})

		t.Run("Error", func(t *testing.T) {
			val := true
			result := &val
			result = StringToBoolPtr[bool]("fuck you")
			if result != nil {
				t.Errorf("StringToBoolPtr() = %v, want nil", result)
			}
		})
	})

	t.Run("StringTemplate", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("hello ${template_1}, i am ${template_2}, nice to meet $[template_3], and ${template_4}", map[string]string{
				"template_1": "value_1",
				"template_2": "value_2",
				"template_3": "value_3",
				"template_4": "value_4",
			})

			want := "hello value_1, i am value_2, nice to meet value_3, and value_4"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("Include Not ASCII", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("你好 ${template_1}, 我是 ${template_2}, nice to meet $[template_3], and $[template_4]", map[string]string{
				"template_1": "世界",
				"template_2": "value_2",
				"template_3": "value_3",
				"template_4": "value_4",
			})

			want := "你好 世界, 我是 value_2, nice to meet value_3, and value_4"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("LongVariableName", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("hello, ${this_is_a_very_very_very_very_very_very_very_very_very_looooooooooong_template}", map[string]string{
				"this_is_a_very_very_very_very_very_very_very_very_very_looooooooooong_template": "value",
			})

			want := "hello, value"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("EndWithSignal", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("hello, ${template}! {", map[string]string{
				"template": "value",
			})

			want := "hello, value! {"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("SignalNotFound", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("hello, $template!", map[string]string{
				"tempalte": "value",
			})

			want := "hello, $template!"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("VariableNotExist", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("hello, ${template_1}! nice to meet ${template_2}", map[string]string{
				"template_1": "value",
			})

			want := "hello, value! nice to meet ${template_2}"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("NilArguments", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("hello, ${template_1}! nice to meet ${template_2}", nil)

			want := "hello, ${template_1}! nice to meet ${template_2}"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("Mismatch Signals", func(t *testing.T) {
			tmp := NewStringTemplateWithMap("hello ${template_1], i am $[template_2}", map[string]string{
				"template_1": "value_1",
				"template_2": "value_2",
			})

			want := "hello ${template_1], i am $[template_2}"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("CustomSignal", func(t *testing.T) {
			tmp := NewStringTemplateWithCustomSignals("hello $<template_1>, i am ${template_2}, nice to meet $[template_3], and $[template_4]",
				map[string]string{
					"template_1": "value_1",
					"template_2": "value_2",
					"template_3": "value_3",
					"template_4": "value_4",
				}, "<", ">",
			)

			want := "hello value_1, i am ${template_2}, nice to meet $[template_3], and $[template_4]"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("CustomSignalButEmpty", func(t *testing.T) {
			tmp := NewStringTemplateWithCustomSignals("hello ${template_1}, i am ${template_2}, nice to meet $[template_3], and $[template_4]",
				map[string]string{
					"template_1": "value_1",
					"template_2": "value_2",
					"template_3": "value_3",
					"template_4": "value_4",
				}, "", "",
			)

			want := "hello value_1, i am value_2, nice to meet value_3, and value_4"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("TemplateWithStructure", func(t *testing.T) {
			type testStruct struct {
				Template1 string `cpc:"key:template_1"`
				Template2 string `cpc:"key:template_2,omitempty"`
				Template3 string `cpc:"key:template_3,default:hello"`
			}

			tmp := NewStringTemplate("hello ${template_1}, i am ${template_2}, i say ${template_3}", &testStruct{
				Template1: "value_1",
			})
			want := "hello value_1, i am ${template_2}, i say hello"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("TemplateWithStructureButGiveNil", func(t *testing.T) {
			tmp := NewStringTemplate("hello ${template_1}, i am ${template_2}, i say {template_3}", nil)
			want := "hello ${template_1}, i am ${template_2}, i say {template_3}"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("TemplateWithEmbeddingStructure", func(t *testing.T) {
			type embeddingStruct struct {
				Template1 string `cpc:"key:template_1"`
			}

			type testStruct struct {
				Template1 string           `cpc:"key:template_1"`
				Template2 string           `cpc:"key:template_2,omitempty"`
				Template3 string           `cpc:"key:template_3,default:hello"`
				SubStruct *embeddingStruct `cpc:"key:sub_struct"`
			}

			tmp := NewStringTemplate("hello ${template_1}, i am ${template_2}, i say ${template_3}, and i have ${sub_struct.template_1}", &testStruct{
				Template1: "value_1",
				SubStruct: &embeddingStruct{
					Template1: "pencil",
				},
			})
			want := "hello value_1, i am ${template_2}, i say hello, and i have pencil"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("TemplateWithEmbeddingStructureButGiveNil", func(t *testing.T) {
			type embeddingStruct struct {
				Template1 string `cpc:"key:template_1"`
			}

			type testStruct struct {
				Template1 string           `cpc:"key:template_1"`
				Template2 string           `cpc:"key:template_2,omitempty"`
				Template3 string           `cpc:"key:template_3,default:hello"`
				SubStruct *embeddingStruct `cpc:"key:sub_struct"`
			}

			tmp := NewStringTemplate("hello ${template_1}, i am ${template_2}, i say ${template_3}, and i have ${sub_struct.template_1}", &testStruct{
				Template1: "value_1",
				SubStruct: nil,
			})
			want := "hello value_1, i am ${template_2}, i say hello, and i have ${sub_struct.template_1}"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("TemplateWithStructureButGiveMap", func(t *testing.T) {
			tmp := NewStringTemplate("hello ${template_1}, i am ${template_2}, i say ${template_3}", map[string]string{
				"template_1": "value_1",
				"template_3": "hello",
			})
			want := "hello value_1, i am ${template_2}, i say hello"
			assert.Equal(t, want, tmp.Parse())
		})

		t.Run("TemplateWithStructureButGiveNilMap", func(t *testing.T) {
			var mp map[string]string = nil
			var m any = mp
			tmp := NewStringTemplate("hello ${template_1}, i am ${template_2}, i say ${template_3}", m)
			want := "hello ${template_1}, i am ${template_2}, i say ${template_3}"
			assert.Equal(t, want, tmp.Parse())
		})
	})
}

func TestNumber(t *testing.T) {
	t.Run("IntToString", func(t *testing.T) {
		result := IntToString(123)
		if result != "123" {
			t.Errorf("IntToString() = %v, want 123", result)
		}
	})

	t.Run("Int64ToString", func(t *testing.T) {
		result := Int64ToString(int64(123))
		if result != "123" {
			t.Errorf("Int64ToString() = %v, want 123", result)
		}
	})

	t.Run("UintToString", func(t *testing.T) {
		result := UintToString(uint(123))
		if result != "123" {
			t.Errorf("UintToString() = %v, want 123", result)
		}
	})

	t.Run("Uint64ToString", func(t *testing.T) {
		result := Uint64ToString(uint64(123))
		if result != "123" {
			t.Errorf("Uint64ToString() = %v, want 123", result)
		}
	})

	t.Run("Float32ToString", func(t *testing.T) {
		result := Float32ToString(float32(123.456))
		if result != "123.456" {
			t.Errorf("Float32ToString() = %v, want 123.456", result)
		}
	})

	t.Run("Float64ToString", func(t *testing.T) {
		result := Float64ToString(123.456)
		if result != "123.456" {
			t.Errorf("Float64ToString() = %v, want 123.456", result)
		}
	})
}

func TestMultiply(t *testing.T) {
	t.Run("NumberOrStringValueGetString", func(t *testing.T) {
		testCases := []struct {
			Case  string
			Value any
			Want  string
		}{
			{
				Case:  "StringType",
				Value: "string",
				Want:  "string",
			},
			{
				Case:  "IntType",
				Value: 123,
				Want:  "123",
			},
			{
				Case:  "Int8Type",
				Value: int8(123),
				Want:  "123",
			},
			{
				Case:  "Int16Type",
				Value: int16(123),
				Want:  "123",
			},
			{
				Case:  "Int32Type",
				Value: int32(123),
				Want:  "123",
			},
			{
				Case:  "Int64Type",
				Value: int64(123),
				Want:  "123",
			},
			{
				Case:  "UintType",
				Value: uint(233),
				Want:  "233",
			},
			{
				Case:  "Uint8Type",
				Value: uint8(233),
				Want:  "233",
			},
			{
				Case:  "Uint16Type",
				Value: uint16(233),
				Want:  "233",
			},
			{
				Case:  "Uint32Type",
				Value: uint32(233),
				Want:  "233",
			},
			{
				Case:  "Uint64Type",
				Value: uint64(233),
				Want:  "233",
			},
			{
				Case:  "Float32Type",
				Value: float32(123.456),
				Want:  "123.456",
			},
			{
				Case:  "Float64Type",
				Value: 123.456,
				Want:  "123.456",
			},
			{
				Case:  "BoolType",
				Value: true,
				Want:  "true",
			},
			{
				Case:  "OtherType",
				Value: nil,
				Want:  "",
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.Case, func(t *testing.T) {
				if NumberOrStringValueGetString(testCase.Value) != testCase.Want {
					t.Errorf("NumberOrStringValueGetString() = %v, want %v", NumberOrStringValueGetString(testCase.Value), testCase.Want)
				}
			})
		}
	})

	t.Run("NumberOrStringValueGetInt", func(t *testing.T) {
		testCases := []struct {
			Case  string
			Value any
			Want  int
		}{
			{
				Case:  "StringType",
				Value: "string",
				Want:  0,
			},
			{
				Case:  "IntType",
				Value: 123,
				Want:  123,
			},
			{
				Case:  "Int8Type",
				Value: int8(123),
				Want:  123,
			},
			{
				Case:  "Int16Type",
				Value: int16(123),
				Want:  123,
			},
			{
				Case:  "Int32Type",
				Value: int32(123),
				Want:  123,
			},
			{
				Case:  "Int64Type",
				Value: int64(123),
				Want:  123,
			},
			{
				Case:  "UintType",
				Value: uint(233),
				Want:  233,
			},
			{
				Case:  "Uint8Type",
				Value: uint8(233),
				Want:  233,
			},
			{
				Case:  "Uint16Type",
				Value: uint16(233),
				Want:  233,
			},
			{
				Case:  "Uint32Type",
				Value: uint32(233),
				Want:  233,
			},
			{
				Case:  "Uint64Type",
				Value: uint64(233),
				Want:  233,
			},
			{
				Case:  "Float32Type",
				Value: float32(123.456),
				Want:  123,
			},
			{
				Case:  "Float64Type",
				Value: 123.456,
				Want:  123,
			},
			{
				Case:  "BoolType",
				Value: true,
				Want:  1,
			},
			{
				Case:  "BoolType2",
				Value: false,
				Want:  0,
			},
			{
				Case:  "OtherType",
				Value: nil,
				Want:  0,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.Case, func(t *testing.T) {
				if NumberOrStringValueGetInt(testCase.Value) != testCase.Want {
					t.Errorf("NumberOrStringValueGetInt() = %v, want %v", NumberOrStringValueGetInt(testCase.Value), testCase.Value)
				}
			})
		}
	})

	t.Run("NumberOrStringValueGetFloat", func(t *testing.T) {
		testCases := []struct {
			Case  string
			Value any
			Want  float64
		}{
			{
				Case:  "StringType",
				Value: "string",
				Want:  0,
			},
			{
				Case:  "IntType",
				Value: 123,
				Want:  123,
			},
			{
				Case:  "Int8Type",
				Value: int8(123),
				Want:  123,
			},
			{
				Case:  "Int16Type",
				Value: int16(123),
				Want:  123,
			},
			{
				Case:  "Int32Type",
				Value: int32(123),
				Want:  123,
			},
			{
				Case:  "Int64Type",
				Value: int64(123),
				Want:  123,
			},
			{
				Case:  "UintType",
				Value: uint(233),
				Want:  233,
			},
			{
				Case:  "Uint8Type",
				Value: uint8(233),
				Want:  233,
			},
			{
				Case:  "Uint16Type",
				Value: uint16(233),
				Want:  233,
			},
			{
				Case:  "Uint32Type",
				Value: uint32(233),
				Want:  233,
			},
			{
				Case:  "Uint64Type",
				Value: uint64(233),
				Want:  233,
			},
			{
				Case:  "Float32Type",
				Value: float32(123),
				Want:  123,
			},
			{
				Case:  "Float64Type",
				Value: 123.456,
				Want:  123.456,
			},
			{
				Case:  "BoolType",
				Value: true,
				Want:  1,
			},
			{
				Case:  "BoolType2",
				Value: false,
				Want:  0,
			},
			{
				Case:  "OtherType",
				Value: nil,
				Want:  0,
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.Case, func(t *testing.T) {
				if NumberOrStringValueGetFloat(testCase.Value) != testCase.Want {
					t.Errorf("NumberOrStringValueGetFloat() = %v, want %v", NumberOrStringValueGetFloat(testCase.Value), testCase.Want)
				}
			})
		}
	})
}

func TestDefault(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		Nil[int]()
	})

	t.Run("NilFunction", func(t *testing.T) {
		NilFunction[int]()()
	})

	t.Run("Ptr", func(t *testing.T) {
		b := true
		bPtr := Ptr(&b)
		if **bPtr != b {
			t.Errorf("Ptr() = %v, want %v", **bPtr, b)
		}
	})
}

func TestArray(t *testing.T) {
	t.Run("Reverse", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			arr := []int{1, 2, 3, 4, 5}
			rev := []int{5, 4, 3, 2, 1}
			for i, v := range ReverseArray(arr) {
				if rev[i] != v {
					t.Errorf("ReverseArray[%d] = %v, want %v", i, v, rev[i])
				}
			}
		})

		t.Run("Empty", func(t *testing.T) {
			arr := []int{}
			rev := ReverseArray(arr)
			if rev == nil || len(rev) != 0 {
				t.Errorf("ReverseArray = %v, want %v", rev, arr)
			}
		})
	})

	t.Run("Contains", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		if !ContainsArray(arr, 3) {
			t.Errorf("ContainsArray() = %v, want %v", ContainsArray(arr, 3), true)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
		unique := []int{1, 2, 3, 4, 5}
		if len(UniqueArray(arr)) != len(unique) {
			t.Errorf("UniqueArray() = %v, want %v", len(UniqueArray(arr)), len(unique))
		}

		for i, v := range UniqueArray(arr) {
			if unique[i] != v {
				t.Errorf("UniqueArray[%d] = %v, want %v", i, v, unique[i])
			}
		}
	})

	t.Run("Remove", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		removed := []int{1, 2, 3, 5}
		for i, v := range RemoveArray(arr, 3, 4) {
			if removed[i] != v {
				t.Errorf("RemoveArray[%d] = %v, want %v", i, v, removed[i])
			}
		}
	})

	t.Run("Index", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		if IndexArray(arr, 3) != 2 {
			t.Errorf("IndexArray() = %v, want %v", IndexArray(arr, 3), 2)
		}
	})

	t.Run("Shuffle", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
		clone := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
		shuffled := ShuffleArray(arr, 0, len(arr))
		if len(shuffled) != len(arr) {
			t.Errorf("ShuffleArray() = %v, want %v", len(shuffled), len(arr))
		}

		same := true
		for i, v := range shuffled {
			if v != clone[i] {
				same = false
				break
			}
		}

		if same {
			t.Errorf("ShuffleArray() = %v, want not %v", shuffled, clone)
		}
	})

	t.Run("Merge", func(t *testing.T) {
		arr1 := []int{1, 2, 3, 4, 5}
		arr2 := []int{6, 7, 8, 9, 10}
		arr3 := []int{11, 12, 13, 14, 15}
		merged := MergeArrays(arr1, arr2, arr3)
		if len(merged) != 15 {
			t.Errorf("MergeArrays() = %v, want %v", len(merged), 15)
		}

		for i, v := range merged {
			if i < 5 && v != arr1[i] {
				t.Errorf("MergeArrays[%d] = %v, want %v", i, v, arr1[i])
			} else if i >= 5 && i < 10 && v != arr2[i-5] {
				t.Errorf("MergeArrays[%d] = %v, want %v", i, v, arr2[i-5])
			} else if i >= 10 && v != arr3[i-10] {
				t.Errorf("MergeArrays[%d] = %v, want %v", i, v, arr3[i-10])
			}
		}
	})

	t.Run("Filter", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		filtered := FilterArray(arr, func(v int) bool {
			return v%2 == 0
		})
		if len(filtered) != 2 {
			t.Errorf("FilterArray() = %v, want %v", len(filtered), 2)
		}

		for i, v := range filtered {
			if v != arr[i*2+1] {
				t.Errorf("FilterArray[%d] = %v, want %v", i, v, arr[i*2+1])
			}
		}
	})

	t.Run("Sort", func(t *testing.T) {
		arr := []int{5, 4, 3, 2, 1}
		sorted := SortArray(arr, func(a, b int) bool {
			return a < b
		})
		for i, v := range sorted {
			if v != i+1 {
				t.Errorf("SortArray[%d] = %v, want %v", i, v, i+1)
			}
		}
	})

	t.Run("LastOf", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		if LastOfArray(arr) != 5 {
			t.Errorf("LastOfArray() = %v, want %v", LastOfArray(arr), 5)
		}

		if LastOfArray([]int{}) != 0 {
			t.Errorf("LastOfArray() = %v, want %v", LastOfArray([]int{}), 0)
		}
	})

	t.Run("TopN", func(t *testing.T) {
		t.Run("TopN-30Percent", func(t *testing.T) {
			array := []int{10, 80, 30, 90, 40, 50, 70, 20, 60, 100}
			percentN := 30
			result := TopNArray(array, percentN, func(a, b int) bool {
				return a < b
			})
			if result != 30 {
				t.Errorf("TopNArray() = %v, want %v", result, 30)
			}
		})

		t.Run("EmptyArray", func(t *testing.T) {
			var array []int
			percentN := 50
			result := TopNArray(array, percentN, func(a, b int) bool {
				return a < b
			})
			var want int
			if result != want {
				t.Errorf("TopNArray() = %v, want %v", result, want)
			}
		})

		t.Run("SingleElement", func(t *testing.T) {
			array := []int{42}
			percentN := 50
			result := TopNArray(array, percentN, func(a, b int) bool {
				return a < b
			})
			want := 42
			if result != want {
				t.Errorf("TopNArray() = %v, want %v", result, want)
			}
		})

		t.Run("PercentN-LessThanZero", func(t *testing.T) {
			array := []int{10, 20}
			percentN := -10
			result := TopNArray(array, percentN, func(a, b int) bool {
				return a < b
			})
			want := 10
			if result != want {
				t.Errorf("TopNArray() = %v, want %v", result, want)
			}
		})

		t.Run("PercentN-GreaterThan100", func(t *testing.T) {
			array := []int{10, 20}
			percentN := 150
			result := TopNArray(array, percentN, func(a, b int) bool {
				return a < b
			})
			want := 20
			if result != want {
				t.Errorf("TopNArray() = %v, want %v", result, want)
			}
		})

		t.Run("ArrayWithDuplicates", func(t *testing.T) {
			array := []int{10, 10, 20, 20, 30, 30}
			percentN := 50
			result := TopNArray(array, percentN, func(a, b int) bool {
				return a < b
			})
			want := 20
			if result != want {
				t.Errorf("TopNArray() = %v, want %v", result, want)
			}
		})

		t.Run("TestWithFloats", func(t *testing.T) {
			array := []float64{1.1, 2.2, 3.3, 4.4, 5.5}
			percentN := 40
			result := TopNArray(array, percentN, func(a, b float64) bool {
				return a < b
			})
			want := 2.2
			if result != want {
				t.Errorf("TopNArray() = %v, want %v", result, want)
			}
		})
	})
}

func TestReflect(t *testing.T) {
	t.Run("CheckStruct", func(t *testing.T) {
		t.Run("CheckSuccess", func(t *testing.T) {
			type successStruct struct {
				Name string         `vc:"key:name,required"`
				Age  int            `vc:"key:age,required"`
				st   *successStruct `vc:"key:st"`
			}

			s := &successStruct{
				Name: "test",
				Age:  18,
				st: &successStruct{
					Name: "test2",
					Age:  19,
					st:   nil,
				},
			}

			if CheckStruct(s) != "" {
				t.Errorf("CheckStruct() = %v, want %v", CheckStruct(s), "")
			}
		})

		t.Run("CheckFailed", func(t *testing.T) {
			type subStruct struct {
				Name string `vc:"key:name,required"`
				Age  int    `vc:"key:age"`
			}

			type failedStruct struct {
				Name string    `vc:"key:name_u,required"`
				Age  int       `vc:"key:age_u,required"`
				st   subStruct `vc:"key:sub"`
			}

			f := &failedStruct{
				Name: "1",
				Age:  18,
				st: subStruct{
					Age: 1,
				},
			}

			if CheckStruct(f) != "sub.name" {
				t.Errorf("CheckStruct() = %v, want %v", CheckStruct(f), "st.name")
			}
		})
	})
}

func TestUnsafe(t *testing.T) {
	t.Run("StringToByte", func(t *testing.T) {
		UnsafeStringToBytes("")
		UnsafeStringToBytes("114514")
	})

	t.Run("ByteToString", func(t *testing.T) {
		UnsafeBytesToString([]byte{})
		UnsafeBytesToString([]byte("114514"))
	})
}
