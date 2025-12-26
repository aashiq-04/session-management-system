package compliance

// Top-level SOC2 overview
type SOC2Overview struct {
	OrganizationID string

	AccessControl     AccessControlMetrics
	Authentication    AuthenticationMetrics
	SystemOperations SystemOperationsMetrics
	SecurityIncidents SecurityIncidentMetrics
}

// CC1 – Control Environment
type AccessControlMetrics struct {
	TotalUsers     int
	AdminUsers     int
	LastRoleChange string
}

// CC6 – Logical Access
type AuthenticationMetrics struct {
	MFAEnabledPercentage float64
	FailedLogins30Days   int
	AuthzDenials30Days   int
}

// CC7 – System Operations
type SystemOperationsMetrics struct {
	ActiveSessions     int
	NewDevices30Days   int
	UnresolvedAlerts   int
}

// Security Incidents (SOC2 evidence)
type SecurityIncidentMetrics struct {
	HighSeverityAlerts int
	MediumSeverityAlerts int
}
