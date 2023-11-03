package smtp

type DialSmtpServerError struct {
	Err error
}

func (d *DialSmtpServerError) Error() string {
	return "dial smtp server error: " + d.Err.Error()
}

func NewDialSmtpServerError(err error) *DialSmtpServerError {
	return &DialSmtpServerError{Err: err}
}

type InitMailContentError struct {
	Content []byte
	Err     error
}

func (i *InitMailContentError) Error() string {
	return "init mail content error: " + i.Err.Error() + ", content: " + string(i.Content)
}

func NewInitMailContentError(content []byte, err error) *InitMailContentError {
	return &InitMailContentError{Content: content, Err: err}
}
