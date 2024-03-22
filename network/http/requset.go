package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// RequestBuilder is used to build a http request.
type RequestBuilder interface {
	WithContext(ctx context.Context) RequestBuilder
	WithMethod(method Method) RequestBuilder
	WithPath(path string) RequestBuilder
	WithPathFormat(format string, args ...any) RequestBuilder
	WithPathTemplate(template string, args map[string]string) RequestBuilder
	WithQuery(key, value string) RequestBuilder
	WithQueryIgnoreEmptyValue(key, value string) RequestBuilder
	WithHeader(key, value string) RequestBuilder
	WithCookie(key, value string) RequestBuilder
	WithBody(body io.Reader) RequestBuilder
	WithJsonBody(body any) RequestBuilder
	WithUserAgent(userAgent UserAgent) RequestBuilder
	WithBearerToken(token string) RequestBuilder
	WithAccept(accept ContentType) RequestBuilder
	WithContentType(contentType ContentType) RequestBuilder
	Clone() RequestBuilder
	Build() (request *http.Request, err error)
}

type requestBuilder struct {
	ctx context.Context

	method  string
	path    string
	queries *url.Values

	object  any
	body    *bytes.Buffer
	cookies map[string]*http.Cookie
	headers map[string]string
}

func (b *requestBuilder) WithContext(ctx context.Context) RequestBuilder {
	b.ctx = trace.FromContext(ctx)
	return b
}

func (b *requestBuilder) WithMethod(method Method) RequestBuilder {
	b.method = method
	return b
}

func (b *requestBuilder) WithPath(path string) RequestBuilder {
	b.path = path
	return b
}

func (b *requestBuilder) WithPathFormat(format string, args ...any) RequestBuilder {
	return b.WithPath(fmt.Sprintf(format, args...))
}

func (b *requestBuilder) WithPathTemplate(template string, args map[string]string) RequestBuilder {
	b.path = values.NewStringTemplate(template, args).Parse()
	return b
}

func (b *requestBuilder) WithQuery(key, value string) RequestBuilder {
	if b.queries == nil {
		b.queries = &url.Values{}
	}
	b.queries.Add(key, value)
	return b
}

func (b *requestBuilder) WithQueryIgnoreEmptyValue(key, value string) RequestBuilder {
	if value != "" {
		return b.WithQuery(key, value)
	}
	return b
}

func (b *requestBuilder) WithHeader(key, value string) RequestBuilder {
	if b.headers == nil {
		b.headers = map[string]string{}
	}
	b.headers[key] = value
	return b
}

func (b *requestBuilder) WithCookie(key, value string) RequestBuilder {
	if b.cookies == nil {
		b.cookies = map[string]*http.Cookie{}
	}
	b.cookies[key] = &http.Cookie{
		Name:  key,
		Value: value,
	}
	return b
}

func (b *requestBuilder) WithBody(body io.Reader) RequestBuilder {
	b.body, b.object = &bytes.Buffer{}, nil
	_, _ = io.Copy(b.body, body)
	return b
}

func (b *requestBuilder) WithJsonBody(body any) RequestBuilder {
	b.object, b.body = body, &bytes.Buffer{}
	return b.WithContentType(ContentTypeJson)
}

func (b *requestBuilder) WithUserAgent(userAgent UserAgent) RequestBuilder {
	return b.WithHeader("User-Agent", userAgent)
}

func (b *requestBuilder) WithBearerToken(token string) RequestBuilder {
	return b.WithHeader("Authorization", values.BuildStrings("Bearer ", token))
}

func (b *requestBuilder) WithAccept(accept ContentType) RequestBuilder {
	return b.WithHeader("Accept", accept)
}

func (b *requestBuilder) WithContentType(contentType ContentType) RequestBuilder {
	return b.WithHeader("Content-Type", contentType)
}

func (b *requestBuilder) Clone() RequestBuilder {
	builder := &requestBuilder{
		ctx:     b.ctx,
		method:  b.method,
		path:    b.path,
		headers: b.headers,
	}

	if b.queries != nil {
		builder.queries = &url.Values{}
		*builder.queries = *b.queries
	}

	if b.cookies != nil {
		builder.cookies = map[string]*http.Cookie{}
		for k, v := range b.cookies {
			// copy cookie
			cookie := &http.Cookie{
				Name:       v.Name,
				Value:      v.Value,
				Path:       v.Path,
				Domain:     v.Domain,
				Expires:    v.Expires,
				RawExpires: v.RawExpires,
				MaxAge:     v.MaxAge,
				Secure:     v.Secure,
				HttpOnly:   v.HttpOnly,
				SameSite:   v.SameSite,
				Raw:        v.Raw,
				Unparsed:   v.Unparsed,
			}
			builder.cookies[k] = cookie
		}
	}

	if b.object != nil {
		builder.object = b.object
	}

	if b.body != nil {
		builder.body = &bytes.Buffer{}
		_, _ = io.Copy(builder.body, b.body)
	}

	return builder
}

func (b *requestBuilder) Build() (request *http.Request, err error) {
	// parse path
	parsed, parseErr := url.Parse(b.path)
	if parseErr != nil {
		return nil, fmt.Errorf("parse path failed: %w", parseErr)
	}

	if b.queries != nil && len(*b.queries) > 0 {
		parsed.RawQuery = b.queries.Encode()
	}

	// encode json body
	if b.object != nil {
		objectBytes, marshalErr := json.Marshal(b.object)
		if marshalErr != nil {
			return nil, fmt.Errorf("marshal json body failed: %w", marshalErr)
		}

		b.body = bytes.NewBuffer(objectBytes)
	}

	// create request
	request, err = http.NewRequest(b.method, parsed.String(), bytes.NewReader(b.body.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// set context
	request = request.WithContext(b.ctx)

	// set headers
	for k, v := range b.headers {
		request.Header.Set(k, v)
	}

	// set cookies
	for _, cookie := range b.cookies {
		request.AddCookie(cookie)
	}

	return request, nil
}

// NewRequestBuilder returns a new RequestBuilder, which is used to build a http request.
// example:
//
//	builder := NewRequestBuilder().
//		WithMethod(http.MethodGet).
//		WithScheme("https").
//		WithHostname("www.baidu.com").
//		WithPath("/").
//		WithQuery("key", "value").
//		WithHeader("Content-Type", "application/json").
//		WithCookie("key", "value")
//	request, err := builder.Build()
//
// then you can use the request to send a http request
func NewRequestBuilder() RequestBuilder {
	rbq := &requestBuilder{
		ctx:     trace.NewContext(),
		method:  "",
		path:    "",
		queries: &url.Values{},
		object:  nil,
		body:    &bytes.Buffer{},
		cookies: map[string]*http.Cookie{},
		headers: map[string]string{},
	}
	return rbq.WithUserAgent(AliothClient)
}

// MultipartBodyBuilder is used to build a multipart body for a request
type MultipartBodyBuilder interface {
	WithFile(fileKey, fileName string, payload io.Reader) MultipartBodyBuilder
	WithForm(key, value string) MultipartBodyBuilder
	Build() (body io.Reader, contentType string, err error)
}

type multipartBodyBuilder struct {
	fileKey  string
	fileName string
	payload  io.Reader
	forms    map[string]string
}

func (m *multipartBodyBuilder) WithFile(fileKey, fileName string, payload io.Reader) MultipartBodyBuilder {
	m.fileKey = fileKey
	m.fileName = fileName
	m.payload = payload
	return m
}

func (m *multipartBodyBuilder) WithForm(key, value string) MultipartBodyBuilder {
	m.forms[key] = value
	return m
}

func (m *multipartBodyBuilder) Build() (body io.Reader, contentType string, err error) {
	// if payload is nil, return an empty buffer
	if m.payload == nil {
		return &bytes.Buffer{}, "", nil
	}

	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	// write form fields
	for key, value := range m.forms {
		writeErr := writer.WriteField(key, value)
		if writeErr != nil {
			return nil, "", fmt.Errorf("write form field failed: %w", writeErr)
		}
	}

	// write file payload
	fileWriter, createErr := writer.CreateFormFile(m.fileKey, m.fileName)
	if createErr != nil {
		return nil, "", fmt.Errorf("create form file failed: %w", createErr)
	}

	// copy file payload to writer
	_, copyErr := io.Copy(fileWriter, m.payload)
	if copyErr != nil {
		return nil, "", fmt.Errorf("copy file payload failed: %w", copyErr)
	}

	// close writer
	closeErr := writer.Close()
	if closeErr != nil {
		return nil, "", fmt.Errorf("close writer failed: %w", closeErr)
	}

	return buffer, writer.FormDataContentType(), nil
}

// NewMultipartBodyBuilder returns a new MultipartBodyBuilder, which is used to build a multipart body for a request
// example:
//
//	builder := NewMultipartBodyBuilder().WithFile("file", "file.txt", fileReader).WithForm("key", "value")
//	body, contentType, err := builder.Build()
//	if err != nil {
//		panic(err)
//	}
//
//	request := NewRequest().WithBody(body).WithHeader("Content-Type", contentType)
//
// then you can use the request to send a multipart request
func NewMultipartBodyBuilder() MultipartBodyBuilder {
	return &multipartBodyBuilder{
		fileKey:  "",
		fileName: "",
		payload:  nil,
		forms:    map[string]string{},
	}
}
