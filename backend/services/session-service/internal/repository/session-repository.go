package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/aashiq-04/session-management-system/backend/services/session-service/internal/models"
	// "github.com/google/uuid"
)

// SessionRepository handles database operations for sessions
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// GetUserSessions retrieves all sessions for a user
func (r *SessionRepository) GetUserSessions(userID string, includeInactive bool) ([]models.SessionWithDevice, error) {
	query := `
		SELECT 
			s.id, s.user_id, s.device_id, s.refresh_token, s.ip_address, 
			s.user_agent, s.location_country, s.location_city, s.latitude, s.longitude,
			s.is_active, s.expires_at, s.created_at, s.revoked_at,
			d.device_name, d.device_type, d.os, d.browser
		FROM sessions s
		LEFT JOIN devices d ON s.device_id = d.id
		WHERE s.user_id = $1
	`
	
	if !includeInactive {
		query += " AND s.is_active = true AND s.expires_at > NOW()"
	}
	
	query += " ORDER BY s.created_at DESC"
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()
	
	var sessions []models.SessionWithDevice
	for rows.Next() {
		var s models.SessionWithDevice
		err := rows.Scan(
			&s.ID, &s.UserID, &s.DeviceID, &s.RefreshToken, &s.IPAddress,
			&s.UserAgent, &s.LocationCountry, &s.LocationCity, &s.Latitude, &s.Longitude,
			&s.IsActive, &s.ExpiresAt, &s.CreatedAt, &s.RevokedAt,
			&s.DeviceName, &s.DeviceType, &s.OS, &s.Browser,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, s)
	}
	
	return sessions, nil
}

// GetSessionByID retrieves a specific session
func (r *SessionRepository) GetSessionByID(sessionID string) (*models.SessionWithDevice, error) {
	query := `
		SELECT 
			s.id, s.user_id, s.device_id, s.refresh_token, s.ip_address, 
			s.user_agent, s.location_country, s.location_city, s.latitude, s.longitude,
			s.is_active, s.expires_at, s.created_at, s.revoked_at,
			d.device_name, d.device_type, d.os, d.browser
		FROM sessions s
		LEFT JOIN devices d ON s.device_id = d.id
		WHERE s.id = $1
	`
	
	var s models.SessionWithDevice
	err := r.db.QueryRow(query, sessionID).Scan(
		&s.ID, &s.UserID, &s.DeviceID, &s.RefreshToken, &s.IPAddress,
		&s.UserAgent, &s.LocationCountry, &s.LocationCity, &s.Latitude, &s.Longitude,
		&s.IsActive, &s.ExpiresAt, &s.CreatedAt, &s.RevokedAt,
		&s.DeviceName, &s.DeviceType, &s.OS, &s.Browser,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	
	return &s, nil
}

// RevokeSession revokes a specific session
func (r *SessionRepository) RevokeSession(sessionID string) error {
	query := `
		UPDATE sessions
		SET is_active = false, revoked_at = $1
		WHERE id = $2
	`
	
	result, err := r.db.Exec(query, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}
	
	return nil
}

// RevokeAllSessions revokes all sessions for a user
func (r *SessionRepository) RevokeAllSessions(userID string, exceptSessionID string) (int64, error) {
	var query string
	var result sql.Result
	var err error
	
	if exceptSessionID != "" {
		query = `
			UPDATE sessions
			SET is_active = false, revoked_at = $1
			WHERE user_id = $2 AND id != $3 AND is_active = true
		`
		result, err = r.db.Exec(query, time.Now(), userID, exceptSessionID)
	} else {
		query = `
			UPDATE sessions
			SET is_active = false, revoked_at = $1
			WHERE user_id = $2 AND is_active = true
		`
		result, err = r.db.Exec(query, time.Now(), userID)
	}
	
	if err != nil {
		return 0, fmt.Errorf("failed to revoke sessions: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	return rowsAffected, nil
}

// GetUserDevices retrieves all devices for a user
func (r *SessionRepository) GetUserDevices(userID string) ([]models.Device, error) {
	query := `
		SELECT id, user_id, device_fingerprint, device_name, device_type, 
		       os, browser, is_trusted, first_seen_at, last_seen_at, created_at
		FROM devices
		WHERE user_id = $1
		ORDER BY last_seen_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	defer rows.Close()
	
	var devices []models.Device
	for rows.Next() {
		var d models.Device
		err := rows.Scan(
			&d.ID, &d.UserID, &d.DeviceFingerprint, &d.DeviceName, &d.DeviceType,
			&d.OS, &d.Browser, &d.IsTrusted, &d.FirstSeenAt, &d.LastSeenAt, &d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, d)
	}
	
	return devices, nil
}

// TrustDevice marks a device as trusted
func (r *SessionRepository) TrustDevice(deviceID string) error {
	query := `UPDATE devices SET is_trusted = true WHERE id = $1`
	
	result, err := r.db.Exec(query, deviceID)
	if err != nil {
		return fmt.Errorf("failed to trust device: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("device not found")
	}
	
	return nil
}

// GetSessionStats retrieves session statistics for a user
func (r *SessionRepository) GetSessionStats(userID string) (*models.SessionStats, error) {
	stats := &models.SessionStats{}
	
	// Get total and active session counts
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN is_active = true AND expires_at > NOW() THEN 1 END) as active
		FROM sessions
		WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(&stats.TotalSessions, &stats.ActiveSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to get session counts: %w", err)
	}
	
	// Get device counts
	query = `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN is_trusted = true THEN 1 END) as trusted
		FROM devices
		WHERE user_id = $1
	`
	err = r.db.QueryRow(query, userID).Scan(&stats.TotalDevices, &stats.TrustedDevices)
	if err != nil {
		return nil, fmt.Errorf("failed to get device counts: %w", err)
	}
	
	// Get last login info
	query = `
		SELECT created_at, 
		       COALESCE(location_city || ', ' || location_country, location_country, 'Unknown') as location
		FROM sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	err = r.db.QueryRow(query, userID).Scan(&stats.LastLogin, &stats.LastLoginLocation)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get last login: %w", err)
	}
	
	// Get recent locations (last 5 unique locations)
	query = `
		SELECT location
		FROM (
			SELECT DISTINCT ON (location_country, location_city)
				COALESCE(location_city || ', ' || location_country, location_country, 'Unknown') as location,
				created_at
			FROM sessions
			WHERE user_id = $1 AND location_country IS NOT NULL
			ORDER BY location_country, location_city, created_at DESC
		) AS unique_locations
		ORDER BY created_at DESC
		LIMIT 5
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent locations: %w", err)
	}
	defer rows.Close()
	
	stats.RecentLocations = []string{}
	for rows.Next() {
		var location string
		if err := rows.Scan(&location); err != nil {
			continue
		}
		stats.RecentLocations = append(stats.RecentLocations, location)
	}
	
	return stats, nil
}

// CreateAuditLog creates an audit log entry
func (r *SessionRepository) CreateAuditLog(log *models.AuditLog) error {
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