package http

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	Expires  string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func NewBasicCookie(name, value string) *Cookie {
	return &Cookie{Name: name, Value: value}
}
