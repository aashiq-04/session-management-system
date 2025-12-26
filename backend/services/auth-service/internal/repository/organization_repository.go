package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/models"
)

type OrganizationRepository struct {
	db *sql.DB
}

// Organization represents an organization in the system


func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// GetPrimaryOrganizationID returns the primary organization for a user
func (r *OrganizationRepository) GetPrimaryOrganizationID(userID string) (string, error) {
	const query = `
		SELECT organization_id
		FROM user_organizations
		WHERE user_id = $1
		  AND is_primary = true
		LIMIT 1;
	`

	var orgID string
	err := r.db.QueryRow(query, userID).Scan(&orgID)
	if err == sql.ErrNoRows {
		return "", errors.New("no primary organization found")
	}
	if err != nil {
		return "", err
	}

	return orgID, nil
}

// CreateOrganization creates a new organization
func (r *OrganizationRepository) CreateOrganization(org *models.Organization) error {
	const query = `
		INSERT INTO organizations (id, name, slug, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, org.ID, org.Name, org.Slug, org.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	return nil
}

// AddUserToOrganization adds a user to an organization
func (r *OrganizationRepository) AddUserToOrganization(userID, orgID string, isPrimary bool) error {
	const query = `
		INSERT INTO user_organizations (user_id, organization_id, is_primary, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, userID, orgID, isPrimary, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}

	return nil
}
