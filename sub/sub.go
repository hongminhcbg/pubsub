package sub

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type Subscriber struct {
	client *pubsub.Client
	handler func(ctx context.Context, m *pubsub.Message)
	subId string
}

func NewSubscriber(client *pubsub.Client, handler func(ctx context.Context, m *pubsub.Message), subId string) *Subscriber {
	return &Subscriber{
		client:  client,
		handler: handler,
		subId:   subId,
	}
}

func (s *Subscriber) Start() error {
	sub := s.client.Subscription(s.subId)
	return  sub.Receive(context.Background(), s.handler)
}


func (s *Subscriber) Close() error {
	return s.client.Close()
}