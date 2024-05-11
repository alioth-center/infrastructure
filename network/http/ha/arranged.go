// Package ha: http arranged handler chain

package ha

import (
	"reflect"

	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/gin-gonic/gin"

	"github.com/alioth-center/infrastructure/network/http"
)

type Arranged interface {
	UniqueName() string
	Priority() int
	RequestType() reflect.Type
	ResponseType() reflect.Type
}

type ArrangedChain[request, response any] interface {
	Arranged
	Chain() http.Chain[request, response]
	Request() request
	Response() response
}

type arrangedChain[request, response any] struct {
	uniqueName string
	priority   int
	chain      http.Chain[request, response]
	req        request
	res        response
}

func (c arrangedChain[request, response]) UniqueName() string {
	return c.uniqueName
}

func (c arrangedChain[request, response]) Priority() int {
	return c.priority
}

func (c arrangedChain[request, response]) Chain() http.Chain[request, response] {
	return c.chain
}

func (c arrangedChain[request, response]) RequestType() reflect.Type {
	return reflect.TypeOf(c.req)
}

func (c arrangedChain[request, response]) ResponseType() reflect.Type {
	return reflect.TypeOf(c.res)
}

func (c arrangedChain[request, response]) Request() request {
	return c.req
}

func (c arrangedChain[request, response]) Response() response {
	return c.res
}

func NewArrangedChain[request, response any](name string, chain http.Chain[request, response]) ArrangedChain[request, response] {
	return &arrangedChain[request, response]{
		uniqueName: name,
		chain:      chain,
		req:        values.Nil[request](),
		res:        values.Nil[response](),
		priority:   -1,
	}
}

func NewArrangedChainWithPriority[request, response any](name string, priority int, chain http.Chain[request, response]) ArrangedChain[request, response] {
	return &arrangedChain[request, response]{
		uniqueName: name,
		chain:      chain,
		req:        values.Nil[request](),
		res:        values.Nil[response](),
		priority:   priority,
	}
}

type ArrangedEndPoint interface {
	UniqueName() string
	ParseConfig(opts EndPointConfig) (endpoint http.EndPointInterface, err error)
}

type arrangedEndPoint[request, response any] struct {
	*http.EndPoint[request, response]
	uniqueName  string
	middlewares gin.HandlersChain
	handlers    http.Chain[request, response]
	req         request
	res         response
}

func (e *arrangedEndPoint[request, response]) UniqueName() string {
	return e.uniqueName
}

func (e *arrangedEndPoint[request, response]) ParseConfig(opts EndPointConfig) (endpoint http.EndPointInterface, err error) {
	mds := gin.HandlersChain{}
	for _, name := range opts.Middlewares {
		middleware, getErr := GetMiddleware(name)
		if getErr != nil {
			return nil, getErr
		}

		mds = append(mds, middleware...)
	}

	var hs []http.Chain[request, response]
	for _, name := range opts.Chains {
		arranged, parseErr := GetArrangedChain(name, e.req, e.res)
		if parseErr != nil {
			return nil, parseErr
		}

		hs = append(hs, arranged.Chain())
	}
	values.SortArray(hs, func(a http.Chain[request, response], b http.Chain[request, response]) bool {
		return len(a) > len(b)
	})
	if len(hs) == 0 {
		return nil, ErrEmptyChain
	}

	handler := hs[0]
	for i := 1; i < len(hs); i++ {
		handler = append(handler, hs[i]...)
	}

	e.EndPoint = http.NewEndPointBuilder[request, response]().
		SetRouter(http.NewRouter(opts.Path)).
		SetAllowMethods(opts.Methods...).
		SetHandlerChain(handler).
		SetGinMiddlewares(mds...).
		SetNecessaryHeaders(opts.Headers.Necessary...).
		SetAdditionalHeaders(opts.Headers.Additional...).
		SetNecessaryParams(opts.Paths.Necessary...).
		SetAdditionalParams(opts.Paths.Additional...).
		SetNecessaryQueries(opts.Queries.Necessary...).
		SetAdditionalQueries(opts.Queries.Additional...).
		SetNecessaryCookies(opts.Cookies.Necessary...).
		SetAdditionalCookies(opts.Cookies.Additional...).
		Build()

	return e, nil
}
