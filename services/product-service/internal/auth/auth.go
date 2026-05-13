package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/o1egl/paseto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AccessTokenClaims struct {
	Subject   string    `json:"sub"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

func ExtractClaims(ctx context.Context, secret string) (*AccessTokenClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, errors.New("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == authHeader[0] {
		return nil, errors.New("invalid authorization format")
	}

	key := make([]byte, 32)
	copy(key, []byte(secret))

	var claims AccessTokenClaims
	pasetoV2 := paseto.NewV2()
	err := pasetoV2.Decrypt(token, key, &claims, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if time.Now().UTC().After(claims.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

type contextKey string

const ClaimsKey contextKey = "claims"

func UnaryInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Public methods (optional: add list of public methods here)
		// For now, we attempt to extract claims but don't fail if missing, 
		// letting individual handlers decide if they require auth.
		
		claims, err := ExtractClaims(ctx, secret)
		if err == nil {
			ctx = context.WithValue(ctx, ClaimsKey, claims)
		}

		return handler(ctx, req)
	}
}

func GetClaims(ctx context.Context) (*AccessTokenClaims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*AccessTokenClaims)
	return claims, ok
}
