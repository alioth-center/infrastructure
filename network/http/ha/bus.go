package ha

import (
	"errors"
	"reflect"

	"github.com/alioth-center/infrastructure/network/http"
	"github.com/gin-gonic/gin"

	"github.com/alioth-center/infrastructure/utils/concurrency"
)

var (
	endpoints   = concurrency.NewHashMap[string, ArrangedEndPoint](concurrency.HashMapNodeOptionSmallSize)
	handlers    = concurrency.NewHashMap[string, Arranged](concurrency.HashMapNodeOptionSmallSize)
	middlewares = concurrency.NewHashMap[string, gin.HandlersChain](concurrency.HashMapNodeOptionSmallSize)

	ErrEndpointNotFound   = errors.New("endpoint not found")
	ErrChainNotFound      = errors.New("chain not found")
	ErrTypeNotMatched     = errors.New("type not matched")
	ErrEmptyChain         = errors.New("empty chain")
	ErrMiddlewareNotFound = errors.New("middleware not found")

	// PriorityDefault is the default priority for chains, which is the lowest priority
	PriorityDefault = -1
	// PriorityLow is the second-lowest priority
	PriorityLow = 0
	// PriorityNormal is the normal priority
	PriorityNormal = 1 << 16
	// PriorityHigh is the second-highest priority
	PriorityHigh = 1 << 32
	// PriorityEmergency is the highest priority
	PriorityEmergency = 1 << 48
)

func RegisterArrangedChain[request, response any](name string, fns ...http.Handler[request, response]) {
	chain := NewArrangedChain[request, response](name, http.NewChain(fns...))
	handlers.Set(chain.UniqueName(), chain)
}

func RegisterPriorityArrangedChain[request, response any](name string, priority int, fns ...http.Handler[request, response]) {
	chain := NewArrangedChainWithPriority[request, response](name, priority, http.NewChain(fns...))
	handlers.Set(chain.UniqueName(), chain)
}

func GetArrangedChain[request, response any](name string, req request, res response) (chain ArrangedChain[request, response], err error) {
	arranged, got := handlers.Get(name)
	if !got {
		return nil, ErrChainNotFound
	}

	sameRequest := arranged.RequestType().AssignableTo(reflect.TypeOf(req))
	sameResponse := arranged.ResponseType().AssignableTo(reflect.TypeOf(res))
	if !sameRequest || !sameResponse {
		return nil, ErrTypeNotMatched
	}

	converted, success := arranged.(ArrangedChain[request, response])
	if !success {
		return nil, ErrTypeNotMatched
	}
	return converted, nil
}

func RegisterMiddleware(name string, middleware gin.HandlersChain) {
	middlewares.Set(name, middleware)
}

func GetMiddleware(name string) (middleware gin.HandlersChain, err error) {
	m, got := middlewares.Get(name)
	if !got {
		return nil, ErrMiddlewareNotFound
	}
	return m, nil
}

func RegisterEndPoint[request, response any](name string, allocRequest request, allocResponse response) {
	ep := &arrangedEndPoint[request, response]{
		uniqueName: name,
		req:        allocRequest,
		res:        allocResponse,
	}

	endpoints.Set(name, ep)
}

func GetEndPoint(name string) (endpoint ArrangedEndPoint, err error) {
	ep, got := endpoints.Get(name)
	if !got {
		return nil, ErrEndpointNotFound
	}
	return ep, nil
}
