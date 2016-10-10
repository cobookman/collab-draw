package main

import (
	"net/http"
	"time"
	"strconv"
)

// Creates a new drawing canvas.
func CreateCanvas(r *http.Request) (interface{}, error) {
	canvas := Canvas{}
	err := canvas.Create()
	return canvas, err
}

// Gets a canvas by specified Id.
func GetCanvas(r *http.Request) (interface{}, error) {
	id := r.FormValue("id")
	canvas := Canvas{}
	err := canvas.Get(id)
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
	err = canvases.GetAll(activeSince, limit)
	return canvases, err
}
