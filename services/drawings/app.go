package main

import (
	"cloud.google.com/go/pubsub"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	upgrader = websocket.Upgrader{}
	gopath   = os.Getenv("GOPATH")
	mq       *MessagingQueue
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

	r.Path("/canvas").
		HandlerFunc(WebsocketMiddleware(UserCanvasSocket))

	s.Handle("/", r)

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

		err := http.ListenAndServe("0.0.0.0:8080", s)
		log.Print("Failed to serve 8080", err)
	}()

	// Create blocking thread for http server to bypass appengine proxy
	go func() {
		defer wg.Done()

		err := http.ListenAndServe("0.0.0.0:65080", s)
		log.Print("Failed to serve 65080", err)
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

type WebsocketApi func(r *http.Request, c *websocket.Conn, mq *MessagingQueue)
type IncomingDrawingHandler func(drawing Drawing) error

func WebsocketMiddleware(f WebsocketApi) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// upgrader.Upgrade auto writes http errors on failure
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer ws.Close()
		f(r, ws, mq)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func DrawingMessageMiddleware(handler IncomingDrawingHandler) IncomingMessageHandler {
	return func(msg *pubsub.Message) error {
		drawing := Drawing{}
		if err := drawing.Unmarshal(msg.Data); err != nil {
			return err
		}

		return handler(drawing)
	}
}
