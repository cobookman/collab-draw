package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"log"
	"os"
)

var (
	upstreamTopic *pubsub.Topic
	DrawingKind   = "Drawing"
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
	X int `json:"y" datastore:"y"`
	Y int `json:"x" datastore:"x"`
}

type Drawing struct {
	// a uuid for the give ndrawing
	DrawingID string `json:"drawingId" datastore:"drawingId"`

	// where this drawing belongs
	CanvasID string `json:"canvasId" datastore:"canvasId"`

	// Points to be connected by lines in order that makes up the drawing
	Points []Point `json:"points" datastore:"points"`
}

func (d *Drawing) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Drawing) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, d)
}

func (d *Drawing) Forward(ctx context.Context) error {
	// if no drawing id, generate a unique id
	if len(d.DrawingID) == 0 {
		d.DrawingID = uuid.NewV4().String()
	}

	data, err := d.Marshal()
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
		Filter("canvasId =", canvas.ID)
	drawings := make([]Drawing, 0, 100)
	keys, err := dc.GetAll(ctx, q, &drawings)
	if err != nil {
		return nil, err
	}

	// Attach Drawing Key Encoding to each obj for future ref.
	for i, key := range keys {
		drawings[i].DrawingID = key.Name
	}
	return drawings, nil
}
