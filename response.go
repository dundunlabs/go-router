package gorouter

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	statusCode int
	b          []byte
	body       any
}

func (r *Response) Status(statusCode int) *Response {
	r.statusCode = statusCode
	return r
}

func (r *Response) Body() any {
	return r.body
}

func (r *Response) SendString(s string) {
	r.body = s
	r.b = []byte(s)
}

func (r *Response) SendJSON(v any) error {
	r.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.body = v
	r.b = b
	return nil
}

func (r *Response) MustSendJSON(v any) {
	if err := r.SendJSON(v); err != nil {
		panic(err)
	}
}
