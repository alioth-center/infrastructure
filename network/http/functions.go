package http

import "github.com/alioth-center/infrastructure/utils/values"

// ParseJsonResponse parse json response from http response
// example:
//
//	type User struct {
//		Name string `json:"name"`
//	}
//
//	func main() {
//		client := NewClient()
//		response, executeErr := client.ExecuteRequest[User](
//			NewRequest().WithPath("https://example.com/api/user").WithMethod(GET),
//		)
//		if executeErr != nil {
//			panic(executeErr)
//		}
//	}
//
// then it will return a User object
func ParseJsonResponse[T any](resp ResponseParser, executeErr error) (data T, err error) {
	if executeErr != nil {
		return values.Nil[T](), executeErr
	}

	if unmarshalErr := resp.BindJson(&data); unmarshalErr != nil {
		return values.Nil[T](), unmarshalErr
	}

	return data, nil
}
