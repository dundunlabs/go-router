package main

import (
	"errors"
	"fmt"
	"net/http"

	gorouter "github.com/dundunlabs/go-router"
)

type ApiResponse struct {
	Error any `json:"error"`
	Data  any `json:"data"`
}

// error handling and response formatting
func apiHandler(next gorouter.HandlerFunc) gorouter.HandlerFunc {
	return func(req *gorouter.Request, res *gorouter.Response) {
		// error handling
		defer func() {
			if err := recover(); err != nil {
				res.Status(http.StatusInternalServerError).SendJSON(ApiResponse{
					Error: fmt.Sprint(err),
				})
			}
		}()

		next(req, res)

		// response formatting
		data := res.Body()
		res.SendJSON(ApiResponse{Data: data})
	}
}

var routes = []gorouter.Route{
	{
		Path:       "/api",
		Middleware: apiHandler,
		Children: []gorouter.Route{
			{
				Path:   "/hello",
				Method: http.MethodGet,
				Handler: func(req *gorouter.Request, res *gorouter.Response) {
					res.SendString("Hello World!")
				},
			},
			{
				Path:   "/error",
				Method: http.MethodGet,
				Handler: func(req *gorouter.Request, res *gorouter.Response) {
					panic(errors.New("an error occurred!"))
				},
			},
			{
				Path:   "/users",
				Method: http.MethodPost,
				Handler: func(req *gorouter.Request, res *gorouter.Response) {
					var m map[string]any
					req.MustParseBody().MustBind(&m)
					res.MustSendJSON(m)
				},
			},
			{
				Path:   "/users/:id",
				Method: http.MethodGet,
				Handler: func(req *gorouter.Request, res *gorouter.Response) {
					params := req.Params()
					res.MustSendJSON(params)
				},
			},
		},
	},
	{
		Path:   "/*",
		Method: http.MethodGet,
		Handler: func(req *gorouter.Request, res *gorouter.Response) {
			res.SendString("Hello from \"" + req.Route() + "\" route")
		},
	},
}

func main() {
	router := gorouter.New(routes)
	http.ListenAndServe(":8080", router)
}
