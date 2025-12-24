package rbac

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/models"
	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/repository"
	"github.com/google/uuid"
)

// AuditAdapter adapts RBAC audit events to existing DB-based audit logs
type AuditAdapter struct {
	repo *repository.UserRepository
}

func NewAuditAdapter(repo *repository.UserRepository) *AuditAdapter {
	return &AuditAdapter{repo: repo}
}

// LogEvent satisfies the AuditClient interface
func (a *AuditAdapter) LogEvent(ctx context.Context, event AuditEvent) error {
	log := &models.AuditLog{
		ID:            uuid.New().String(),
		EventType:     event.EventType,
		EventCategory: event.EventCategory,
		Severity:      event.Severity,
		Success:       event.Success,
		CreatedAt:     time.Now(),
	}

	if event.UserID != "" {
		log.UserID = &event.UserID
	}
	if event.OrganizationID != "" {
		// keep in metadata until org_id column is added
		log.Metadata = mapToJSON(map[string]interface{}{
			"organization_id": event.OrganizationID,
			"permission":      event.Metadata["permission"],
			"reason":          event.Metadata["reason"],
		})
	}

	return a.repo.CreateAuditLog(log)
}

// helper
func mapToJSON(m map[string]interface{}) *string {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	s := string(b)
	return &s
}
