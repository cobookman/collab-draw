package models

import (
	"cloud.google.com/go/datastore"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"time"
)

var (
	CanvasKind string = "Canvas"
)

type Canvas struct {
	ID      string    `datastore:"-" json:"id"`
	Created time.Time `datastore:"created" json:"created"`
}

func GetCanvas(ctx context.Context, id string) (*Canvas, error) {
	c := &Canvas{ID: id}
	if err := DatastoreClient().Get(ctx, c.Key(), c); err != nil {
		return nil, err
	}

	return c, nil
}

func GetCanvases(ctx context.Context, createdSince time.Time, limit int) ([]*Canvas, error) {
	q := datastore.NewQuery(CanvasKind).
		Filter("created >", createdSince).
		Limit(limit)
	var canvases []*Canvas
	keys, err := DatastoreClient().GetAll(ctx, q, &canvases)
	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		canvases[i].ID = k.Name
	}
	return canvases, nil
}

func NewCanvas(ctx context.Context) (*Canvas, error) {
	c := &Canvas{
		Created: time.Now().UTC(),
	}
	c.assignID()
	_, err := DatastoreClient().Put(ctx, c.Key(), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Canvas) assignID() {
	c.ID = uuid.NewV4().String()
}

func (c Canvas) Key() *datastore.Key {
	return datastore.NameKey(CanvasKind, c.ID, nil)
}
