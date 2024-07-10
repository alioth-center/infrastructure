package logger

import (
	"encoding/json"
	"time"
)

func defaultMarshaller(fields Fields) []byte {
	entry := fields.Export()
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

	bytes, _ := json.Marshal(entry)
	return append(bytes, '\n')
}
