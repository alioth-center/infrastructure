package http

import "github.com/alioth-center/infrastructure/utils/values"

func ParseJsonResponse[T any](resp Response) (data T, err error) {
	if resp.Error() != nil {
		return values.Nil[T](), resp.Error()
	} else if unmarshalErr := resp.BindJsonBody(&data); unmarshalErr != nil {
		return values.Nil[T](), unmarshalErr
	} else {
		return data, nil
	}
}
