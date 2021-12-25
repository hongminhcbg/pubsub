package pub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
)

type Publisher struct {
	topic *pubsub.Topic
}

func NewPublisher(topic *pubsub.Topic) *Publisher {
	return &Publisher{topic: topic}
}

func (p *Publisher) PublishMessage(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx := context.Background()
	resp := p.topic.Publish(ctx, &pubsub.Message{
		Data:            b,
	})

	_, err = resp.Get(ctx)
	return err
}