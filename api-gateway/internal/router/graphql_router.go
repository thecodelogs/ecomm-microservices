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

	// REST Handlers
	authHandler := handler.NewAuthHandler(userClient.Auth)
	userHandler := handler.NewUserHandler(userClient.User, userClient.Addr)
	productHandler := handler.NewProductHandler(productClient.Product, productClient.Category, s3Storage)
	adminHandler := handler.NewAdminHandler(userClient.User)

	// API Group
	api := r.Group("/api")
	{
		// Public Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/admin/login", authHandler.AdminLogin)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
		}

		// User Panel / Customer Routes
		users := api.Group("/users")
		users.Use(authMiddleware.RequireAuth())
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)

			// Addresses
			users.GET("/me/addresses", userHandler.ListAddresses)
			users.POST("/me/addresses", userHandler.CreateAddress)
			users.GET("/me/addresses/:id", userHandler.GetAddress)
			users.PUT("/me/addresses/:id/default", userHandler.SetDefaultAddress)
			users.DELETE("/me/addresses/:id", userHandler.DeleteAddress)
		}

		// Public Catalog Endpoints (User Panel API)
		api.GET("/products", productHandler.ListProducts)
		api.GET("/categories", productHandler.ListCategories)

		// Admin Catalog Endpoints
		catalogAdmin := api.Group("/categories")
		catalogAdmin.Use(authMiddleware.RequireAuth(), middleware.AdminOnly())
		{
			catalogAdmin.POST("", productHandler.CreateCategory)
		}

		// Admin User Management Group
		adminUsers := api.Group("/admin/users")
		adminUsers.Use(authMiddleware.RequireAuth(), middleware.AdminOnly())
		{
			adminUsers.GET("", adminHandler.ListUsers)
			adminUsers.GET("/:id", adminHandler.GetUser)
			adminUsers.PUT("/:id/status", adminHandler.UpdateUserStatus)
			adminUsers.DELETE("/:id", adminHandler.DeleteUser)
		}
	}

	return r
}
