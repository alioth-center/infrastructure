package cls

import (
	"fmt"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/generate"
	"github.com/alioth-center/infrastructure/utils/timezone"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFieldParser(t *testing.T) {
	log := Logger{
		opts: Config{
			Locale: timezone.LocationUTCEast3Dot5,
		},
	}
	t.Log(log.prepareStructure(logger.NewFields(trace.NewContext()).WithMessage("hello")))
}

func TestFieldProcess(t *testing.T) {
	log := Logger{
		opts: Config{
			Locale: timezone.LocationUTCEast3Dot5,
		},
	}
	src := "aaaa-bbbb-cccc-dddd"
	dst := log.prepareField(src)
	if dst != strings.ReplaceAll(src, "-", "_") {
		t.Error("field process failed")
	}
}

func TestClsLogger(t *testing.T) {
	endpoint := os.Getenv("CLS_ENDPOINT")
	accessKeyID := os.Getenv("CLS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("CLS_ACCESS_KEY_SECRET")
	topic := os.Getenv("CLS_TOPIC_ID")
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || topic == "" {
		t.Skip("missing cls environment variables")
	}

	cfg := Config{
		Locale:     timezone.LocationShanghai,
		Service:    "infrastructure",
		Instance:   "sunist",
		SecretKey:  accessKeySecret,
		SecretID:   accessKeyID,
		Endpoint:   endpoint,
		TopicID:    topic,
		MaxRetries: 3,
		LogLocal:   true,
		LogLevel:   logger.LevelDebug,
	}
	fallback := logger.Default()
	cls, err := NewClsLogger(cfg, fallback)
	if err != nil {
		t.Fatal(err)
	}

	custom := map[string]any{
		"username": "sunist",
		"age":      14,
		"institution": map[string]any{
			"name": "alioth-center",
			"url":  "https://www.alioth.center",
			"projects": []string{
				"infrastructure",
				"restoration",
				"authentic",
				"akasha-terminal",
			},
		},
		"tags": []string{
			"test",
		},
		"key": generate.RandomBase62(16),
	}

	ctx := trace.NewContext()
	fmt.Printf("%#v\n", ctx)

	cls.Debug(logger.NewFields(ctx).WithMessage(generate.RandomBase62(2)).WithData(&custom))
	cls.Info(logger.NewFields(ctx).WithMessage(generate.RandomBase62(4)).WithData(&custom))
	cls.Warn(logger.NewFields(ctx).WithMessage(generate.RandomBase62(6)).WithData(&custom))
	cls.Error(logger.NewFields(ctx).WithMessage(generate.RandomBase62(8)).WithData(&custom))
	cls.Fatal(logger.NewFields(ctx).WithMessage(generate.RandomBase62(10)).WithData(&custom))

	time.Sleep(5 * time.Second)
}
