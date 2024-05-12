package ha

import (
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/network/http"
)

func ParseConfig(config EngineConfig) (engine *http.Engine, err error) {
	engine = http.NewEngine(config.Path)
	group, parseErr := ParseRouter(config.Router)
	if parseErr != nil {
		return nil, parseErr
	}

	engine.AddEndPoints(group)

	if !config.Serving {
		return engine, nil
	}

	if config.Block {
		return engine, engine.Serve(config.Bind)
	}

	exitChan := make(chan struct{})
	exit.Register(func(sig string) string {
		exitChan <- struct{}{}
		return "arranged http engine exited"
	}, "exit arranged http engine")
	engine.ServeAsync(config.Bind, exitChan)

	return engine, nil
}

func ParseRouter(config RouterConfig) (group *http.EndpointGroup, err error) {
	group = http.NewEndPointGroup(config.Path)

	for _, ep := range config.EndPoints {
		arranged, gotErr := GetEndPoint(ep.Name)
		if gotErr != nil {
			return nil, gotErr
		}

		endpoint, parseErr := arranged.ParseConfig(ep)
		if parseErr != nil {
			return nil, parseErr
		}

		group.AddEndPoints(endpoint)
	}

	for _, sub := range config.SubRouters {
		subGroup, parseErr := ParseRouter(sub)
		if parseErr != nil {
			return nil, parseErr
		}

		group.AddEndPoints(subGroup)
	}

	return group, nil
}
