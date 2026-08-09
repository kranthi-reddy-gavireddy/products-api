package brokerclient

import (
	"context"
	"fmt"

	"github.com/kranthi-reddy-gavireddy/message-broker/internal/listeners"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func New(brokers []string, topic string, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: reader}
}

func (c *Consumer) Start(ctx context.Context, handler listeners.Handler) {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			continue
		}
		handler(msg.Value)
	}
	// Process the message (e.g., call a handler function)
}
