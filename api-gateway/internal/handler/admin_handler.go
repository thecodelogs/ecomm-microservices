package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminClient userpb.UserServiceClient
}

func NewAdminHandler(adminClient userpb.UserServiceClient) *AdminHandler {
	return &AdminHandler{adminClient: adminClient}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	search := c.Query("search")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	_, err := h.adminClient.DeleteUser(ctx, &userpb.DeleteUserRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
