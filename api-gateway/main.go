package main

import (
	"log"

	"github.com/manojnegi/ecommerce/api-gateway/internal/middleware"

	"github.com/manojnegi/ecommerce/api-gateway/internal/router"

	"github.com/manojnegi/ecommerce/api-gateway/internal/config"

	"github.com/manojnegi/ecommerce/api-gateway/internal/client"
)

func main() {
	cfg := config.Load()

	// Connect to user-service
	userClient, err := client.NewUserClient(cfg.UserServiceURL)
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}
	defer userClient.Close()

	// Setup middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret, userClient.User)

	// Setup router
	r := router.Setup(userClient, authMiddleware)

	// Start server
	log.Printf("API Gateway running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
