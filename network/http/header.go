package http

type RequestHeader struct {
	Accept         string
	AcceptEncoding string
	AcceptLanguage string
	UserAgent      string
	ContentType    string
	ContentLength  int
	Origin         string
	Referer        string
	Authorization  string
	ApiKey         string
}
