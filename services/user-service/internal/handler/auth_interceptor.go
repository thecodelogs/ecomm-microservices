package handler

import (
	"context"
	"net"
	"strings"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Context key type (unexported to prevent collisions)
type contextKey string

const (
	ClientIPKey contextKey = "client_ip"
	TokenKey    contextKey = "token"
	UserIDKey   contextKey = "user_id"
	RoleKey     contextKey = "role"
)

// ExtractClientIP gets the real client IP from gRPC context
func ExtractClientIP(ctx context.Context) string {
	// Try X-Forwarded-For first (if behind API Gateway / load balancer)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			ips := strings.Split(xff[0], ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
		// Try X-Real-Ip
		if xri := md.Get("x-real-ip"); len(xri) > 0 {
			return xri[0]
		}
	}

	// Fall back to gRPC peer address
	if p, ok := peer.FromContext(ctx); ok {
		if addr, ok := p.Addr.(*net.TCPAddr); ok {
			return addr.IP.String()
		}
		// Parse "IP:port" format
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			return host
		}
		return p.Addr.String()
	}

	return "unknown"
}

// AuthInterceptor validates PASETO token and injects user info into context
func AuthInterceptor(authSvc *service.AuthService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract and inject client IP
		clientIP := ExtractClientIP(ctx)
		ctx = context.WithValue(ctx, ClientIPKey, clientIP)

		// Public methods that don't need auth
		publicMethods := map[string]bool{
			"/userpb.AuthService/Register":     true,
			"/userpb.AuthService/Login":        true,
			"/userpb.AuthService/AdminLogin":   true,
			"/userpb.AuthService/RefreshToken": true,
		}

		// Skip auth for public methods
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract authorization header
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		// Validate token and extract claims
		userID, role, err := authSvc.ValidateToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// Inject all auth info into context
		ctx = context.WithValue(ctx, TokenKey, token)
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		return handler(ctx, req)
	}
}

// Helper functions to extract values from context

func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

func GetRole(ctx context.Context) string {
	if role, ok := ctx.Value(RoleKey).(string); ok {
		return role
	}
	return ""
}

func GetClientIP(ctx context.Context) string {
	if ip, ok := ctx.Value(ClientIPKey).(string); ok {
		return ip
	}
	return "unknown"
}

func GetToken(ctx context.Context) string {
	if token, ok := ctx.Value(TokenKey).(string); ok {
		return token
	}
	return ""
}
