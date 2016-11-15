package main

import (
	"cloud.google.com/go/datastore"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

var (
	dc   *datastore.Client
	Kind string = "Canvas"
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
	ID      string    `datastore:"-" json:"id"`
	Created time.Time `datastore:"created" json:"created"`
}

type Canvases struct {
	Canvases []Canvas
}

func GetCanvas(ctx context.Context, id string) (*Canvas, error) {
	k := datastore.NameKey(Kind, id, nil)
	c := new(Canvas)
	if err := dc.Get(ctx, k, c); err != nil {
		return nil, err
	}
	c.ID = id
	return c, nil
}

func NewCanvas(ctx context.Context) (*Canvas, error) {
	c := new(Canvas)
	c.Created = time.Now()
	c.ID = uuid.NewV4().String()

	k := datastore.NameKey(Kind, c.ID, nil)
	k, err := dc.Put(ctx, k, c)

	return c, err
}

func (cs *Canvases) GetAll(ctx context.Context, activeSince time.Time, limit int) error {
	q := datastore.NewQuery(Kind).
		Order("-Timestamp").
		Limit(limit).
		Filter("activeSince >", activeSince)

	_, err := dc.GetAll(ctx, q, &cs.Canvases)
	return err
}
