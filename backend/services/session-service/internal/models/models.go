package models

import (
	"time"
)

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

// SessionWithDevice combines session and device information
type SessionWithDevice struct {
	Session
	DeviceName *string `db:"device_name"`
	DeviceType *string `db:"device_type"`
	OS         *string `db:"os"`
	Browser    *string `db:"browser"`
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

// SessionStats represents session statistics
type SessionStats struct {
	TotalSessions    int
	ActiveSessions   int
	TotalDevices     int
	TrustedDevices   int
	LastLogin        *time.Time
	LastLoginLocation *string
	RecentLocations  []string
}