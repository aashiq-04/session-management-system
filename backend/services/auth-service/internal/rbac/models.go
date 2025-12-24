package rbac

import "time"

// Role represents an org-scoped role
type Role struct {
	ID             string
	OrganizationID string
	Name           string
	IsActive       bool
	IsSystem       bool
	CreatedAt      time.Time
}

// Permission represents a global permission
type Permission struct {
	ID          string
	Name        string
	Description string
	CreatedAt  time.Time
}

// UserRole is the assignment of a role to a user within an org
type UserRole struct {
	ID             string
	UserID         string
	RoleID         string
	OrganizationID string
	AssignedBy     *string
	AssignedAt     time.Time
}
