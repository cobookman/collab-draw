package models

import (
	"golang.org/x/net/context"
	"testing"
)

func TestSubscription(t *testing.T) {
	sub := NewSubscription("SomeCanvasID", "SomeHostTopicID")
	if sub.Key().Name != "SomeCanvasID.SomeHostTopicID" {
		t.Fatal("Key is not pointing to right resource")
	}
	ctx := context.Background()
	if err := sub.Subscribe(ctx); err != nil {
		t.Fatal(err)
	}
	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatal(err)
	}
}
