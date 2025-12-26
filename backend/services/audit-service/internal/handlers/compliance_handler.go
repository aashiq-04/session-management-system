package handlers

import (
	"context"

	"github.com/aashiq-04/session-management-system/backend/services/audit-service/internal/compliance"
	compb "github.com/aashiq-04/session-management-system/backend/services/audit-service/proto/compliance"
)

type ComplianceHandler struct {
	compb.UnimplementedComplianceServiceServer
	service *compliance.Service
}

func NewComplianceHandler(service *compliance.Service) *ComplianceHandler {
	return &ComplianceHandler{service: service}
}

func (h *ComplianceHandler) GetSOC2Overview(
	ctx context.Context,
	req *compb.SOC2Request,
) (*compb.SOC2OverviewResponse, error) {

	overview, err := h.service.GetSOC2Overview(req.OrganizationId)
	if err != nil {
		return nil, err
	}

	return &compb.SOC2OverviewResponse{
		OrganizationId: overview.OrganizationID,
		AccessControl: &compb.AccessControlMetrics{
			TotalUsers: int32(overview.AccessControl.TotalUsers),
			AdminUsers: int32(overview.AccessControl.AdminUsers),
		},
		Authentication: &compb.AuthenticationMetrics{
			MfaEnabledPercentage: overview.Authentication.MFAEnabledPercentage,
			FailedLogins_30Days:  int32(overview.Authentication.FailedLogins30Days),
			AuthzDenials_30Days:  int32(overview.Authentication.AuthzDenials30Days),
		},
		SystemOperations: &compb.SystemOperationsMetrics{
			ActiveSessions:   int32(overview.SystemOperations.ActiveSessions),
			UnresolvedAlerts: int32(overview.SystemOperations.UnresolvedAlerts),
		},
		SecurityIncidents: &compb.SecurityIncidentMetrics{
			HighSeverityAlerts:   int32(overview.SecurityIncidents.HighSeverityAlerts),
			MediumSeverityAlerts: int32(overview.SecurityIncidents.MediumSeverityAlerts),
		},
	}, nil
}
