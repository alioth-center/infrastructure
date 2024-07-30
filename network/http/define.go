package http

import (
	"os"

	"github.com/gin-gonic/gin"
)

const (
	ModeRelease = "RELEASE"
)

var mode = os.Getenv("AC_MODE")

func SetMode(m string) {
	mode = m
	gin.SetMode(gin.ReleaseMode)
}

const (
	defaultTraceHeaderKey  = "Ac-Request-Id"
	defaultErrorContextKey = "ac-error"
)

var (
	traceHeaderKey  = defaultTraceHeaderKey
	errorContextKey = defaultErrorContextKey
)

func TraceHeaderKey() string {
	return traceHeaderKey
}

func SetTraceHeaderKey(key string) {
	if traceHeaderKey == defaultTraceHeaderKey {
		traceHeaderKey = key
	}
}

func ErrorContextKey() string {
	return errorContextKey
}

func SetErrorContextKey(key string) {
	if errorContextKey == defaultErrorContextKey {
		errorContextKey = key
	}
}

type Method = string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	OPTIONS Method = "OPTIONS"
	HEAD    Method = "HEAD"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	TRACE   Method = "TRACE"
	CONNECT Method = "CONNECT"
	PATCH   Method = "PATCH"
)

type MethodCode = int32

const (
	CodeGet MethodCode = 1 << iota
	CodePost
	CodeOptions
	CodeHead
	CodePut
	CodeDelete
	CodeTrace
	CodeConnect
	CodePatch
)

var (
	MethodMap = map[Method]MethodCode{
		GET:     CodeGet,
		POST:    CodePost,
		OPTIONS: CodeOptions,
		HEAD:    CodeHead,
		PUT:     CodePut,
		DELETE:  CodeDelete,
		TRACE:   CodeTrace,
		CONNECT: CodeConnect,
		PATCH:   CodePatch,
	}

	MethodCodeMap = map[MethodCode]Method{
		CodeGet:     GET,
		CodePost:    POST,
		CodeOptions: OPTIONS,
		CodeHead:    HEAD,
		CodePut:     PUT,
		CodeDelete:  DELETE,
		CodeTrace:   TRACE,
		CodeConnect: CONNECT,
		CodePatch:   PATCH,
	}
)

type UserAgent = string

const (
	Curl         UserAgent = "curl/7.64.1"
	Postman      UserAgent = "PostmanRuntime/7.26.8"
	ChromeOSX    UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
	Safari       UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15"
	Firefox      UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0"
	AliothClient UserAgent = "alioth-http-client/1.0.0"
)

type ContentType = string

const (
	ContentTypeJson       ContentType = "application/json"
	ContentTypeXml        ContentType = "application/xml"
	ContentTypeYaml       ContentType = "application/yaml"
	ContentTypeForm       ContentType = "application/x-www-form-urlencoded"
	ContentTypeFileStream ContentType = "application/octet-stream"
	ContentTypeMultipart  ContentType = "multipart/form-data"
	ContentTypeTextPlain  ContentType = "text/plain"
	ContentTypeTextHtml   ContentType = "text/html"
	ContentTypeTextJson   ContentType = "text/json"
	ContentTypeTextXml    ContentType = "text/xml"
	ContentTypeTextYaml   ContentType = "text/yaml"
	ContentTypeTextCsv    ContentType = "text/csv"
	ContentTypeImagePng   ContentType = "image/png"
	ContentTypeImageJpeg  ContentType = "image/jpeg"
	ContentTypeImageGif   ContentType = "image/gif"
)

type Status = int

const (
	StatusContinue              Status = 100
	StatusSwitchingProtocols    Status = 101
	StatusProcessing            Status = 102
	StatusEarlyHints            Status = 103
	StatusOK                    Status = 200
	StatusCreated               Status = 201
	StatusAccepted              Status = 202
	StatusNonAuthoritative      Status = 203
	StatusNoContent             Status = 204
	StatusResetContent          Status = 205
	StatusPartialContent        Status = 206
	StatusMultiStatus           Status = 207
	StatusAlreadyReported       Status = 208
	StatusIMUsed                Status = 226
	StatusMultipleChoices       Status = 300
	StatusMovedPermanently      Status = 301
	StatusFound                 Status = 302
	StatusSeeOther              Status = 303
	StatusNotModified           Status = 304
	StatusUseProxy              Status = 305
	StatusTemporaryRedirect     Status = 307
	StatusPermanentRedirect     Status = 308
	StatusBadRequest            Status = 400
	StatusUnauthorized          Status = 401
	StatusPaymentRequired       Status = 402
	StatusForbidden             Status = 403
	StatusNotFound              Status = 404
	StatusMethodNotAllowed      Status = 405
	StatusNotAcceptable         Status = 406
	StatusProxyAuthRequired     Status = 407
	StatusRequestTimeout        Status = 408
	StatusConflict              Status = 409
	StatusGone                  Status = 410
	StatusLengthRequired        Status = 411
	StatusPreconditionFailed    Status = 412
	StatusRequestEntityToo      Status = 413
	StatusRequestURITooLong     Status = 414
	StatusUnsupportedMedia      Status = 415
	StatusRequestedRangeNot     Status = 416
	StatusExpectationFailed     Status = 417
	StatusTeapot                Status = 418
	StatusMisdirectedRequest    Status = 421
	StatusUnprocessable         Status = 422
	StatusLocked                Status = 423
	StatusFailedDependency      Status = 424
	StatusTooEarly              Status = 425
	StatusUpgradeRequired       Status = 426
	StatusPreconditionRequired  Status = 428
	StatusTooManyRequests       Status = 429
	StatusRequestHeaderFields   Status = 431
	StatusUnavailableForLegal   Status = 451
	StatusInternalServerError   Status = 500
	StatusNotImplemented        Status = 501
	StatusBadGateway            Status = 502
	StatusServiceUnavailable    Status = 503
	StatusGatewayTimeout        Status = 504
	StatusHTTPVersionNot        Status = 505
	StatusVariantAlsoNegotiates Status = 506
	StatusInsufficientStorage   Status = 507
	StatusLoopDetected          Status = 508
	StatusNotExtended           Status = 510
	StatusNetworkAuthentication Status = 511
)

type FrameworkErrorCode = int

const (
	ErrorCodeMissingRequiredQuery  FrameworkErrorCode = 4001
	ErrorCodeMissingRequiredParam  FrameworkErrorCode = 4002
	ErrorCodeMissingRequiredHeader FrameworkErrorCode = 4003
	ErrorCodeMissingRequiredCookie FrameworkErrorCode = 4004
	ErrorCodeInvalidRequestBody    FrameworkErrorCode = 4005
	ErrorCodeBadRequestBody        FrameworkErrorCode = 4006

	ErrorCodeResourceNotFound FrameworkErrorCode = 4041

	ErrorCodeMethodNotSupported FrameworkErrorCode = 4051
	ErrorCodeMethodNotAllowed   FrameworkErrorCode = 4052

	ErrorCodeInternalErrorOccurred FrameworkErrorCode = 5001
	ErrorCodePanicErrorRecovered   FrameworkErrorCode = 5002
)

type FrameworkResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	RequestID    string `json:"request_id"`
	Data         any    `json:"data,omitempty"`
}

type HeaderEnum = string

const (
	HeaderAuthorization  HeaderEnum = "Authorization"
	HeaderContentType    HeaderEnum = "Content-Type"
	HeaderUserAgent      HeaderEnum = "User-Agent"
	HeaderAccept         HeaderEnum = "Accept"
	HeaderEncoding       HeaderEnum = "Accept-Encoding"
	HeaderAcceptLanguage HeaderEnum = "Accept-Language"
	HeaderContentLength  HeaderEnum = "Content-Length"
)

type NoBody = struct{}

type NoResponse = struct{}
