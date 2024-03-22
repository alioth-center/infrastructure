package http

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

// ResponseParser is used to parse a http response.
type ResponseParser interface {
	RawResponse() *http.Response
	RawRequest() *http.Request
	RawBody() []byte
	Context() context.Context
	Status() (code int, message string)
	BindJson(receiver any) (bindErr error)
	BindXml(receiver any) (bindErr error)
	BindCustom(receiver any, decoder func(reader io.Reader, receiver any) error) (bindErr error)
	BindHeader(fields ...string) (header map[string][]string)
	BindCookie(fields ...string) (cookies map[string]*http.Cookie)
}

type simpleParser struct {
	raw *http.Response
	buf *bytes.Buffer

	headers map[string][]string
	cookies map[string]*http.Cookie
}

func (p *simpleParser) RawResponse() *http.Response {
	return p.raw
}

func (p *simpleParser) RawRequest() *http.Request {
	return p.raw.Request
}

func (p *simpleParser) RawBody() []byte {
	return p.buf.Bytes()
}

func (p *simpleParser) Context() context.Context {
	return p.raw.Request.Context()
}

func (p *simpleParser) Status() (code int, message string) {
	return p.raw.StatusCode, p.raw.Status
}

func (p *simpleParser) BindJson(receiver any) (bindErr error) {
	return json.NewDecoder(p.buf).Decode(receiver)
}

func (p *simpleParser) BindXml(receiver any) (bindErr error) {
	return xml.NewDecoder(p.buf).Decode(receiver)
}

func (p *simpleParser) BindCustom(receiver any, decoder func(reader io.Reader, receiver any) error) (bindErr error) {
	return decoder(p.buf, receiver)
}

func (p *simpleParser) BindHeader(fields ...string) (header map[string][]string) {
	header = map[string][]string{}
	for _, field := range fields {
		header[field] = p.headers[field]
	}

	return header
}

func (p *simpleParser) BindCookie(fields ...string) (cookies map[string]*http.Cookie) {
	cookies = map[string]*http.Cookie{}
	for _, field := range fields {
		cookies[field] = p.cookies[field]
	}

	return cookies
}

func NewSimpleResponseParser(r *http.Response) ResponseParser {
	// read response body
	buf := &bytes.Buffer{}
	if r != nil && r.Body != nil {
		payloadBytes, _ := io.ReadAll(r.Body)
		buf.Write(payloadBytes)
		r.Body = io.NopCloser(bytes.NewReader(payloadBytes))
	}

	// read response headers
	cookies := map[string]*http.Cookie{}
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie
	}

	// return response parser
	return &simpleParser{
		raw:     r,
		buf:     buf,
		headers: r.Header,
		cookies: cookies,
	}
}
