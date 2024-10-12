package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/gin-gonic/gin"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
)

func TestHttpClient(t *testing.T) {
	type Request struct {
		Msg string `json:"msg" xml:"msg"`
		Val int64  `json:"val" xml:"val"`
	}
	type Back struct {
		Json Request `json:"json" xml:"json"`
	}

	t.Run("RequestBuilder", func(t *testing.T) {
		builder := NewRequestBuilder().
			WithContext(trace.NewContext()).
			WithMethod(POST).
			WithPath("https://echo.apifox.com/post").
			WithPathFormat("https://echo.apifox.com/post?args1=%v", 1).
			WithPathTemplate("https://${hostname}/post?args1=$[args1]", map[string]string{"hostname": "echo.apifox.com", "args1": "1"}).
			WithQuery("args2", "2").
			WithHeader("testH", "testV").
			WithCookie("testC", "testV").
			WithBody(bytes.NewBufferString("test")).
			WithJsonBody(Request{Msg: "test", Val: 1}).
			WithUserAgent(AliothClient).
			WithBearerToken("114514").
			WithAccept(ContentTypeJson).
			WithContentType(ContentTypeJson).
			Clone()

		client := NewSimpleClient()
		res, err := client.ExecuteRequest(builder)
		if err != nil {
			t.Fatal(err)
		}

		receiver := map[string]any{}
		t.Log(res.BindJson(&receiver))
		t.Log(receiver)
	})

	t.Run("ResponseParser", func(t *testing.T) {
		client := NewSimpleClient()
		response, err := client.ExecuteRequest(NewRequestBuilder().
			WithPath("https://echo.apifox.com/post").
			WithMethod(POST).
			WithJsonBody(Request{Msg: "test", Val: 1}),
		)
		if err != nil {
			t.Fatal(err)
		}

		response.RawResponse()
		response.RawRequest()
		response.RawBody()
		response.Context()
		response.Status()
		t.Log(response.BindJson(&Back{}))
		t.Log(response.BindHeader("Content-Type"))
		t.Log(response.BindCookie("testC"))
		t.Log(response.BindCustom(map[string]any{}, func(reader io.Reader, receiver any) error {
			return nil
		}))
		t.Log(response.BindXml(&map[string]any{}))
	})

	t.Run("LoggerClient", func(t *testing.T) {
		client := NewLoggerClient(logger.Default())
		// client := NewSimpleClient()
		response, err := client.ExecuteRequest(NewRequestBuilder().
			WithPath("https://echo.apifox.com/post?fuck=you").
			WithMethod(POST).
			WithJsonBody(map[string]string{"fuck": "you"}).
			WithHeader("Content-Type", ContentTypeJson).
			WithAccept(ContentTypeJson),
		)
		if err != nil {
			t.Fatal(err)
		}

		receiver := map[string]any{}
		t.Log(response.BindJson(&receiver))
		t.Log(receiver)
	})

	t.Run("MockClient", func(t *testing.T) {
		client := NewMockClientWithLogger(logger.Default(), &MockOptions{
			Trigger: func(req *http.Request) bool { return true },
			Handler: func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"msg":"hello"}`)),
					Request:    req,
				}
			},
		})

		response, err := ParseJsonResponse[map[string]any](
			client.ExecuteRequest(NewRequestBuilder().WithPath("https://fuck.yo")),
		)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(response)
	})
}

func TestHttpServer(t *testing.T) {
	// gin.SetMode(gin.ReleaseMode)
	type Request struct {
		Msg string `json:"msg" vc:"key:msg,required"`
	}
	type Response struct {
		Msg string `json:"msg"`
	}

	t.Run("EndPointOpts", func(t *testing.T) {
		engine := NewEngine("/test")
		engine.AddEndPoints(
			NewEndPointWithOpts[Request, Response](
				WithRouterOpts[Request, Response](engine.BaseRouter().Group("/echo/:name")),
				WithAllowedMethodsOpts[Request, Response](GET, POST),
				WithChainOpts[Request, Response](NewChain[Request, Response](
					func(ctx Context[Request, Response]) {
						ctx.SetResponse(Response{Msg: ctx.PathParams().GetString("name") + ctx.Request().Msg})
						ctx.SetStatusCode(StatusOK)
					},
				)),
				WithHeaderOpts[Request, Response](map[string]bool{
					"Content-Type": true,
				}),
				WithParamOpts[Request, Response](map[string]bool{
					"name": true,
				}),
				WithQueryOpts[Request, Response](map[string]bool{
					"admin": true,
				}),
				WithCookieOpts[Request, Response](map[string]bool{
					"test": true,
				}),
			),
		)

		ex := make(chan struct{}, 1)
		ec := engine.ServeAsync("0.0.0.0:8080", ex)
		go func() {
			for {
				select {
				case e := <-ec:
					t.Errorf("http handlers occur error: %v", e)
				}
			}
		}()

		// 等待服务启动
		time.Sleep(time.Millisecond * 100)

		response, executeErr := NewLoggerClient(logger.Default()).ExecuteRequest(
			NewRequestBuilder().
				WithPath("http://localhost:8080/test/echo/sunist").
				WithQuery("admin", "1").
				WithMethod(POST).
				WithCookie("test", "test").
				WithHeader("Authorization", ContentTypeJson).
				WithJsonBody(&Request{Msg: ""}),
		)
		if executeErr != nil {
			t.Fatal(executeErr)
		}

		receiver := FrameworkResponse{}
		if !StatusCodeIs4XX(response) {
			t.Fatal("status code is not 4XX")
		}
		t.Log(response.RawResponse().StatusCode)
		t.Log(response.BindJson(&receiver))
		t.Log(receiver)

		ex <- struct{}{}
	})

	t.Run("EndPointBuilder", func(t *testing.T) {
		engine := NewEngine("/test1")
		echoGroup := NewEndPointGroup("/echo")
		nameHandler := NewEndPointBuilder[Request, Response]().
			SetHandlerChain(
				NewChain(func(ctx Context[Request, Response]) {
					ctx.SetResponse(Response{Msg: ctx.PathParams().GetString("name") + ctx.Request().Msg})
					ctx.SetStatusCode(StatusOK)
				}),
			).
			SetRouter(NewRouter("/:name")).
			SetAllowMethods(GET, POST).
			SetNecessaryHeaders("Content-Type").
			SetNecessaryParams("name").
			SetNecessaryQueries("admin").
			SetNecessaryCookies("test").
			Build()
		echoGroup.AddEndPoints(nameHandler)
		engine.AddEndPoints(echoGroup)

		ex := make(chan struct{}, 1)
		ec := engine.ServeAsync("0.0.0.0:8081", ex)
		go func() {
			for {
				select {
				case e := <-ec:
					t.Errorf("http handlers occur error: %v", e)
				}
			}
		}()

		// 等待服务启动
		time.Sleep(time.Millisecond * 100)

		response, executeErr := NewLoggerClient(logger.Default()).ExecuteRequest(
			NewRequestBuilder().
				WithPath("http://localhost:8081/test1/echo/sunist").
				WithQuery("admin", "1").
				WithMethod(POST).
				WithCookie("test", "test").
				WithHeader("Authorization", ContentTypeJson).
				WithHeader("Ac-Request-Id", "114514").
				WithJsonBody(&Request{Msg: ""}),
		)
		if executeErr != nil {
			t.Fatal(executeErr)
		}

		receiver := FrameworkResponse{}
		if !StatusCodeIs4XX(response) {
			t.Fatal("status code is not 4XX")
		}
		t.Log(response.BindJson(&receiver))
		if receiver.RequestID != "114514" {
			t.Fatal("request id is not equal")
		}
		t.Log(response.RawResponse().StatusCode)
		t.Log(receiver)

		ex <- struct{}{}
	})
}

func TestHttpFunctions(t *testing.T) {
	okResponse := &simpleParser{
		raw: &http.Response{
			Status:     "ok",
			StatusCode: 200,
		},
	}
	badResponse := &simpleParser{
		raw: &http.Response{
			Status:     "bad",
			StatusCode: 400,
		},
	}
	errorResponse := &simpleParser{
		raw: &http.Response{
			Status:     "error",
			StatusCode: 500,
		},
	}

	t.Run("StatusCode2XX", func(t *testing.T) {
		if !StatusCodeIs2XX(okResponse) {
			t.Error("StatusCodeIs2XX should return true")
		}
	})

	t.Run("StatusCode4XX", func(t *testing.T) {
		if !StatusCodeIs4XX(badResponse) {
			t.Error("StatusCodeIs4XX should return true")
		}
	})

	t.Run("StatusCode5XX", func(t *testing.T) {
		if !StatusCodeIs5XX(errorResponse) {
			t.Error("StatusCodeIs5XX should return true")
		}
	})

	t.Run("CheckStatusCode", func(t *testing.T) {
		if CheckStatusCode(okResponse) {
			t.Error("CheckStatusCode should return false when empty want list")
		}
		if !CheckStatusCode(okResponse, StatusOK) {
			t.Error("CheckStatusCode should return true when 200")
		}
		if CheckStatusCode(okResponse, StatusBadRequest, StatusBadGateway) {
			t.Error("CheckStatusCode should return false when 400 and 502")
		}
	})
}

func TestAcContext(t *testing.T) {
	type Request struct {
		Msg string `json:"msg"`
	}
	type Response struct {
		Msg string `json:"msg"`
	}

	t.Run("ContextLifecycle", func(t *testing.T) {
		ctx := NewContext[Request, Response]().(*acContext[Request, Response])

		// Test reset
		ctx.setHandlers(NewChain[Request, Response](
			func(cc Context[Request, Response]) { ctx.Next() },
			func(cc Context[Request, Response]) { ctx.Next() },
		))
		ctx.Next()
		if ctx.idx <= 2 {
			t.Fatal("Expected idx to be greater than 2 after Next")
		}
		ctx.reset()
		if ctx.idx != -1 {
			t.Fatal("Expected idx to be -1 after reset")
		}

		// Test setting handlers
		chain := NewChain[Request, Response](
			func(ctx Context[Request, Response]) { ctx.Next() },
		)
		ctx.setHandlers(chain)
		if len(ctx.h) != 1 {
			t.Fatal("Expected chain length to be 1")
		}

		// Test Next
		ctx.Next()
		if ctx.idx <= 1 {
			t.Fatal("Expected idx to be grater than 1 after Next")
		}

		// Test Abort
		ctx.Abort()
		if ctx.idx != 1 {
			t.Fatal("Expected idx to be 1 after Abort")
		}
		if !ctx.IsAborted() {
			t.Fatal("Expected IsAborted to return true")
		}

		// Test context methods
		rawReq, _ := http.NewRequest("GET", "http://example.com", nil)
		ctx.SetRawRequest(rawReq)
		if ctx.RawRequest() != rawReq {
			t.Fatal("Expected RawRequest to return the set raw request")
		}

		queryParams := Params{"key": "value"}
		ctx.SetQueryParams(queryParams)
		if ctx.QueryParams().GetString("key") != "value" {
			t.Fatal("Expected QueryParams to return the set query params")
		}

		pathParams := Params{"id": "123"}
		ctx.SetPathParams(pathParams)
		if ctx.PathParams().GetString("id") != "123" {
			t.Fatal("Expected PathParams to return the set path params")
		}

		headerParams := Params{"Authorization": "Bearer token"}
		ctx.SetHeaderParams(headerParams)
		if ctx.HeaderParams().GetString("Authorization") != "Bearer token" {
			t.Fatal("Expected HeaderParams to return the set header params")
		}

		cookieParams := Params{"session": "abc123"}
		ctx.SetCookieParams(cookieParams)
		if ctx.CookieParams().GetString("session") != "abc123" {
			t.Fatal("Expected CookieParams to return the set cookie params")
		}

		extraParams := Params{"extra": "data"}
		ctx.SetExtraParams(extraParams)
		if ctx.ExtraParams().GetString("extra") != "data" {
			t.Fatal("Expected ExtraParams to return the set extra params")
		}

		req := Request{Msg: "hello"}
		ctx.SetRequest(req)
		if ctx.Request().Msg != "hello" {
			t.Fatal("Expected Request to return the set request")
		}

		resp := Response{Msg: "world"}
		ctx.SetResponse(resp)
		if ctx.Response().Msg != "world" {
			t.Fatal("Expected Response to return the set response")
		}

		ctx.SetResponseHeader("Content-Type", "application/json")
		if ctx.ResponseHeaders().GetString("Content-Type") != "application/json" {
			t.Fatal("Expected ResponseHeaders to return the set response headers")
		}

		cookie := Cookie{Name: "session", Value: "abc123"}
		ctx.SetResponseSetCookie(cookie)
		if len(ctx.ResponseSetCookies()) != 1 || ctx.ResponseSetCookies()[0].Name != "session" {
			t.Fatal("Expected ResponseSetCookies to return the set response cookies")
		}

		ctx.SetStatusCode(StatusOK)
		if ctx.StatusCode() != StatusOK {
			t.Fatal("Expected StatusCode to return the set status code")
		}

		err := &http.ProtocolError{ErrorString: "test error"}
		ctx.SetError(err)
		if !errors.Is(err, ctx.Error()) {
			t.Fatal("Expected Error to return the set error")
		}

		testKey := "testKey"
		testValue := "testValue"
		ctx.SetValue(testKey, testValue)
		if ctx.Value(testKey) != testValue {
			t.Fatal("Expected Value to return the set value")
		}

		// Test context interface methods
		if _, ok := ctx.Deadline(); ok {
			t.Fatal("Expected Deadline to return false")
		}
		if ctx.Done() != nil {
			t.Fatal("Expected Done to return nil")
		}
		if ctx.Err() != nil {
			t.Fatal("Expected Err to return nil")
		}
		if ctx.Value(testKey) != testValue {
			t.Fatal("Expected Value to return the set value")
		}

		// test client ip
		ctx.raw.Header.Set("X-Forwarded-For", "1.1.1.1")
		if ctx.ClientIP() != "1.1.1.1" {
			t.Fatal("Expected ClientIP to return the set client ip")
		}
		ctx.raw.Header.Del("X-Forwarded-For")
		ctx.raw.Header.Set("X-Real-IP", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		if ctx.ClientIP() != "2001:0db8:85a3:0000:0000:8a2e:0370:7334" {
			t.Fatal("Expected ClientIP to return the set client ip")
		}
		ctx.raw.Header.Del("X-Real-IP")
		ctx.raw.RemoteAddr = "114.51.41.191:9810"
		if ctx.ClientIP() != "114.51.41.191" {
			t.Fatal("Expected ClientIP to return the set client ip" + ctx.ClientIP())
		}
		ctx.raw.RemoteAddr = ""
		ctx.extraParams[RemoteIPKey] = "127.0.0.1"
		if ctx.ClientIP() != "127.0.0.1" {
			t.Fatal("Expected ClientIP to return the set client ip")
		}
		ctx.extraParams[RemoteIPKey] = ""
		if ctx.ClientIP() != "" {
			t.Fatal("Expected ClientIP to return empty string")
		}
	})
}

func TestContextOpts(t *testing.T) {
	type Request struct {
		Msg string `json:"msg"`
	}
	type Response struct {
		Msg string `json:"msg"`
	}

	t.Run("WithContext", func(t *testing.T) {
		rawCtx := context.WithValue(context.Background(), "key", "value")
		ctx := NewContext[Request, Response](WithContext[Request, Response](rawCtx)).(*acContext[Request, Response])
		if ctx.Value("key") != "value" {
			t.Fatal("Expected context value to be set")
		}
	})

	t.Run("WithQueryParams", func(t *testing.T) {
		cc := acContext[Request, Response]{}
		cc.QueryParams()
		NewContext[Request, Response](WithQueryParams[Request, Response](nil))

		params := Params{"key": "value"}
		ctx := NewContext[Request, Response](WithQueryParams[Request, Response](params)).(*acContext[Request, Response])
		if ctx.QueryParams().GetString("key") != "value" {
			t.Fatal("Expected query params to be set")
		}
	})

	t.Run("WithPathParams", func(t *testing.T) {
		cc := acContext[Request, Response]{}
		cc.PathParams()
		NewContext[Request, Response](WithPathParams[Request, Response](nil))

		params := Params{"id": "123"}
		ctx := NewContext[Request, Response](WithPathParams[Request, Response](params)).(*acContext[Request, Response])
		if ctx.PathParams().GetString("id") != "123" {
			t.Fatal("Expected path params to be set")
		}
	})

	t.Run("WithHeaderParams", func(t *testing.T) {
		cc := acContext[Request, Response]{}
		cc.HeaderParams()
		NewContext[Request, Response](WithHeaderParams[Request, Response](nil))

		params := Params{"Authorization": "Bearer token"}
		ctx := NewContext[Request, Response](WithHeaderParams[Request, Response](params)).(*acContext[Request, Response])
		if ctx.HeaderParams().GetString("Authorization") != "Bearer token" {
			t.Fatal("Expected header params to be set")
		}
	})

	t.Run("WithCookieParams", func(t *testing.T) {
		cc := acContext[Request, Response]{}
		cc.CookieParams()
		NewContext[Request, Response](WithCookieParams[Request, Response](nil))

		params := Params{"session": "abc123"}
		ctx := NewContext[Request, Response](WithCookieParams[Request, Response](params)).(*acContext[Request, Response])
		if ctx.CookieParams().GetString("session") != "abc123" {
			t.Fatal("Expected cookie params to be set")
		}
	})

	t.Run("WithExtraParams", func(t *testing.T) {
		cc := acContext[Request, Response]{}
		cc.ExtraParams()

		cd := acContext[Request, Response]{}
		cd.SetExtraParam("extra", "data")
		cd.SetResponseHeader("Content-Type", "application/json")
		cd.SetResponseSetCookie(Cookie{Name: "session", Value: "abc123"})

		NewContext[Request, Response](WithExtraParams[Request, Response](nil))

		params := Params{"extra": "data"}
		ctx := NewContext[Request, Response](WithExtraParams[Request, Response](params)).(*acContext[Request, Response])
		if ctx.ExtraParams().GetString("extra") != "data" {
			t.Fatal("Expected extra params to be set")
		}
	})

	t.Run("WithRawRequest", func(t *testing.T) {
		rawReq, _ := http.NewRequest("GET", "http://example.com", nil)
		ctx := NewContext[Request, Response](WithRawRequest[Request, Response](rawReq)).(*acContext[Request, Response])
		if ctx.RawRequest() != rawReq {
			t.Fatal("Expected raw request to be set")
		}
	})

	t.Run("WithRequest", func(t *testing.T) {
		req := Request{Msg: "hello"}
		ctx := NewContext[Request, Response](WithRequest[Request, Response](req)).(*acContext[Request, Response])
		if ctx.Request().Msg != "hello" {
			t.Fatal("Expected request to be set")
		}
	})

	t.Run("WithRequestHeader", func(t *testing.T) {
		headers := RequestHeader{Authorization: "Bearer token"}
		ctx := NewContext[Request, Response](WithRequestHeader[Request, Response](headers)).(*acContext[Request, Response])
		if ctx.NormalHeaders().Authorization != "Bearer token" {
			t.Fatal("Expected request header to be set")
		}
	})

	t.Run("WithResponse", func(t *testing.T) {
		resp := Response{Msg: "world"}
		ctx := NewContext[Request, Response](WithResponse[Request, Response](resp)).(*acContext[Request, Response])
		if ctx.Response().Msg != "world" {
			t.Fatal("Expected response to be set")
		}
	})

	t.Run("WithResponseHeaders", func(t *testing.T) {
		NewContext[Request, Response](WithResponseHeaders[Request, Response](nil))

		headers := Params{"Content-Type": "application/json"}
		ctx := NewContext[Request, Response](WithResponseHeaders[Request, Response](headers)).(*acContext[Request, Response])
		if ctx.ResponseHeaders().GetString("Content-Type") != "application/json" {
			t.Fatal("Expected response headers to be set")
		}
	})

	t.Run("WithSetHeaders", func(t *testing.T) {
		NewContext[Request, Response](WithSetHeaders[Request, Response](nil))

		headers := Params{"Content-Type": "application/json"}
		ctx := NewContext[Request, Response](WithSetHeaders[Request, Response](headers)).(*acContext[Request, Response])
		if ctx.ResponseHeaders().GetString("Content-Type") != "application/json" {
			t.Fatal("Expected set headers to be set")
		}
	})

	t.Run("WithSetCookies", func(t *testing.T) {
		NewContext[Request, Response](WithSetCookies[Request, Response](nil))

		cookies := []Cookie{{Name: "session", Value: "abc123"}}
		ctx := NewContext[Request, Response](WithSetCookies[Request, Response](cookies)).(*acContext[Request, Response])
		if len(ctx.ResponseSetCookies()) != 1 || ctx.ResponseSetCookies()[0].Name != "session" {
			t.Fatal("Expected set cookies to be set")
		}
	})

	t.Run("WithStatusCode", func(t *testing.T) {
		ctx := NewContext[Request, Response](WithStatusCode[Request, Response](StatusOK)).(*acContext[Request, Response])
		if ctx.StatusCode() != StatusOK {
			t.Fatal("Expected status code to be set")
		}
	})

	t.Run("WithError", func(t *testing.T) {
		err := &http.ProtocolError{ErrorString: "test error"}
		ctx := NewContext[Request, Response](WithError[Request, Response](err)).(*acContext[Request, Response])
		if ctx.Error() != err {
			t.Fatal("Expected error to be set")
		}
	})
}

func TestMethodList(t *testing.T) {
	ml := methodList(0)
	if ml.isAllowedAll() {
		t.Fatal("methodList should not allow all methods initially")
	}
	if !ml.allowAll().isAllowedAll() {
		t.Fatal("methodList should allow all methods after calling allowAll")
	}
	if ml.allowNone().isAllowedAll() {
		t.Fatal("methodList should allow no methods after calling allowNone")
	}

	m := methodList(0).allowMethod(GET, POST)
	if _, allowed := m.isAllowed("GET"); !allowed {
		t.Fatal("methodList should allow GET method")
	}
	if _, allowed := m.isAllowed("DELETE"); allowed {
		t.Fatal("methodList should not allow DELETE method")
	}
	methods := m.allowedMethods()
	if len(methods) != 2 {
		t.Fatal("methodList allowed methods are incorrect")
	}
	if !values.ContainsArray(methods, "GET") || !values.ContainsArray(methods, "POST") {
		t.Fatal("methodList allowed methods are incorrect")
	}
}

func TestParams(t *testing.T) {
	params := Params{
		"string": "test",
		"int":    "123",
		"uint":   "123",
		"float":  "123.45",
		"bool":   "true",
	}

	if params.GetString("string") != "test" {
		t.Fatal("Params GetString failed")
	}
	if params.GetInt("int") != 123 {
		t.Fatal("Params GetInt failed")
	}
	if params.GetUint("uint") != 123 {
		t.Fatal("Params GetUint failed")
	}
	if params.GetFloat("float") != 123.45 {
		t.Fatal("Params GetFloat failed")
	}
	if !params.GetBool("bool") {
		t.Fatal("Params GetBool failed")
	}
}

func TestHttpFunction(t *testing.T) {
	okResponse := &simpleParser{
		raw: &http.Response{
			Status:     "ok",
			StatusCode: 200,
		},
	}
	badResponse := &simpleParser{
		raw: &http.Response{
			Status:     "bad",
			StatusCode: 400,
		},
	}
	errorResponse := &simpleParser{
		raw: &http.Response{
			Status:     "error",
			StatusCode: 500,
		},
	}

	t.Run("StatusCode2XX", func(t *testing.T) {
		if !StatusCodeIs2XX(okResponse) {
			t.Error("StatusCodeIs2XX should return true")
		}
	})

	t.Run("StatusCode4XX", func(t *testing.T) {
		if !StatusCodeIs4XX(badResponse) {
			t.Error("StatusCodeIs4XX should return true")
		}
	})

	t.Run("StatusCode5XX", func(t *testing.T) {
		if !StatusCodeIs5XX(errorResponse) {
			t.Error("StatusCodeIs5XX should return true")
		}
	})

	t.Run("CheckStatusCode", func(t *testing.T) {
		if CheckStatusCode(okResponse) {
			t.Error("CheckStatusCode should return false when empty want list")
		}
		if !CheckStatusCode(okResponse, StatusOK) {
			t.Error("CheckStatusCode should return true when 200")
		}
		if CheckStatusCode(okResponse, StatusBadRequest, StatusBadGateway) {
			t.Error("CheckStatusCode should return false when 400 and 502")
		}
	})
}

func TestHttpErrors(t *testing.T) {
	_ = UnsupportedAcceptError{}.Error()
	_ = UnsupportedContentTypeError{}.Error()
	_ = MethodNotAllowedError{}.Error()
	_ = UnsupportedMethodError{}.Error()
	_ = NecessaryHeaderMissingError{}.Error()
	_ = NecessaryQueryMissingError{}.Error()
	_ = ServerAlreadyServingError{}.Error()
	_ = NecessaryCookieMissingError{}.Error()
	_ = ContentTypeMismatchError{}.Error()
}

func TestHttpDefines(t *testing.T) {
	SetTraceHeaderKey(TraceHeaderKey())
	SetErrorContextKey(ErrorContextKey())
}

func TestHttpCookies(t *testing.T) {
	NewBasicCookie("name", "value")
}

func TestHttpChains(t *testing.T) {
	c := NewChain[any, any]()
	c = NewChain[any, any]([]Handler[any, any]{}...)
	c.AddHandlerBack()
	c.AddHandlerFront()
	c.Execute(NewContext[any, any]())
}

func TestRouter(t *testing.T) {
	// Test NewRouter with empty base
	r := NewRouter("")
	if r.FullRouterPath() != "" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "", r.FullRouterPath())
	}
	if r.BaseRouterPath() != "" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "", r.BaseRouterPath())
	}

	// Test NewRouter with base '/api'
	r = NewRouter("/api")
	if r.FullRouterPath() != "/api" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api", r.FullRouterPath())
	}
	if r.BaseRouterPath() != "/api" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "/api", r.BaseRouterPath())
	}

	// Test NewRouter with base 'api'
	r = NewRouter("api")
	if r.FullRouterPath() != "/api" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api", r.FullRouterPath())
	}
	if r.BaseRouterPath() != "/api" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "/api", r.BaseRouterPath())
	}

	// Test Group method
	subRouter := r.Group("/v1")
	if subRouter.FullRouterPath() != "/api/v1" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1", subRouter.FullRouterPath())
	}
	if subRouter.BaseRouterPath() != "/v1" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "/v1", subRouter.BaseRouterPath())
	}

	subRouter2 := subRouter.Group("user")
	if subRouter2.FullRouterPath() != "/api/v1/user" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1/user", subRouter2.FullRouterPath())
	}
	if subRouter2.BaseRouterPath() != "/user" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "user", subRouter2.BaseRouterPath())
	}

	subRouter3 := subRouter2.Group("info/")
	if subRouter3.FullRouterPath() != "/api/v1/user/info" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1/user/info", subRouter3.FullRouterPath())
	}
	if subRouter3.BaseRouterPath() != "/info" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "info", subRouter3.BaseRouterPath())
	}

	// Test Group with empty sub path
	subRouter4 := subRouter3.Group("")
	if subRouter4.FullRouterPath() != "/api/v1/user/info" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1/user/info", subRouter4.FullRouterPath())
	}
	if subRouter4.BaseRouterPath() != "/info" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "info", subRouter4.BaseRouterPath())
	}

	// Test Extend method
	r = NewRouter("/api")
	subRouter = NewRouter("/v1")
	subRouter.Extend(r)
	if subRouter.FullRouterPath() != "/api/v1" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1", subRouter.FullRouterPath())
	}
	if subRouter.BaseRouterPath() != "/v1" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "/v1", subRouter.BaseRouterPath())
	}

	// Test extending an already extended router
	subRouter2 = NewRouter("/v2")
	subRouter2.Extend(r)
	subRouter2.Extend(subRouter) // should not change anything
	if subRouter2.FullRouterPath() != "/api/v2" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v2", subRouter2.FullRouterPath())
	}
	if subRouter2.BaseRouterPath() != "/v2" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "/v2", subRouter2.BaseRouterPath())
	}

	// Test FullRouterPath
	r = NewRouter("/api")
	if r.FullRouterPath() != "/api" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api", r.FullRouterPath())
	}

	subRouter = r.Group("/v1")
	if subRouter.FullRouterPath() != "/api/v1" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1", subRouter.FullRouterPath())
	}

	subRouter2 = subRouter.Group("user")
	if subRouter2.FullRouterPath() != "/api/v1/user" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1/user", subRouter2.FullRouterPath())
	}

	subRouter3 = subRouter2.Group("info/")
	if subRouter3.FullRouterPath() != "/api/v1/user/info" {
		t.Errorf("expected FullRouterPath '%s' but got '%s'", "/api/v1/user/info", subRouter3.FullRouterPath())
	}

	// Test BaseRouterPath
	r = NewRouter("/api")
	if r.BaseRouterPath() != "/api" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "/api", r.BaseRouterPath())
	}

	subRouter = r.Group("/v1")
	if subRouter.BaseRouterPath() != "/v1" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "/v1", subRouter.BaseRouterPath())
	}

	subRouter2 = subRouter.Group("user")
	if subRouter2.BaseRouterPath() != "/user" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "user", subRouter2.BaseRouterPath())
	}

	subRouter3 = subRouter2.Group("info/")
	if subRouter3.BaseRouterPath() != "/info" {
		t.Errorf("expected BaseRouterPath '%s' but got '%s'", "info", subRouter3.BaseRouterPath())
	}
}

func TestEndPointBuilder(t *testing.T) {
	type request struct{}
	type response struct{}

	// Test NewBasicEndPoint
	ep := NewBasicEndPoint[request, response](GET, nil, nil)
	if ep == nil {
		t.Errorf("NewBasicEndPoint returned nil")
	}

	// Test SetAllowMethods
	builder := &EndPointBuilder[request, response]{}
	builder.SetAllowMethods(GET, POST)
	if len(builder.options) != 1 {
		t.Errorf("SetAllowMethods failed, expected 1 option, got %d", len(builder.options))
	}

	// Test SetNecessaryParams
	builder.SetNecessaryParams("param1", "param2")
	if len(builder.options) != 2 {
		t.Errorf("SetNecessaryParams failed, expected 2 options, got %d", len(builder.options))
	}

	// Test SetAdditionalParams
	builder.SetAdditionalParams("param3", "param4")
	if len(builder.options) != 3 {
		t.Errorf("SetAdditionalParams failed, expected 3 options, got %d", len(builder.options))
	}

	// Test SetParams
	builder.SetParams(map[string]bool{"param5": true, "param6": false})
	if len(builder.options) != 4 {
		t.Errorf("SetParams failed, expected 4 options, got %d", len(builder.options))
	}

	// Test SetNecessaryQueries
	builder.SetNecessaryQueries("query1", "query2")
	if len(builder.options) != 5 {
		t.Errorf("SetNecessaryQueries failed, expected 5 options, got %d", len(builder.options))
	}

	// Test SetAdditionalQueries
	builder.SetAdditionalQueries("query3", "query4")
	if len(builder.options) != 6 {
		t.Errorf("SetAdditionalQueries failed, expected 6 options, got %d", len(builder.options))
	}

	// Test SetQueries
	builder.SetQueries(map[string]bool{"query5": true, "query6": false})
	if len(builder.options) != 7 {
		t.Errorf("SetQueries failed, expected 7 options, got %d", len(builder.options))
	}

	// Test SetNecessaryHeaders
	builder.SetNecessaryHeaders("header1", "header2")
	if len(builder.options) != 8 {
		t.Errorf("SetNecessaryHeaders failed, expected 8 options, got %d", len(builder.options))
	}

	// Test SetAdditionalHeaders
	builder.SetAdditionalHeaders("header3", "header4")
	if len(builder.options) != 9 {
		t.Errorf("SetAdditionalHeaders failed, expected 9 options, got %d", len(builder.options))
	}

	// Test SetHeaders
	builder.SetHeaders(map[string]bool{"header5": true, "header6": false})
	if len(builder.options) != 10 {
		t.Errorf("SetHeaders failed, expected 10 options, got %d", len(builder.options))
	}

	// Test SetNecessaryCookies
	builder.SetNecessaryCookies("cookie1", "cookie2")
	if len(builder.options) != 11 {
		t.Errorf("SetNecessaryCookies failed, expected 11 options, got %d", len(builder.options))
	}

	// Test SetAdditionalCookies
	builder.SetAdditionalCookies("cookie3", "cookie4")
	if len(builder.options) != 12 {
		t.Errorf("SetAdditionalCookies failed, expected 12 options, got %d", len(builder.options))
	}

	// Test SetCookies
	builder.SetCookies(map[string]bool{"cookie5": true, "cookie6": false})
	if len(builder.options) != 13 {
		t.Errorf("SetCookies failed, expected 13 options, got %d", len(builder.options))
	}

	// Test SetRouter
	builder.SetRouter(nil)
	if len(builder.options) != 14 {
		t.Errorf("SetRouter failed, expected 14 options, got %d", len(builder.options))
	}

	// Test SetGinMiddlewares
	builder.SetGinMiddlewares(nil)
	if len(builder.options) != 15 {
		t.Errorf("SetGinMiddlewares failed, expected 15 options, got %d", len(builder.options))
	}

	// Test SetHandlerChain
	builder.SetHandlerChain(nil)
	if len(builder.options) != 16 {
		t.Errorf("SetHandlerChain failed, expected 16 options, got %d", len(builder.options))
	}

	// Test SetCustomPreprocessors
	builder.SetCustomPreprocessors(nil)
	if len(builder.options) != 17 {
		t.Errorf("SetCustomPreprocessors failed, expected 17 options, got %d", len(builder.options))
	}
}

func TestEngine(t *testing.T) {
	// Test NewEngine
	engine := NewEngine("/base")
	if engine == nil {
		t.Fatal("NewEngine returned nil")
	}
	if engine.core == nil {
		t.Fatal("Engine core is nil")
	}
	if engine.baseRouter == nil {
		t.Fatal("Engine baseRouter is nil")
	}
	if len(engine.endpoints) != 0 {
		t.Fatalf("Expected 0 endpoints, got %d", len(engine.endpoints))
	}
	if len(engine.middlewares) != 0 {
		t.Fatalf("Expected 0 middlewares, got %d", len(engine.middlewares))
	}

	// Test AddMiddlewares
	mw1 := func(c *gin.Context) {}
	mw2 := func(c *gin.Context) {}
	engine.AddMiddlewares(mw1, mw2)
	if len(engine.middlewares) != 2 {
		t.Fatalf("Expected 2 middlewares, got %d", len(engine.middlewares))
	}

	// Test registerEndpoints
	engine.registerEndpoints()
	if len(engine.core.Handlers) != 4 { // 1 for traceContext, 2 for added middlewares
		t.Fatalf("Expected 4 handlers, got %d", len(engine.core.Handlers))
	}

	// Test traceContext
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	engine.traceContext(ctx)
	if ctx.GetString(trace.ContextKey()) == "" {
		t.Fatal("Expected trace ID in context, got empty string")
	}

	// Test defaultHandler
	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)
	engine.defaultHandler(ctx)
}
