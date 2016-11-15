package main

import (
	"errors"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
)

// Creates a new drawing canvas.
func HandleCreateCanvas(r *http.Request) (interface{}, error) {
	ctx := context.Background()
	canvas, err := NewCanvas(ctx)
	return canvas, err
}

// Gets a canvas by specified Id.
func HandleGetCanvas(r *http.Request) (interface{}, error) {
	id := r.FormValue("id")
	if len(id) == 0 {
		return nil, errors.New("Please specify a canvas Id")
	}
	ctx := context.Background()
	canvas, err := GetCanvas(ctx, id)
	return canvas, err
}

// Lists active canvases by active since rfc3339 timestamp and up to the given limit.
func HandleListCanvases(r *http.Request) (interface{}, error) {
	activeSince, err := time.Parse(time.RFC3339, r.FormValue("activeSince"))
	if err != nil {
		return nil, err
	}
	limit := 0
	if len(r.FormValue("limit")) != 0 {
		var err error
		limit, err = strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			return nil, err
		}
	}

	canvases := Canvases{}

	ctx := context.Background()
	err = canvases.GetAll(ctx, activeSince, limit)
	return canvases, err
}

