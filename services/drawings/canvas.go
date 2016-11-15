package main

import (
	"cloud.google.com/go/datastore"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"strings"
	"time"
)

var (
	CanvasKind    string = "Canvas"
	CanvasSubKind string = "CanvasSub"
)

type Canvas struct {
	ID      string    `datastore:"-" json:"id"`
	Created time.Time `datastore:"created" json:"created"`
}

type Canvases struct {
	Canvases []Canvas
}

type CanvasSubscription struct {
	CanvasId string `json:"canvasId" datastore:"canvasId"`
	TopicId  string `json:"topic" datastore:"topicId"`
}

func GetCanvas(ctx context.Context, id string) (*Canvas, error) {
	k := datastore.NameKey(CanvasKind, id, nil)
	c := new(Canvas)
	if err := datastoreClient().Get(ctx, k, c); err != nil {
		return nil, err
	}

	// attach the canvas's id for future ref.
	c.ID = id
	return c, nil
}

func CreateCanvas(ctx context.Context) (*Canvas, error) {
	id := uuid.NewV4().String()
	k := datastore.NameKey(CanvasKind, id, nil)
	c := &Canvas{
		Created: time.Now(),
	}

	c.Created = time.Now().UTC()
	k, err := datastoreClient().Put(ctx, k, c)

	if err != nil {
		return nil, err
	}

	c.ID = k.Name
	return c, nil
}

// Records the pubsub topic name in a datastore record so that other servers
// know who to notify when a new drawing is added.
func AddCanvasSubscription(ctx context.Context, canvasId string, topicName string) error {
	name := canvasId + "." + topicName
	k := datastore.NameKey(CanvasSubKind, name, nil)
	sub := new(CanvasSubscription)
	sub.CanvasId = canvasId
	sub.TopicId = topicName

	_, err := datastoreClient().Put(ctx, k, sub)
	return err
}

func RemoveCanvasSubscription(ctx context.Context, canvasId string, topicName string) error {
	name := canvasId + "." + topicName
	k := datastore.NameKey(CanvasSubKind, name, nil)
	return dc.Delete(ctx, k)
}

func GetCanvasSubscriptions(ctx context.Context, canvasId string) ([]string, error) {
	namePrefix := canvasId + "."

	// Get all records with a key that starts with namePrefix
	// \ufffd = largest possible utf8 character
	q := datastore.NewQuery(CanvasSubKind).
		Filter("__key__ >", namePrefix).
		Filter("__key__ <", namePrefix+"\ufffd").
		KeysOnly()

	subNames := make([]string, 0, 10)
	for t := datastoreClient().Run(ctx, q); ; {
		var sub CanvasSubscription
		k, err := t.Next(&sub)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		subName := strings.Replace(k.Name, namePrefix, "", 1)
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
	keys, err := datastoreClient().GetAll(ctx, q, &cs.Canvases)

	// Attach Id to each canvas obj for future ref.
	for i, key := range keys {
		cs.Canvases[i].ID = key.Encode()
	}
	return cs, err
}
