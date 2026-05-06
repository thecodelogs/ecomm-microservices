package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminClient userpb.UserServiceClient
}

func NewAdminHandler(adminClient userpb.UserServiceClient) *AdminHandler {
	return &AdminHandler{adminClient: adminClient}
}

// GET /api/admin/users — List all users (paginated, searchable)
// func (h *AdminHandler) ListUsers(c *gin.Context) {
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
// 	status := c.Query("status")
// 	search := c.Query("search")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	resp, err := h.adminClient.ListUsers(ctx, &userpb.ListUsersRequest{
// 		Page:     int32(page),
// 		PageSize: int32(pageSize),
// 		Status:   status,
// 		Search:   search,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	search := c.Query("search")

	// 1. Extract the Authorization header from the incoming HTTP request
	authHeader := c.GetHeader("Authorization")

	// 2. Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Inject the token into gRPC Metadata
	// The key must be lowercase "authorization" for most gRPC interceptors
	md := metadata.Pairs("authorization", authHeader)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 4. Call the client with the ENRICHED context
	resp, err := h.adminClient.ListUsers(ctx, &userpb.ListUsersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Status:   status,
		Search:   search,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GET /api/admin/users/:id — Get specific user
func (h *AdminHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.adminClient.GetUser(ctx, &userpb.GetUserRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PUT /api/admin/users/:id/status — Update user status
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"` // active, suspended, deleted
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.adminClient.UpdateUserStatus(ctx, &userpb.UpdateUserStatusRequest{
		UserId: userID,
		Status: req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DELETE /api/admin/users/:id — Soft delete user
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.adminClient.DeleteUser(ctx, &userpb.DeleteUserRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
