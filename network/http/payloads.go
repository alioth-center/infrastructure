package http

import (
	"encoding/json"
	"encoding/xml"

	"github.com/alioth-center/infrastructure/utils/values"
	"gopkg.in/yaml.v3"
)

var processors = map[string]func(in []byte, receiver any) error{
	"":                 json.Unmarshal, // default use json processor
	"application/json": json.Unmarshal,
	"text/json":        json.Unmarshal,
	"application/yaml": yaml.Unmarshal,
	"text/yaml":        yaml.Unmarshal,
	"application/xml":  xml.Unmarshal,
	"text/xml":         xml.Unmarshal,
}

func defaultPayloadProcessor[request any](contentType string, payload []byte, marshaller func(in []byte, receiver any) error) (request, error) {
	if marshaller == nil {
		defaultMarshaller, supported := processors[contentType]
		if !supported {
			return values.Nil[request](), UnsupportedContentTypeError{ContentType: contentType}
		}
		if defaultMarshaller == nil {
			return values.Nil[request](), UnsupportedContentTypeError{ContentType: contentType}
		}

		marshaller = defaultMarshaller
	}

	req := values.Nil[request]()
	unmarshalErr := marshaller(payload, &req)
	if unmarshalErr != nil {
		return values.Nil[request](), unmarshalErr
	}

	return req, nil
}

func SetPayloadProcessor(contentType string, marshaller func(in []byte, receiver any) error) {
	processors[contentType] = marshaller
}
