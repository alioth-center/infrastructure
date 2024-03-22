package http

import (
	"bytes"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"testing"
	"time"
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
		//client := NewSimpleClient()
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
	gin.SetMode(gin.ReleaseMode)
	type Request struct {
		Msg string `json:"msg" vc:"key:msg,required"`
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
	t.Log(response.RawResponse().StatusCode)
	t.Log(response.BindJson(&receiver))
	t.Log(receiver)

	ex <- struct{}{}
}
