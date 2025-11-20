package models

import (
	"time"
)

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

// SecurityAlert represents a detected security anomaly
type SecurityAlert struct {
	ID              string     `db:"id"`
	UserID          string     `db:"user_id"`
	AlertType       string     `db:"alert_type"`
	Severity        string     `db:"severity"`
	Description     string     `db:"description"`
	Metadata        *string    `db:"metadata"` // JSON string
	IPAddress       *string    `db:"ip_address"`
	LocationCountry *string    `db:"location_country"`
	LocationCity    *string    `db:"location_city"`
	IsResolved      bool       `db:"is_resolved"`
	ResolvedAt      *time.Time `db:"resolved_at"`
	CreatedAt       time.Time  `db:"created_at"`
}

// EventCount represents count of events by type
type EventCount struct {
	EventType string `db:"event_type"`
	Count     int    `db:"count"`
}

// DailyActivity represents daily login activity
type DailyActivity struct {
	Date             string `db:"date"`
	LoginCount       int    `db:"login_count"`
	FailedLoginCount int    `db:"failed_login_count"`
}

// ComplianceReport represents a compliance summary
type ComplianceReport struct {
	TotalEvents        int
	SuccessfulLogins   int
	FailedLogins       int
	SessionRevocations int
	SecurityAlerts     int
	MFAEvents          int
	EventBreakdown     []EventCount
	TopLocations       []string
}

// ActivitySummary represents user activity summary
type ActivitySummary struct {
	TotalLogins         int
	FailedLoginAttempts int
	UniqueDevices       int
	UniqueLocations     int
	DailyActivity       []DailyActivity
}