package openai

import "github.com/alioth-center/infrastructure/utils/values"

type ResponseStatusError struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
}

func (e *ResponseStatusError) Error() string {
	return values.BuildStrings("openai response status code unexpected: ", e.Status, "(", values.IntToString(e.StatusCode), ")")
}
