package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type CustomConfig struct {
		Port    int    `yaml:"port" xml:"port" json:"port"`
		Backend string `yaml:"backend" xml:"backend" json:"backend"`
	}
	mp := CustomConfig{
		Port:    8080,
		Backend: "https://api.your.com",
	}

	t.Run("JsonConfig", func(t *testing.T) {
		jsonMp := CustomConfig{}
		jsonBytes := []byte(`{"port":8080,"backend":"https://api.your.com"}`)
		e := loadJsonConfigWithKeys(&jsonMp, jsonBytes)
		if e != nil {
			t.Error(e)
		}
		if !reflect.DeepEqual(jsonMp, mp) {
			t.Errorf("jsonMp: %+v", jsonMp)
		}
	})

	t.Run("YamlConfig", func(t *testing.T) {
		yamlMp := CustomConfig{}
		yamlBytes := []byte("port: 8080\nbackend: https://api.your.com")
		e := loadYamlConfigWithKeys(&yamlMp, yamlBytes)
		if e != nil {
			t.Error(e)
		}
		if !reflect.DeepEqual(yamlMp, mp) {
			t.Errorf("yamlMp: %+v", yamlMp)
		}
	})

	t.Run("XmlConfig", func(t *testing.T) {
		xmlMp := CustomConfig{}
		xmlBytes := []byte("<CustomConfig><port>8080</port><backend>https://api.your.com</backend></CustomConfig>")
		e := loadXmlConfigWithKeys(&xmlMp, xmlBytes)
		if e != nil {
			t.Error(e)
		}
		if !reflect.DeepEqual(xmlMp, mp) {
			t.Errorf("xmlMp: %+v", xmlMp)
		}
	})
}

type CustomConfig struct {
	Port      int    `yaml:"port" xml:"port" json:"port"`
	Backend   string `yaml:"backend" xml:"backend" json:"backend"`
	ExtraData Extra  `yaml:"extra" xml:"extra" json:"extra"`
}

type Extra struct {
	Username string `yaml:"username" xml:"username" json:"username"`
	Password string `yaml:"password" xml:"password" json:"password"`
}

func TestLoadConfigWithKeys(t *testing.T) {
	exMp := Extra{
		Username: "admin",
		Password: "123456",
	}

	t.Run("JsonConfig", func(t *testing.T) {
		jsonMp := Extra{}
		jsonBytes := []byte(`{"port":8080,"backend":"https://api.your.com","extra":{"username":"admin","password":"123456"}}`)
		e := loadJsonConfigWithKeys(&jsonMp, jsonBytes, "extra")
		if e != nil {
			t.Error(e)
		}
		if !reflect.DeepEqual(jsonMp, exMp) {
			t.Errorf("jsonMp: %+v", jsonMp)
		}
	})

	t.Run("YamlConfig", func(t *testing.T) {
		yamlMp := Extra{}
		yamlBytes := []byte("port: 8080\nbackend: https://api.your.com\nextra:\n    username: admin\n    password: \"123456\"")
		e := loadYamlConfigWithKeys(&yamlMp, yamlBytes, "extra")
		if e != nil {
			t.Error(e)
		}
		if !reflect.DeepEqual(yamlMp, exMp) {
			t.Errorf("yamlMp: %+v", yamlMp)
		}
	})
}
