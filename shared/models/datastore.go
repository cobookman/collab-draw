package models

import (
	"sync"
	"log"
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

var (
	dc *datastore.Client
	dOnce sync.Once
)

func createClient() {
	ctx := context.Background()
	var err error
	dc, err = datastore.NewClient(ctx, ProjectID())
	if err != nil {
		log.Fatal(err)
	}
}

// Get a datastore client. Threadsafe
func DatastoreClient() *datastore.Client {
	if dc == nil {
		// Only create the client one time, as we re-use
		// our datastore client across threads
		dOnce.Do(createClient)
	}
	return dc
}
