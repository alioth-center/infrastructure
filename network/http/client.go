package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"gopkg.in/yaml.v3"
)

type Client interface {
	ExecuteRequest(request RequestBuilder) (response ResponseParser, err error)
	ExecuteRawRequest(request *http.Request) (response *http.Response, err error)
}

type simpleClient struct {
	cli *http.Client
}

func (c simpleClient) ExecuteRequest(request RequestBuilder) (response ResponseParser, err error) {
	rawReq, buildReqErr := request.Build()
	if buildReqErr != nil {
		return nil, buildReqErr
	}

	rawResp, doReqErr := c.ExecuteRawRequest(rawReq)
	if doReqErr != nil {
		return nil, doReqErr
	}

	return NewSimpleResponseParser(rawResp), nil
}

func (c simpleClient) ExecuteRawRequest(request *http.Request) (response *http.Response, err error) {
	return c.cli.Do(request)
}

func NewSimpleClient() Client {
	return &simpleClient{
		cli: &http.Client{},
	}
}

type loggerClient struct {
	log logger.Logger
	cli *http.Client
}

type RequestLoggingFields struct {
	Url         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"content_type"`
	Body        any    `json:"body"`
}

type ResponseLoggingFields struct {
	StatusCode  int    `json:"status_code"`
	ContentType string `json:"content_type"`
	Body        any    `json:"body"`
}

func (c loggerClient) copyReader(reader io.Reader) (io.ReadCloser, *bytes.Buffer) {
	bytePayload, _ := io.ReadAll(reader)
	return io.NopCloser(bytes.NewReader(bytePayload)), bytes.NewBuffer(bytePayload)
}

func (c loggerClient) buildRequestLoggingFields(req *http.Request) *RequestLoggingFields {
	contentType := req.Header.Get("Content-Type")
	requestFields := RequestLoggingFields{
		Url:         req.URL.String(),
		Method:      req.Method,
		ContentType: contentType,
	}

	// only text body will be logged
	switch {
	case strings.Contains(contentType, ContentTypeJson), strings.Contains(contentType, ContentTypeTextJson):
		var loggingBody *bytes.Buffer
		req.Body, loggingBody = c.copyReader(req.Body)
		decodeErr := json.NewDecoder(loggingBody).Decode(&requestFields.Body)
		if decodeErr != nil {
			requestFields.Body = loggingBody.String()
		}
	case strings.Contains(contentType, ContentTypeYaml), strings.Contains(contentType, ContentTypeTextYaml):
		var loggingBody *bytes.Buffer
		req.Body, loggingBody = c.copyReader(req.Body)
		decodeErr := yaml.NewDecoder(loggingBody).Decode(&requestFields.Body)
		if decodeErr != nil {
			requestFields.Body = loggingBody.String()
		}
	case strings.Contains(contentType, ContentTypeTextPlain), contentType == "":
		var loggingBody *bytes.Buffer
		req.Body, loggingBody = c.copyReader(req.Body)
		requestFields.Body = loggingBody.String()
	default:
		// not text body, set payload to empty buffer
		requestFields.Body = &bytes.Buffer{}
	}

	return &requestFields
}

func (c loggerClient) buildResponseLoggingFields(res *http.Response) *ResponseLoggingFields {
	contentType := res.Header.Get("Content-Type")
	responseFields := ResponseLoggingFields{
		StatusCode:  res.StatusCode,
		ContentType: contentType,
	}

	// only text body will be logged
	switch {
	case strings.Contains(contentType, ContentTypeJson), strings.Contains(contentType, ContentTypeTextJson):
		var loggingBody *bytes.Buffer
		res.Body, loggingBody = c.copyReader(res.Body)
		decodeErr := json.NewDecoder(loggingBody).Decode(&responseFields.Body)
		if decodeErr != nil {
			responseFields.Body = loggingBody.String()
		}
	case strings.Contains(contentType, ContentTypeYaml), strings.Contains(contentType, ContentTypeTextYaml):
		var loggingBody *bytes.Buffer
		res.Body, loggingBody = c.copyReader(res.Body)
		decodeErr := yaml.NewDecoder(loggingBody).Decode(&responseFields.Body)
		if decodeErr != nil {
			responseFields.Body = loggingBody.String()
		}
	case strings.Contains(contentType, ContentTypeTextPlain), strings.Contains(contentType, ContentTypeTextHtml), contentType == "":
		var loggingBody *bytes.Buffer
		res.Body, loggingBody = c.copyReader(res.Body)
		responseFields.Body = loggingBody.String()
	default:
		// not text body, set payload to empty buffer
		responseFields.Body = &bytes.Buffer{}
	}

	return &responseFields
}

func (c loggerClient) ExecuteRequest(request RequestBuilder) (response ResponseParser, err error) {
	rawReq, buildReqErr := request.Build()
	if buildReqErr != nil {
		return nil, buildReqErr
	}

	rawResp, doReqErr := c.ExecuteRawRequest(rawReq)
	if doReqErr != nil {
		return nil, doReqErr
	}

	return NewSimpleResponseParser(rawResp), nil
}

func (c loggerClient) ExecuteRawRequest(request *http.Request) (response *http.Response, err error) {
	// log request before execute
	if c.log != nil {
		if tid := request.Header.Get(TraceHeaderKey()); tid != "" {
			request = request.WithContext(trace.Context(request.Context(), tid))
		}
		fields := logger.NewFields(request.Context()).
			WithMessage("http request executed").
			WithData(c.buildRequestLoggingFields(request)).
			WithField("subsystem_scope", "http_request")
		c.log.Info(fields)
	}

	// execute request
	rawResp, doReqErr := c.cli.Do(request)
	if doReqErr != nil {
		if c.log != nil {
			fields := logger.NewFields(request.Context()).
				WithMessage("http request failed").
				WithData(doReqErr).
				WithField("subsystem_scope", "http_error")
			c.log.Error(fields)
		}
		return nil, doReqErr
	}

	// log response after execute
	if c.log != nil {
		fields := logger.NewFields(request.Context()).
			WithMessage("http response received").
			WithData(c.buildResponseLoggingFields(rawResp)).
			WithField("subsystem_scope", "http_response")
		c.log.Info(fields)
	}

	// return response
	return rawResp, nil
}

func NewLoggerClient(log logger.Logger) Client {
	return &loggerClient{
		log: log,
		cli: &http.Client{},
	}
}

type MockOptions struct {
	Trigger func(req *http.Request) bool
	Handler func(req *http.Request) *http.Response
}

type mockClient struct {
	loggerClient
	opts []*MockOptions
}

func (c mockClient) ExecuteRawRequest(request *http.Request) (response *http.Response, err error) {
	// if logger enable, log request
	if c.log != nil {
		if tid := request.Header.Get(TraceHeaderKey()); tid != "" {
			request = request.WithContext(trace.Context(request.Context(), tid))
		}
		fields := logger.NewFields(request.Context()).
			WithMessage("http request executed").
			WithData(c.buildRequestLoggingFields(request)).
			WithField("subsystem_scope", "http_request")
		c.log.Info(fields)
	}

	// find mock response
	mocked := false
	for _, opt := range c.opts {
		if opt.Trigger(request) {
			response = opt.Handler(request)
			mocked = true
			break
		}
	}

	// if mock response not found, execute request
	if !mocked {
		response, err = c.cli.Do(request)
		if err != nil {
			if c.log != nil {
				fields := logger.NewFields(request.Context()).
					WithMessage("http request failed").
					WithData(err).
					WithField("subsystem_scope", "http_error")
				c.log.Error(fields)
			}
			return nil, err
		}
	}

	// if logger enable, log response
	if c.log != nil {
		fields := logger.NewFields(request.Context()).
			WithMessage("http response received").
			WithData(c.buildResponseLoggingFields(response)).
			WithField("subsystem_scope", "http_response")
		if mocked {
			fields = fields.WithField("mock_status", true)
		} else {
			fields = fields.WithField("mock_status", false)
		}
		c.log.Info(fields)
	}

	// return response
	return response, nil
}

func (c mockClient) ExecuteRequest(request RequestBuilder) (response ResponseParser, err error) {
	rawReq, buildReqErr := request.Build()
	if buildReqErr != nil {
		return nil, buildReqErr
	}

	rawResp, doReqErr := c.ExecuteRawRequest(rawReq)
	if doReqErr != nil {
		return nil, doReqErr
	}

	return NewSimpleResponseParser(rawResp), nil
}

// NewMockClientWithLogger create a mock http client with logger.
// example:
//
//	client := NewMockClientWithLogger(logger.NewLogger())
//
// then
//
//	response, err := client.ExecuteRawRequest(&http.Request{})
//
// will return a mock response.
func NewMockClientWithLogger(log logger.Logger, opts ...*MockOptions) Client {
	if opts == nil {
		opts = []*MockOptions{}
	}

	return &mockClient{
		loggerClient: loggerClient{
			log: log,
			cli: &http.Client{},
		},
		opts: opts,
	}
}

// NewMockClient create a mock http client.
// example:
//
//	client := NewMockClient(
//		&MockOptions{
//			Trigger: func(req *http.Request) bool {
//				return true
//			},
//			Handler: func(req *http.Request) *http.Response {
//				return &http.Response{
//					StatusCode: 200,
//					Body:       io.NopCloser(strings.NewReader("mock response")),
//				}
//			},
//		},
//	)
//
// then
//
//	response, err := client.ExecuteRawRequest(&http.Request{})
//
// will return a mock response.
func NewMockClient(opts ...*MockOptions) Client {
	if opts == nil {
		opts = []*MockOptions{}
	}

	return &mockClient{
		loggerClient: loggerClient{
			log: nil,
			cli: &http.Client{},
		},
		opts: opts,
	}
}
