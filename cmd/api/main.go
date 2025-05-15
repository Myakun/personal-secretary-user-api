package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Myakun/personal-secretary-user-api/internal/api"
	"github.com/Myakun/personal-secretary-user-api/internal/config"
	"github.com/Myakun/personal-secretary-user-api/pkg/env"
	loggerPkg "github.com/Myakun/personal-secretary-user-api/pkg/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	envFile := flag.String("env_file", ".env", "Path to environment file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Fatalf("Error loading env file: %s", err)
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	appEnv, err := env.FromString(cfg.Env)
	if err != nil {
		log.Fatalf("Error parsing env: %s", err)
	}

	logger, err := loggerPkg.NewLogger(appEnv)
	if err != nil {
		log.Fatalf("Error initializing logger: %s", err)
	}

	logger.Info("Connecting to MongoDB...")

	// Connect to MongoDB
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.Mongo.User, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.Database)
	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(5 * time.Second)
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("Failed to open mongo: ", err)
	}

	// Verify MongoDB connection
	maxRetries := 5
	retryDelay := time.Second
	for i := 0; i < maxRetries; i++ {
		ctxPing, cancelPing := context.WithTimeout(ctx, 2*time.Second)
		err = mongoClient.Ping(ctxPing, nil)
		cancelPing()
		if err == nil {
			break
		}

		msg := fmt.Sprintf("Mongo ping failed (try %d/%d): %v", i+1, maxRetries, err)
		logger.Fatal(msg)

		time.Sleep(retryDelay)
	}
	if err != nil {
		msg := fmt.Sprintf("Failed to ping mongo after %d retries: %v", maxRetries, err)
		logger.Fatal(msg)
	}

	logger.Info("MongoDB connection established")

	// Disconnect from MongoDB
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect mongo: ", err)
		} else {
			logger.Info("MongoDB connection closed")
		}
	}()

	apiPort := cfg.Api.Port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", apiPort),
		Handler: api.GetRouter(),
	}

	// Start server
	go func() {
		logger.Info("Starting HTTP server on port ", apiPort)
		if err = server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Fatal("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error: ", err)
	}

	logger.Info("HTTP server shut down cleanly")
}
