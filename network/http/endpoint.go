package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/gin-gonic/gin"
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

type EndPointInterface interface {
	bindRouter(base *gin.RouterGroup, father Router)
}

type EndpointGroup struct {
	groupUrl  string
	endpoints []EndPointInterface
}

func (e *EndpointGroup) bindRouter(base *gin.RouterGroup, father Router) {
	sub := father.Group(e.groupUrl)
	sub.Extend(father)
	for i := range e.endpoints {
		e.endpoints[i].bindRouter(base, sub)
	}
}

func (e *EndpointGroup) AddEndPoints(endpoints ...EndPointInterface) {
	if len(endpoints) == 0 {
		return
	}

	e.endpoints = append(e.endpoints, endpoints...)
}

// EndPoint is the interface for http endpoint.
//
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
//
// or using the builder like:
//
//	ep := NewEndPointBuilder[request, response]().
//		SetAllowMethods(http.GET, http.POST).
//		SetNecessaryParams("name").
//		SetAdditionalQueries("id").
//		SetAdditionalHeaders("Content-Type").
//		SetAdditionalCookies("session").
//		Build()
type EndPoint[request any, response any] struct {
	router         Router
	chain          Chain[request, response]
	allowMethods   methodList
	customRender   bool
	parsingHeaders map[string]bool
	parsingQueries map[string]bool
	parsingParams  map[string]bool
	parsingCookies map[string]bool
	preprocessors  []EndpointPreprocessor[request, response]
	ginMiddlewares []gin.HandlerFunc
}

func (ep *EndPoint[request, response]) bindRouter(router *gin.RouterGroup, base Router) {
	// no routers to bind, do nothing
	if router == nil || ep.router == nil {
		return
	}

	// try to bind father router
	ep.router.Extend(base)
	routerPath := ep.router.FullRouterPath()

	// final router path fix, must begin with a slash and end without a slash
	if routerPath == "" {
		routerPath = "/"
	}

	// inject gin middlewares
	var fns []gin.HandlerFunc
	if len(ep.ginMiddlewares) > 0 {
		fns = append(fns, ep.ginMiddlewares...)
	}
	fns = append(fns, ep.Serve)

	// bind router
	for _, method := range ep.allowMethods.allowedMethods() {
		switch method {
		case http.MethodGet:
			router.GET(routerPath, fns...)
		case http.MethodPost:
			router.POST(routerPath, fns...)
		case http.MethodPut:
			router.PUT(routerPath, fns...)
		case http.MethodDelete:
			router.DELETE(routerPath, fns...)
		case http.MethodPatch:
			router.PATCH(routerPath, fns...)
		case http.MethodHead:
			router.HEAD(routerPath, fns...)
		case http.MethodOptions:
			router.OPTIONS(routerPath, fns...)
		}
	}
}

func (ep *EndPoint[request, response]) SetGinMiddlewares(middlewares ...gin.HandlerFunc) {
	if len(middlewares) > 0 {
		ep.ginMiddlewares = middlewares
	}
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
	// attach extra params
	extraParams := Params{}
	tid, tracedCtx := trace.TransformContext(ctx)
	extraParams[TraceIDKey] = tid
	extraParams[RemoteIPKey] = ctx.ClientIP()
	extraParams[RequestTimeKey] = values.Int64ToString(time.Now().UnixMilli())

	defer func() {
		if recovered := concurrency.RecoverErr(recover()); recovered != nil {
			errResponse := &FrameworkResponse{
				ErrorCode:    ErrorCodePanicErrorRecovered,
				ErrorMessage: values.BuildStrings("internal error: ", recovered.Error()),
				RequestID:    tid,
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)

			if mode != ModeRelease {
				// recovered error should be panic with stack trace when not in release mode
				panic(fmt.Sprintf("error: %v\nwith stack:\n%s", recovered.Error(), trace.Stack(1)))
			}
		}
	}()

	// copy request body
	content, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(content))

	// init context and preprocess
	context := NewContext[request, response](WithContext[request, response](tracedCtx), WithExtraParams[request, response](extraParams), WithRawRequest[request, response](ctx.Request))
	preprocessors := DefaultPreprocessors[request, response]()
	if len(ep.preprocessors) > 0 {
		preprocessors = ep.preprocessors
	}
	for _, preprocessor := range preprocessors {
		preprocessor(ep, ctx, context.(PreprocessedContext[request, response]))
	}
	if ctx.IsAborted() {
		return
	}

	// write request body back
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(content))

	// execute endpoint handler chain
	ep.chain.Execute(context)

	// internal error occurred
	if context.Error() != nil {
		errResponse := &FrameworkResponse{
			ErrorCode:    ErrorCodeInternalErrorOccurred,
			ErrorMessage: values.BuildStrings("internal error: ", context.Error().Error()),
			RequestID:    tid,
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)
		return
	} else if len(ctx.Errors) > 0 {
		errResponse := &FrameworkResponse{
			ErrorCode:    ErrorCodeInternalErrorOccurred,
			ErrorMessage: ctx.Errors.String(),
			RequestID:    tid,
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)
		return
	}

	// set response
	for key, value := range context.ResponseHeaders() {
		ctx.Header(key, value)
	}
	for _, cookie := range context.ResponseSetCookies() {
		ctx.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}

	// write response
	if ep.customRender {
		// enable no render or customRender options, write nothing to response writer
		ctx.Status(context.StatusCode())
		return
	}

	// disable custom render options, use default json render
	outData, encodeErr := json.Marshal(context.Response())
	if encodeErr != nil {
		errResponse := &FrameworkResponse{
			ErrorCode:    ErrorCodeInternalErrorOccurred,
			ErrorMessage: values.BuildStrings("internal error: ", encodeErr.Error()),
			RequestID:    tid,
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errResponse)
		return
	}
	ctx.Data(context.StatusCode(), ContentTypeJson, outData)
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

// WithGinMiddlewaresOpts sets the gin middlewares for EndPoint, nil is allowed, but it will not work.
// example:
//
//	ep := NewEndPointWithOpts[request, response](
//		WithGinMiddlewaresOpts[request, response](middlewares...),
//	)
func WithGinMiddlewaresOpts[request any, response any](middlewares ...gin.HandlerFunc) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.SetGinMiddlewares(middlewares...)
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

func WithCustomPreprocessors[request any, response any](preprocessors ...EndpointPreprocessor[request, response]) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.preprocessors = preprocessors
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

func WithCustomRender[request any, response any](enable bool) EndPointOptions[request, response] {
	return func(ep *EndPoint[request, response]) {
		ep.customRender = enable
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

type EndPointBuilder[request, response any] struct {
	options []EndPointOptions[request, response]
}

// SetAllowMethods sets the allowed methods for EndPoint.
//
// example:
//
//	ep := NewEndPointBuilder[request, response]().
//		SetAllowMethods(http.GET, http.POST, http.PUT, http.DELETE, http.PATCH).
//		Build()
func (eb *EndPointBuilder[request, response]) SetAllowMethods(methods ...Method) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithAllowedMethodsOpts[request, response](methods...))
	return eb
}

// SetNecessaryParams sets the necessary params for EndPoint.
//
// example:
//
//	ep := NewEndPointBuilder[request, response]().
//		SetNecessaryParams("name").
//		Build()
func (eb *EndPointBuilder[request, response]) SetNecessaryParams(params ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range params {
		args[params[i]] = true
	}
	eb.options = append(eb.options, WithParamOpts[request, response](args))
	return eb
}

// SetAdditionalParams sets the additional params for EndPoint.
//
// example:
//
//	ep := NewEndPointBuilder[request, response]().
//		SetAdditionalParams("name").
//		Build()
func (eb *EndPointBuilder[request, response]) SetAdditionalParams(params ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range params {
		args[params[i]] = false
	}
	eb.options = append(eb.options, WithParamOpts[request, response](args))
	return eb
}

// SetParams sets the params for EndPoint.
//
// example:
//
//	ep := NewEndPointBuilder[request, response]().
//		SetParams(map[string]bool{"name": true}).
//		Build()
func (eb *EndPointBuilder[request, response]) SetParams(params map[string]bool) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithParamOpts[request, response](params))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetNecessaryQueries(queries ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range queries {
		args[queries[i]] = true
	}
	eb.options = append(eb.options, WithQueryOpts[request, response](args))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetAdditionalQueries(queries ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range queries {
		args[queries[i]] = false
	}
	eb.options = append(eb.options, WithQueryOpts[request, response](args))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetQueries(queries map[string]bool) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithQueryOpts[request, response](queries))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetNecessaryHeaders(headers ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range headers {
		args[headers[i]] = true
	}
	eb.options = append(eb.options, WithHeaderOpts[request, response](args))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetAdditionalHeaders(headers ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range headers {
		args[headers[i]] = false
	}
	eb.options = append(eb.options, WithHeaderOpts[request, response](args))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetHeaders(headers map[string]bool) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithHeaderOpts[request, response](headers))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetNecessaryCookies(cookies ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range cookies {
		args[cookies[i]] = true
	}
	eb.options = append(eb.options, WithCookieOpts[request, response](args))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetAdditionalCookies(cookies ...string) *EndPointBuilder[request, response] {
	args := map[string]bool{}
	for i := range cookies {
		args[cookies[i]] = false
	}
	eb.options = append(eb.options, WithCookieOpts[request, response](args))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetCookies(cookies map[string]bool) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithCookieOpts[request, response](cookies))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetRouter(router Router) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithRouterOpts[request, response](router))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetGinMiddlewares(middlewares ...gin.HandlerFunc) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithGinMiddlewaresOpts[request, response](middlewares...))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetHandlerChain(chain Chain[request, response]) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithChainOpts[request, response](chain))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetCustomPreprocessors(preprocessors ...EndpointPreprocessor[request, response]) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithCustomPreprocessors[request, response](preprocessors...))
	return eb
}

func (eb *EndPointBuilder[request, response]) SetCustomRender(enable bool) *EndPointBuilder[request, response] {
	eb.options = append(eb.options, WithCustomRender[request, response](enable))
	return eb
}

func (eb *EndPointBuilder[request, response]) Build() *EndPoint[request, response] {
	return NewEndPointWithOpts[request, response](eb.options...)
}

func NewEndPointBuilder[request, response any]() *EndPointBuilder[request, response] {
	return &EndPointBuilder[request, response]{options: []EndPointOptions[request, response]{}}
}

func NewEndPointGroup(group string, endpoints ...EndPointInterface) *EndpointGroup {
	return &EndpointGroup{
		groupUrl:  group,
		endpoints: endpoints,
	}
}
