package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/kranthi-reddy-gavireddy/message-broker/internal/listeners"
	brokerclient "github.com/kranthi-reddy-gavireddy/message-broker/pkg/broker-client"
)

type App struct {
	cancel context.CancelFunc
	context.Context
}

func New() *App {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	app := &App{
		cancel:  cancel,
		Context: ctx,
	}
	return app
}

func (a *App) Stop() {
	a.cancel()
}

func (a *App) Start() {
	fmt.Println("Starting message broker...")

	topics := listeners.GetTopics()

	for _, topic := range topics {
		handler, exists := listeners.GetHandler(topic)
		if !exists {
			fmt.Printf("No handler found for topic: %s\n", topic)
			continue
		}

		consumer := brokerclient.New([]string{"localhost:9092"}, topic, "message-broker-group")

		go consumer.Start(a.Context, handler)
	}

	<-a.Done()
	fmt.Println("Shutting down message broker...")
}
