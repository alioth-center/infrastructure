package ha

import (
	"testing"
	"time"

	"github.com/alioth-center/infrastructure/config"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/network/http/hm"
	"github.com/gin-gonic/gin"
)

func TestBusWithParser(t *testing.T) {
	RegisterArrangedChain("echo-chain", http.EchoHandler[map[string]any]())
	RegisterPriorityArrangedChain("tracing-request", PriorityEmergency, hm.TracingRequestMiddleware(logger.Default(), map[string]any{}, map[string]any{}))
	RegisterEndPoint("echo-main", map[string]any{}, map[string]any{})
	RegisterEndPoint("echo-v1", map[string]any{}, map[string]any{})
	RegisterMiddleware("tracing-log", gin.HandlersChain{hm.TracingRawRequestMiddleware(logger.Default())})

	cfg := EngineConfig{}
	e := config.LoadConfig(&cfg, "./test.yml")
	if e != nil {
		t.Fatal(e)
	}

	_, err := ParseConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)
	response, executeErr := http.NewLoggerClient(logger.Default()).ExecuteRequest(
		http.NewRequestBuilder().
			WithPath("http://localhost:8003/api/echo").
			WithQuery("admin", "1").
			WithMethod(http.POST).
			WithCookie("test", "value").
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

	time.Sleep(100 * time.Millisecond)
}
