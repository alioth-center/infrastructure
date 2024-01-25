package http

import "github.com/alioth-center/infrastructure/utils/values"

type UnsupportedContentTypeError struct {
	ContentType string
}

func (e UnsupportedContentTypeError) Error() string {
	return "unsupported content type: " + e.ContentType
}

type UnsupportedAcceptError struct {
	Accept string
}

func (e UnsupportedAcceptError) Error() string {
	return "unsupported accept: " + e.Accept
}

type MethodNotAllowedError struct {
	Method string
}

func (e MethodNotAllowedError) Error() string {
	return "method not allowed: " + e.Method
}

type UnsupportedMethodError struct {
	Method string
}

func (e UnsupportedMethodError) Error() string {
	return "unsupported method: " + e.Method
}

type NecessaryHeaderMissingError struct {
	Header string
}

func (e NecessaryHeaderMissingError) Error() string {
	return "necessary header missing: " + e.Header
}

type NecessaryQueryMissingError struct {
	Query string
}

func (e NecessaryQueryMissingError) Error() string {
	return "necessary query missing: " + e.Query
}

type ServerAlreadyServingError struct {
	Address string
}

func (e ServerAlreadyServingError) Error() string {
	return "http server is already serving at " + e.Address
}

type NecessaryCookieMissingError struct {
	Cookie string
}

func (e NecessaryCookieMissingError) Error() string {
	return "necessary cookie missing: " + e.Cookie
}

type ContentTypeMismatchError struct {
	Expected string
	Actual   string
}

func (e ContentTypeMismatchError) Error() string {
	return values.BuildStrings("content type mismatch, expected: ", e.Expected, ", actual: ", e.Actual)
}
