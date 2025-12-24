package rbac

import "context"

type Repository interface {
	HasPermission(
		ctx context.Context,
		userID string,
		orgID string,
		permission string,
	) (bool, error)
}
