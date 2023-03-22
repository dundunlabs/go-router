package gorouter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testResponse(method string, target string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, body)
	router := New([]Route{
		{

			Path:   "/401",
			Method: http.MethodGet,
			Handler: func(r *Request, w *Response) {
				w.Status(http.StatusUnauthorized)
			},
		},
		{
			Path:   "/string",
			Method: http.MethodGet,
			Handler: func(r *Request, w *Response) {
				w.SendString("Hello World!")
			},
		},
		{
			Path:   "/json",
			Method: http.MethodGet,
			Handler: func(r *Request, w *Response) {
				w.MustSendJSON(map[string]string{"foo": "bar"})
			},
		},
	})
	router.ServeHTTP(w, req)
	return w
}

func TestStatus(t *testing.T) {
	w := testResponse(http.MethodGet, "/401", nil)
	if want, got := http.StatusUnauthorized, w.Result().StatusCode; want != got {
		t.Errorf("got %d, wanted %d", got, want)
	}
}

func TestSendString(t *testing.T) {
	w := testResponse(http.MethodGet, "/string", nil)
	if want, got := "text/plain; charset=utf-8", w.Header().Get("Content-Type"); want != got {
		t.Errorf("got %s, wanted %s", got, want)
	}
	if want, got := "Hello World!", w.Body.String(); want != got {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestSendJSON(t *testing.T) {
	w := testResponse(http.MethodGet, "/json", nil)
	if want, got := "application/json", w.Header().Get("Content-Type"); want != got {
		t.Errorf("got %s, wanted %s", got, want)
	}
	if want, got := "{\"foo\":\"bar\"}", w.Body.String(); want != got {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
