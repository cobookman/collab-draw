package main

import (
	"cloud.google.com/go/datastore"
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
	Id string `json:"id"`
}

type Canvases struct {
	Canvases []Canvas
}

func (c *Canvas) Get(ctx context.Context, id string) error {
	k, err := datastore.DecodeKey(id)
	if err != nil {
		return err
	}

	err = dc.Get(ctx, k, c)
	return err
}

func (c *Canvas) Create(ctx context.Context) error {
	k := datastore.NewIncompleteKey(ctx, Kind, nil)
	k, err := dc.Put(ctx, k, c)
	if err == nil {
		c.Id = k.Encode()
	}
	return err
}

func (cs *Canvases) GetAll(ctx context.Context, activeSince time.Time, limit int) error {
	q := datastore.NewQuery(Kind).
		Order("-Timestamp").
		Limit(limit).
		Filter("activeSince >", activeSince)

	_, err := dc.GetAll(ctx, q, cs.Canvases)
	return err
}
