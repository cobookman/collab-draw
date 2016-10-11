package main

import (
	"log"
	"net/http"
	"time"
	"strconv"

	"golang.org/x/net/context"
	"github.com/gorilla/websocket"
)

// Creates a new drawing canvas.
func CreateCanvas(r *http.Request) (interface{}, error) {
	canvas := Canvas{}

	ctx := context.Background()
	err := canvas.Create(ctx)
	return canvas, err
}

// Gets a canvas by specified Id.
func GetCanvas(r *http.Request) (interface{}, error) {
	id := r.FormValue("id")
	canvas := Canvas{}

	ctx := context.Background()
	err := canvas.Get(ctx, id)
	return canvas, err
}

// Lists active canvases by active since rfc3339 timestamp and up to the given limit.
func ListCanvases(r *http.Request) (interface{}, error) {
	activeSince, err := time.Parse(time.RFC3339, r.FormValue("activeSince"))
	if err != nil {
		return nil, err
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		return nil, err
	}

	canvases := Canvases{}

	ctx := context.Background()
	err = canvases.GetAll(ctx, activeSince, limit)
	return canvases, err
}

func ListenCanvas(r *http.Request, c *websocket.Conn) {
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}

