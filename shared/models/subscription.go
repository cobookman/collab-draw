package models

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

var (
	SubscriptionKind string = "Subscription"
)

type Subscription struct {
	CanvasID      string `json:"canvasId" datastore:"canvasId"`
	ServerTopicID string `json:"serverTopicId" datastore:"serverTopicId"`
}

func NewSubscription(canvasID string) *Subscription {
	return &Subscription{
		CanvasID:      canvasID,
		ServerTopicID: DownstreamTopic().ID(),
	}
}

func (sub Subscription) Key() *datastore.Key {
	name := sub.CanvasID + "." + sub.ServerTopicID
	k := datastore.NameKey(SubscriptionKind, name, nil)
	return k
}

func (sub Subscription) Subscribe(ctx context.Context) error {
	k := sub.Key()
	_, err := DatastoreClient().Put(ctx, k, &sub)
	return err
}

func (sub Subscription) Unsubscribe(ctx context.Context) error {
	k := sub.Key()
	return DatastoreClient().Delete(ctx, k)
}
