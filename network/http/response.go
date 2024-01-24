package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response interface {
	Error() error
	Result() (status int, body []byte, err error)
	StringResult() (status int, body string, err error)
	BindJsonResult(receiver any) (status int, err error)
	BindXmlResult(receiver any) (status int, err error)
	GetHeader(key string, defaultVal ...string) (val string)
	GetBearerToken() (token string)
	GetCookie(cookieName string) (ck *http.Cookie)
	GetStatusCode() int
	GetBody() (body []byte)
	BindJsonBody(receiver any) (err error)
	BindXmlBody(receiver any) (err error)
}

type response struct {
	cks map[string]*http.Cookie
	bd  []byte
	r   *http.Response
	e   error
}

func (r *response) Error() error {
	return r.e
}

func (r *response) GetHeader(key string, defaultVal ...string) (val string) {
	if r.e != nil {
		return ""
	}

	if v := r.r.Header.Get(key); v != "" {
		return v
	}

	if len(defaultVal) > 0 {
		return defaultVal[0]
	}

	return ""
}

func (r *response) GetBearerToken() (token string) {
	if r.e != nil {
		return ""
	}

	return strings.TrimPrefix(r.GetHeader("Authorization"), "Bearer ")
}

func (r *response) GetCookie(cookieName string) (ck *http.Cookie) {
	if r.e != nil {
		return nil
	}

	if r.cks == nil {
		r.cks = map[string]*http.Cookie{}
		for _, cookie := range r.r.Cookies() {
			r.cks[cookie.Name] = cookie
		}
	}

	return r.cks[cookieName]
}

func (r *response) GetStatusCode() int {
	if r.e != nil {
		return -1
	}

	return r.r.StatusCode
}

func (r *response) GetBody() (body []byte) {
	if r.e != nil {
		return []byte{}
	}

	return r.bd
}

func (r *response) BindJsonBody(receiver any) (err error) {
	if r.e != nil {
		return r.e
	}

	contentType := r.r.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "json") || contentType == "" {
		return json.NewDecoder(bytes.NewBuffer(r.bd)).Decode(receiver)
	}

	return fmt.Errorf("content type is not json: %s", contentType)
}

func (r *response) BindXmlBody(receiver any) (err error) {
	if r.e != nil {
		return r.e
	}

	contentType := r.r.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(contentType), "xml") || contentType == "" {
		return xml.NewDecoder(bytes.NewBuffer(r.bd)).Decode(receiver)
	}

	return fmt.Errorf("content type is not xml: %s", contentType)
}

func (r *response) Result() (status int, body []byte, err error) {
	return r.GetStatusCode(), r.GetBody(), r.Error()
}

func (r *response) StringResult() (status int, body string, err error) {
	return r.GetStatusCode(), string(r.GetBody()), r.Error()
}

func (r *response) BindJsonResult(receiver any) (status int, err error) {
	return r.GetStatusCode(), r.BindJsonBody(receiver)
}

func (r *response) BindXmlResult(receiver any) (status int, err error) {
	return r.GetStatusCode(), r.BindXmlBody(receiver)
}

func NewResponse(r *http.Response, e error) Response {
	if r == nil {
		return &response{r: nil, e: e, bd: []byte{}, cks: nil}
	}

	body, _ := io.ReadAll(r.Body)
	return &response{r: r, e: e, bd: body, cks: nil}
}
