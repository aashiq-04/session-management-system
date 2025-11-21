package clients

import (
	"fmt"
	"log"

	authpb "github.com/aashiq-04/session-management-system/backend/gateway/proto/auth"
	auditpb "github.com/aashiq-04/session-management-system/backend/gateway/proto/audit"
	sessionpb "github.com/aashiq-04/session-management-system/backend/gateway/proto/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClients holds all gRPC client connections
type GRPCClients struct {
	AuthClient    authpb.AuthServiceClient
	SessionClient sessionpb.SessionServiceClient
	AuditClient   auditpb.AuditServiceClient
}

// NewGRPCClients creates and initializes all gRPC clients
func NewGRPCClients(authURL, sessionURL, auditURL string) (*GRPCClients, error) {
	// Connect to Auth Service
	authConn, err := grpc.Dial(authURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}
	log.Printf("Connected to Auth Service at %s", authURL)

	// Connect to Session Service
	sessionConn, err := grpc.Dial(sessionURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session service: %w", err)
	}
	log.Printf("Connected to Session Service at %s", sessionURL)

	// Connect to Audit Service
	auditConn, err := grpc.Dial(auditURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to audit service: %w", err)
	}
	log.Printf("Connected to Audit Service at %s", auditURL)

	return &GRPCClients{
		AuthClient:    authpb.NewAuthServiceClient(authConn),
		SessionClient: sessionpb.NewSessionServiceClient(sessionConn),
		AuditClient:   auditpb.NewAuditServiceClient(auditConn),
	}, nil
}