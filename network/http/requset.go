package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/trace"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type Request interface {
	SetUrl(path string) Request
	SetMethod(method string) Request
	SetBearerToken(token string) Request
	SetUserAgent(userAgent string) Request
	SetAccept(accept string) Request
	SetHeader(key, value string) Request
	SetJsonBody(ptr any) Request
	SetMultiPartBodyWithErrorHandler(multipartField, multipartName string, multiWriter io.Reader, others map[string]string, handler func(err error) (ignore bool)) Request
	SetMultiPartBody(multipartField, multipartName string, multiWriter io.Reader, others map[string]string) Request
	SetCookie(cookieKey, cookieValue string) Request
	build() (req *http.Request, err error)
}

type request struct {
	ctx     context.Context
	method  string
	path    *url.URL
	headers map[string]string
	cookies map[string]string
	body    io.Reader
}

func (r *request) SetUrl(path string) Request {
	parsed, _ := url.Parse(path)
	r.path = parsed
	return r
}

func (r *request) SetMethod(method string) Request {
	r.method = method
	return r
}

func (r *request) SetBearerToken(token string) Request {
	r.headers["Authorization"] = "Bearer " + token
	return r
}

func (r *request) SetUserAgent(userAgent string) Request {
	r.headers["User-Agent"] = userAgent
	return r
}

func (r *request) SetAccept(accept string) Request {
	r.headers["Accept"] = accept
	return r
}

func (r *request) SetHeader(key, value string) Request {
	r.headers[key] = value
	return r
}

func (r *request) SetJsonBody(ptr any) Request {
	if ptr == nil {
		return r
	}

	if payload, marshalErr := json.Marshal(ptr); marshalErr == nil {
		r.SetHeader("Content-Type", "application/json")
		r.body = bytes.NewReader(payload)
	}
	return r
}

func (r *request) SetMultiPartBodyWithErrorHandler(multipartField, multipartName string, multiWriter io.Reader, others map[string]string, handler func(err error) (ignore bool)) Request {
	body := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(body)

	if others != nil {
		for k, v := range others {
			if writeFieldErr := bodyWriter.WriteField(k, v); writeFieldErr != nil {
				if handler(writeFieldErr) {
					continue
				} else {
					return r
				}
			}
		}
	}

	// 创建文件部分
	formFilePart, createFormFileErr := bodyWriter.CreateFormFile(multipartField, multipartName)
	if createFormFileErr != nil && !handler(createFormFileErr) {
		return r
	}

	// 复制文件数据到表单
	_, copyDataErr := io.Copy(formFilePart, multiWriter)
	if copyDataErr != nil && !handler(copyDataErr) {
		return r
	}

	// 关闭multipart写入器以写入结尾分隔符
	closeWriterErr := bodyWriter.Close()
	if closeWriterErr != nil && !handler(closeWriterErr) {
		return r
	}

	// 设置内容类型为multipart form数据的内容类型
	r.SetHeader("Content-Type", bodyWriter.FormDataContentType())

	// 设置请求体
	r.body = body

	return r
}

func (r *request) SetMultiPartBody(multipartField, multipartName string, multiWriter io.Reader, others map[string]string) Request {
	return r.SetMultiPartBodyWithErrorHandler(multipartField, multipartName, multiWriter, others, func(_ error) (ignore bool) {
		return true
	})
}

func (r *request) SetCookie(cookieKey, cookieValue string) Request {
	r.cookies[cookieKey] = cookieValue
	return r
}

func (r *request) SetContext(ctx context.Context) Request {
	_, tracedCtx := trace.GetTraceID(ctx)
	r.ctx = tracedCtx
	return r
}

func (r *request) build() (req *http.Request, err error) {
	if req, err = http.NewRequest(r.method, r.path.String(), r.body); err != nil {
		return nil, fmt.Errorf("build request failed: %w", err)
	} else {
		req = req.WithContext(r.ctx)
		for k, v := range r.headers {
			req.Header.Set(k, v)
		}
		for k, v := range r.cookies {
			cookie := &http.Cookie{
				Name:  k,
				Value: v,
			}
			req.AddCookie(cookie)
		}
		return req, nil
	}
}

func NewRequest() Request {
	return &request{
		ctx:     trace.NewContextWithTraceID(),
		method:  "",
		path:    &url.URL{},
		headers: map[string]string{},
		cookies: map[string]string{},
		body:    &bytes.Buffer{},
	}
}

func NewGetRequest(url string) Request {
	return NewRequest().SetMethod(http.MethodGet).SetUrl(url)
}

func NewPostRequest(url string, payload any) Request {
	return NewRequest().SetMethod(http.MethodPost).SetUrl(url).SetJsonBody(payload)
}
