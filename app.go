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

	apiV1.Path("/devices").Methods("GET").
		HandlerFunc(RestfulMiddleware(deviceListHandler))

	http.Handle("/", r)
	appengine.Main()
}

type RestfulApi func(r *http.Request) (interface {}, error)

func Jsonify(v interface{}, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func RestfulMiddleware(f RestfulApi) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := f(r)
		var status int
		var out interface{}
		if err != nil {
			status = http.StatusInternalServerError
			out = err
		} else {
			status = http.StatusOK
			out = v
		}
		Jsonify(out, w, status)
	}
}

type Device struct {
	Name string
}

type DeviceList struct {
	Devices []Device
}

func deviceListHandler(r *http.Request) (interface{}, error) {
	devices := DeviceList{
		Devices: []Device{
			Device{
				Name: "Colin",
			},
		},
	}
	return devices, nil
}
