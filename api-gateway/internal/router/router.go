package router

import (
	"github.com/manojnegi/ecommerce/api-gateway/internal/client"

	"github.com/manojnegi/ecommerce/api-gateway/internal/handler"

	"github.com/manojnegi/ecommerce/api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(userClient *client.UserClient, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Health checks (no auth)
	health := handler.NewHealthHandler()
	r.GET("/health", health.Check)
	r.GET("/ready", health.Ready)

	// ── Public Auth Routes ──
	authHandler := handler.NewAuthHandler(userClient.Auth)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/admin/login", authHandler.AdminLogin)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
	}

	// ── User Routes (authenticated) ──
	userHandler := handler.NewUserHandler(userClient.User, userClient.Addr)
	user := r.Group("/api/users/me")
	user.Use(authMiddleware.RequireAuth())
	{
		user.GET("", userHandler.GetProfile)
		user.PUT("", userHandler.UpdateProfile)

		// Addresses
		user.GET("/addresses", userHandler.ListAddresses)
		user.POST("/addresses", userHandler.CreateAddress)
		user.GET("/addresses/:id", userHandler.GetAddress)
		user.PUT("/addresses/:id/default", userHandler.SetDefaultAddress)
		user.DELETE("/addresses/:id", userHandler.DeleteAddress)
	}

	// ── Admin Routes (admin only) ──
	adminHandler := handler.NewAdminHandler(userClient.Admin)
	admin := r.Group("/api/admin")
	admin.Use(authMiddleware.RequireAuth(), middleware.AdminOnly())
	{
		admin.GET("/users", adminHandler.ListUsers)
		admin.GET("/users/:id", adminHandler.GetUser)
		admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
	}

	return r
}
