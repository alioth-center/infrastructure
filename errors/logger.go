package errors

type EmptyLogEntryError struct{}

func (e EmptyLogEntryError) Error() string {
	return "entry is nil or empty"
}

func NewEmptyLogEntryError() error {
	return EmptyLogEntryError{}
}

type FileWriterWriteToDirectoryError struct {
	Path string
}

func (e FileWriterWriteToDirectoryError) Error() string {
	return "write to directory: " + e.Path
}

func NewFileWriterWriteToDirectoryError(path string) error {
	return FileWriterWriteToDirectoryError{Path: path}
}
