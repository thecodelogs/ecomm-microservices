package handler

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Role context key
type roleKey struct{}

// WithRole injects role into context (used after auth)
func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey{}, role)
}

// GetRole extracts role from context
func GetRole(ctx context.Context) string {
	if r, ok := ctx.Value(roleKey{}).(string); ok {
		return r
	}
	return ""
}

// RequireRole returns interceptor that checks for required role
func RequireRole(roles ...string) grpc.UnaryServerInterceptor {
	allowed := make(map[string]bool)
	for _, r := range roles {
		allowed[r] = true
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		role := GetRole(ctx)
		if !allowed[role] {
			return nil, status.Error(codes.PermissionDenied, "insufficient privileges")
		}
		return handler(ctx, req)
	}
}
