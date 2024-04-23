package http

import "github.com/alioth-center/infrastructure/utils/values"

// ParseJsonResponse parse json response from http response
//
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

// StatusCodeIs2XX check if the status code of the response is 2XX
//
// example:
//
//	func main() {
//		client := NewClient()
//		response, executeErr := client.ExecuteRequest(
//			NewRequest().WithPath("https://example.com/api/user").WithMethod(GET),
//		)
//		if executeErr != nil {
//			panic(executeErr)
//		}
//		if !StatusCodeIs2XX(response) {
//			panic("status code is not 2XX")
//		}
//	}
func StatusCodeIs2XX(resp ResponseParser) bool {
	code := resp.RawResponse().StatusCode
	return code >= 200 && code < 300
}

// StatusCodeIs4XX check if the status code of the response is 4XX
//
// example:
//
//	func main() {
//		client := NewClient()
//		response, executeErr := client.ExecuteRequest(
//			NewRequest().WithPath("https://example.com/api/user").WithMethod(GET),
//		)
//		if executeErr != nil {
//			panic(executeErr)
//		}
//		if !StatusCodeIs4XX(response) {
//			panic("status code is not 4XX")
//		}
//	}
func StatusCodeIs4XX(resp ResponseParser) bool {
	code := resp.RawResponse().StatusCode
	return code >= 400 && code < 500
}

// StatusCodeIs5XX check if the status code of the response is 5XX
//
// example:
//
//	func main() {
//		client := NewClient()
//		response, executeErr := client.ExecuteRequest(
//			NewRequest().WithPath("https://example.com/api/user").WithMethod(GET),
//		)
//		if executeErr != nil {
//			panic(executeErr)
//		}
//		if !StatusCodeIs5XX(response) {
//			panic("status code is not 5XX")
//		}
//	}
func StatusCodeIs5XX(resp ResponseParser) bool {
	code := resp.RawResponse().StatusCode
	return code >= 500 && code < 600
}

// CheckStatusCode check if the status code of the response is in the given status code list
//
// example:
//
//	func main() {
//		client := NewClient()
//		response, executeErr := client.ExecuteRequest(
//			NewRequest().WithPath("https://example.com/api/user").WithMethod(GET),
//		)
//		if executeErr != nil {
//			panic(executeErr)
//		}
//		if !CheckStatusCode(response, 200, 201, 202) {
//			panic("status code is not 200, 201 or 202")
//		}
//	}
func CheckStatusCode(resp ResponseParser, want ...Status) bool {
	if len(want) == 0 {
		return false
	}

	code := resp.RawResponse().StatusCode
	return values.ContainsArray(want, code)
}
