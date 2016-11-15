package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	gopath = os.Getenv("GOPATH")
)

func main() {
	s := http.NewServeMux()
	r := mux.NewRouter()

	ae := r.PathPrefix("/_ah").Subrouter()
	api := r.PathPrefix("/api").Subrouter()
	apiV1 := api.PathPrefix("/v1").Subrouter()

	apiV1.Path("/canvas").Methods("POST").
		HandlerFunc(RestfulMiddleware(HandleCreateCanvas))

	apiV1.Path("/canvas").Methods("GET").
		HandlerFunc(RestfulMiddleware(HandleGetCanvas))

	apiV1.Path("/canvases").Methods("GET").
		HandlerFunc(RestfulMiddleware(HandleListCanvases))

	// Required for appengine flex to measure the health of service
	ae.Path("/health").Methods("GET", "POST").
		HandlerFunc(healthCheck)

	// Serve static assets. Note that gorilla matches in order of route order
	// So this should be the last route added to our http server
	static_assets_path := gopath + "/src/" + "github.com/cobookman/collabdraw/services/default/public"
	log.Print("Static assets serving from: ", static_assets_path)
	// r.PathPrefix("/*.html").Handler(http.FileServer(http.Dir(static_assets_path)))
	r.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exts := []string{"html", "js", "css", "png", "json"}
		fext := filepath.Ext(r.URL.Path)
		log.Print(r.URL.Path, fext)
		for _, ext := range exts {
			if fext == "."+ext {
				http.ServeFile(w, r, static_assets_path+"/"+r.URL.Path)
				return
			}
		}

		http.ServeFile(w, r, static_assets_path+"/index.html")
	}))

	s.Handle("/", r)
	err := http.ListenAndServe("0.0.0.0:8080", handlers.CORS()(s))
	log.Print("Failed to serve 8080", err)
}

type RestfulApi func(r *http.Request) (interface{}, error)

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

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
