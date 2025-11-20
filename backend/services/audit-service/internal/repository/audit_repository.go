package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/aashiq-04/session-management-system/backend/services/audit-service/internal/models"
)

// AuditRepository handles database operations for audit logs
type AuditRepository struct {
	db *sql.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// CreateAuditLog creates a new audit log entry
func (r *AuditRepository) CreateAuditLog(log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, session_id, device_id, event_type, event_category,
		                        severity, ip_address, user_agent, location_country, location_city,
		                        metadata, success, failure_reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	
	_, err := r.db.Exec(
		query,
		log.ID,
		log.UserID,
		log.SessionID,
		log.DeviceID,
		log.EventType,
		log.EventCategory,
		log.Severity,
		log.IPAddress,
		log.UserAgent,
		log.LocationCountry,
		log.LocationCity,
		log.Metadata,
		log.Success,
		log.FailureReason,
		log.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	
	return nil
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (r *AuditRepository) GetUserAuditLogs(userID string, limit, offset int, eventCategory, severity string, successOnly bool) ([]models.AuditLog, int, error) {
	// Build query with filters
	query := `
		SELECT id, user_id, session_id, device_id, event_type, event_category,
		       severity, ip_address, user_agent, location_country, location_city,
		       metadata, success, failure_reason, created_at
		FROM audit_logs
		WHERE user_id = $1
	`
	
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1`
	
	args := []interface{}{userID}
	argIndex := 2
	
	if eventCategory != "" {
		query += fmt.Sprintf(" AND event_category = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND event_category = $%d", argIndex)
		args = append(args, eventCategory)
		argIndex++
	}
	
	if severity != "" {
		query += fmt.Sprintf(" AND severity = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND severity = $%d", argIndex)
		args = append(args, severity)
		argIndex++
	}
	
	if successOnly {
		query += fmt.Sprintf(" AND success = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND success = $%d", argIndex)
		args = append(args, true)
		argIndex++
	}
	
	query += " ORDER BY created_at DESC"
	
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, limit, offset)
	}
	
	// Get total count
	var totalCount int
	err := r.db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}
	
	// Get logs
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()
	
	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.SessionID, &log.DeviceID,
			&log.EventType, &log.EventCategory, &log.Severity,
			&log.IPAddress, &log.UserAgent, &log.LocationCountry, &log.LocationCity,
			&log.Metadata, &log.Success, &log.FailureReason, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}
	
	return logs, totalCount, nil
}

// GetAuditLogsByEvent retrieves audit logs by event type
func (r *AuditRepository) GetAuditLogsByEvent(eventType string, limit, offset int, startDate, endDate time.Time) ([]models.AuditLog, int, error) {
	query := `
		SELECT id, user_id, session_id, device_id, event_type, event_category,
		       severity, ip_address, user_agent, location_country, location_city,
		       metadata, success, failure_reason, created_at
		FROM audit_logs
		WHERE event_type = $1 AND created_at >= $2 AND created_at <= $3
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`
	
	countQuery := `
		SELECT COUNT(*)
		FROM audit_logs
		WHERE event_type = $1 AND created_at >= $2 AND created_at <= $3
	`
	
	// Get total count
	var totalCount int
	err := r.db.QueryRow(countQuery, eventType, startDate, endDate).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}
	
	// Get logs
	rows, err := r.db.Query(query, eventType, startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()
	
	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(
			&log.ID, &log.UserID, &log.SessionID, &log.DeviceID,
			&log.EventType, &log.EventCategory, &log.Severity,
			&log.IPAddress, &log.UserAgent, &log.LocationCountry, &log.LocationCity,
			&log.Metadata, &log.Success, &log.FailureReason, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}
	
	return logs, totalCount, nil
}

// CreateSecurityAlert creates a new security alert
func (r *AuditRepository) CreateSecurityAlert(alert *models.SecurityAlert) error {
	query := `
		INSERT INTO security_alerts (id, user_id, alert_type, severity, description,
		                             metadata, ip_address, location_country, location_city,
		                             is_resolved, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.Exec(
		query,
		alert.ID,
		alert.UserID,
		alert.AlertType,
		alert.Severity,
		alert.Description,
		alert.Metadata,
		alert.IPAddress,
		alert.LocationCountry,
		alert.LocationCity,
		alert.IsResolved,
		alert.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create security alert: %w", err)
	}
	
	return nil
}

// GetSecurityAlerts retrieves security alerts for a user
func (r *AuditRepository) GetSecurityAlerts(userID string, includeResolved bool, severity string) ([]models.SecurityAlert, error) {
	query := `
		SELECT id, user_id, alert_type, severity, description, metadata,
		       ip_address, location_country, location_city,
		       is_resolved, resolved_at, created_at
		FROM security_alerts
		WHERE user_id = $1
	`
	
	args := []interface{}{userID}
	argIndex := 2
	
	if !includeResolved {
		query += " AND is_resolved = false"
	}
	
	if severity != "" {
		query += fmt.Sprintf(" AND severity = $%d", argIndex)
		args = append(args, severity)
	}
	
	query += " ORDER BY created_at DESC"
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query security alerts: %w", err)
	}
	defer rows.Close()
	
	var alerts []models.SecurityAlert
	for rows.Next() {
		var alert models.SecurityAlert
		err := rows.Scan(
			&alert.ID, &alert.UserID, &alert.AlertType, &alert.Severity,
			&alert.Description, &alert.Metadata, &alert.IPAddress,
			&alert.LocationCountry, &alert.LocationCity,
			&alert.IsResolved, &alert.ResolvedAt, &alert.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan security alert: %w", err)
		}
		alerts = append(alerts, alert)
	}
	
	return alerts, nil
}

// ResolveSecurityAlert marks a security alert as resolved
func (r *AuditRepository) ResolveSecurityAlert(alertID string) error {
	query := `
		UPDATE security_alerts
		SET is_resolved = true, resolved_at = $1
		WHERE id = $2
	`
	
	result, err := r.db.Exec(query, time.Now(), alertID)
	if err != nil {
		return fmt.Errorf("failed to resolve security alert: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("security alert not found")
	}
	
	return nil
}

// GetComplianceReport generates a compliance report
func (r *AuditRepository) GetComplianceReport(userID string, startDate, endDate time.Time) (*models.ComplianceReport, error) {
	report := &models.ComplianceReport{}
	
	// Get total events
	query := `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3`
	err := r.db.QueryRow(query, userID, startDate, endDate).Scan(&report.TotalEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to get total events: %w", err)
	}
	
	// Get successful logins
	query = `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1 AND event_type = 'user_login' AND success = true AND created_at >= $2 AND created_at <= $3`
	err = r.db.QueryRow(query, userID, startDate, endDate).Scan(&report.SuccessfulLogins)
	if err != nil {
		return nil, fmt.Errorf("failed to get successful logins: %w", err)
	}
	
	// Get failed logins
	query = `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1 AND event_type = 'login_failed' AND created_at >= $2 AND created_at <= $3`
	err = r.db.QueryRow(query, userID, startDate, endDate).Scan(&report.FailedLogins)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed logins: %w", err)
	}
	
	// Get session revocations
	query = `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1 AND event_type IN ('session_revoked', 'all_sessions_revoked') AND created_at >= $2 AND created_at <= $3`
	err = r.db.QueryRow(query, userID, startDate, endDate).Scan(&report.SessionRevocations)
	if err != nil {
		return nil, fmt.Errorf("failed to get session revocations: %w", err)
	}
	
	// Get security alerts
	query = `SELECT COUNT(*) FROM security_alerts WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3`
	err = r.db.QueryRow(query, userID, startDate, endDate).Scan(&report.SecurityAlerts)
	if err != nil {
		return nil, fmt.Errorf("failed to get security alerts: %w", err)
	}
	
	// Get MFA events
	query = `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1 AND event_type IN ('mfa_enabled', 'mfa_verified') AND created_at >= $2 AND created_at <= $3`
	err = r.db.QueryRow(query, userID, startDate, endDate).Scan(&report.MFAEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to get MFA events: %w", err)
	}
	
	// Get event breakdown
	query = `
		SELECT event_type, COUNT(*) as count
		FROM audit_logs
		WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY event_type
		ORDER BY count DESC
		LIMIT 10
	`
	rows, err := r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get event breakdown: %w", err)
	}
	defer rows.Close()
	
	report.EventBreakdown = []models.EventCount{}
	for rows.Next() {
		var ec models.EventCount
		if err := rows.Scan(&ec.EventType, &ec.Count); err != nil {
			continue
		}
		report.EventBreakdown = append(report.EventBreakdown, ec)
	}
	
	// Get top locations
	query = `
		SELECT DISTINCT COALESCE(location_city || ', ' || location_country, location_country, 'Unknown') as location
		FROM audit_logs
		WHERE user_id = $1 AND created_at >= $2 AND created_at <= $3 AND location_country IS NOT NULL
		ORDER BY created_at DESC
		LIMIT 5
	`
	rows, err = r.db.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get top locations: %w", err)
	}
	defer rows.Close()
	
	report.TopLocations = []string{}
	for rows.Next() {
		var location string
		if err := rows.Scan(&location); err != nil {
			continue
		}
		report.TopLocations = append(report.TopLocations, location)
	}
	
	return report, nil
}

// GetActivitySummary generates an activity summary
func (r *AuditRepository) GetActivitySummary(userID string, days int) (*models.ActivitySummary, error) {
	summary := &models.ActivitySummary{}
	startDate := time.Now().AddDate(0, 0, -days)
	
	// Get total logins
	query := `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1 AND event_type = 'user_login' AND success = true AND created_at >= $2`
	err := r.db.QueryRow(query, userID, startDate).Scan(&summary.TotalLogins)
	if err != nil {
		return nil, fmt.Errorf("failed to get total logins: %w", err)
	}
	
	// Get failed login attempts
	query = `SELECT COUNT(*) FROM audit_logs WHERE user_id = $1 AND event_type = 'login_failed' AND created_at >= $2`
	err = r.db.QueryRow(query, userID, startDate).Scan(&summary.FailedLoginAttempts)
	if err != nil {
		return nil, fmt.Errorf("failed to get failed logins: %w", err)
	}
	
	// Get unique devices
	query = `SELECT COUNT(DISTINCT device_id) FROM audit_logs WHERE user_id = $1 AND device_id IS NOT NULL AND created_at >= $2`
	err = r.db.QueryRow(query, userID, startDate).Scan(&summary.UniqueDevices)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique devices: %w", err)
	}
	
	// Get unique locations
	query = `SELECT COUNT(DISTINCT location_country) FROM audit_logs WHERE user_id = $1 AND location_country IS NOT NULL AND created_at >= $2`
	err = r.db.QueryRow(query, userID, startDate).Scan(&summary.UniqueLocations)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique locations: %w", err)
	}
	
	// Get daily activity
	query = `
		SELECT 
			DATE(created_at) as date,
			COUNT(CASE WHEN event_type = 'user_login' AND success = true THEN 1 END) as login_count,
			COUNT(CASE WHEN event_type = 'login_failed' THEN 1 END) as failed_login_count
		FROM audit_logs
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`
	rows, err := r.db.Query(query, userID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily activity: %w", err)
	}
	defer rows.Close()
	
	summary.DailyActivity = []models.DailyActivity{}
	for rows.Next() {
		var activity models.DailyActivity
		var date time.Time
		if err := rows.Scan(&date, &activity.LoginCount, &activity.FailedLoginCount); err != nil {
			continue
		}
		activity.Date = date.Format("2006-01-02")
		summary.DailyActivity = append(summary.DailyActivity, activity)
	}
	
	return summary, nil
}