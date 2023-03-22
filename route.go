package gorouter

type HandlerFunc func(req *Request, res *Response)
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Route struct {
	Path       string
	Method     string
	Handler    HandlerFunc
	Middleware MiddlewareFunc
	Children   []Route
}
