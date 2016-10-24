package main

import (
	"strings"
	"cloud.google.com/go/datastore"
	"errors"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

var (
	dc   *datastore.Client
	CanvasKind string = "Canvas"
	CanvasSubKind string = "CanvasSub"
)

func init() {
	ctx := context.Background()
	projectID := os.Getenv("GCLOUD_DATASET_ID")

	var err error
	dc, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
}

type Canvas struct {
	Created time.Time
}

type Canvases struct {
	Canvases []Canvas
}

type CanvasSubscription struct {
	CanvasId string `json:"canvasId"`
	TopicId string `json:"topic"`
}

func GetCanvas(ctx context.Context, id string) (*Canvas, error) {
	k, err := datastore.DecodeKey(id)
	if err != nil {
		return nil, err
	}

	if k.Kind() != CanvasKind {
		return nil, errors.New("Invalid key")
	}

	c := new(Canvas)
	err = dc.Get(ctx, k, c)
	return c, err
}

func CreateCanvas(ctx context.Context) (*Canvas, error) {
	k := datastore.NewIncompleteKey(ctx, CanvasKind, nil)
	c := new(Canvas)
	c.Created = time.Now().UTC()
	k, err := dc.Put(ctx, k, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Records the pubsub topic name in a datastore record so that other servers
// know who to notify when a new drawing is added.
func AddCanvasSubscription(ctx context.Context, canvasId string, topicName string) error {
	name := canvasId + "." + topicName
	k := datastore.NewKey(ctx, CanvasSubKind, name, 0, nil)
	sub := new(CanvasSubscription)
	sub.CanvasId = canvasId
	sub.TopicId = topicName

	_, err := dc.Put(ctx, k, sub)
	return err
}

func RemoveCanvasSubscription(ctx context.Context, canvasId string, topicName string) error {
	name := canvasId + "." + topicName
	k := datastore.NewKey(ctx, CanvasSubKind, name, 0, nil)
	return dc.Delete(ctx, k)
}

func GetCanvasSubscriptions(ctx context.Context, canvasId string) ([]string, error) {
	namePrefix := canvasId + "."

	// Get all records with a key that starts with namePrefix
	// \ufffd = largest possible utf8 character
	q := datastore.NewQuery(CanvasSubKind).
		Filter("__key__ >", namePrefix).
		Filter("__key__ <", namePrefix + "\ufffd").
		KeysOnly()

	subNames := make([]string, 0, 10)
	for t := dc.Run(ctx, q); ; {
		var sub CanvasSubscription
		k, err := t.Next(&sub)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		subName := strings.Replace(k.Name(), namePrefix, "", 1)
		subNames = append(subNames, subName)
	}
	return subNames, nil
}

func GetAllCanvases(ctx context.Context, activeSince time.Time, limit int) (*Canvases, error) {
	q := datastore.NewQuery(CanvasKind).
		Order("-Timestamp").
		Limit(limit).
		Filter("activeSince >", activeSince)

	cs := new(Canvases)
	_, err := dc.GetAll(ctx, q, &cs.Canvases)
	return cs, err
}
