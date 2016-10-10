package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func main() {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiV1 := apiRouter.PathPrefix("/v1").Subrouter()

	apiV1.Path("/canvas").Methods("POST").
		HandlerFunc(RestfulMiddleware(CreateCanvas))

	apiV1.Path("/canvas").Methods("GET").
		HandlerFunc(RestfulMiddleware(GetCanvas))

	apiV1.Path("/canvases").Methods("GET").
		HandlerFunc(RestfulMiddleware(ListCanvases))

	http.Handle("/", r)
	appengine.Main()
}

type RestfulApi func(r *http.Request) (interface {}, error)

func Jsonify(v interface{}, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type ErrorResp struct {
	Error string `json:"error"`
}

func RestfulMiddleware(f RestfulApi) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var status int
		var out interface{}

		v, err := f(r)
		if err != nil {
			status = http.StatusInternalServerError
			out = ErrorResp{
				Error: err.Error(),
			}
		} else {
			status = http.StatusOK
			out = v
		}

		Jsonify(out, w, status)
	}
}

