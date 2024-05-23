package http

import (
	"bytes"
	"context"
	"github.com/alioth-center/infrastructure/utils/values"
	"io"
	"net/http"
	"testing"
	"time"

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
		if ctx.Error() != err {
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
