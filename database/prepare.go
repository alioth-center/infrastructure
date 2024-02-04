package database

import "strings"

type Template interface {
	Prepare() (sql string, err error)
}

type template struct {
	template string
	argsMap  map[string]string
}

// escape escapes the raw string for SQL statement
func (t *template) escape(raw string) string {
	builder := strings.Builder{}
	builder.Grow(len(raw) * 2) //预分配足够的空间
	for _, r := range raw {
		switch r {
		case '\'', '\\', '"':
			builder.WriteRune('\\')
			builder.WriteRune(r)
		default:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// findVariable locates the variable within the template and returns it along with its length
func (t *template) findVariable(template []rune, index int, endSignal rune) (variable string, endIndex int, err error) {
	startIndex, length := index, len(template)
	for index < length && template[index] != endSignal {
		if template[index] == ' ' || template[index] == '\t' || template[index] == '\n' || template[index] == '\r' {
			return "", 0, NewTemplateFormatError(template[startIndex:])
		}
		index++
	}
	if index-startIndex > 32 { // variable name too long
		return "", 0, NewTemplateFormatError(template[startIndex:])
	}
	if index == length { // did not find the endSignal
		return "", 0, NewTemplateFormatError(template[startIndex:])
	}
	return string(template[startIndex:index]), index, nil
}

func (t *template) Prepare() (string, error) {
	templates := []rune(t.template)
	builder, length := strings.Builder{}, len(templates)
	builder.Grow(length)

	var index, prefixLength int
	var endSignal rune

	for index < length {
		current := templates[index]
		switch current {
		case '$': // Potential start of a variable
			if index+1 < length {
				next := templates[index+1]
				if next == '{' || next == '[' {
					prefixLength = 2
					if next == '{' {
						endSignal = '}'
					} else {
						endSignal = ']'
					}
				} else {
					builder.WriteRune(current)
					index++
					continue
				}
			} else {
				builder.WriteRune(current)
				index++
				continue
			}
		case '{', '[': // Start of a variable without $
			prefixLength = 1
			endSignal = '}'
			if current == '[' {
				endSignal = ']'
			}
		default:
			builder.WriteRune(current)
			index++
			continue
		}

		variable, endIndex, err := t.findVariable(templates, index+prefixLength, endSignal)
		if err != nil {
			return "", err
		}

		value, exist := t.argsMap[variable]
		if !exist {
			// Variable not found, keep the original placeholder
			builder.WriteString(string(templates[index : endIndex+1]))
		} else if value != "" {
			// Variable found, replace the placeholder with the value, but if the value is empty, ignore it
			builder.WriteString(t.escape(value))
		}

		index = endIndex + 1 // Move past the end of the variable
	}

	return builder.String(), nil
}

func NewTemplate(tmp string, argsMap map[string]string) Template {
	t := &template{
		template: tmp,
		argsMap:  argsMap,
	}

	return t
}
