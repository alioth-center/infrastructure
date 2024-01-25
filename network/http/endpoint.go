package http

import (
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	TraceIDKey     = "trace-id"
	RemoteIPKey    = "remote-ip"
	RequestTimeKey = "request-time"
)

type methodList int32

func (ml methodList) isAllowed(method string) (exist bool, allow bool) {
	code, hasMethod := MethodMap[method]
	if !hasMethod {
		return false, false
	}

	return true, (int32(ml) & code) != 0
}

func (ml methodList) isAllowedAll() bool {
	return int32(ml) == 0x1ff
}

func (ml methodList) allowNone() methodList {
	return methodList(0)
}

func (ml methodList) allowAll() methodList {
	return methodList(0x1ff)
}

func (ml methodList) allowMethod(methods ...Method) methodList {
	for _, method := range methods {
		if code, hasMethod := MethodMap[method]; hasMethod {
			ml |= methodList(code)
		}
	}

	return ml
}

func (ml methodList) allowedMethods() []string {
	var methods []string
	for code, method := range MethodCodeMap {
		if int32(ml)&code != 0 {
			methods = append(methods, method)
		}
	}

	return methods
}

type Params map[string]string

func (p Params) GetString(key string) string {
	return p[key]
}

func (p Params) GetInt(key string) int {
	return values.StringToInt(p[key], 0)
}

func (p Params) GetUint(key string) uint {
	return values.StringToUint(p[key], uint(0))
}

func (p Params) GetFloat(key string) float64 {
	return values.StringToFloat64(p[key], 0.0)
}

func (p Params) GetBool(key string) bool {
	return values.StringToBool(p[key], false)
}

func defaultRequestProcessor[request any](raw *http.Request) (request, error) {
	payloads, readErr := io.ReadAll(raw.Body)
	if readErr != nil {
		return values.Nil[request](), readErr
	}

	if len(payloads) == 0 {
		return values.Nil[request](), nil
	}

	contentType := raw.Header.Get("Content-Type")
	req, err := defaultPayloadProcessor[request](contentType, payloads, nil)
	if err != nil {
		return values.Nil[request](), err
	}

	return req, nil
}

func defaultResponseProcessor[response any](resp response, status int, err error) (response, int, error) {
	return resp, status, err
}

type EndPointInterface interface {
	bindRouter(base *gin.RouterGroup)
	fullRouterPath() string
}

type RequestProcessor[request any] func(raw *http.Request) (request, error)

type ResponseProcessor[response any] func(resp response, status int, err error) (response, int, error)

// EndPoint is the interface for http endpoint.
// example:
//
//	ep := NewBasicEndPoint[request, response](http.GET, router, chain)
//	ep.AddParsingHeaders("Content-Type", true)
//	ep.AddParsingQueries("id", true)
//	ep.AddParsingParams("name", true)
//
// then
//
//	GET /api/v1/user/:name?id=1
//	Content-Type: application/json
//
// then in your handler
//
//	func handler(ctx Context[request, response]) {
//		ctx.HeaderParams().GetString("Content-Type") // application/json
//		ctx.QueryParams().GetInt("id") // 1
//		ctx.PathParams().GetString("name") // name
//
//		// do something
//	}
type EndPoint[request any, response any] struct {
	router            Router
	requestProcessor  RequestProcessor[request]
	responseProcessor ResponseProcessor[response]
	chain             Chain[request, response]
	allowMethods      methodList
	parsingHeaders    map[string]bool
	parsingQueries    map[string]bool
	parsingParams     map[string]bool
	parsingCookies    map[string]bool
}

func (ep *EndPoint[request, response]) bindRouter(base *gin.RouterGroup) {
	if base == nil || ep.router == nil {
		return
	}

	routerPath := ep.router.FullRouterPath()
	if routerPath == "" {
		routerPath = "/"
	}

	for _, method := range ep.allowMethods.allowedMethods() {
		switch method {
		case http.MethodGet:
			base.GET(routerPath, ep.Serve)
		case http.MethodPost:
			base.POST(routerPath, ep.Serve)
		case http.MethodPut:
			base.PUT(routerPath, ep.Serve)
		case http.MethodDelete:
			base.DELETE(routerPath, ep.Serve)
		case http.MethodPatch:
			base.PATCH(routerPath, ep.Serve)
		case http.MethodHead:
			base.HEAD(routerPath, ep.Serve)
		case http.MethodOptions:
			base.OPTIONS(routerPath, ep.Serve)
		}
	}
}

func (ep *EndPoint[request, response]) fullRouterPath() string {
	return ep.router.FullRouterPath()
}

func (ep *EndPoint[request, response]) preprocessRequestMethod(ctx *gin.Context) (checked bool, status int, err error) {
	existMethod, allowMethod := ep.allowMethods.isAllowed(ctx.Request.Method)
	if !existMethod {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, values.Nil[response]())
		status, err = http.StatusMethodNotAllowed, UnsupportedMethodError{ctx.Request.Method}
		return false, status, err
	}
	if !allowMethod {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, values.Nil[response]())
		status, err = http.StatusMethodNotAllowed, MethodNotAllowedError{ctx.Request.Method}
		return false, status, err
	}

	return true, http.StatusOK, nil
}

func (ep *EndPoint[request, response]) preprocessRequestHeader(ctx *gin.Context) (checked bool, headers Params, status int, err error) {
	checked, headers, status, err = true, map[string]string{}, http.StatusOK, nil
	for key, necessary := range ep.parsingHeaders {
		if necessary && ctx.GetHeader(key) == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, values.Nil[response]())
			checked, status, err = false, http.StatusBadRequest, NecessaryHeaderMissingError{key}
		}

		headers[key] = ctx.GetHeader(key)
	}

	return checked, headers, status, err
}

func (ep *EndPoint[request, response]) preprocessNormalRequestHeader(ctx *gin.Context) (headers RequestHeader) {
	headers = RequestHeader{
		Accept:         ctx.GetHeader("Accept"),
		AcceptEncoding: ctx.GetHeader("Accept-Encoding"),
		AcceptLanguage: ctx.GetHeader("Accept-Language"),
		UserAgent:      ctx.GetHeader("User-Agent"),
		ContentType:    ctx.GetHeader("Content-Type"),
		ContentLength:  values.StringToInt(ctx.GetHeader("Content-Length"), 0),
		Origin:         ctx.GetHeader("Origin"),
		Referer:        ctx.GetHeader("Referer"),
		Authorization:  ctx.GetHeader("Authorization"),
		ApiKey:         ctx.GetHeader("X-API-Key"),
	}

	return headers
}

func (ep *EndPoint[request, response]) preprocessRequestQuery(ctx *gin.Context) (checked bool, queries Params, status int, err error) {
	checked, queries, status, err = true, map[string]string{}, http.StatusOK, nil
	for key, necessary := range ep.parsingQueries {
		value := ctx.Query(key)
		if necessary && value == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, values.Nil[response]())
			checked, status, err = false, http.StatusBadRequest, NecessaryQueryMissingError{key}
		}

		queries[key] = value
	}

	return checked, queries, status, err
}

func (ep *EndPoint[request, response]) preprocessRequestParam(ctx *gin.Context) (checked bool, params Params, status int, err error) {
	checked, params, status, err = true, map[string]string{}, http.StatusOK, nil
	for key, necessary := range ep.parsingParams {
		value := ctx.Param(key)
		if strings.HasPrefix(value, "/") {
			value = value[1:]
		}
		if strings.HasSuffix(value, "/") {
			value = value[:len(value)-1]
		}

		if necessary && value == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, values.Nil[response]())
			checked, status, err = false, http.StatusBadRequest, NecessaryQueryMissingError{key}
		}

		params[key] = value
	}

	return checked, params, status, err
}

func (ep *EndPoint[request, response]) preprocessRequestCookie(ctx *gin.Context) (checked bool, cookies Params, status int, err error) {
	checked, cookies, status, err = true, map[string]string{}, http.StatusOK, nil
	for key, necessary := range ep.parsingCookies {
		value, err := ctx.Cookie(key)
		if necessary && err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, values.Nil[response]())
			checked, status, err = false, http.StatusBadRequest, NecessaryCookieMissingError{key}
		}

		cookies[key] = value
	}

	return checked, cookies, status, err
}

func (ep *EndPoint[request, response]) preprocessRequestBody(ctx *gin.Context) (req request, status int, err error) {
	preprocessor := ep.requestProcessor
	if preprocessor == nil {
		preprocessor = defaultRequestProcessor[request]
	}

	req, err = preprocessor(ctx.Request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, values.Nil[response]())
		return values.Nil[request](), http.StatusBadRequest, err
	}

	return req, http.StatusOK, nil
}

func (ep *EndPoint[request, response]) SetRequestProcessor(processor RequestProcessor[request]) {
	if processor == nil {
		processor = func(raw *http.Request) (request, error) {
			return values.Nil[request](), nil
		}
	}

	ep.requestProcessor = processor
}

func (ep *EndPoint[request, response]) SetResponseProcessor(processor ResponseProcessor[response]) {
	if processor == nil {
		processor = func(resp response, status int, err error) (response, int, error) {
			return resp, http.StatusOK, nil
		}
	}

	ep.responseProcessor = processor
}

func (ep *EndPoint[request, response]) SetHandlerChain(chain Chain[request, response]) {
	ep.chain = NewChain[request, response](chain...)
}

func (ep *EndPoint[request, response]) SetAllowedMethods(methods ...Method) {
	ep.allowMethods = ep.allowMethods.allowMethod(methods...)
}

func (ep *EndPoint[request, response]) AddParsingHeaders(key string, necessary bool) {
	if ep.parsingHeaders == nil {
		ep.parsingHeaders = map[string]bool{}
	}

	ep.parsingHeaders[key] = necessary
}

func (ep *EndPoint[request, response]) AddParsingQueries(key string, necessary bool) {
	if ep.parsingQueries == nil {
		ep.parsingQueries = map[string]bool{}
	}

	ep.parsingQueries[key] = necessary
}

func (ep *EndPoint[request, response]) AddParsingParams(key string, necessary bool) {
	if ep.parsingParams == nil {
		ep.parsingParams = map[string]bool{}
	}

	ep.parsingParams[key] = necessary
}

func (ep *EndPoint[request, response]) AddParsingCookies(key string, necessary bool) {
	if ep.parsingCookies == nil {
		ep.parsingCookies = map[string]bool{}
	}

	ep.parsingCookies[key] = necessary
}

func (ep *EndPoint[request, response]) Serve(ctx *gin.Context) {
	var (
		finalResponse   response
		finalStatusCode int
		finalError      error
		finalHeaders    Params = map[string]string{}
		finalCookies    []Cookie
	)

	defer func() {
		preprocessor := ep.responseProcessor
		if preprocessor == nil {
			preprocessor = defaultResponseProcessor[response]
		}

		// preprocessing http response
		resp, status, err := preprocessor(finalResponse, finalStatusCode, finalError)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			return
		}

		for key, value := range finalHeaders {
			ctx.Header(key, value)
		}
		for _, cookie := range finalCookies {
			ctx.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		}
		ctx.JSON(status, resp)
	}()

	var (
		headerParams Params = map[string]string{}
		queryParams  Params = map[string]string{}
		pathParams   Params = map[string]string{}
		extraParams  Params = map[string]string{}
		cookieParams Params = map[string]string{}
		normalHeader RequestHeader
		bodyPayload  request
	)

	// attach extra params
	tid, tracedCtx := trace.GetTraceID(ctx)
	extraParams[TraceIDKey] = tid
	extraParams[RemoteIPKey] = ctx.ClientIP()
	extraParams[RequestTimeKey] = values.Int64ToString(time.Now().UnixMilli())

	// preprocessing http request
	if checked, status, err := ep.preprocessRequestMethod(ctx); !checked {
		finalResponse, finalStatusCode, finalError = values.Nil[response](), status, err
		return
	}
	if checked, headers, status, err := ep.preprocessRequestHeader(ctx); !checked {
		finalResponse, finalStatusCode, finalError = values.Nil[response](), status, err
		return
	} else {
		headerParams = headers
	}
	if checked, queries, status, err := ep.preprocessRequestQuery(ctx); !checked {
		finalResponse, finalStatusCode, finalError = values.Nil[response](), status, err
		return
	} else {
		queryParams = queries
	}
	if checked, params, status, err := ep.preprocessRequestParam(ctx); !checked {
		finalResponse, finalStatusCode, finalError = values.Nil[response](), status, err
		return
	} else {
		pathParams = params
	}
	if checked, cookies, status, err := ep.preprocessRequestCookie(ctx); !checked {
		finalResponse, finalStatusCode, finalError = values.Nil[response](), status, err
		return
	} else {
		cookieParams = cookies
	}
	if req, status, err := ep.preprocessRequestBody(ctx); err != nil {
		finalResponse, finalStatusCode, finalError = values.Nil[response](), status, err
		return
	} else {
		bodyPayload = req
	}
	normalHeader = ep.preprocessNormalRequestHeader(ctx)

	// init context
	context := NewContext[request, response](
		WithContext[request, response](tracedCtx),
		WithQueryParams[request, response](queryParams),
		WithPathParams[request, response](pathParams),
		WithHeaderParams[request, response](headerParams),
		WithExtraParams[request, response](extraParams),
		WithCookieParams[request, response](cookieParams),
		WithRawRequest[request, response](ctx.Request),
		WithRequest[request, response](bodyPayload),
		WithRequestHeader[request, response](normalHeader),
	)

	ep.chain.Execute(context)
	finalResponse, finalStatusCode, finalError, finalHeaders, finalCookies =
		context.Response(), context.StatusCode(), context.Error(), context.ResponseHeaders(), context.ResponseSetCookies()
}

// EndPointOptions is the options for EndPoint.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithRouterOpts[request, response](router),
//		WithChainOpts[request, response](chain),
//		WithHeaderOpts[request, response](headers),
//		WithQueryOpts[request, response](queries),
//		WithParamOpts[request, response](params),
//		WithRequestProcessorOpts[request, response](processor),
//		WithResponseProcessorOpts[request, response](processor),
//		WithAllowedMethodsOpts[request, response](methods...),
//	)
type EndPointOptions[request any, response any] func(ep *EndPoint[request, response])

// WithRouterOpts sets the router for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithRouterOpts[request, response](router),
//	)
func WithRouterOpts[request any, response any](router Router) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.router = router
	}
}

// WithChainOpts sets the chain for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithChainOpts[request, response](chain),
//	)
func WithChainOpts[request any, response any](chain Chain[request, response]) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.SetHandlerChain(chain)
	}
}

// WithHeaderOpts sets the headers for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithHeaderOpts[request, response](headers),
//	)
func WithHeaderOpts[request any, response any](headers map[string]bool) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		if headers == nil {
			return
		}

		for key, necessary := range headers {
			ep.AddParsingHeaders(key, necessary)
		}
	}
}

// WithQueryOpts sets the queries for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithQueryOpts[request, response](queries),
//	)
func WithQueryOpts[request any, response any](queries map[string]bool) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		if queries == nil {
			return
		}

		for key, necessary := range queries {
			ep.AddParsingQueries(key, necessary)
		}
	}
}

// WithParamOpts sets the params for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithParamOpts[request, response](params),
//	)
func WithParamOpts[request any, response any](params map[string]bool) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		if params == nil {
			return
		}

		for key, necessary := range params {
			ep.AddParsingParams(key, necessary)
		}
	}
}

// WithCookieOpts sets the cookies for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithCookieOpts[request, response](cookies),
//	)
func WithCookieOpts[request any, response any](cookies map[string]bool) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		if cookies == nil {
			return
		}

		for key, necessary := range cookies {
			ep.AddParsingCookies(key, necessary)
		}
	}
}

// WithRequestProcessorOpts sets the request processor for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithRequestProcessorOpts[request, response](processor),
//	)
func WithRequestProcessorOpts[request any, response any](processor RequestProcessor[request]) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.SetRequestProcessor(processor)
	}
}

// WithResponseProcessorOpts sets the response processor for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithResponseProcessorOpts[request, response](processor),
//	)
func WithResponseProcessorOpts[request any, response any](processor ResponseProcessor[response]) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.SetResponseProcessor(processor)
	}
}

// WithAllowedMethodsOpts sets the allowed methods for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithAllowedMethodsOpts[request, response](http.GET, http.POST, http.PUT, http.DELETE, http.PATCH),
//	)
func WithAllowedMethodsOpts[request any, response any](methods ...Method) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.SetAllowedMethods(methods...)
	}
}

// NewEndPointWithOpts creates a new EndPoint with options.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithRouterOpts[request, response](router),
//		WithChainOpts[request, response](chain),
//		WithHeaderOpts[request, response](headers),
//		WithQueryOpts[request, response](queries),
//		WithParamOpts[request, response](params),
//		WithRequestProcessorOpts[request, response](processor),
//		WithResponseProcessorOpts[request, response](processor),
//		WithAllowedMethodsOpts[request, response](methods...),
//	)
func NewEndPointWithOpts[request any, response any](opts ...EndPointOptions[request, response]) *EndPoint[request, response] {
	ep := &EndPoint[request, response]{}

	for _, opt := range opts {
		if opt != nil {
			opt(ep)
		}
	}

	return ep
}

// NewBasicEndPoint creates a new EndPoint with basic options. It's a shortcut for NewEndPointWithOpts.
// example:
//
//	ep := NewBasicEndPoint[request, response](http.GET, router, chain)
func NewBasicEndPoint[request any, response any](method Method, router Router, chain Chain[request, response]) *EndPoint[request, response] {
	return NewEndPointWithOpts[request, response](
		WithAllowedMethodsOpts[request, response](method),
		WithRouterOpts[request, response](router),
		WithChainOpts[request, response](chain),
	)
}
