package http

import (
	"github.com/alioth-center/infrastructure/utils/values"
	"strings"
)

// Router is the interface that wraps the basic methods of a router.
type Router interface {
	// Group returns a new router with the given sub path. sub path must begin with a slash and end without a slash.
	//
	// example:
	//	router := NewRouter("/api")
	//	subRouter := router.Group("/v1") // subRouter.FullRouterPath() == "/api/v1"
	//	subRouter2 := subRouter.Group("user") // subRouter2.FullRouterPath() == "/api/v1/user"
	//	subRouter3 := subRouter2.Group("info/") // subRouter3.FullRouterPath() == "/api/v1/user/info"
	Group(sub string) Router

	// Extend returns a new router with the given father router.
	//
	// example:
	//	router := NewRouter("/api")
	//	subRouter := NewRouter("/v1")
	//	subRouter.Extend(router) // subRouter2.FullRouterPath() == "/api/v1"
	Extend(father Router) Router

	// FullRouterPath returns the full path of the router.
	//
	// example:
	//	router := NewRouter("/api")
	//	subRouter := router.Group("/v1") // subRouter.FullRouterPath() == "/api/v1"
	FullRouterPath() string

	// BaseRouterPath returns the base path of the router.
	//
	// example:
	//	router := NewRouter("/api")
	//	subRouter := router.Group("/v1") // subRouter.BaseRouterPath() == "/v1"
	BaseRouterPath() string
}

type router struct {
	father   *router
	children map[string]*router

	content string
}

func (r *router) Group(sub string) Router {
	if sub == "" || sub == "." || sub == "/" {
		return r
	}

	if !strings.HasPrefix(sub, "/") {
		sub = "/" + sub
	}

	if strings.HasSuffix(sub, "/") {
		sub = strings.TrimSuffix(sub, "/")
	}

	nr := &router{father: r, children: map[string]*router{}, content: sub}

	if r.children == nil {
		r.children = map[string]*router{}
	}
	r.children[sub] = nr

	return nr
}

func (r *router) Extend(father Router) Router {
	if father == nil {
		return r
	}

	if r.father != nil {
		return r
	}

	router, ok := father.(*router)
	if !ok {
		return r
	}

	r.father = router
	return r
}

func (r *router) FullRouterPath() string {
	var fathers []string
	for currentPtr := r; currentPtr != nil; currentPtr = currentPtr.father {
		fathers = append(fathers, currentPtr.content)
	}

	return values.BuildStrings(values.ReverseArray(fathers)...)
}

func (r *router) BaseRouterPath() string {
	return r.content
}

func NewRouter(base string) Router {
	if base == "" {
		return &router{
			father:   nil,
			children: map[string]*router{},
			content:  "",
		}
	}

	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	if strings.HasSuffix(base, "/") {
		base = strings.TrimSuffix(base, "/")
	}

	r := &router{
		father:   nil,
		children: map[string]*router{},
		content:  base,
	}

	return r
}
