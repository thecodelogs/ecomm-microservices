package main

import (
	"log"

	"github.com/manojnegi/ecommerce/api-gateway/internal/client"
	"github.com/manojnegi/ecommerce/api-gateway/internal/config"
	"github.com/manojnegi/ecommerce/api-gateway/internal/middleware"
	"github.com/manojnegi/ecommerce/api-gateway/internal/router"
	"github.com/manojnegi/ecommerce/api-gateway/internal/storage"
)

func main() {
	cfg := config.Load()

	// Connect to user-service
	userClient, err := client.NewUserClient(cfg.UserServiceURL)
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}
	defer userClient.Close()

	// Connect to product-service
	productClient, err := client.NewProductClient(cfg.ProductServiceURL)
	if err != nil {
		log.Fatalf("failed to connect to product-service: %v", err)
	}
	defer productClient.Close()

	// Auth middleware
	authMiddleware := middleware.NewAuthMiddleware(
		cfg.JWTSecret,
		userClient.User,
	)

	// S3 Storage
	s3Storage, err := storage.NewS3Storage(
		cfg.AWSAccessKey,
		cfg.AWSSecretKey,
		cfg.AWSRegion,
		cfg.S3Bucket,
	)
	if err != nil {
		log.Fatalf("failed to initialize s3 storage: %v", err)
	}

	// Setup GraphQL router
	r := router.SetupGraphQL(
		userClient,
		productClient,
		authMiddleware,
		s3Storage,
		cfg,
	)

	// Start server
	log.Printf("GraphQL API Gateway running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
