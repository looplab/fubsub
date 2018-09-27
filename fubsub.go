package fubsub

import (
	"context"
	"errors"
	"io"
	"math/rand"

	"cloud.google.com/go/pubsub"
)

var _ = io.ReadWriteCloser(&Fubsub{})

// Fubsub is a io.ReadWriteCloser for GCP pubsub.
type Fubsub struct {
	topic *pubsub.Topic
	sub   *pubsub.Subscription
}

// New creates a new Fubsub.
func New(projectID, name string) (*Fubsub, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.New("could not create pubsub client: " + err.Error())
	}

	// Get or create the topic.
	topicName := "fubsub_" + name
	topic := client.Topic(topicName)
	if ok, err := topic.Exists(ctx); err != nil {
		return nil, errors.New("could not get topic name: " + err.Error())
	} else if !ok {
		if topic, err = client.CreateTopic(ctx, topicName); err != nil {
			return nil, errors.New("could not create topic: " + err.Error())
		}
	}

	// Get or create the subscription.
	subscriptionID := topicName + "_" + randString(8)
	sub := client.Subscription(subscriptionID)
	if ok, err := sub.Exists(ctx); err != nil {
		return nil, errors.New("could not get subscription name: " + err.Error())
	} else if !ok {
		if sub, err = client.CreateSubscription(ctx, subscriptionID,
			pubsub.SubscriptionConfig{Topic: topic},
		); err != nil {
			return nil, errors.New("could not create subscription: " + err.Error())
		}
	}

	return &Fubsub{
		topic: topic,
		sub:   sub,
	}, nil
}

// Read implements the Read method of the io.Reader interface.
func (f *Fubsub) Read(p []byte) (n int, err error) {
	// NOTE: Not at all meant to call Receive every time like this...
	ctx, cancel := context.WithCancel(context.Background())
	if err = f.sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		copy(p, msg.Data)
		n = len(msg.Data)
		msg.Ack()
		// Hack to close the receive loop after each single message.
		cancel()
	}); err == context.Canceled {
		err = nil
	}
	return
}

// Write implements the Write method of the io.Writer interface.
func (f *Fubsub) Write(p []byte) (n int, err error) {
	ctx := context.Background()
	r := f.topic.Publish(ctx, &pubsub.Message{Data: p})
	if _, err = r.Get(ctx); err == nil {
		n = len(p)
	}
	return
}

// Close implements the Close method of the io.Closer interface.
func (f *Fubsub) Close() error {
	ctx := context.Background()
	if err := f.sub.Delete(ctx); err != nil {
		return errors.New("could not delete subscription: " + err.Error())
	}
	if err := f.topic.Delete(ctx); err != nil {
		return errors.New("could not delete topic: " + err.Error())
	}
	return nil
}

func randString(n int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
