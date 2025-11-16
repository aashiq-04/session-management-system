package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	IsActive     bool      `db:"is_active"`
	MFAEnabled   bool      `db:"mfa_enabled"`
	MFASecret    *string   `db:"mfa_secret"` // pointer to handle NULL
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Device represents a device that has accessed the system
type Device struct {
	ID                string    `db:"id"`
	UserID            string    `db:"user_id"`
	DeviceFingerprint string    `db:"device_fingerprint"`
	DeviceName        *string   `db:"device_name"`
	DeviceType        *string   `db:"device_type"`
	OS                *string   `db:"os"`
	Browser           *string   `db:"browser"`
	IsTrusted         bool      `db:"is_trusted"`
	FirstSeenAt       time.Time `db:"first_seen_at"`
	LastSeenAt        time.Time `db:"last_seen_at"`
	CreatedAt         time.Time `db:"created_at"`
}

// Session represents an active user session
type Session struct {
	ID              string     `db:"id"`
	UserID          string     `db:"user_id"`
	DeviceID        string     `db:"device_id"`
	RefreshToken    string     `db:"refresh_token"`
	IPAddress       string     `db:"ip_address"`
	UserAgent       *string    `db:"user_agent"`
	LocationCountry *string    `db:"location_country"`
	LocationCity    *string    `db:"location_city"`
	Latitude        *float64   `db:"latitude"`
	Longitude       *float64   `db:"longitude"`
	IsActive        bool       `db:"is_active"`
	ExpiresAt       time.Time  `db:"expires_at"`
	CreatedAt       time.Time  `db:"created_at"`
	RevokedAt       *time.Time `db:"revoked_at"`
}

// AuditLog represents a security event in the system
type AuditLog struct {
	ID              string     `db:"id"`
	UserID          *string    `db:"user_id"`
	SessionID       *string    `db:"session_id"`
	DeviceID        *string    `db:"device_id"`
	EventType       string     `db:"event_type"`
	EventCategory   string     `db:"event_category"`
	Severity        string     `db:"severity"`
	IPAddress       *string    `db:"ip_address"`
	UserAgent       *string    `db:"user_agent"`
	LocationCountry *string    `db:"location_country"`
	LocationCity    *string    `db:"location_city"`
	Metadata        *string    `db:"metadata"` // JSON string
	Success         bool       `db:"success"`
	FailureReason   *string    `db:"failure_reason"`
	CreatedAt       time.Time  `db:"created_at"`
}

// MFABackupCode represents a backup code for MFA recovery
type MFABackupCode struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	CodeHash  string     `db:"code_hash"`
	IsUsed    bool       `db:"is_used"`
	UsedAt    *time.Time `db:"used_at"`
	CreatedAt time.Time  `db:"created_at"`
}