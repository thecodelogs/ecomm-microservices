package handler

import (
	"context"
	"net/http"
	"time"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"github.com/manojnegi/ecommerce/api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
	addrClient userpb.AddressServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient, addrClient userpb.AddressServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
		addrClient: addrClient,
	}
}

// GET /api/users/me — Get own profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.userClient.GetProfile(ctx, &userpb.GetProfileRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PUT /api/users/me — Update own profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req userpb.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserId = userID

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.userClient.UpdateProfile(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ── Addresses ──

// GET /api/users/me/addresses
func (h *UserHandler) ListAddresses(c *gin.Context) {
	userID := middleware.GetUserID(c)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.addrClient.ListAddresses(ctx, &userpb.ListAddressesRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// POST /api/users/me/addresses
func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req userpb.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserId = userID

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.addrClient.CreateAddress(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GET /api/users/me/addresses/:id
func (h *UserHandler) GetAddress(c *gin.Context) {
	addressID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.addrClient.GetAddress(ctx, &userpb.GetAddressRequest{AddressId: addressID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PUT /api/users/me/addresses/:id/default
func (h *UserHandler) SetDefaultAddress(c *gin.Context) {
	userID := middleware.GetUserID(c)
	addressID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	_, err := h.addrClient.SetDefaultAddress(ctx, &userpb.SetDefaultAddressRequest{
		UserId:    userID,
		AddressId: addressID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "default address updated"})
}

// DELETE /api/users/me/addresses/:id
func (h *UserHandler) DeleteAddress(c *gin.Context) {
	addressID := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	_, err := h.addrClient.DeleteAddress(ctx, &userpb.DeleteAddressRequest{AddressId: addressID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
}
