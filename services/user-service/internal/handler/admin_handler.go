package handler

import (
	"context"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdminLogin — same as Login but requires admin role in token
func (h *AuthHandler) AdminLogin(ctx context.Context, req *userpb.AdminLoginRequest) (*userpb.AuthResponse, error) {
	tokens, user, err := h.authSvc.Login(ctx, req.Email, req.Password, "")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Verify user has admin role
	roles, err := h.userSvc.GetUserRoles(ctx, user.ID)
	if err != nil || !contains(roles, "admin") {
		return nil, status.Error(codes.PermissionDenied, "admin access required")
	}

	return &userpb.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Unix(),
		User:         toProtoUser(user),
	}, nil
}

// ListUsers — admin only
func (h *AuthHandler) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.UserList, error) {
	// Role already checked by interceptor
	users, total, err := h.userSvc.ListUsers(ctx, req.Page, req.PageSize, req.Status, req.Search)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbUsers []*userpb.User
	for _, u := range users {
		// Create a copy of u to safely take its address in the loop
		uCopy := u
		pbUsers = append(pbUsers, toProtoUser(&uCopy))
	}

	return &userpb.UserList{
		Users: pbUsers,
		Total: total,
		Page:  req.Page,
	}, nil
}

// GetUser — admin only
func (h *AuthHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	user, err := h.userSvc.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return toProtoUser(user), nil
}

// UpdateUserStatus — admin only
func (h *AuthHandler) UpdateUserStatus(ctx context.Context, req *userpb.UpdateUserStatusRequest) (*userpb.User, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := h.userSvc.UpdateStatus(ctx, userID, req.Status); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.GetUser(ctx, &userpb.GetUserRequest{UserId: req.UserId})
}

// DeleteUser — soft delete, admin only
func (h *AuthHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.Empty, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := h.userSvc.SoftDelete(ctx, userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.Empty{}, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
