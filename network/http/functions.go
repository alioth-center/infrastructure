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
//		request := NewRequest().SetUrl("https://example.com/api/user").SetAccept(ContentTypeJson)
//		responseData, parseErr := ParseJsonResponse[User](NewClient().ExecuteRequest(request))
//		if parseErr != nil {
//			panic(parseErr)
//		}
//
// }
func ParseJsonResponse[T any](resp Response) (data T, err error) {
	if resp.Error() != nil {
		return values.Nil[T](), resp.Error()
	}
	if unmarshalErr := resp.BindJsonBody(&data); unmarshalErr != nil {
		return values.Nil[T](), unmarshalErr
	}

	return data, nil
}

// ParseXmlResponse parse xml response from http response
// example:
//
//	type User struct {
//		Name string `xml:"name"`
//	}
//
//	func main() {
//		request := NewRequest().SetUrl("https://example.com/api/user").SetAccept(ContentTypeXml)
//		responseData, parseErr := ParseXmlResponse[User](NewClient().ExecuteRequest(request))
//		if parseErr != nil {
//			panic(parseErr)
//		}
//
// }
func ParseXmlResponse[T any](resp Response) (data T, err error) {
	if resp.Error() != nil {
		return values.Nil[T](), resp.Error()
	}
	if unmarshalErr := resp.BindXmlBody(&data); unmarshalErr != nil {
		return values.Nil[T](), unmarshalErr
	}

	return data, nil
}
