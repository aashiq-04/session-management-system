package graph

import (
	"github.com/aashiq-04/session-management-system/backend/gateway/clients"
)

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