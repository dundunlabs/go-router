package gorouter

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(req *http.Request, test func(*Request)) {
	router := New([]Route{
		{

			Path:   "/test/:a/:b/:c",
			Method: http.MethodGet,
			Handler: func(w *Response, r *Request) {
				test(r)
			},
		},
		{

			Path:   "/",
			Method: http.MethodPost,
			Handler: func(w *Response, r *Request) {
				test(r)
			},
		},
	})
	router.ServeHTTP(httptest.NewRecorder(), req)
}

func TestParams(t *testing.T) {
	testRequest(httptest.NewRequest(http.MethodGet, "/test/1/2/3", nil), func(r *Request) {
		params := r.Params()
		want := map[string]string{
			"a": "1",
			"b": "2",
			"c": "3",
		}
		for k, v := range params {
			if want[k] != v {
				t.Errorf("got %s, wanted %s", v, want[k])
			}
		}
	})
}

func TestParam(t *testing.T) {
	testRequest(httptest.NewRequest(http.MethodGet, "/test/1/2/3", nil), func(r *Request) {
		want := map[string]string{
			"a": "1",
			"b": "2",
			"c": "3",
		}
		for k, v := range want {
			if r.Param(k) != v {
				t.Errorf("got %s, wanted %s", v, want[k])
			}
		}
	})
}

func TestBody(t *testing.T) {
	body := bytes.NewBufferString("{\"id\":1,\"profile\":{\"name\":\"Foo Bar\"}}")
	want := map[string]any{
		"id": 1,
		"profile": map[string]any{
			"name": "Foo Bar",
		},
	}
	testRequest(httptest.NewRequest(http.MethodPost, "/", body), func(r *Request) {
		type Profile struct {
			Name string `json:"name"`
		}

		type Body struct {
			ID      int     `json:"id"`
			Profile Profile `json:"profile"`
		}

		var body Body

		r.MustParseBody().Bind(&body)

		if w, g := want["id"], body.ID; w != g {
			t.Errorf("got %d, wanted %d", g, w)
		}
		if w, g := want["profile"].(map[string]any)["name"], body.Profile.Name; w != g {
			t.Errorf("got %s, wanted %s", g, w)
		}
	})
}
