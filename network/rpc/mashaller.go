package rpc

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Marshaller interface {
	Marshal(request any) ([]byte, error)
	Unmarshal(data []byte, response any) error
}

type jsonMarshaller struct{}

func (j jsonMarshaller) Marshal(request any) ([]byte, error) {
	return json.Marshal(request)
}

func (j jsonMarshaller) Unmarshal(data []byte, response any) error {
	return json.Unmarshal(data, response)
}

func NewJsonMarshaller() Marshaller {
	return jsonMarshaller{}
}

type yamlMarshaller struct{}

func (y yamlMarshaller) Marshal(request any) ([]byte, error) {
	return yaml.Marshal(request)
}

func (y yamlMarshaller) Unmarshal(data []byte, response any) error {
	return yaml.Unmarshal(data, response)
}

func NewYamlMarshaller() Marshaller {
	return yamlMarshaller{}
}
