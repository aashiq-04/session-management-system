package graph

import (
	"context"
	"net/http"
	"net"
	"github.com/aashiq-04/session-management-system/backend/gateway/clients"
)
func extractRealIP(r *http.Request) string {
    // X-Forwarded-For (proxy)
    if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
        return forwarded
    }

    // X-Real-IP (nginx)
    if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
        return realIP
    }

    // Remote address fallback
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err == nil {
        return ip
    }

    return "0.0.0.0"
}


func getIPFromContext(ctx context.Context) string {
    req, ok := ctx.Value("httpRequest").(*http.Request)
    if !ok || req == nil {
        return "0.0.0.0"
    }

    // Call the shared helper from main.go
    return extractRealIP(req)
}

// Resolver is the main resolver that holds all dependencies
type Resolver struct {
	Clients   *clients.GRPCClients
	JWTSecret string
}

// NewResolver creates a new resolver instance
func NewResolver(clients *clients.GRPCClients, jwtSecret string) *Resolver {
	return &Resolver{
		Clients:   clients,
		JWTSecret: jwtSecret,
	}
}