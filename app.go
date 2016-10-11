package main

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"sync"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{}
	gopath = os.Getenv("GOPATH")
)

func main() {
	log.Print("Starting up server")
	sHttp := http.NewServeMux()
	sWs := http.NewServeMux()

	r := mux.NewRouter()
	ae := r.PathPrefix("/_ah").Subrouter()
	api := r.PathPrefix("/api").Subrouter()
	apiV1 := api.PathPrefix("/v1").Subrouter()

	apiV1.Path("/canvas").Methods("POST").
		HandlerFunc(RestfulMiddleware(CreateCanvas))

	apiV1.Path("/canvas").Methods("GET").
		HandlerFunc(RestfulMiddleware(GetCanvas))

	apiV1.Path("/canvases").Methods("GET").
		HandlerFunc(RestfulMiddleware(ListCanvases))

	// Required for appengine flex to measure the health of service
	ae.Path("/health").Methods("GET", "POST").
		HandlerFunc(healthCheck)

	// Serve static assets. Note that gorilla matches in order of route order
	// So this should be the last route added
	static_assets_path := gopath + "/src/" + "github.com/cobookman/collabdraw/public"
	log.Print("Static assets serving from: ", static_assets_path)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(static_assets_path)))
	sHttp.Handle("/", r)

	// Add websocket routes
	sWs.Handle("/canvas", WebsocketMiddleware(ListenCanvas))

	// Serve both http servers
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := http.ListenAndServe("0.0.0.0:8080", sHttp)
		log.Print("Failed to serve 8080", err)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err := http.ListenAndServe("0.0.0.0:65080", sWs)
		log.Print("Failed to serve 65080", err)
		wg.Done()
	}()

	log.Print("Started up both servers")
	wg.Wait()
	//appengine.Main()
}

type RestfulApi func(r *http.Request) (interface{}, error)
type WebsocketApi func(r *http.Request, c *websocket.Conn)

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

func WebsocketMiddleware(f WebsocketApi) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Jsonify(ErrorResp{Error: err.Error()}, w, http.StatusInternalServerError)
			return
		}

		defer ws.Close()
		f(r, ws)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
