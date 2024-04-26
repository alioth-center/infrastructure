package cls

import (
	"fmt"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/encrypt"
	"github.com/alioth-center/infrastructure/utils/timezone"
	"github.com/alioth-center/infrastructure/utils/values"
	tcls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

var (
	clients = concurrency.NewHashMap[string, *tcls.AsyncProducerClient](concurrency.HashMapNodeOptionSmallSize)
)

type client struct {
	instance *tcls.AsyncProducerClient
	fallback logger.Logger
	opts     Config
}

func (c *client) init() error {
	clientKey := encrypt.HashMD5(values.BuildStrings(c.opts.Endpoint, c.opts.TopicID))
	if ist, ok := clients.Get(clientKey); ok && ist != nil {
		// got exist cls client instance, reuse it
		c.instance = ist
		return nil
	}

	// create new cls client instance
	cfg := tcls.GetDefaultAsyncProducerClientConfig()
	cfg.Endpoint = c.opts.Endpoint
	cfg.AccessKeyID = c.opts.SecretID
	cfg.AccessKeySecret = c.opts.SecretKey
	instance, initErr := tcls.NewAsyncProducerClient(cfg)
	if initErr != nil {
		return fmt.Errorf("failed to create cls client: %w", initErr)
	}
	if instance == nil {
		return fmt.Errorf("failed to create cls client: nil instance")
	}

	// start cls client
	instance.Start()
	exit.Register(c.exit, values.BuildStrings("cls client exit: ", c.opts.TopicID))

	clients.Set(clientKey, instance)
	c.instance = instance
	return nil
}

func (c *client) Success(_ *tcls.Result) {}

func (c *client) Fail(result *tcls.Result) {
	c.fallback.Error(logger.NewFields(trace.NewContext()).
		WithMessage("failed to send log to cls").
		WithTraceID(result.GetRequestId()).
		WithField("error_code", result.GetErrorCode()).
		WithField("error_message", result.GetErrorMessage()).
		WithData(result.GetReservedAttempts()))
}

func (c *client) exit(sig string) string {
	if c.instance != nil {
		e := c.instance.Close(1000)
		if e != nil {
			return values.BuildStrings("failed to close cls client: ", e.Error())
		}
	}

	return values.BuildStrings("cls client closed with signal: ", sig)
}

func (c *client) execute(fields map[string]string) {
	log := tcls.NewCLSLog(timezone.NowInZeroTimeUnix(), fields)
	go func(fields map[string]string) {
		var e error = nil
		for i := 0; i < c.opts.MaxRetries; i++ {
			e = c.instance.SendLog(c.opts.TopicID, log, c)
			if e == nil {
				return
			}
		}

		c.fallback.Error(logger.NewFields(trace.NewContextWithTid(fields["tid"])).WithMessage("cls fallback log").WithData(fields))
	}(fields)
}

func newClsClient(opts Config, fallback logger.Logger) (cli *client, err error) {
	cli = &client{
		fallback: fallback,
		opts:     opts,
	}
	initErr := cli.init()
	if initErr != nil {
		return nil, fmt.Errorf("failed to create cls client: %w", initErr)
	}

	return cli, nil
}
