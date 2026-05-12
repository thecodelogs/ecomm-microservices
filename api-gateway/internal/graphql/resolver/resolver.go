package resolver

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/manojnegi/ecommerce/api-gateway/internal/client"
	"github.com/manojnegi/ecommerce/api-gateway/internal/config"

	"github.com/manojnegi/ecommerce/api-gateway/internal/middleware"
	"github.com/manojnegi/ecommerce/api-gateway/internal/storage"
)

// Resolver is the root resolver
type Resolver struct {
	UserClient     *client.UserClient
	ProductClient  *client.ProductClient
	AuthMiddleware *middleware.AuthMiddleware
	S3Storage      *storage.S3Storage
	Config         *config.Config
}

func New(
	userClient *client.UserClient,
	productClient *client.ProductClient,
	authMiddleware *middleware.AuthMiddleware,
	s3Storage *storage.S3Storage,
	cfg *config.Config,
) *Resolver {
	return &Resolver{
		UserClient:     userClient,
		ProductClient:  productClient,
		AuthMiddleware: authMiddleware,
		S3Storage:      s3Storage,
		Config:         cfg,
	}
}

// Helper to extract Gin context
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContext")
	if ginContext == nil {
		return nil, fmt.Errorf("gin context not found")
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("invalid gin context type")
	}
	return gc, nil
}


func intPtr(i int) *int { return &i }
