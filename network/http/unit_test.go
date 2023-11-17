package http

import (
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
	} else {
		t.Logf("%+v", responseData)
	}
}
