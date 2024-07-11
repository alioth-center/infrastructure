package config

import (
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig 从指定路径加载配置到 receiver 中，通过文件扩展名自动识别配置文件类型，支持 yaml、json、xml
func LoadConfig(receiver any, path string) (err error) {
	bytesOfConfig, readConfigFileErr := readConfigFile(path)
	if readConfigFileErr != nil {
		return fmt.Errorf("read config file error: %w", readConfigFileErr)
	}

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		return loadYamlConfigWithKeys(receiver, bytesOfConfig)
	case ".json":
		return loadJsonConfigWithKeys(receiver, bytesOfConfig)
	case ".xml":
		return loadXmlConfigWithKeys(receiver, bytesOfConfig)
	default:
		return ErrUnSupportedConfigExtension
	}
}

// LoadConfigWithKeys 从指定路径加载配置到 receiver 中，只加载指定的配置项，通过文件扩展名自动识别配置文件类型，支持 yaml、json
func LoadConfigWithKeys(receiver any, path string, keys ...string) (err error) {
	bytesOfConfig, readConfigFileErr := readConfigFile(path)
	if readConfigFileErr != nil {
		return fmt.Errorf("read config file error: %w", readConfigFileErr)
	}

	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		return loadYamlConfigWithKeys(receiver, bytesOfConfig, keys...)
	case ".json":
		return loadJsonConfigWithKeys(receiver, bytesOfConfig, keys...)
	default:
		return ErrUnSupportedConfigExtension
	}
}

func LoadEmbedConfig(receiver any, fs embed.FS, name string) (err error) {
	bytesOfConfig, readConfigFileErr := fs.ReadFile(name)
	if readConfigFileErr != nil {
		return fmt.Errorf("read config file error: %w", readConfigFileErr)
	}

	switch filepath.Ext(name) {
	case ".yaml", ".yml":
		return loadYamlConfigWithKeys(receiver, bytesOfConfig)
	case ".json":
		return loadJsonConfigWithKeys(receiver, bytesOfConfig)
	case ".xml":
		return loadXmlConfigWithKeys(receiver, bytesOfConfig)
	default:
		return ErrUnSupportedConfigExtension
	}
}

func LoadEmbedConfigWithKeys(receiver any, fs embed.FS, name string, keys ...string) (err error) {
	bytesOfConfig, readConfigFileErr := fs.ReadFile(name)
	if readConfigFileErr != nil {
		return fmt.Errorf("read config file error: %w", readConfigFileErr)
	}

	switch filepath.Ext(name) {
	case ".yaml", ".yml":
		return loadYamlConfigWithKeys(receiver, bytesOfConfig, keys...)
	case ".json":
		return loadJsonConfigWithKeys(receiver, bytesOfConfig, keys...)
	default:
		return ErrUnSupportedConfigExtension
	}
}

func WriteConfig(path string, object any) (err error) {
	switch filepath.Ext(path) {
	case ".yaml", ".yml":
		return writeYamlConfig(path, object)
	case ".json":
		return writeJsonConfig(path, object)
	case ".xml":
		return writeXmlConfig(path, object)
	default:
		return ErrUnSupportedConfigExtension
	}
}

// readConfigFile 读取配置文件
func readConfigFile(path string, bytes ...[]byte) (content []byte, err error) {
	if len(bytes) > 0 {
		return bytes[0], nil
	}

	fileInfo, se := os.Stat(path)
	if se != nil {
		return nil, fmt.Errorf("stat config file error: %w", se)
	}

	if fileInfo.IsDir() {
		return nil, ErrConfigFilePathIsDir
	}

	f, ofe := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o666)
	if ofe != nil {
		return nil, fmt.Errorf("open config file error: %w", ofe)
	}

	bytesOfConfig, rfe := io.ReadAll(f)
	if rfe != nil {
		return nil, fmt.Errorf("read config file error: %w", rfe)
	}

	cfe := f.Close()
	if cfe != nil {
		return nil, fmt.Errorf("close config file error: %w", cfe)
	}

	return bytesOfConfig, nil
}

// loadConfigWithKeys 加载配置文件，只加载指定的配置项
func loadConfigWithKeys(receiver any, bytesOfConfig []byte, unmarshalFunc func(data []byte, receiver any) error, marshalFunc func(object any) (data []byte, err error), keys ...string) (err error) {
	var (
		bytesBuffer  = bytesOfConfig
		objectBuffer = map[string]any{}
	)

	for _, key := range keys {
		unmarshalErr := unmarshalFunc(bytesBuffer, &objectBuffer)
		if unmarshalErr != nil {
			return fmt.Errorf("unmarshal config buffer error: %w", unmarshalErr)
		}

		contentBuffer, existKey := objectBuffer[key]
		if !existKey {
			return ErrConfigContentNotExists
		}

		marshalBytes, marshalErr := marshalFunc(contentBuffer)
		if marshalErr != nil {
			return fmt.Errorf("marshal config buffer error: %w", marshalErr)
		}

		objectBuffer = map[string]any{}
		bytesBuffer = marshalBytes
	}

	unmarshalErr := unmarshalFunc(bytesBuffer, receiver)
	if unmarshalErr != nil {
		return fmt.Errorf("unmarshal config buffer error: %w", unmarshalErr)
	}

	return nil
}

// loadJsonConfigWithKeys 加载 json 配置文件，只加载指定的配置项
func loadJsonConfigWithKeys(receiver any, bytesOfConfig []byte, keys ...string) (err error) {
	return loadConfigWithKeys(receiver, bytesOfConfig, json.Unmarshal, json.Marshal, keys...)
}

// loadYamlConfigWithKeys 加载 yaml 配置文件，只加载指定的配置项
func loadYamlConfigWithKeys(receiver any, bytesOfConfig []byte, keys ...string) (err error) {
	return loadConfigWithKeys(receiver, bytesOfConfig, yaml.Unmarshal, yaml.Marshal, keys...)
}

// loadXmlConfigWithKeys 加载 xml 配置文件，只加载指定的配置项
func loadXmlConfigWithKeys(receiver any, bytesOfConfig []byte, keys ...string) (err error) {
	return loadConfigWithKeys(receiver, bytesOfConfig, xml.Unmarshal, xml.Marshal, keys...)
}

func writeJsonConfig(path string, object any) (err error) {
	bytesOfConfig, marshalErr := json.Marshal(object)
	if marshalErr != nil {
		return fmt.Errorf("marshal config object error: %w", marshalErr)
	}

	writeErr := writeFile(path, bytesOfConfig)
	if writeErr != nil {
		return fmt.Errorf("write config file error: %w", writeErr)
	}

	return nil
}

func writeYamlConfig(path string, object any) (err error) {
	bytesOfConfig, marshalErr := yaml.Marshal(object)
	if marshalErr != nil {
		return fmt.Errorf("marshal config object error: %w", marshalErr)
	}

	writeErr := writeFile(path, bytesOfConfig)
	if writeErr != nil {
		return fmt.Errorf("write config file error: %w", writeErr)
	}

	return nil
}

func writeXmlConfig(path string, object any) (err error) {
	bytesOfConfig, marshalErr := xml.Marshal(object)
	if marshalErr != nil {
		return fmt.Errorf("marshal config object error: %w", marshalErr)
	}

	writeErr := writeFile(path, bytesOfConfig)
	if writeErr != nil {
		return fmt.Errorf("write config file error: %w", writeErr)
	}

	return nil
}

func writeFile(p string, data []byte) (err error) {
	f, createErr := os.Create(p)
	if createErr != nil {
		return fmt.Errorf("create config file error: %w", createErr)
	}

	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			err = fmt.Errorf("close config file error: %w", closeErr)
		}
	}()
	_, writeErr := f.Write(data)
	if writeErr != nil {
		return fmt.Errorf("write file error: %w", writeErr)
	}

	return nil
}
