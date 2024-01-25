package http

import (
	"github.com/alioth-center/infrastructure/utils/values"
	"strings"
)

type Router interface {
	Group(sub string) Router
	FullRouterPath() string
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

	r := &router{
		father:   nil,
		children: map[string]*router{},
		content:  base,
	}

	return r
}
