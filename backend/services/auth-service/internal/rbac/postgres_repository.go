package rbac

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

// Constructor used by main.go
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// HasPermission checks if a user has a permission in an organization
// FAIL-CLOSED: any error → deny
func (r *PostgresRepository) HasPermission(
	ctx context.Context,
	userID string,
	orgID string,
	permission string,
) (bool, error) {

	const query = `
		SELECT 1
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		JOIN role_permissions rp ON r.id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		WHERE ur.user_id = $1
		  AND ur.organization_id = $2
		  AND r.is_active = true
		  AND p.name = $3
		LIMIT 1;
	`

	var exists int
	err := r.db.QueryRowContext(ctx, query, userID, orgID, permission).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
