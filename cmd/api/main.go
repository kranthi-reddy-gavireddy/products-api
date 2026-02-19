package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"products-api/internal/config"
	"products-api/internal/database"
	"products-api/internal/handlers"
	"products-api/internal/repository"
	"products-api/internal/routes"
	"products-api/internal/server"
	"products-api/internal/services"
	"products-api/logger"
	"products-api/message-broker/listeners"
	"products-api/utils/cache"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger := logger.New()

	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// Stop message processors first
	fiberServer.StopMessageProcessors()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		logger.Errorf("Server forced to shutdown with error: %v", err)
	}

	logger.Infof("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	logger := logger.New()

	logger.Infof("Product API Service Starting...")

	err := config.LoadConfig()
	if err != nil {
		logger.Errorf("Error loading config: %v", err)
		return
	}

	err = cache.Set()
	if err != nil {
		logger.Errorf("Error setting up cache: %v", err)
		return
	}

	server := server.New()
	server.RegisterFiberRoutes()

	dbInstance := database.New().GetDB()

	productRepo := repository.New(dbInstance)
	prodcutService := services.New(productRepo)
	productHandler := handlers.New(prodcutService)
	productRoutes := routes.New(productHandler)
	productRoutes.RegisterRoutes(server.App)

	server.AddMessageProcessor("http://localstack:4566/000000000000/OrderCreatedTopic", func(msg *types.Message) error {
		return listeners.HandleOrderCreation(msg, prodcutService)
	})

	server.StartMessageProcessors()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		err := server.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
