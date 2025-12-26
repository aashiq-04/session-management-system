package compliance

import "math"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSOC2Overview(orgID string) (*SOC2Overview, error) {

	totalUsers, err := s.repo.CountTotalUsers(orgID)
	if err != nil {
		return nil, err
	}

	adminUsers, err := s.repo.CountAdminUsers(orgID)
	if err != nil {
		return nil, err
	}

	failedLogins, _ := s.repo.CountFailedLoginsLast30Days(orgID)
	authzDenials, _ := s.repo.CountAuthzDenialsLast30Days(orgID)
	activeSessions, _ := s.repo.CountActiveSessions(orgID)
	unresolvedAlerts, _ := s.repo.CountUnresolvedAlerts(orgID)

	var mfaPercentage float64
	if totalUsers > 0 {
		// Placeholder – we’ll refine MFA query in Step 5B
		mfaPercentage = math.Round((float64(adminUsers) / float64(totalUsers)) * 100)
	}

	return &SOC2Overview{
		OrganizationID: orgID,
		AccessControl: AccessControlMetrics{
			TotalUsers: totalUsers,
			AdminUsers: adminUsers,
		},
		Authentication: AuthenticationMetrics{
			MFAEnabledPercentage: mfaPercentage,
			FailedLogins30Days:   failedLogins,
			AuthzDenials30Days:   authzDenials,
		},
		SystemOperations: SystemOperationsMetrics{
			ActiveSessions:   activeSessions,
			UnresolvedAlerts: unresolvedAlerts,
		},
		SecurityIncidents: SecurityIncidentMetrics{
			HighSeverityAlerts:   unresolvedAlerts, // refined later
			MediumSeverityAlerts: 0,
		},
	}, nil
}
