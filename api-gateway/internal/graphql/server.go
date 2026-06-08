package graphql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/manojnegi/ecommerce/api-gateway/internal/client"
	"github.com/manojnegi/ecommerce/api-gateway/internal/config"
	"github.com/manojnegi/ecommerce/api-gateway/internal/graphql/generated"
	"github.com/manojnegi/ecommerce/api-gateway/internal/graphql/resolver"
	"github.com/manojnegi/ecommerce/api-gateway/internal/middleware"
	"github.com/manojnegi/ecommerce/api-gateway/internal/storage"
)

type Server struct {
	userClient     *client.UserClient
	productClient  *client.ProductClient
	authMiddleware *middleware.AuthMiddleware
	s3Storage      *storage.S3Storage
	cfg            *config.Config
	resolver       *resolver.Resolver
}

func NewServer(
	userClient *client.UserClient,
	productClient *client.ProductClient,
	authMiddleware *middleware.AuthMiddleware,
	s3Storage *storage.S3Storage,
	cfg *config.Config,
) *Server {
	res := resolver.New(userClient, productClient, authMiddleware, s3Storage, cfg)
	return &Server{
		userClient:     userClient,
		productClient:  productClient,
		authMiddleware: authMiddleware,
		s3Storage:      s3Storage,
		cfg:            cfg,
		resolver:       res,
	}
}

func (s *Server) Handler() gin.HandlerFunc {
	config := generated.Config{Resolvers: s.resolver}

	// @auth directive: ensures user is authenticated
	config.Directives.Auth = func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
		gc, err := resolver.GinContextFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("internal error: gin context missing")
		}

		// The AuthMiddleware already populates "userID" if the token is valid
		userID, exists := gc.Get("userID")
		if !exists || userID == "" {
			return nil, fmt.Errorf("access denied: unauthenticated (check Authorization header)")
		}

		return next(ctx)
	}

	// @admin directive: ensures user has ADMIN role
	config.Directives.Admin = func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
		gc, err := resolver.GinContextFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("internal error: gin context missing")
		}

		role, exists := gc.Get("role")
		log.Printf("DEBUG: Admin Check - exists: %v, role: %v (%T)", exists, role, role)

		roleStr, ok := role.(string)
		if !exists {
			return nil, fmt.Errorf("access denied: missing role in context (unauthenticated?)")
		}
		if !ok {
			return nil, fmt.Errorf("access denied: role is not a string (%T)", role)
		}
		if !strings.EqualFold(roleStr, "admin") {
			return nil, fmt.Errorf("access denied: admin only (your role: %s)", roleStr)
		}

		return next(ctx)
	}

	// Create gqlgen handler
	h := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	return func(c *gin.Context) {
		// Pass gin context to graphql context
		ctx := context.WithValue(c.Request.Context(), "GinContext", c)
		c.Request = c.Request.WithContext(ctx)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
