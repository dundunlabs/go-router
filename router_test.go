package gorouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var tRoutes = []Route{
	{
		Path:   "/ping",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r Request) {
			w.Write([]byte("pong"))
		},
	},
	{
		Path: "/api",
		Middleware: func(next HandlerFunc) HandlerFunc {
			return func(w http.ResponseWriter, r Request) {
				w.Header().Set("Content-Type", "application/json")
				next(w, r)
			}
		},
		Children: []Route{
			{
				Path:   "/hello",
				Method: http.MethodPost,
				Middleware: func(next HandlerFunc) HandlerFunc {
					return func(w http.ResponseWriter, r Request) {
						w.Header().Set("X-Test", "test")
						next(w, r)
					}
				},
				Handler: func(w http.ResponseWriter, r Request) {
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
		Handler: func(w http.ResponseWriter, r Request) {
			w.Header().Set("X-Params", fmt.Sprintf("%s", r.Params()))
			w.Write([]byte(r.Route()))
		},
	},
	{
		Path:   "/users/:userId/blogs/:blogId",
		Method: http.MethodGet,
		Handler: func(w http.ResponseWriter, r Request) {
			w.Header().Set("X-Params", fmt.Sprintf("%s", r.Params()))
			w.Write([]byte(r.Route()))
		},
	},
	{
		Path: "/posts/:postId",
		Children: []Route{
			{
				Path:   "/comments/:commentId",
				Method: http.MethodGet,
				Handler: func(w http.ResponseWriter, r Request) {
					w.Header().Set("X-Params", fmt.Sprintf("%s", r.Params()))
					w.Write([]byte(r.Route()))
				},
			},
		},
	},
	{
		Path: "/:foo",
		Children: []Route{
			{
				Path:   "/foo",
				Method: http.MethodGet,
				Handler: func(w http.ResponseWriter, r Request) {
					w.Header().Set("X-Params", fmt.Sprintf("%s", r.Params()))
					w.Write([]byte(r.Route()))
				},
			},
		},
	},
	{
		Path: "/:bar",
		Children: []Route{
			{
				Path:   "/bar",
				Method: http.MethodGet,
				Handler: func(w http.ResponseWriter, r Request) {
					w.Header().Set("X-Params", fmt.Sprintf("%s", r.Params()))
					w.Write([]byte(r.Route()))
				},
			},
		},
	},
}

var tRouter = New(tRoutes)

type ExpectedResponse struct {
	body       string
	statusCode int
	expect     func(http.ResponseWriter)
}

type RouteTest struct {
	req *http.Request
	res ExpectedResponse
}

func testRoute(t *testing.T, rt RouteTest) {
	w := httptest.NewRecorder()
	tRouter.ServeHTTP(w, rt.req)

	if code, want := w.Result().StatusCode, rt.res.statusCode; code != want {
		t.Errorf("got %d, wanted %d", code, want)
	}

	if text, want := w.Body.String(), rt.res.body; text != want {
		t.Errorf("got %s, wanted %s", text, want)
	}

	if expect := rt.res.expect; expect != nil {
		expect(w)
	}
}

func TestRoute(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/ping", nil),
		ExpectedResponse{"pong", 200, nil},
	})
}

func TestNestedRoute(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodPost, "/api/hello", bytes.NewBufferString("{\"name\": \"John Doe\"}")),
		ExpectedResponse{"John Doe", 200, nil},
	})
}

func TestDynamicRoutes(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/users/1", "/users/:userId"},
		{"/users/1/blogs/2", "/users/:userId/blogs/:blogId"},
	}

	for _, test := range tests {
		t.Run("Path="+test.path, func(t *testing.T) {
			testRoute(t, RouteTest{
				httptest.NewRequest(http.MethodGet, test.path, nil),
				ExpectedResponse{test.want, 200, nil},
			})
		})
	}
}

func TestNestedDynamicRoute(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/posts/1/comments/2", nil),
		ExpectedResponse{"/posts/:postId/comments/:commentId", 200, nil},
	})
}

func TestRoutePath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/foo/foo", "/:foo/foo"},
		{"/bar/bar", "/:bar/bar"},
	}

	for _, test := range tests {
		t.Run("Path="+test.path, func(t *testing.T) {
			testRoute(t, RouteTest{
				httptest.NewRequest(http.MethodGet, test.path, nil),
				ExpectedResponse{test.want, 200, nil},
			})
		})
	}
}

func TestRequestParams(t *testing.T) {
	tests := []struct {
		path   string
		body   string
		params string
	}{
		{"/users/1", "/users/:userId", "map[userId:1]"},
		{"/posts/1/comments/2", "/posts/:postId/comments/:commentId", "map[commentId:2 postId:1]"},
	}

	for _, test := range tests {
		t.Run("Path="+test.path, func(t *testing.T) {
			testRoute(t, RouteTest{
				httptest.NewRequest(http.MethodGet, test.path, nil),
				ExpectedResponse{test.body, 200, func(w http.ResponseWriter) {
					if result := w.Header().Get("X-Params"); result != test.params {
						t.Errorf("Expected params: %s, got: %s", test.params, result)
					}
				}},
			})
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/api/hello", nil),
		ExpectedResponse{"", 405, nil},
	})
}

func TestNotFound(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodGet, "/hello", nil),
		ExpectedResponse{"404 page not found\n", 404, nil},
	})
}

func TestMiddleware(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodPost, "/api/hello", bytes.NewBufferString("{\"name\": \"John Doe\"}")),
		ExpectedResponse{"John Doe", 200, func(w http.ResponseWriter) {
			if result, want := w.Header().Get("Content-Type"), "application/json"; result != want {
				t.Errorf("Expected content-type: %s, got: %s", want, result)
			}
		}},
	})
}

func TestNestedMiddleware(t *testing.T) {
	testRoute(t, RouteTest{
		httptest.NewRequest(http.MethodPost, "/api/hello", bytes.NewBufferString("{\"name\": \"John Doe\"}")),
		ExpectedResponse{"John Doe", 200, func(w http.ResponseWriter) {
			if result, want := w.Header().Get("Content-Type"), "application/json"; result != want {
				t.Errorf("Expected content-type: %s, got: %s", want, result)
			}

			if result, want := w.Header().Get("X-Test"), "test"; result != want {
				t.Errorf("Expected 'X-Test' header: %s, got: %s", want, result)
			}
		}},
	})
}
