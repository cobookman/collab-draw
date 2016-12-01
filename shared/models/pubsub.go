package models

import (
	"cloud.google.com/go/pubsub"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"log"
	"sync"
)

var (
	pc         *pubsub.Client
	upstream   *pubsub.Topic
	downstream *pubsub.Topic
	pOnce      sync.Once
)

// Instantiates our pubsub client.
func createPubsubClient() {
	ctx := context.Background()
	var err error
	pc, err = pubsub.NewClient(ctx, ProjectID())
	if err != nil {
		log.Fatal(err)
	}
}

// Get a pubsub client, threadsafe.
func PubsubClient() *pubsub.Client {
	if pc == nil {
		pOnce.Do(createPubsubClient)
	}
	return pc
}

// Reference to upstream topic, will create if does not exist.
func UpstreamTopic() *pubsub.Topic {
	if upstream == nil {
		ctx := context.Background()
		upstream, _ = PubsubClient().CreateTopic(ctx, UpstreamTopicName())
	}
	return upstream
}

// Reference to downstream topic, will create if does not exist.
func DownstreamTopic() *pubsub.Topic {
	if downstream == nil {
		ctx := context.Background()
		downstream, _ = PubsubClient().CreateTopic(ctx, "topic-"+uuid.NewV4().String())
	}
	return downstream
}
