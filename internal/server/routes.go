package server

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)

	s.App.Get("/health", s.healthHandler)

	s.App.Post("/notify", s.notifyHandler)

	s.App.Get("/events", s.eventsHandler)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) notifyHandler(c *fiber.Ctx) error {
	if s.sns == nil {
		return c.Status(500).JSON(fiber.Map{"error": "SNS client not initialized"})
	}

	// Parse request body for topic ARN and message
	var payload struct {
		TopicArn string `json:"topicArn"`
		Message  string `json:"message"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if payload.TopicArn == "" {
		return c.Status(400).JSON(fiber.Map{"error": "topicArn is required"})
	}

	if payload.Message == "" {
		payload.Message = "Default notification message"
	}

	_, err := s.sns.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: &payload.TopicArn,
		Message:  &payload.Message,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Notification sent"})
}

func (s *FiberServer) eventsHandler(c *fiber.Ctx) error {
	if s.sqs == nil {
		return c.Status(500).JSON(fiber.Map{"error": "SQS client not initialized"})
	}

	// Receive messages from queue
	queueURL := "http://localstack:4566/000000000000/MyQueue" // For LocalStack; use env for real AWS

	result, err := s.sqs.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     0,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	messages := []fiber.Map{}
	for _, msg := range result.Messages {
		messages = append(messages, fiber.Map{
			"messageId": *msg.MessageId,
			"body":      *msg.Body,
		})

		// Delete the message after processing
		_, delErr := s.sqs.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: msg.ReceiptHandle,
		})
		if delErr != nil {
			log.Printf("Failed to delete message: %v", delErr)
		}
	}

	return c.JSON(fiber.Map{"events": messages})
}
