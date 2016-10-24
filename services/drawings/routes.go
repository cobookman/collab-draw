package main

import (
	"errors"
	"sync"
	"reflect"
	"log"
	"net/http"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/gorilla/websocket"
)

var (
	userSockets = make(map[string][]*websocket.Conn)
	lock = sync.RWMutex{}
)


type SocketErr struct {
	Error error
	Message string
}

type SocketMsg struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

func SocketErrf(c *websocket.Conn, err error, msg string) {
	se := SocketErr{
		Error: err,
		Message: msg,
	}
	b, _ := json.Marshal(se)
	c.WriteMessage(websocket.TextMessage, b)
}

func SocketRespf(c *websocket.Conn, v interface{}) {
	sm := SocketMsg{
		Type: reflect.TypeOf(v).Name(),
		Data: v,
	}
	b, _ := json.Marshal(sm)
	c.WriteMessage(websocket.TextMessage, b)
}

// adds a websocket conn
func addSub(canvasID string, c *websocket.Conn) {
	lock.Lock()
	defer lock.Unlock()

	socs := userSockets[canvasID]
	userSockets[canvasID] = append(socs, c)
}

// removes the specified websocket conn, and gives a count of how many
// remaining subs are on the server. returns -1 if element not removed
func removeSub(canvasID string, c *websocket.Conn) (int, error) {
	lock.Lock()
	defer lock.Unlock()

	cs := userSockets[canvasID]
	for i := 0; i < len(cs); i++ {
		if cs[i] == c {
			cs = append(cs[:i], cs[i+1:]...)
			userSockets[canvasID] = cs
			return len(cs), nil
		}
	}
	return -1, errors.New("Failed to find the websocket obj subscribed to given canvas")
}

// Locks the usersocket map, and returns a copy of websocket connections
func getSubs(canvasID string) []*websocket.Conn {
	lock.RLock()
	defer lock.RUnlock()
	cs := userSockets[canvasID]
	csCopy := make([]*websocket.Conn, len(cs))
	copy(csCopy, cs)
	return csCopy
}

// Cleanup a subscriptions on websocket close
func cleanup(canvasID string, c *websocket.Conn, mq *MessagingQueue) {
	ctx := context.Background()

	remainingSubs, err := removeSub(canvasID, c)
	if err != nil {
		log.Print(err)
	}

	if remainingSubs == 0 {
		RemoveCanvasSubscription(ctx, canvasID, mq.Topic.ID())
	}
}

func UserCanvasSocket(r *http.Request, c *websocket.Conn, mq *MessagingQueue) {
	ctx := context.Background()
	canvasID := r.FormValue("id")

	// add a subscription
	addSub(canvasID, c)
	if err := AddCanvasSubscription(ctx, canvasID, mq.Topic.ID()); err != nil {
		log.Print(err)
		SocketErrf(c, err, "Failed to add subscription to canvas updates")
		return
	}

	// when done, cleanup
	defer cleanup(canvasID, c, mq)

	// get canvas & send to end user
	canvas, err := GetCanvas(ctx, canvasID)
	if err != nil {
		log.Print(err)
		SocketErrf(c, err, "Failed to get specified canvas")
		return
	}
	SocketRespf(c, canvas)

	// Deal with incoming drawings
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
		// TODO(bookman): Route incoming drawing 
	}
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
		SocketRespf(subs[i], drawing)
	}

	return nil
}
