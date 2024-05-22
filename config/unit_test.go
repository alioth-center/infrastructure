package config

import (
	"embed"
	"os"
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

//go:embed testdata/*
var testFS embed.FS

func TestLoadConfigViaGPT(t *testing.T) {
	type CustomConfig struct {
		Port    int    `yaml:"port" xml:"port" json:"port"`
		Backend string `yaml:"backend" xml:"backend" json:"backend"`
	}
	mp := CustomConfig{
		Port:    8080,
		Backend: "https://api.your.com",
	}

	t.Run("JsonConfig", func(t *testing.T) {
		jsonMps := CustomConfig{}
		err := LoadConfig(&jsonMps, "testdata/config.json")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(jsonMps, mp) {
			t.Errorf("jsonMps: %+v", jsonMps)
		}

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
		yamlMps := CustomConfig{}
		err := LoadConfig(&yamlMps, "testdata/config.yaml")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(yamlMps, mp) {
			t.Errorf("yamlMps: %+v", yamlMps)
		}

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
		xmlMps := CustomConfig{}
		err := LoadConfig(&xmlMps, "testdata/config.xml")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(xmlMps, mp) {
			t.Errorf("xmlMps: %+v", xmlMps)
		}

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

	t.Run("UnsupportedExtension", func(t *testing.T) {
		txtMps := CustomConfig{}
		e := LoadConfig(&txtMps, "testdata/config.txt")
		if e == nil {
			t.Error("expected error for unsupported extension, got nil")
		}

		txtMp := CustomConfig{}
		err := LoadConfig(&txtMp, "config.txt")
		if err == nil {
			t.Error("expected error for unsupported extension, got nil")
		}
	})
}

func TestLoadConfigWithKeysViaGPT(t *testing.T) {
	exMp := Extra{
		Username: "admin",
		Password: "123456",
	}

	t.Run("JsonConfig", func(t *testing.T) {
		jsonMps := Extra{}
		err := LoadConfigWithKeys(&jsonMps, "testdata/config_with_keys.json", "extra")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(jsonMps, exMp) {
			t.Errorf("jsonMps: %+v", jsonMps)
		}

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
		yamlMps := Extra{}
		err := LoadConfigWithKeys(&yamlMps, "testdata/config_with_keys.yaml", "extra")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(yamlMps, exMp) {
			t.Errorf("yamlMps: %+v", yamlMps)
		}

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

	t.Run("KeyNotFound", func(t *testing.T) {
		jsonMp := Extra{}
		jsonBytes := []byte(`{"port":8080,"backend":"https://api.your.com"}`)
		err := loadJsonConfigWithKeys(&jsonMp, jsonBytes, "nonexistent")
		if err == nil {
			t.Error("expected error for nonexistent key, got nil")
		}
	})

	t.Run("InvalidJson", func(t *testing.T) {
		jsonMp := Extra{}
		jsonBytes := []byte(`{"port":8080,"backend":}`)
		err := loadJsonConfigWithKeys(&jsonMp, jsonBytes, "extra")
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		txtMp := Extra{}
		err := LoadConfigWithKeys(&txtMp, "testdata/config.txt", "extra")
		if err == nil {
			t.Error("expected error for invalid format, got nil")
		}
	})
}

func TestLoadEmbedConfig(t *testing.T) {
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
		err := LoadEmbedConfig(&jsonMp, testFS, "testdata/config.json")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(jsonMp, mp) {
			t.Errorf("jsonMp: %+v", jsonMp)
		}
	})

	t.Run("YamlConfig", func(t *testing.T) {
		yamlMp := CustomConfig{}
		err := LoadEmbedConfig(&yamlMp, testFS, "testdata/config.yaml")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(yamlMp, mp) {
			t.Errorf("yamlMp: %+v", yamlMp)
		}
	})

	t.Run("XmlConfig", func(t *testing.T) {
		xmlMp := CustomConfig{}
		err := LoadEmbedConfig(&xmlMp, testFS, "testdata/config.xml")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(xmlMp, mp) {
			t.Errorf("xmlMp: %+v", xmlMp)
		}
	})

	t.Run("UnsupportedExtension", func(t *testing.T) {
		txtMp := CustomConfig{}
		err := LoadEmbedConfig(&txtMp, testFS, "testdata/config.txt")
		if err == nil {
			t.Error("expected error for unsupported extension, got nil")
		}
	})
}

func TestLoadEmbedConfigWithKeys(t *testing.T) {
	exMp := Extra{
		Username: "admin",
		Password: "123456",
	}

	t.Run("JsonConfig", func(t *testing.T) {
		jsonMp := Extra{}
		err := LoadEmbedConfigWithKeys(&jsonMp, testFS, "testdata/config_with_keys.json", "extra")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(jsonMp, exMp) {
			t.Errorf("jsonMp: %+v", jsonMp)
		}
	})

	t.Run("YamlConfig", func(t *testing.T) {
		yamlMp := Extra{}
		err := LoadEmbedConfigWithKeys(&yamlMp, testFS, "testdata/config_with_keys.yaml", "extra")
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(yamlMp, exMp) {
			t.Errorf("yamlMp: %+v", yamlMp)
		}
	})

	t.Run("KeyNotFound", func(t *testing.T) {
		jsonMp := Extra{}
		err := LoadEmbedConfigWithKeys(&jsonMp, testFS, "testdata/config.json", "nonexistent")
		if err == nil {
			t.Error("expected error for nonexistent key, got nil")
		}
	})

	t.Run("InvalidJson", func(t *testing.T) {
		jsonMp := Extra{}
		err := LoadEmbedConfigWithKeys(&jsonMp, testFS, "testdata/invalid.json", "extra")
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})
}

func TestReadConfigFile(t *testing.T) {
	t.Run("FileNotExist", func(t *testing.T) {
		_, err := readConfigFile("nonexistent.yaml")
		if err == nil {
			t.Error("expected error for nonexistent file, got nil")
		}
	})

	t.Run("IsDir", func(t *testing.T) {
		_, err := readConfigFile(".")
		if err == nil {
			t.Error("expected error for directory, got nil")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		emptyFilePath := "testdata/empty.yaml"
		_, err := os.Create(emptyFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(emptyFilePath)

		content, err := readConfigFile(emptyFilePath)
		if err != nil {
			t.Error(err)
		}
		if len(content) != 0 {
			t.Errorf("expected empty content, got %s", content)
		}
	})
}
