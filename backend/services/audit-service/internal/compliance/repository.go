package compliance

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// --- CC1: Access Control ---

func (r *Repository) CountTotalUsers(orgID string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM user_organizations
		WHERE organization_id = $1
	`, orgID).Scan(&count)
	return count, err
}

func (r *Repository) CountAdminUsers(orgID string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT ur.user_id)
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.organization_id = $1
		  AND r.name IN ('ORG_ADMIN', 'ADMIN')
	`, orgID).Scan(&count)
	return count, err
}

// --- CC6: Logical Access ---

func (r *Repository) CountFailedLoginsLast30Days(orgID string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM audit_logs
		WHERE event_type = 'login_failed'
		  AND created_at >= NOW() - INTERVAL '30 days'
	`).Scan(&count)
	return count, err
}

func (r *Repository) CountAuthzDenialsLast30Days(orgID string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM audit_logs
		WHERE event_type = 'AUTHZ_DENIED'
		  AND created_at >= NOW() - INTERVAL '30 days'
	`).Scan(&count)
	return count, err
}

// --- CC7: System Operations ---

func (r *Repository) CountActiveSessions(orgID string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM sessions
		WHERE is_active = true
	`).Scan(&count)
	return count, err
}

func (r *Repository) CountUnresolvedAlerts(orgID string) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM security_alerts
		WHERE is_resolved = false
	`).Scan(&count)
	return count, err
}
