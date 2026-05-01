package handler

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// ExtractClientIP gets the real client IP from gRPC context
func ExtractClientIP(ctx context.Context) string {
	// Try X-Forwarded-For first (if behind API Gateway / load balancer)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			// X-Forwarded-For: client, proxy1, proxy2
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
