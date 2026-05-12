package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
	"google.golang.org/grpc/metadata"
)

// PASETO claims (must match user-service)
type AccessTokenClaims struct {
	Subject   string    `json:"sub"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

type AuthMiddleware struct {
	pasetoV2     *paseto.V2
	symmetricKey []byte
	userClient   userpb.UserServiceClient
}

func NewAuthMiddleware(secret string, userClient userpb.UserServiceClient) *AuthMiddleware {
	key := make([]byte, 32)
	copy(key, []byte(secret))

	return &AuthMiddleware{
		pasetoV2:     paseto.NewV2(),
		symmetricKey: key,
		userClient:   userClient,
	}
}

// Authenticate middleware attempts to authenticate the user but does not abort if no token is provided.
func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.Next()
			return
		}

		// Validate locally
		claims, err := a.validateToken(token)
		if err != nil {
			// Fallback to gRPC
			claims, err = a.validateViaGRPC(c.Request.Context(), token)
			if err != nil {
				// Invalid token, but we don't abort, just don't set userID
				c.Next()
				return
			}
		}

		// Set user info in context
		c.Set("userID", claims.Subject)
		c.Set("role", claims.Role)
		c.Set("token", token)

		// Inject token into gRPC outgoing context
		md := metadata.Pairs("authorization", "Bearer "+token)
		newCtx := metadata.NewOutgoingContext(c.Request.Context(), md)
		c.Request = c.Request.WithContext(newCtx)

		c.Next()
	}
}

// RequireAuth Gin middleware that validates PASETO token and aborts if missing
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := a.validateToken(token)
		if err != nil {
			claims, err = a.validateViaGRPC(c.Request.Context(), token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
		}

		c.Set("userID", claims.Subject)
		c.Set("role", claims.Role)
		c.Set("token", token)

		md := metadata.Pairs("authorization", "Bearer "+token)
		newCtx := metadata.NewOutgoingContext(c.Request.Context(), md)
		c.Request = c.Request.WithContext(newCtx)

		c.Next()
	}
}

func (a *AuthMiddleware) validateToken(token string) (*AccessTokenClaims, error) {
	var claims AccessTokenClaims
	err := a.pasetoV2.Decrypt(token, a.symmetricKey, &claims, nil)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if time.Now().UTC().After(claims.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func (a *AuthMiddleware) validateViaGRPC(ctx context.Context, token string) (*AccessTokenClaims, error) {
	resp, err := a.userClient.ValidateToken(ctx, &userpb.ValidateTokenRequest{Token: token})
	if err != nil || !resp.Valid {
		return nil, errors.New("invalid token")
	}

	return &AccessTokenClaims{
		Subject: resp.UserId,
		Role:    resp.Role,
	}, nil
}

// Extract user ID from context (call after RequireAuth)
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("userID")
	return userID.(string)
}

// Extract role from context
func GetRole(c *gin.Context) string {
	role, _ := c.Get("role")
	return role.(string)
}
