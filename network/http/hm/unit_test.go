package hm

import (
	"testing"
	"time"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
)

func TestLoggingMiddlewares(t *testing.T) {
	engine := http.NewEngine("/test")
	echoGroup := http.NewEndPointGroup("/echo")
	handler := http.NewEndPointBuilder[map[string]any, map[string]any]().
		SetHandlerChain(
			http.NewChain(
				TracingRequestMiddleware(logger.Default(), map[string]any{}, map[string]any{}),
				http.EchoHandler[map[string]any](),
			),
		).
		SetRouter(http.NewRouter("/:sub")).
		SetAllowMethods(http.GET, http.POST).
		SetNecessaryHeaders("Content-Type").
		SetNecessaryParams("sub").
		SetGinMiddlewares(TracingRawRequestMiddleware(logger.Default())).
		Build()
	echoGroup.AddEndPoints(handler)
	engine.AddEndPoints(echoGroup)

	engine.ServeAsync("0.0.0.0:8001", make(chan struct{}))
	time.Sleep(time.Millisecond * 100)

	response, executeErr := http.NewLoggerClient(logger.Default()).ExecuteRequest(
		http.NewRequestBuilder().
			WithPath("http://localhost:8001/test/echo/template").
			WithQuery("admin", "1").
			WithMethod(http.POST).
			WithCookie("test", "value").
			WithHeader("Ac-Request-Id", "114514").
			WithJsonBody(map[string]any{
				"int":    1,
				"string": "hello",
				"float":  114.514,
				"bool":   true,
			}),
	)
	if executeErr != nil {
		t.Fatal(executeErr)
	}

	receiver := map[string]any{}
	t.Log(response.BindJson(&receiver))
	t.Log(receiver)
}
