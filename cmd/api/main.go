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

	"github.com/Myakun/personal-secretary-user-api/internal/config"
	"github.com/Myakun/personal-secretary-user-api/internal/delivery/api"
	registerHandler "github.com/Myakun/personal-secretary-user-api/internal/delivery/api/handler/register"
	mongoUserRepo "github.com/Myakun/personal-secretary-user-api/internal/infrastructure/repository/mongo/user"
	registerPresentationPkg "github.com/Myakun/personal-secretary-user-api/internal/presentation/user/registration"
	userUseCasePkg "github.com/Myakun/personal-secretary-user-api/internal/usecase/user"
	"github.com/Myakun/personal-secretary-user-api/pkg/env"
	pkgLogger "github.com/Myakun/personal-secretary-user-api/pkg/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	envFile := flag.String("env_file", ".env", "Path to environment file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Fatalf("Err loading env file: %s", err)
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Err loading config: %s", err)
	}

	appEnv, err := env.FromString(cfg.Env)
	if err != nil {
		log.Fatalf("Err parsing env: %s", err)
	}

	logger, err := pkgLogger.NewLogger(appEnv)
	if err != nil {
		log.Fatalf("Err initializing logger: %s", err)
	}

	logger.Info("Connecting to MongoDB...")

	// Connect to MongoDB
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.Mongo.User, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.Database)
	clientOptions := mongoOptions.Client().ApplyURI(uri).SetConnectTimeout(5 * time.Second)
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

	mongoDb := mongoClient.Database(cfg.Mongo.Database)

	// Repositories
	userRepo := mongoUserRepo.NewUserRepository(mongoDb.Collection("users"), logger)

	// Use cases
	userUseCase := userUseCasePkg.NewUserUseCase(logger, userRepo)

	// Presentations
	registerPresentation := registerPresentationPkg.NewUserRegistration(logger, userUseCase)

	// API handlers
	routerHandlers := &api.RouterHandlers{
		RegisterHandler: registerHandler.NewRegisterHandler(logger, registerPresentation),
	}

	// Prepare server
	apiPort := cfg.Api.Port
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", apiPort),
		Handler: api.GetRouter(routerHandlers),
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
