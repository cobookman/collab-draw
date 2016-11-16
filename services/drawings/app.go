package main

import (
	"encoding/json"
	"cloud.google.com/go/pubsub"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var (
	gopath = os.Getenv("GOPATH")
	mq     *MessagingQueue
)

func main() {
	log.Print("Starting up messaging queue")
	ctx := context.Background()
	projectID := os.Getenv("GCLOUD_PROJECT")
	topicName := "topic-" + uuid.NewV4().String()
	subName := "sub-" + uuid.NewV4().String()
	var err error
	mq, err = NewMessagingQueue(ctx, projectID, topicName, subName)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Starting up server")
	s := http.NewServeMux()

	r := mux.NewRouter()

	ae := r.PathPrefix("/_ah").Subrouter()

	// Required for appengine flex to measure the health of service
	ae.Path("/health").Methods("GET", "POST").
		HandlerFunc(healthCheck)

	r.Path("/ipaddr").
		HandlerFunc(hostIpHandler)

	r.Path("/canvas").
		HandlerFunc(ServerPushMiddleware(UserCanvasDrawingPush))

	r.Path("/drawing").Methods("POST").
		HandlerFunc(RestfulMiddleware(HandleIncomingDrawing))

	s.Handle("/", r)

	// cleanup on os kill signal
	go func() {
		c := make(chan os.Signal, 10)
		signal.Notify(c, os.Interrupt)
		// block until we get a kill signal
		<-c
		if err := mq.Cleanup(ctx); err != nil {
			log.Fatal(err)
		} else {
			log.Print("Cleaned up before exit")
		}
		os.Exit(0)
	}()

	// Start servers. If any server crashes the wait group is decrimented,
	// and the main server will quit, causing our docker manager to restart the
	// docker instance.
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Start a server that handles incoming appengine proxy calls
	// Used for health checks, and to tell us when we're about to
	// be drained
	go func() {
		defer wg.Done()

		err := http.ListenAndServe("0.0.0.0:8080", handlers.CORS()(s))
		log.Print("Failed to serve 8080", err)
	}()

	// Create thread for messaging queue worker
	go func() {
		defer wg.Done()

		// Start a single worker
		if err := mq.OnMessage(DrawingMessageMiddleware(OnIncomingDrawing)); err != nil {
			log.Print("Worker failed", err)
		}
	}()
	log.Print("Started up all servers")

	// Will block until wg.Done() is called
	wg.Wait()

	// Some service died, clean ourselves up then kill
	if err := mq.Cleanup(ctx); err != nil {
		log.Fatal(err)
	}
}

type ServerPushHandler func(w http.ResponseWriter, r *http.Request, mq *MessagingQueue)
func ServerPushMiddleware(f ServerPushHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		// Set Server push headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// https://www.nginx.com/resources/wiki/start/topics/examples/x-accel/#x-accel-buffering
		w.Header().Set("X-Accel-Buffering", "no")
		f(w, r, mq)
	}
}

type IncomingDrawingHandler func(drawing Drawing) error
func DrawingMessageMiddleware(handler IncomingDrawingHandler) IncomingMessageHandler {
	return func(msg *pubsub.Message) error {
		drawing := Drawing{}
		if err := drawing.Unmarshal(msg.Data); err != nil {
			return err
		}

		return handler(drawing)
	}
}

type RestfulApi func(r *http.Request) (interface{}, error)
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

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(out)
        }
}


func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func hostIpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"wsip\":\"%s\"}", HostIp()+":65080")
}

