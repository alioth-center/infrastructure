package errors

type UnsupportedConfigExtensionError struct {
	Extension string
}

func (e UnsupportedConfigExtensionError) Error() string {
	return "unsupported config extension: " + e.Extension
}

func NewUnSupportedConfigExtensionError(extension string) error {
	return UnsupportedConfigExtensionError{Extension: extension}
}

type ConfigFilePathIsDirError struct {
	Path string
}

func (e ConfigFilePathIsDirError) Error() string {
	return "config file path is a directory: " + e.Path
}

func NewConfigFilepathIsDirError(path string) error {
	return ConfigFilePathIsDirError{Path: path}
}

type ConfigContentNotExistsError struct {
	Key string
}

func (e ConfigContentNotExistsError) Error() string {
	return "config content not exists: " + e.Key
}

func NewConfigContentNotExistsError(key string) error {
	return ConfigContentNotExistsError{Key: key}
}
