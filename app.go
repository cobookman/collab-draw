package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"google.golang.org/appengine"

)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)


func main() {
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

	apiV1.Path("/canvas/ws").Methods("GET").
		HandlerFunc(WebsocketMiddleware(ListenCanvas))

	// Required for appengine flex to measure the health of service
	ae.Path("/health").Methods("GET", "POST").
		HandlerFunc(healthCheck)

	// Serve static assets. Note that gorilla matches in order of route order
	// So this should be the last route added
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	http.Handle("/", r)
	appengine.Main()
}

type RestfulApi func(r *http.Request) (interface {}, error)
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
	return func (w http.ResponseWriter, r *http.Request) {
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
