package main

import (
	"log"
	"os"
	"golang.org/x/net/context"
	"cloud.google.com/go/datastore"
)

var (
	dc *datastore.Client
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

func datastoreClient() *datastore.Client {
	return dc
}
