package main

import (
	"time"
	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)

type MessagingQueue struct {
	Topic        *pubsub.Topic
	Subscription *pubsub.Subscription
}

// Creates a new subscription to a messaging queue
func NewMessagingQueue(ctx context.Context, projectID string, topicName string, subName string) (*MessagingQueue, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic, _ := client.CreateTopic(ctx, topicName)
	sub, err := client.CreateSubscription(ctx, subName, topic, 10*time.Second, nil)
	if err != nil {
		return nil, err
	}

	mq := new(MessagingQueue)
	mq.Topic = topic
	mq.Subscription = sub
	return mq, nil
}

func (mq *MessagingQueue) Cleanup(ctx context.Context) error {
	if mq.Subscription != nil {
		if err := mq.Subscription.Delete(ctx); err != nil {
			return err
		}
		mq.Subscription = nil
	}
	if mq.Topic != nil {
		if err := mq.Topic.Delete(ctx); err != nil {
			return err
		}
		mq.Topic = nil
	}
	return nil
}

type IncomingMessageHandler func(msg *pubsub.Message) error

func (mq MessagingQueue) OnMessage(f IncomingMessageHandler) error {
	ctx := context.Background()
	it, err := mq.Subscription.Pull(ctx)
	if err != nil {
		return err
	}
	defer it.Stop()

	// You might want to use a worker pool to pull down the messages
	// to maximize messages/s processed. For now we are going to run
	// this in 1 thread.
	for {
		msg, err := it.Next()
		if err != nil {
			return err
		}
		if err := f(msg); err != nil {
			msg.Done(false)
		} else {
			msg.Done(true)
		}

	}
}
