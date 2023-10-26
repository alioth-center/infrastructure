package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v3"
)

type CustomConfig struct {
	Port    int    `yaml:"port" xml:"port" json:"port"`
	Backend string `yaml:"backend" xml:"backend" json:"backend"`
	Extra   `yaml:"extra" xml:"extra" json:"extra"`
}

type Extra struct {
	Username string `yaml:"username" xml:"username" json:"username"`
	Password string `yaml:"password" xml:"password" json:"password"`
}

func main() {
	object := CustomConfig{
		Port:    8080,
		Backend: "http://api.your.com",
		Extra: Extra{
			Username: "admin",
			Password: "123456",
		},
	}
	yamlBytes, _ := yaml.Marshal(&object)
	println(string(yamlBytes))

	jsonBytes, _ := json.Marshal(&object)
	println(string(jsonBytes))

	xmlBytes, e := xml.Marshal(&object)
	println(string(xmlBytes), e)
	fmt.Println(string(xmlBytes))
}
