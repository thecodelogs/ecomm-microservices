package router

import (
	"github.com/gin-gonic/gin"
	"github.com/manojnegi/ecommerce/api-gateway/internal/client"
	"github.com/manojnegi/ecommerce/api-gateway/internal/config"
	"github.com/manojnegi/ecommerce/api-gateway/internal/graphql"
	"github.com/manojnegi/ecommerce/api-gateway/internal/handler"
	"github.com/manojnegi/ecommerce/api-gateway/internal/middleware"
	"github.com/manojnegi/ecommerce/api-gateway/internal/storage"
)

func SetupGraphQL(
	userClient *client.UserClient,
	productClient *client.ProductClient,
	authMiddleware *middleware.AuthMiddleware,
	s3Storage *storage.S3Storage,
	cfg *config.Config,
) *gin.Engine {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Health checks (REST fallback for infrastructure)
	health := handler.NewHealthHandler()
	r.GET("/health", health.Check)
	r.GET("/ready", health.Ready)

	// GraphQL Server
	gqlServer := graphql.NewServer(userClient, productClient, authMiddleware, s3Storage, cfg)

	// GraphQL endpoint
	r.POST("/graphql", authMiddleware.Authenticate(), gqlServer.Handler())
	r.GET("/graphql", authMiddleware.Authenticate(), gqlServer.Handler()) // For queries via GET (useful for caching)

	// GraphQL Playground (disable in production)
	if cfg.Environment != "production" {
		r.GET("/playground", gqlServer.PlaygroundHandler())
	}

	return r
}
