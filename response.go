package gorouter

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	statusCode int
}

func (r *Response) Status(statusCode int) *Response {
	r.statusCode = statusCode
	return r
}

func (r *Response) End() {
	if r.statusCode > 0 {
		r.WriteHeader(r.statusCode)
	}
}

func (r *Response) SendString(s string) error {
	r.End()
	_, err := r.Write([]byte(s))
	return err
}

func (r *Response) MustSendString(s string) {
	if err := r.SendString(s); err != nil {
		panic(err)
	}
}

func (r *Response) SendJSON(v any) error {
	r.Header().Set("Content-Type", "application/json")
	r.End()
	return json.NewEncoder(r).Encode(v)
}

func (r *Response) MustSendJSON(v any) {
	if err := r.SendJSON(v); err != nil {
		panic(err)
	}
}
