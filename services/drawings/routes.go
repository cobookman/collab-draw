package main

import (
	"errors"
	"github.com/cobookman/collabdraw/shared/models"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	conns = NewSSEConns()
)

// Cleanup a subscriptions on websocket close
func cleanup(canvasID string, c *SSEConn) {
	ctx := context.Background()

	remaining, err := conns.Remove(canvasID, c)
	if err != nil {
		log.Print(err)
	}

	if remaining == 0 {
		sub := models.NewSubscription(canvasID)
		sub.Unsubscribe(ctx)
	}
}

func UserCanvasDrawingPush(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn := NewSSEConn(&w)

	canvasID := r.FormValue("canvasId")
	if len(canvasID) == 0 {
		log.Print("missing canvasId query parameter.")
		conn.WriteErr(nil, "missing canvasId query parameter.")
		return
	}

	conns.Add(canvasID, conn)
	defer cleanup(canvasID, conn)

	// subscribe to drawing updates
	sub := models.NewSubscription(canvasID)
	if err := sub.Subscribe(ctx); err != nil {
		log.Print(err)
		conn.WriteErr(err, "Failed to add subscription to canvas updates")
		return
	}

	// get canvas & send to end user
	canvas, err := models.GetCanvas(ctx, canvasID)
	if err != nil {
		log.Print(err)
		conn.WriteErr(err, "Failed to get specified canvas")
		return
	}
	conn.WriteMsg(canvas)

	// get canvas's drawings & send to end user
	drawings, err := models.GetDrawings(ctx, canvas)
	if err != nil {
		log.Print(err)
		conn.WriteErr(err, "Failed to get previous drawings for canvas")
		return
	}

	// Send up the drawings one at a time
	for _, d := range drawings {
		conn.WriteMsg(d)
	}

	// Keep informing our client that we are still alive
	// And handle socket closure
	closenotify := w.(http.CloseNotifier).CloseNotify()
	for {
		select {
		case <-closenotify:
			return
		default:
			conn.WriteAlive()
		}
		time.Sleep(5 * time.Second)
	}
}

func HandleIncomingDrawing(r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	d := models.Drawing{}
	if err := d.Unmarshal(body); err != nil {
		return nil, err
	}

	// Give the drawing a universal unique id
	d.AssignID()

	// Forward drawing to microservices
	ctx := context.Background()
	if err := d.Forward(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

// Route drawing to relavent clients.
func OnIncomingDrawing(drawing models.Drawing) error {
	writers := conns.GetWriters(drawing.CanvasID)
	if len(writers) == 0 {
		return errors.New("No subscribers of canvasID: " + drawing.CanvasID)
	}

	for i := 0; i < len(writers); i++ {
		writer := writers[i]
		writer.WriteMsg(drawing)
	}

	return nil
}
