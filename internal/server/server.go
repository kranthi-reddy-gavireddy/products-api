package server

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gofiber/fiber/v2"

	"products-api/internal/database"
	"products-api/internal/services"
)

type FiberServer struct {
	*fiber.App
	db         database.Service
	sns        *sns.Client
	sqs        *sqs.Client
	product    *services.ProductService
	processors []*MessageProcessor
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

type MessageProcessor struct {
	queueURL string
	handler  func(msg *types.Message) error
}

func New() *FiberServer {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load AWS config: %v", err)
		// Continue without SNS if config fails
	}

	dbSvc := database.New()
	ctx, cancel := context.WithCancel(context.Background())
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "products-api",
			AppName:      "products-api",
		}),

		db:     dbSvc,
		ctx:    ctx,
		cancel: cancel,
	}

	if cfg.Region != "" {
		server.sns = sns.NewFromConfig(cfg)
		server.sqs = sqs.NewFromConfig(cfg)
	}

	return server
}

// AddMessageProcessor adds a new message processor for a queue
func (s *FiberServer) AddMessageProcessor(queueURL string, handler func(msg *types.Message) error) {
	processor := &MessageProcessor{
		queueURL: queueURL,
		handler:  handler,
	}
	s.processors = append(s.processors, processor)
}

// StartMessageProcessors starts all registered message processors
func (s *FiberServer) StartMessageProcessors() {
	for _, processor := range s.processors {
		s.wg.Add(1)
		go s.processQueue(processor)
	}
}

// StopMessageProcessors stops all message processors
func (s *FiberServer) StopMessageProcessors() {
	s.cancel()
	s.wg.Wait()
}

// processQueue continuously processes messages from a queue
func (s *FiberServer) processQueue(processor *MessageProcessor) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("Stopping message processor for queue: %s", processor.queueURL)
			return
		default:
			// Receive messages
			result, err := s.sqs.ReceiveMessage(s.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            &processor.queueURL,
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20, // Long polling
				VisibilityTimeout:   30, // 30 seconds to process
			})

			if err != nil {
				log.Printf("Error receiving messages from %s: %v", processor.queueURL, err)
				time.Sleep(5 * time.Second) // Back off on error
				continue
			}

			// Process messages
			for _, msg := range result.Messages {
				if err := processor.handler(&msg); err != nil {
					log.Printf("Error processing message %s: %v", *msg.MessageId, err)
					// Don't delete the message if processing failed
					continue
				}

				// Delete the message after successful processing
				_, delErr := s.sqs.DeleteMessage(s.ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &processor.queueURL,
					ReceiptHandle: msg.ReceiptHandle,
				})
				if delErr != nil {
					log.Printf("Failed to delete message %s: %v", *msg.MessageId, delErr)
				}
			}
		}
	}
}
