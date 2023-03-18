package gorouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var bRoutes = []Route{
	{
		Path:    "/",
		Method:  http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {},
	},
	{
		Path:    "/ping",
		Method:  http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {},
	},
	{
		Path:    "/users/:id",
		Method:  http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {},
	},
	{
		Path:    "/:resource/:id",
		Method:  http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {},
	},
	{
		Path:    "/:a/:b/:c/:d/:e",
		Method:  http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {},
	},
}

var bRouter = New(bRoutes)

func benchmarkRoute(b *testing.B, r *http.Request) {
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		bRouter.ServeHTTP(w, r)
	}
}

func BenchmarkRootRoute(b *testing.B) {
	benchmarkRoute(b, httptest.NewRequest(http.MethodGet, "/", nil))
}

func BenchmarkSimpleRoute(b *testing.B) {
	benchmarkRoute(b, httptest.NewRequest(http.MethodGet, "/ping", nil))
}

func Benchmark1ParamsRoute(b *testing.B) {
	benchmarkRoute(b, httptest.NewRequest(http.MethodGet, "/users/1", nil))
}

func Benchmark2ParamsRoute(b *testing.B) {
	benchmarkRoute(b, httptest.NewRequest(http.MethodGet, "/blogs/fa221f04-109c-4f8a-8075-495e83e5ba5b", nil))
}

func Benchmark5ParamsRoute(b *testing.B) {
	benchmarkRoute(b, httptest.NewRequest(http.MethodGet, "/1/2/3/4/5", nil))
}
