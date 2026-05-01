package handler

import (
	"context"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	handler "github.com/manojnegi/ecomm-microservices/services/user-service/internal/middleware"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	userpb.UnimplementedAuthServiceServer
	userpb.UnimplementedUserServiceServer
	authSvc *service.AuthService
	userSvc *service.UserService
}

func NewAuthHandler(authSvc *service.AuthService, userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.AuthResponse, error) {
	// Extract real client IP from context
	clientIP := handler.ExtractClientIP(ctx)

	tokens, user, err := h.authSvc.Login(ctx, req.Email, req.Password, clientIP)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userpb.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Unix(),
		User:         toProtoUser(user),
	}, nil
}

// Register also passes IP
func (h *AuthHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	clientIP := handler.ExtractClientIP(ctx)

	user, err := h.authSvc.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName, clientIP)
	if err != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	tokens, _, err := h.authSvc.Login(ctx, req.Email, req.Password, clientIP)
	if err != nil {
		return nil, status.Error(codes.Internal, "registration succeeded but login failed")
	}

	return &userpb.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Unix(),
		User:         toProtoUser(user),
	}, nil
}

// ── AuthService RPCs ──

func (h *AuthHandler) RefreshToken(ctx context.Context, req *userpb.RefreshTokenRequest) (*userpb.TokenResponse, error) {
	clientIP := handler.ExtractClientIP(ctx)
	tokens, err := h.authSvc.RefreshToken(ctx, req.RefreshToken, clientIP)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userpb.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt.Unix(),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *userpb.LogoutRequest) (*userpb.Empty, error) {
	if err := h.authSvc.Logout(ctx, req.RefreshToken); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userpb.Empty{}, nil
}

// ── UserService RPCs ──

func (h *AuthHandler) ValidateToken(ctx context.Context, req *userpb.ValidateTokenRequest) (*userpb.ValidateTokenResponse, error) {
	userID, role, err := h.authSvc.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userpb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID.String(),
		Role:   role,
	}, nil
}

func (h *AuthHandler) GetProfile(ctx context.Context, req *userpb.GetProfileRequest) (*userpb.User, error) {
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

func (h *AuthHandler) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*userpb.User, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	if err := h.userSvc.UpdateProfile(ctx, userID, req.FirstName, req.LastName, req.Phone); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user, err := h.userSvc.GetProfile(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoUser(user), nil
}

// ── Helpers ──

func toProtoUser(u *models.User) *userpb.User {
	return &userpb.User{
		Id:              u.ID.String(),
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Phone:           u.Phone.String,
		AvatarUrl:       u.AvatarURL.String,
		Status:          u.Status,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt.Unix(),
	}
}
