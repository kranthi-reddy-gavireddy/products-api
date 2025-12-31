package server

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gofiber/fiber/v2"

	"products-api/internal/database"
)

type FiberServer struct {
	*fiber.App

	db  database.Service
	sns *sns.Client
	sqs *sqs.Client
}

func New() *FiberServer {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load AWS config: %v", err)
		// Continue without SNS if config fails
	}

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "products-api",
			AppName:      "products-api",
		}),

		db: database.New(),
	}

	if cfg.Region != "" {
		server.sns = sns.NewFromConfig(cfg)
		server.sqs = sqs.NewFromConfig(cfg)
	}

	return server
}
