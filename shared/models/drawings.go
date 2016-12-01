package models

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

var (
	DrawingKind = "Drawing"
)

type Point struct {
	X int `json:"x" datastore:"x"`
	Y int `json:"y" datastore:"y"`
}

type Drawing struct {
	// a uuid for the given drawing
	ID string `json:"id" datastore:"id"`

	// Where this drawing is painted on
	CanvasID string `json:"canvasId" datastore:"canvasId"`

	// The drawing itself is a set of points connected by lines
	Points []Point `json:"points" datastore:"points"`
}

// Encodes a drawing into json representation.
func (d Drawing) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

// Parses a json representation of a drawing into a drawing struct.
func (d *Drawing) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, d)
}

// Forwards the drawing to a pubsub topic
func (d *Drawing) Forward(ctx context.Context) error {
	if len(d.ID) == 0 {
		d.AssignID()
	}

	data, err := d.Marshal()
	if err != nil {
		return err
	}

	msg := pubsub.Message{
		Data: data,
	}

	if _, err := UpstreamTopic().Publish(ctx, &msg); err != nil {
		return err
	}

	return nil
}

// Assigns the drawing a unique id
func (d *Drawing) AssignID() {
	d.ID = uuid.NewV4().String()
}

// Gets a collection of drawings.
func GetDrawings(ctx context.Context, canvas *Canvas) ([]Drawing, error) {
	if canvas == nil {
		return nil, errors.New("canvas cannot be nil")
	}

	q := datastore.NewQuery(DrawingKind).
		Filter("canvasId =", canvas.ID)
	drawings := make([]Drawing, 0, 100)
	keys, err := DatastoreClient().GetAll(ctx, q, &drawings)
	if err != nil {
		return nil, err
	}

	// Attach IDs to Drawings
	for i, key := range keys {
		drawings[i].ID = key.Name
	}
	return drawings, nil
}
