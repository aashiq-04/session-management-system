package clients

import (
	"context"
	"fmt"
	"log"

	auditpb "github.com/aashiq-04/session-management-system/backend/gateway/proto/audit"
	authpb "github.com/aashiq-04/session-management-system/backend/gateway/proto/auth"
	sessionpb "github.com/aashiq-04/session-management-system/backend/gateway/proto/session"
	authzpb "github.com/aashiq-04/session-management-system/backend/services/auth-service/proto/authorization"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	comppb "github.com/aashiq-04/session-management-system/backend/services/audit-service/proto/compliance"
)

// GRPCClients holds all gRPC client connections
type GRPCClients struct {
	AuthClient    authpb.AuthServiceClient
	SessionClient sessionpb.SessionServiceClient
	AuditClient   auditpb.AuditServiceClient
	AuthzClient   authzpb.AuthorizationServiceClient
	ComplianceClient comppb.ComplianceServiceClient
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
	complianceConn, err := grpc.Dial(auditURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to compliance service: %w", err)
	}
	log.Printf("Connected to Compliance Service at %s", auditURL)

	return &GRPCClients{
		AuthClient:    authpb.NewAuthServiceClient(authConn),
		SessionClient: sessionpb.NewSessionServiceClient(sessionConn),
		AuditClient:   auditpb.NewAuditServiceClient(auditConn),
		AuthzClient:   authzpb.NewAuthorizationServiceClient(authConn),
		ComplianceClient: comppb.NewComplianceServiceClient(complianceConn),
		}, nil
}
func (c *GRPCClients) CheckPermission(
	ctx context.Context,
	userID string,
	orgID string,
	permission string,
) bool {

	resp, err := c.AuthzClient.CanUserPerform(ctx, &authzpb.AuthorizationRequest{
		UserId:         userID,
		OrganizationId: orgID,
		Permission:     permission,
	})

	if err != nil {
		log.Printf("[RBAC][ERROR] user=%s perm=%s err=%v", userID, permission, err)
		return false
	}

	log.Printf("[RBAC][DRY-RUN] user=%s perm=%s allowed=%v",
		userID, permission, resp.Allowed)

	return resp.Allowed
}
