package rbac

import (
	"context"
	"errors"
	"time"
)

// AuditClient is a minimal interface to decouple RBAC from audit implementation
// This should match how your auth-service already talks to audit-service
type AuditClient interface {
	LogEvent(ctx context.Context, event AuditEvent) error
}

// AuditEvent represents the minimum data we need for SOC2 authorization logs
type AuditEvent struct {
	EventType     string
	EventCategory string
	Severity      string
	UserID        string
	OrganizationID string
	Success       bool
	Metadata      map[string]interface{}
	Timestamp     time.Time
}

// Service handles authorization decisions
type Service struct {
	repo        Repository
	auditClient AuditClient
}

func NewService(repo Repository, auditClient AuditClient) *Service {
	return &Service{
		repo:        repo,
		auditClient: auditClient,
	}
}

// Authorize checks whether a user can perform an action.
// FAIL-CLOSED: any error results in denial.
func (s *Service) Authorize(
	ctx context.Context,
	userID string,
	orgID string,
	permission string,
) (bool, error) {

	// Basic input validation
	if userID == "" || orgID == "" || permission == "" {
		s.logDenied(ctx, userID, orgID, permission, "invalid_input")
		return false, errors.New("authorization denied")
	}

	allowed, err := s.repo.HasPermission(ctx, userID, orgID, permission)
	if err != nil {
		// DB or internal error → DENY
		s.logDenied(ctx, userID, orgID, permission, "internal_error")
		return false, errors.New("authorization denied")
	}

	if !allowed {
		s.logDenied(ctx, userID, orgID, permission, "permission_not_granted")
		return false, errors.New("authorization denied")
	}

	return true, nil
}

func (s *Service) logDenied(
	ctx context.Context,
	userID string,
	orgID string,
	permission string,
	reason string,
) {
	if s.auditClient == nil {
		// Fail silently: authorization must not depend on audit availability
		return
	}

	_ = s.auditClient.LogEvent(ctx, AuditEvent{
		EventType:      "AUTHZ_DENIED",
		EventCategory:  "authorization",
		Severity:       "warning",
		UserID:         userID,
		OrganizationID: orgID,
		Success:        false,
		Metadata: map[string]interface{}{
			"permission": permission,
			"reason":     reason,
		},
		Timestamp: time.Now().UTC(),
	})
}
