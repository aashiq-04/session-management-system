package handlers

import (
	"context"

	authzpb "github.com/aashiq-04/session-management-system/backend/services/auth-service/proto/authorization"
	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/rbac"
)

type AuthorizationHandler struct {
	authzpb.UnimplementedAuthorizationServiceServer
	rbacService *rbac.Service
}

func NewAuthorizationHandler(rbacService *rbac.Service) *AuthorizationHandler {
	return &AuthorizationHandler{
		rbacService: rbacService,
	}
}

func (h *AuthorizationHandler) CanUserPerform(
	ctx context.Context,
	req *authzpb.AuthorizationRequest,
) (*authzpb.AuthorizationResponse, error) {

	allowed, err := h.rbacService.Authorize(
		ctx,
		req.UserId,
		req.OrganizationId,
		req.Permission,
	)

	if err != nil {
		// Fail-closed: any error → denied
		return &authzpb.AuthorizationResponse{
			Allowed: false,
		}, nil
	}

	return &authzpb.AuthorizationResponse{
		Allowed: allowed,
	}, nil
}
