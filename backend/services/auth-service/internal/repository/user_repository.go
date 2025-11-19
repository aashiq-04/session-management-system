package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, is_active, mfa_enabled, mfa_secret, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.IsActive,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(userID string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, is_active, mfa_enabled, mfa_secret, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.IsActive,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, full_name = $2, is_active = $3, mfa_enabled = $4, 
		    mfa_secret = $5, updated_at = $6
		WHERE id = $7
	`
	
	_, err := r.db.Exec(
		query,
		user.Email,
		user.FullName,
		user.IsActive,
		user.MFAEnabled,
		user.MFASecret,
		time.Now(),
		user.ID,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

// EnableMFA enables multi-factor authentication for a user
func (r *UserRepository) EnableMFA(userID string, secret string) error {
	query := `
		UPDATE users
		SET mfa_enabled = true, mfa_secret = $1, updated_at = $2
		WHERE id = $3
	`
	
	_, err := r.db.Exec(query, secret, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to enable MFA: %w", err)
	}
	
	return nil
}

// CreateDevice inserts a new device into the database
func (r *UserRepository) CreateDevice(device *models.Device) error {
	query := `
		INSERT INTO devices (id, user_id, device_fingerprint, device_name, device_type, 
		                     os, browser, is_trusted, first_seen_at, last_seen_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.Exec(
		query,
		device.ID,
		device.UserID,
		device.DeviceFingerprint,
		device.DeviceName,
		device.DeviceType,
		device.OS,
		device.Browser,
		device.IsTrusted,
		device.FirstSeenAt,
		device.LastSeenAt,
		device.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}
	
	return nil
}

// GetDeviceByFingerprint retrieves a device by its fingerprint
func (r *UserRepository) GetDeviceByFingerprint(fingerprint string) (*models.Device, error) {
	query := `
		SELECT id, user_id, device_fingerprint, device_name, device_type, os, browser,
		       is_trusted, first_seen_at, last_seen_at, created_at
		FROM devices
		WHERE device_fingerprint = $1
	`
	
	device := &models.Device{}
	err := r.db.QueryRow(query, fingerprint).Scan(
		&device.ID,
		&device.UserID,
		&device.DeviceFingerprint,
		&device.DeviceName,
		&device.DeviceType,
		&device.OS,
		&device.Browser,
		&device.IsTrusted,
		&device.FirstSeenAt,
		&device.LastSeenAt,
		&device.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil // Device not found, but not an error
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	
	return device, nil
}

// UpdateDeviceLastSeen updates the last seen timestamp for a device
func (r *UserRepository) UpdateDeviceLastSeen(deviceID string) error {
	query := `UPDATE devices SET last_seen_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), deviceID)
	return err
}

// CreateSession creates a new session
func (r *UserRepository) CreateSession(session *models.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, device_id, refresh_token, ip_address, user_agent,
		                      location_country, location_city, latitude, longitude,
		                      is_active, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	
	_, err := r.db.Exec(
		query,
		session.ID,
		session.UserID,
		session.DeviceID,
		session.RefreshToken,
		session.IPAddress,
		session.UserAgent,
		session.LocationCountry,
		session.LocationCity,
		session.Latitude,
		session.Longitude,
		session.IsActive,
		session.ExpiresAt,
		session.CreatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	
	return nil
}

// GetSessionByRefreshToken retrieves a session by refresh token
func (r *UserRepository) GetSessionByRefreshToken(refreshToken string) (*models.Session, error) {
	query := `
		SELECT id, user_id, device_id, refresh_token, ip_address, user_agent,
		       location_country, location_city, latitude, longitude,
		       is_active, expires_at, created_at, revoked_at
		FROM sessions
		WHERE refresh_token = $1
	`
	
	session := &models.Session{}
	err := r.db.QueryRow(query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.RefreshToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.LocationCountry,
		&session.LocationCity,
		&session.Latitude,
		&session.Longitude,
		&session.IsActive,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.RevokedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	
	return session, nil
}

// CreateAuditLog creates an audit log entry
func (r *UserRepository) CreateAuditLog(log *models.AuditLog) error {
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