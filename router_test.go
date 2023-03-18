package gorouter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var tRoutes = []Route{
	{
		Path:   "/ping",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		},
	},
	{
		Path: "/api",
		Middleware: func(next HandlerFunc) HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				next(w, r)
			}
		},
		Children: []Route{
			{
				Path:   "/hello",
				Method: http.MethodPost,
				Middleware: func(next HandlerFunc) HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Test", "test")
						next(w, r)
					}
				},
				Handler: func(w http.ResponseWriter, r *http.Request) {
					var body map[string]any
					json.NewDecoder(r.Body).Decode(&body)
					w.Write([]byte(body["name"].(string)))
				},
			},
		},
	},
	{
		Path:   "/users/:userId",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.Context().Value("route").(string)))
		},
	},
	{
		Path:   "/users/:userId/blogs/:blogId",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.Context().Value("route").(string)))
		},
	},
	{
		Path: "/posts/:postId",
		Children: []Route{
			{
				Path:   "/comments/:commentId",
				Method: http.MethodGet,
				Handler: func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte(r.Context().Value("route").(string)))
				},
			},
		},
	},
}

var tRouter = New(tRoutes)

type ExpectedResponse struct {
	body       string
	statusCode int
}

type RouteTest struct {
	req *http.Request
	res ExpectedResponse
}

func testRoute(t *testing.T, rt RouteTest) http.ResponseWriter {
	w := httptest.NewRecorder()
	tRouter.ServeHTTP(w, rt.req)

	if code, want := w.Result().StatusCode, rt.res.statusCode; code != want {
		t.Errorf("got %d, wanted %d", code, want)
	}

	if text, want := w.Body.String(), rt.res.body; text != want {
		t.Errorf("got %s, wanted %s", text, want)
	}

	return w
}

func TestRoute(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/ping", nil),
		ExpectedResponse{"pong", 200},
	})
}

func TestNestedRoute(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodPost, "/api/hello", bytes.NewBufferString("{\"name\": \"John Doe\"}")),
		ExpectedResponse{"John Doe", 200},
	})
}

func TestDynamicRoutes(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/users/1", nil),
		ExpectedResponse{"/users/:userId", 200},
	})

	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/users/1/blogs/2", nil),
		ExpectedResponse{"/users/:userId/blogs/:blogId", 200},
	})
}

func TestNestedDynamicRoute(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/posts/1/comments/2", nil),
		ExpectedResponse{"/posts/:postId/comments/:commentId", 200},
	})
}

func TestMethodNotAllowed(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/api/hello", nil),
		ExpectedResponse{"", 405},
	})
}

func TestNotFound(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/hello", nil),
		ExpectedResponse{"404 page not found\n", 404},
	})
}

func TestMiddleware(t *testing.T) {
	w := testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodPost, "/api/hello", bytes.NewBufferString("{\"name\": \"John Doe\"}")),
		ExpectedResponse{"John Doe", 200},
	})

	ctWant := "application/json"
	if ct := w.Header().Get("Content-Type"); ct != ctWant {
		t.Errorf("Expected content-type: %s, got: %s", ctWant, ct)
	}
}

func TestNestedMiddleware(t *testing.T) {
	w := testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodPost, "/api/hello", bytes.NewBufferString("{\"name\": \"John Doe\"}")),
		ExpectedResponse{"John Doe", 200},
	})

	if result, want := w.Header().Get("Content-Type"), "application/json"; result != want {
		t.Errorf("Expected content-type: %s, got: %s", want, result)
	}

	if result, want := w.Header().Get("X-Test"), "test"; result != want {
		t.Errorf("Expected 'X-Test' header: %s, got: %s", want, result)
	}
}
