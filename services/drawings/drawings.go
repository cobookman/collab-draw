package main

import (
	"os"
	"encoding/json"
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"log"
	"errors"
	"cloud.google.com/go/pubsub"
)

var (
	upstreamTopic *pubsub.Topic
	DrawingKind = "Drawing"
)

func init() {
	// setup datastore client
	ctx := context.Background()
	projectID := os.Getenv("GCLOUD_DATASET_ID")

	// setup upstream drawing topic
	upstreamTopicName := os.Getenv("UPSTREAM_DRAWING_TOPIC")
	if len(upstreamTopicName) == 0 {
		log.Fatal(errors.New("need to set env variable UPSTREAM_DRAWING_TOPIC. " +
			  "variable contains the pubsub topic name of where new drawings are sent"))
	}
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	upstreamTopic = client.Topic(upstreamTopicName)
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Drawing struct {
	// where this drawing belongs
	CanvasID  string `json:"canvasId"`

	// Id for keeping track upstream of duplicate drawings
	DrawingID string `json:"drawingId"`

	// Points to be connected by lines in order that makes up the drawing
	Points []Point `json:"points"`
}

func (d *Drawing) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Drawing) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, d)
}

func HandleNewDrawing(ctx context.Context, drawing Drawing) error {
	data, err := drawing.Marshal()
	if err != nil {
		return err
	}

	_, err = upstreamTopic.Publish(ctx, &pubsub.Message{
		Data: data,
	})
	return err
}

func GetDrawings(ctx context.Context, canvas *Canvas) ([]Drawing, error) {
	if canvas == nil {
		return nil, errors.New("canvas cannot be nil")
	}

	q := datastore.NewQuery(DrawingKind).
		Filter("CanvasID =", canvas.ID)
	drawings := make([]Drawing, 0, 10)
	keys, err := dc.GetAll(ctx, q, &drawings)
	if err != nil {
		return nil, err
	}

	// Attach Drawing Key Encoding to each obj for future ref.
	for i, key := range keys {
		drawings[i].DrawingID = key.Encode()
	}
	return drawings, nil
}
