package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/repository"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

type AuthService struct {
	userRepo     *repository.UserRepo
	tokenRepo    *repository.TokenRepo
	pasetoV2     *paseto.V2
	symmetricKey []byte
}

func NewAuthService(userRepo *repository.UserRepo, tokenRepo *repository.TokenRepo, secret string) *AuthService {
	// PASETO requires 32-byte key for local (symmetric) mode
	key := make([]byte, 32)
	copy(key, []byte(secret))

	return &AuthService{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		pasetoV2:     paseto.NewV2(),
		symmetricKey: key,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// ── PASETO Claims Struct ──

type AccessTokenClaims struct {
	Subject   string    `json:"sub"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

type RefreshTokenClaims struct {
	Subject   string    `json:"sub"`
	TokenID   string    `json:"jti"` // unique refresh token ID
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

// ── Public Methods ──

func (s *AuthService) Register(ctx context.Context, email, password, firstName, lastName, clientIP string) (*models.User, error) {
	if _, err := s.userRepo.GetByEmail(ctx, email); err == nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:              uuid.New(),
		Email:           email,
		PasswordHash:    string(hash),
		FirstName:       firstName,
		LastName:        lastName,
		Status:          "active",
		IsEmailVerified: false,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password, clientIP string) (*TokenPair, *models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Check if locked
	if user.LockedUntil.Valid && user.LockedUntil.Time.After(time.Now().UTC()) {
		return nil, nil, errors.New("account temporarily locked due to failed attempts")
	}

	// Check status
	if user.Status != "active" {
		return nil, nil, errors.New("account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		_ = s.userRepo.IncrementFailedLogin(ctx, user.ID)
		return nil, nil, errors.New("invalid credentials")
	}

	// Reset failed attempts and update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user.ID, user.Role, clientIP)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (uuid.UUID, string, error) {
	var claims AccessTokenClaims

	// Decrypt and verify PASETO token
	err := s.pasetoV2.Decrypt(tokenString, s.symmetricKey, &claims, nil)
	if err != nil {
		return uuid.Nil, "", errors.New("invalid token")
	}

	// Check expiration
	if time.Now().UTC().After(claims.ExpiresAt) {
		return uuid.Nil, "", errors.New("token expired")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, "", errors.New("invalid user id")
	}

	// Verify user still exists and is active
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user.Status != "active" {
		return uuid.Nil, "", errors.New("user not found or inactive")
	}

	return userID, claims.Role, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, rawRefreshToken, clientIP string) (*TokenPair, error) {
	var claims RefreshTokenClaims

	// Decrypt refresh token
	err := s.pasetoV2.Decrypt(rawRefreshToken, s.symmetricKey, &claims, nil)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check expiration
	if time.Now().UTC().After(claims.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// Check if token exists in DB and not revoked
	dbToken, err := s.tokenRepo.GetByHash(ctx, rawRefreshToken)
	if err != nil {
		return nil, errors.New("token not found")
	}

	if dbToken.RevokedAt.Valid {
		return nil, errors.New("token has been revoked")
	}

	if dbToken.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("token expired")
	}

	// Fetch user to get current role
	user, err := s.userRepo.GetByID(ctx, dbToken.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new pair
	return s.generateTokenPair(ctx, dbToken.UserID, user.Role, clientIP)
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
	return s.tokenRepo.Revoke(ctx, rawRefreshToken)
}

// ── Private: Generate Token Pair ──

func (s *AuthService) generateTokenPair(ctx context.Context, userID uuid.UUID, role, clientIP string) (*TokenPair, error) {
	now := time.Now().UTC()

	// ── Access Token: 15 minutes ──
	accessExpiry := now.Add(15 * time.Minute)
	accessClaims := AccessTokenClaims{
		Subject:   userID.String(),
		Role:      role,
		IssuedAt:  now,
		ExpiresAt: accessExpiry,
	}

	accessToken, err := s.pasetoV2.Encrypt(s.symmetricKey, accessClaims, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt access token: %w", err)
	}

	// ── Refresh Token: 30 days ──
	refreshTokenID := uuid.New().String()
	refreshExpiry := now.Add(30 * 24 * time.Hour)
	refreshClaims := RefreshTokenClaims{
		Subject:   userID.String(),
		TokenID:   refreshTokenID,
		IssuedAt:  now,
		ExpiresAt: refreshExpiry,
	}

	refreshToken, err := s.pasetoV2.Encrypt(s.symmetricKey, refreshClaims, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	// ── Store refresh token with REAL client IP ──
	deviceInfo := map[string]interface{}{
		"user_agent": extractUserAgent(ctx), // optional helper
		"ip":         clientIP,              // ← real IP from handler
	}
	deviceJSON, _ := json.Marshal(deviceInfo)

	if err := s.tokenRepo.Create(ctx, userID, refreshToken, deviceJSON, refreshExpiry); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Optional: log for audit
	log.Printf("Token generated for user %s from IP %s", userID, clientIP)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry,
	}, nil
}

// Optional: extract user agent from context
func extractUserAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			return ua[0]
		}
	}
	return "unknown"
}
