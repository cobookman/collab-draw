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
	"time"
)

var (
	upgrader = websocket.Upgrader{}
	gopath   = os.Getenv("GOPATH")
)

func main() {
	log.Print("Starting up server")
	s := http.NewServeMux()

	r := mux.NewRouter()

	ae := r.PathPrefix("/_ah").Subrouter()

	// Required for appengine flex to measure the health of service
	ae.Path("/health").Methods("GET", "POST").
		HandlerFunc(healthCheck)

	r.Path("/canvas").
		HandlerFunc(WebsocketMiddleware(ListenCanvas))

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
	go func() {
		defer wg.Done()

		// avoids appengine proxy
		err := http.ListenAndServe("0.0.0.0:65080", s)
		log.Print("Failed to serve 65080", err)
	}()
	go func() {
		defer wg.Done()
		sub, err := subscribeIncomingDrawings()
		if err != nil {
			log.Print("Failed to subscribe to incoming drawing pubsub topic: ", err)
			return
		}

		// Start a single worker
		if err := listenIncomingDrawings(sub, OnIncomingDrawing); err != nil {
			log.Print("Worker failed", err)
		}
	}()

	log.Print("Started up all servers")

	// Will block until wg.Done() is called
	wg.Wait()
}

type WebsocketApi func(r *http.Request, c *websocket.Conn)
type IncomingDrawingHandler func(drawing Drawing) error

func WebsocketMiddleware(f WebsocketApi) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// upgrader.Upgrade auto writes http errors on failure
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer ws.Close()
		f(r, ws)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func subscribeIncomingDrawings() (*pubsub.Subscription, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
	if err != nil {
		return nil, err
	}

	// Create a new unique topic name for this server instance
	topicName := "topic-" + uuid.NewV4().String()

	// Create topic
	topic, _ := client.CreateTopic(ctx, topicName)

	// Add a subscription for our host, using a random unique uuid.
	// We'll give ourselves 10 seconds to respond to a message before its
	// deemed as a failure
	subscriptionName := "subscription-" + uuid.NewV4().String()
	return client.CreateSubscription(ctx, subscriptionName, topic, 10*time.Second, nil)
}

func listenIncomingDrawings(subscription *pubsub.Subscription, handler IncomingDrawingHandler) error {
	ctx := context.Background()
	it, err := subscription.Pull(ctx)
	if err != nil {
		return err
	}
	defer it.Stop()

	for {
		msg, err := it.Next()
		if err != nil {
			return err
		}
		drawing := Drawing{}
		if err := drawing.Unmarshal(msg.Data); err != nil {
			log.Print("Failed to parse data: ", msg.Data)
		}
		if err := handler(drawing); err != nil {
			msg.Done(false)
		} else {
			msg.Done(true)
		}
	}
}
