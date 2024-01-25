package http

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

func TestParsingError(t *testing.T) {
	t.Run("ParsingError:JsonError", func(t *testing.T) {
		resp := NewResponse(nil, errors.New("test error"))
		if _, err := ParseJsonResponse[struct{}](resp); err == nil {
			t.Errorf("parsing error should be occurred")
		}
	})

	t.Run("ParsingError:XmlError", func(t *testing.T) {
		resp := NewResponse(nil, errors.New("test error"))
		if _, err := ParseXmlResponse[struct{}](resp); err == nil {
			t.Errorf("parsing error should be occurred")
		}
	})
}

func TestResponseError(t *testing.T) {
	t.Run("ResponseError:NotJson", func(t *testing.T) {
		contentBytes := []byte("<html><body>test</body></html>")
		httpResp := http.Response{
			StatusCode:    http.StatusOK,
			Header:        http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
			Body:          io.NopCloser(bytes.NewBuffer(contentBytes)),
			ContentLength: int64(len(contentBytes)),
		}

		resp := NewResponse(&httpResp, nil)
		_, err := ParseJsonResponse[struct{}](resp)
		if err == nil {
			t.Errorf("parsing error should be occurred")
		}
		wantErr := ContentTypeMismatchError{
			Expected: ContentTypeJson,
			Actual:   "text/html; charset=utf-8",
		}
		if !errors.As(err, &wantErr) {
			t.Errorf("error is not expected: %+v", resp.Error())
		}
	})

	t.Run("ResponseError:NotXml", func(t *testing.T) {
		contentBytes := []byte("<html><body>test</body></html>")
		httpResp := http.Response{
			StatusCode:    http.StatusOK,
			Header:        http.Header{"Content-Type": []string{"text/html; charset=utf-8"}},
			Body:          io.NopCloser(bytes.NewBuffer(contentBytes)),
			ContentLength: int64(len(contentBytes)),
		}

		resp := NewResponse(&httpResp, nil)
		_, err := ParseXmlResponse[struct{}](resp)
		if err == nil {
			t.Errorf("parsing error should be occurred")
		}
		wantErr := ContentTypeMismatchError{
			Expected: ContentTypeTextXml,
			Actual:   "text/html; charset=utf-8",
		}
		if !errors.As(err, &wantErr) {
			t.Errorf("error is not expected: %+v", err)
		}
	})
}

func TestResponseBind(t *testing.T) {
	t.Run("ResponseBind:Json", func(t *testing.T) {
		contentBytes := []byte("{\"msg\":\"test\"}")
		httpResp := http.Response{
			StatusCode:    http.StatusOK,
			Header:        http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			Body:          io.NopCloser(bytes.NewBuffer(contentBytes)),
			ContentLength: int64(len(contentBytes)),
		}

		resp := NewResponse(&httpResp, nil)
		var data struct {
			Msg string `json:"msg"`
		}

		status, err := resp.BindJsonResult(&data)
		if err != nil {
			t.Errorf("parsing error should not be occurred: %+v", err)
		}
		if data.Msg != "test" {
			t.Errorf("data is not expected: %+v", data)
		}
		if status != StatusOK {
			t.Errorf("status code is not expected: %+v", status)
		}
	})

	t.Run("ResponseBind:Xml", func(t *testing.T) {
		contentBytes := []byte("<xml><msg>test</msg></xml>")
		httpResp := http.Response{
			StatusCode:    http.StatusOK,
			Header:        http.Header{"Content-Type": []string{"application/xml; charset=utf-8"}},
			Body:          io.NopCloser(bytes.NewBuffer(contentBytes)),
			ContentLength: int64(len(contentBytes)),
		}

		resp := NewResponse(&httpResp, nil)
		var data struct {
			Msg string `xml:"msg"`
		}

		status, err := resp.BindXmlResult(&data)
		if err != nil {
			t.Errorf("parsing error should not be occurred: %+v", err)
		}
		if data.Msg != "test" {
			t.Errorf("data is not expected: %+v", data)
		}
		if status != StatusOK {
			t.Errorf("status code is not expected: %+v", status)
		}
	})
}

func TestResponseGet(t *testing.T) {
	httpResp := http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":  []string{"application/json; charset=utf-8"},
			"Set-Cookie":    []string{"test=test"},
			"Authorization": []string{"Bearer test"},
			"Custom":        []string{"test"},
		},
		Body: io.NopCloser(bytes.NewBuffer([]byte("{\"msg\":\"test\"}"))),
	}
	resp := NewResponse(&httpResp, nil)

	status := resp.GetStatusCode()
	if status != StatusOK {
		t.Errorf("status code is not expected: %+v", status)
	}
	auth := resp.GetBearerToken()
	if auth != "test" {
		t.Errorf("auth is not expected: %+v", auth)
	}
	cookie := resp.GetCookie("test")
	if cookie == nil {
		t.Errorf("cookie is not expected: %+v", cookie)
	}
	if cookie.Value != "test" {
		t.Errorf("cookie value is not expected: %+v", cookie)
	}
	header := resp.GetHeader("Custom")
	if header != "test" {
		t.Errorf("header is not expected: %+v", header)
	}
	body := resp.GetBody()
	if string(body) != "{\"msg\":\"test\"}" {
		t.Errorf("body is not expected: %+v", string(body))
	}
	if resp.Error() != nil {
		t.Errorf("error is not expected: %+v", resp.Error())
	}

	resultStatus, stringBody, err := resp.StringResult()
	if err != nil {
		t.Errorf("error is not expected: %+v", err)
	}
	if resultStatus != StatusOK {
		t.Errorf("status code is not expected: %+v", resultStatus)
	}
	if stringBody != "{\"msg\":\"test\"}" {
		t.Errorf("body is not expected: %+v", stringBody)
	}
}
