package main

import (
	"time"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"reflect"
	"sync"
)

var (
	// Using a pointer here as it allows us to keep track of who's who
	// through object pointer addresses
	userSockets = make(map[string][]*http.ResponseWriter)
	lock        = sync.RWMutex{}
)

type SocketErr struct {
	Error   error
	Message string
}

type SocketMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Status struct {
	Message string `json:"msg"`
}

func writeSSE(w http.ResponseWriter, eventType string, data string) {
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
	f, ok := w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func PushErr(w http.ResponseWriter, err error, msg string) {
	se := SocketErr{
		Error:   err,
		Message: msg,
	}
	b, _ := json.Marshal(se)
	writeSSE(w, "error", string(b))
}

func PushResp(w http.ResponseWriter, v interface{}) {
	sm := SocketMsg{
		Type: reflect.TypeOf(v).Name(),
		Data: v,
	}
	b, _ := json.Marshal(sm)
	writeSSE(w, "message", string(b))
}

// adds a websocket conn
func addSub(canvasID string, w *http.ResponseWriter) {
	lock.Lock()
	defer lock.Unlock()

	socs := userSockets[canvasID]
	userSockets[canvasID] = append(socs, w)
}

// removes the specified websocket conn, and gives a count of how many
// remaining subs are on the server. returns -1 if element not removed
func removeSub(canvasID string, w *http.ResponseWriter) (int, error) {
	lock.Lock()
	defer lock.Unlock()

	cs := userSockets[canvasID]
	for i := 0; i < len(cs); i++ {
		if cs[i] == w {
			cs = append(cs[:i], cs[i+1:]...)
			userSockets[canvasID] = cs
			return len(cs), nil
		}
	}
	return -1, errors.New("Failed to find the websocket obj subscribed to given canvas")
}

// Locks the usersocket map, and returns a copy of websocket connections
func getSubs(canvasID string) []*http.ResponseWriter {
	lock.RLock()
	defer lock.RUnlock()
	cs := userSockets[canvasID]
	csCopy := make([]*http.ResponseWriter, len(cs))
	copy(csCopy, cs)
	return csCopy
}

// Cleanup a subscriptions on websocket close
func cleanup(canvasID string, w *http.ResponseWriter, mq *MessagingQueue) {
	ctx := context.Background()

	remainingSubs, err := removeSub(canvasID, w)
	if err != nil {
		log.Print(err)
	}

	if remainingSubs == 0 {
		RemoveCanvasSubscription(ctx, canvasID, mq.Topic.ID())
	}
}

func UserCanvasDrawingPush(w http.ResponseWriter, r *http.Request, mq *MessagingQueue) {
	ctx := context.Background()
	canvasID := r.FormValue("canvasId")
	if len(canvasID) == 0 {
		log.Print("missing canvasId query parameter.")
		PushErr(w, nil, "missing canvasId query parameter.")
		return
	}
	// add a subscription
	addSub(canvasID, &w)
	if err := AddCanvasSubscription(ctx, canvasID, mq.Topic.ID()); err != nil {
		log.Print(err)
		PushErr(w, err, "Failed to add subscription to canvas updates")
		return
	}

	// when done, cleanup ourcanvas subscription
	// we just created
	defer cleanup(canvasID, &w, mq)

	// get canvas & send to end user
	canvas, err := GetCanvas(ctx, canvasID)
	if err != nil {
		log.Print(err)
		PushErr(w, err, "Failed to get specified canvas")
		return
	}
	PushResp(w, canvas)

	// get canvas's drawings & send to end user
	drawings, err := GetDrawings(ctx, canvas)
	if err != nil {
		log.Print(err)
		PushErr(w, err, "Failed to get previous drawings for canvas")
		return
	}

	// Send up the drawings one at a time
	for _, d := range drawings {
		PushResp(w, d)
	}

		PushErr(w, nil, "HI WORLD")

	// Keep informing our client that we are still alive
	closenotify := w.(http.CloseNotifier).CloseNotify()
	for {
		select {
		case <-closenotify:
			return
		default:
			PushResp(w, Status{
				Message: "Still Alive",
			})
		}
		time.Sleep(1 * time.Second)
	}
}

func HandleIncomingDrawing(r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	d := Drawing{}
	if err := d.Unmarshal(body); err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := d.Forward(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

// Called when we get a pubsub message that needs to be sent to some of our
// hosted users
func OnIncomingDrawing(drawing Drawing) error {
	// get users that are on the canvas the drawing belongs to.
	subs := getSubs(drawing.CanvasID)
	if len(subs) == 0 {
		return errors.New("No subscribers of canvasID: " + drawing.CanvasID)
	}

	for i := 0; i < len(subs); i++ {
		PushResp(*subs[i], drawing)
	}

	return nil
}
