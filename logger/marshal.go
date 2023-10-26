package logger

import (
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/errors"
	"strings"
	"time"
)

var (
	JsonMarshaller Marshaller = jsonMarshaller
	TextMarshaller Marshaller = textMarshaller
	CsvMarshaller  Marshaller = csvMarshaller
	TsvMarshaller  Marshaller = tsvMarshaller
	// LokiMarshaller Marshaller = lokiMarshaller 暂未支持
)

type Marshaller func(entry *Entry) ([]byte, error)

func jsonMarshaller(entry *Entry) ([]byte, error) {
	if entry.File == "" {
		entry.File = "unset"
	}
	if entry.Level == "" {
		entry.Level = string(LevelInfo)
	}
	if entry.CallTime == "" {
		entry.CallTime = time.Now().Format(timeFormat)
	}
	if entry.Service == "" {
		entry.Service = "unset"
	}

	bytes, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	} else {
		return append(bytes, '\n'), nil
	}
}

var (
	textMarshallerLevelPlaceholder = map[Level]string{
		LevelDebug: "DEBU",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERRO",
		LevelFatal: "FATA",
		LevelPanic: "PANI",
	}
)

func textMarshaller(entry *Entry) ([]byte, error) {
	if entry.File == "" {
		entry.File = "unset"
	}
	if entry.Level == "" {
		entry.Level = string(LevelInfo)
	}
	if entry.CallTime == "" {
		entry.CallTime = time.Now().Format(timeFormat)
	}
	if entry.Service == "" {
		entry.Service = "unset"
	}

	cleanTextString := func(input string) (cleaned string) {
		cleaned = strings.ReplaceAll(input, "\n", `\n`)
		cleaned = strings.ReplaceAll(cleaned, "\r", `\r`)
		cleaned = strings.ReplaceAll(cleaned, "\t", `\t`)
		cleaned = strings.ReplaceAll(cleaned, `|`, ` `)
		return cleaned
	}

	textBuilder := strings.Builder{}
	textBuilder.WriteString(cleanTextString(entry.File))
	textBuilder.WriteString(" | ")
	textBuilder.WriteString(textMarshallerLevelPlaceholder[Level(entry.Level)])
	textBuilder.WriteString(" | ")
	textBuilder.WriteString(cleanTextString(entry.CallTime))
	textBuilder.WriteString(" | ")
	textBuilder.WriteString(cleanTextString(entry.Service))
	if entry.TraceID != "" {
		textBuilder.WriteString(" | ")
		textBuilder.WriteString(cleanTextString(entry.TraceID))
	}
	if entry.Message != "" {
		textBuilder.WriteString(" | Message: ")
		textBuilder.WriteString(cleanTextString(entry.Message))
	}
	if entry.Data != nil {
		textBuilder.WriteString(" | Data: ")
		textBuilder.WriteString(cleanTextString(fmt.Sprintf("%+v", entry.Data)))
	}
	if entry.Extra != nil {
		textBuilder.WriteString(" | Extra: ")
		for k, v := range entry.Extra {
			textBuilder.WriteString(cleanTextString(fmt.Sprintf("%s=%+v ", k, v)))
		}
	}

	return append([]byte(textBuilder.String()), '\n'), nil
}

func csvMarshaller(entry *Entry) ([]byte, error) {
	if entry.File == "" {
		entry.File = "unset"
	}
	if entry.Level == "" {
		entry.Level = string(LevelInfo)
	}
	if entry.CallTime == "" {
		entry.CallTime = time.Now().Format(timeFormat)
	}
	if entry.Service == "" {
		entry.Service = "unset"
	}

	cleanCsvString := func(input string) (cleaned string) {
		cleaned = strings.ReplaceAll(input, ",", " ")
		cleaned = strings.ReplaceAll(cleaned, "\n", " ")
		cleaned = strings.ReplaceAll(cleaned, "\r", " ")
		cleaned = strings.ReplaceAll(cleaned, "\t", " ")
		return cleaned
	}

	textBuilder := strings.Builder{}
	textBuilder.WriteString(cleanCsvString(entry.File))
	textBuilder.WriteString(",")
	textBuilder.WriteString(cleanCsvString(entry.Level))
	textBuilder.WriteString(",")
	textBuilder.WriteString(cleanCsvString(entry.CallTime))
	textBuilder.WriteString(",")
	textBuilder.WriteString(cleanCsvString(entry.Service))
	textBuilder.WriteString(",")
	textBuilder.WriteString(cleanCsvString(entry.TraceID))
	textBuilder.WriteString(",")
	textBuilder.WriteString(cleanCsvString(entry.Message))
	textBuilder.WriteString(",")
	if entry.Data != nil {
		textBuilder.WriteString(cleanCsvString(fmt.Sprintf("%+v", entry.Data)))
	}
	textBuilder.WriteString(",")
	if entry.Extra != nil {
		for k, v := range entry.Extra {
			textBuilder.WriteString(cleanCsvString(fmt.Sprintf("%s=%+v ", k, v)))
		}
	}

	return append([]byte(textBuilder.String()), '\n'), nil
}

func tsvMarshaller(entry *Entry) ([]byte, error) {
	if entry.File == "" {
		entry.File = "unset"
	}
	if entry.Level == "" {
		entry.Level = string(LevelInfo)
	}
	if entry.CallTime == "" {
		entry.CallTime = time.Now().Format(timeFormat)
	}
	if entry.Service == "" {
		entry.Service = "unset"
	}

	cleanCsvString := func(input string) (cleaned string) {
		cleaned = strings.ReplaceAll(input, ",", " ")
		cleaned = strings.ReplaceAll(cleaned, "\n", " ")
		cleaned = strings.ReplaceAll(cleaned, "\r", " ")
		cleaned = strings.ReplaceAll(cleaned, "\t", " ")
		return
	}

	textBuilder := strings.Builder{}
	textBuilder.WriteString(cleanCsvString(entry.File))
	textBuilder.WriteString("\t")
	textBuilder.WriteString(cleanCsvString(entry.Level))
	textBuilder.WriteString("\t")
	textBuilder.WriteString(cleanCsvString(entry.CallTime))
	textBuilder.WriteString("\t")
	textBuilder.WriteString(cleanCsvString(entry.Service))
	textBuilder.WriteString("\t")
	textBuilder.WriteString(cleanCsvString(entry.TraceID))
	textBuilder.WriteString("\t")
	textBuilder.WriteString(cleanCsvString(entry.Message))
	textBuilder.WriteString("\t")
	if entry.Data != nil {
		textBuilder.WriteString(cleanCsvString(fmt.Sprintf("%+v", entry.Data)))
	}
	textBuilder.WriteString("\t")
	if entry.Extra != nil {
		for k, v := range entry.Extra {
			textBuilder.WriteString(cleanCsvString(fmt.Sprintf("%s=%+v ", k, v)))
		}
	}

	return append([]byte(textBuilder.String()), '\n'), nil
}

func marshalEntry(entry *Entry, marshaller Marshaller) ([]byte, error) {
	if entry == nil {
		return nil, errors.NewEmptyLogEntryError()
	}

	if marshaller == nil {
		marshaller = JsonMarshaller
	}

	return marshaller(entry)
}
