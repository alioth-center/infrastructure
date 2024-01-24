package http

import (
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func TestHttpCall(t *testing.T) {
	type request struct {
		Msg string `json:"msg"`
		Val int64  `json:"val"`
	}
	type back struct {
		Json request `json:"json"`
	}

	responseData, parseErr := ParseJsonResponse[back](NewClient().ExecuteRequest(NewRequest().
		SetUrl("https://echo.apifox.com/post").
		SetMethod(POST).
		SetUserAgent(Curl).
		SetJsonBody(&request{Msg: "HelloWorld", Val: time.Now().Unix()}).
		SetAccept(ContentTypeJson),
	))

	if parseErr != nil {
		t.Fatal(parseErr)
	}

	t.Logf("%+v", responseData)
}

func TestHttpServer(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	type Request struct {
		Msg string `json:"msg"`
	}
	type Response struct {
		Msg string `json:"msg"`
	}

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
	responseData, parseErr := ParseJsonResponse[Response](NewClient().ExecuteRequest(NewRequest().
		SetUrl("http://localhost:8080/test/echo/sunist?admin=1").
		SetMethod(POST).
		SetUserAgent(Curl).
		SetCookie("test", "test").
		SetJsonBody(&Request{Msg: "HelloWorld"}),
	))

	if parseErr != nil {
		t.Fatal(parseErr)
	}
	if responseData.Msg != "sunistHelloWorld" {
		t.Errorf("response data is not expected: %+v", responseData)
	}

	ex <- struct{}{}
}
