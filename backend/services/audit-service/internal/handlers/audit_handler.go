package handlers

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	pb "github.com/aashiq-04/session-management-system/backend/services/audit-service/proto"
	"github.com/aashiq-04/session-management-system/backend/services/audit-service/internal/models"
	"github.com/aashiq-04/session-management-system/backend/services/audit-service/internal/repository"
)

// AuditHandler implements the AuditService gRPC service
type AuditHandler struct {
	pb.UnimplementedAuditServiceServer
	repo *repository.AuditRepository
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(db *sql.DB) *AuditHandler {
	return &AuditHandler{
		repo: repository.NewAuditRepository(db),
	}
}

// CreateAuditLog creates a new audit log entry
func (h *AuditHandler) CreateAuditLog(ctx context.Context, req *pb.CreateAuditLogRequest) (*pb.CreateAuditLogResponse, error) {
	log.Printf("CreateAuditLog request received for event: %s", req.EventType)

	logID := uuid.New().String()
	
	auditLog := &models.AuditLog{
		ID:              logID,
		UserID:          stringToPointer(req.UserId),
		SessionID:       stringToPointer(req.SessionId),
		DeviceID:        stringToPointer(req.DeviceId),
		EventType:       req.EventType,
		EventCategory:   req.EventCategory,
		Severity:        req.Severity,
		IPAddress:       stringToPointer(req.IpAddress),
		UserAgent:       stringToPointer(req.UserAgent),
		LocationCountry: stringToPointer(req.LocationCountry),
		LocationCity:    stringToPointer(req.LocationCity),
		Metadata:        stringToPointer(req.Metadata),
		Success:         req.Success,
		FailureReason:   stringToPointer(req.FailureReason),
		CreatedAt:       time.Now(),
	}

	err := h.repo.CreateAuditLog(auditLog)
	if err != nil {
		log.Printf("Failed to create audit log: %v", err)
		return &pb.CreateAuditLogResponse{
			Success: false,
			Message: "Failed to create audit log",
		}, nil
	}

	return &pb.CreateAuditLogResponse{
		Success: true,
		Message: "Audit log created successfully",
		LogId:   logID,
	}, nil
}

// GetUserAuditLogs retrieves audit logs for a user
func (h *AuditHandler) GetUserAuditLogs(ctx context.Context, req *pb.GetUserAuditLogsRequest) (*pb.GetUserAuditLogsResponse, error) {
	log.Printf("GetUserAuditLogs request received for user: %s", req.UserId)

	limit := int(req.Limit)
	if limit == 0 {
		limit = 50 // Default limit
	}

	offset := int(req.Offset)

	logs, totalCount, err := h.repo.GetUserAuditLogs(
		req.UserId,
		limit,
		offset,
		req.EventCategory,
		req.Severity,
		req.SuccessOnly,
	)

	if err != nil {
		log.Printf("Failed to get user audit logs: %v", err)
		return &pb.GetUserAuditLogsResponse{
			Success: false,
			Message: "Failed to retrieve audit logs",
		}, nil
	}

	// Convert to protobuf format
	pbLogs := make([]*pb.AuditLog, 0, len(logs))
	for _, l := range logs {
		pbLog := &pb.AuditLog{
			Id:              l.ID,
			UserId:          pointerToString(l.UserID),
			SessionId:       pointerToString(l.SessionID),
			DeviceId:        pointerToString(l.DeviceID),
			EventType:       l.EventType,
			EventCategory:   l.EventCategory,
			Severity:        l.Severity,
			IpAddress:       pointerToString(l.IPAddress),
			UserAgent:       pointerToString(l.UserAgent),
			LocationCountry: pointerToString(l.LocationCountry),
			LocationCity:    pointerToString(l.LocationCity),
			Metadata:        pointerToString(l.Metadata),
			Success:         l.Success,
			FailureReason:   pointerToString(l.FailureReason),
			CreatedAt:       l.CreatedAt.Format(time.RFC3339),
		}
		pbLogs = append(pbLogs, pbLog)
	}

	return &pb.GetUserAuditLogsResponse{
		Success:    true,
		Message:    "Audit logs retrieved successfully",
		Logs:       pbLogs,
		TotalCount: int32(totalCount),
	}, nil
}

// GetAuditLogsByEvent retrieves audit logs by event type
func (h *AuditHandler) GetAuditLogsByEvent(ctx context.Context, req *pb.GetAuditLogsByEventRequest) (*pb.GetAuditLogsByEventResponse, error) {
	log.Printf("GetAuditLogsByEvent request received for event: %s", req.EventType)

	limit := int(req.Limit)
	if limit == 0 {
		limit = 50
	}

	offset := int(req.Offset)

	// Parse dates
	startDate := time.Now().AddDate(0, 0, -30) // Default: 30 days ago
	if req.StartDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.StartDate)
		if err == nil {
			startDate = parsed
		}
	}

	endDate := time.Now()
	if req.EndDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.EndDate)
		if err == nil {
			endDate = parsed
		}
	}

	logs, totalCount, err := h.repo.GetAuditLogsByEvent(req.EventType, limit, offset, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get audit logs by event: %v", err)
		return &pb.GetAuditLogsByEventResponse{
			Success: false,
			Message: "Failed to retrieve audit logs",
		}, nil
	}

	// Convert to protobuf format
	pbLogs := make([]*pb.AuditLog, 0, len(logs))
	for _, l := range logs {
		pbLog := &pb.AuditLog{
			Id:              l.ID,
			UserId:          pointerToString(l.UserID),
			SessionId:       pointerToString(l.SessionID),
			DeviceId:        pointerToString(l.DeviceID),
			EventType:       l.EventType,
			EventCategory:   l.EventCategory,
			Severity:        l.Severity,
			IpAddress:       pointerToString(l.IPAddress),
			UserAgent:       pointerToString(l.UserAgent),
			LocationCountry: pointerToString(l.LocationCountry),
			LocationCity:    pointerToString(l.LocationCity),
			Metadata:        pointerToString(l.Metadata),
			Success:         l.Success,
			FailureReason:   pointerToString(l.FailureReason),
			CreatedAt:       l.CreatedAt.Format(time.RFC3339),
		}
		pbLogs = append(pbLogs, pbLog)
	}

	return &pb.GetAuditLogsByEventResponse{
		Success:    true,
		Message:    "Audit logs retrieved successfully",
		Logs:       pbLogs,
		TotalCount: int32(totalCount),
	}, nil
}

// GetSecurityAlerts retrieves security alerts for a user
func (h *AuditHandler) GetSecurityAlerts(ctx context.Context, req *pb.GetSecurityAlertsRequest) (*pb.GetSecurityAlertsResponse, error) {
	log.Printf("GetSecurityAlerts request received for user: %s", req.UserId)

	alerts, err := h.repo.GetSecurityAlerts(req.UserId, req.IncludeResolved, req.Severity)
	if err != nil {
		log.Printf("Failed to get security alerts: %v", err)
		return &pb.GetSecurityAlertsResponse{
			Success: false,
			Message: "Failed to retrieve security alerts",
		}, nil
	}

	// Convert to protobuf format
	pbAlerts := make([]*pb.SecurityAlert, 0, len(alerts))
	unresolvedCount := 0

	for _, a := range alerts {
		if !a.IsResolved {
			unresolvedCount++
		}

		resolvedAt := ""
		if a.ResolvedAt != nil {
			resolvedAt = a.ResolvedAt.Format(time.RFC3339)
		}

		pbAlert := &pb.SecurityAlert{
			Id:              a.ID,
			UserId:          a.UserID,
			AlertType:       a.AlertType,
			Severity:        a.Severity,
			Description:     a.Description,
			Metadata:        pointerToString(a.Metadata),
			IpAddress:       pointerToString(a.IPAddress),
			LocationCountry: pointerToString(a.LocationCountry),
			LocationCity:    pointerToString(a.LocationCity),
			IsResolved:      a.IsResolved,
			ResolvedAt:      resolvedAt,
			CreatedAt:       a.CreatedAt.Format(time.RFC3339),
		}
		pbAlerts = append(pbAlerts, pbAlert)
	}

	return &pb.GetSecurityAlertsResponse{
		Success:         true,
		Message:         "Security alerts retrieved successfully",
		Alerts:          pbAlerts,
		TotalCount:      int32(len(alerts)),
		UnresolvedCount: int32(unresolvedCount),
	}, nil
}

// CreateSecurityAlert creates a new security alert
func (h *AuditHandler) CreateSecurityAlert(ctx context.Context, req *pb.CreateSecurityAlertRequest) (*pb.CreateSecurityAlertResponse, error) {
	log.Printf("CreateSecurityAlert request received: %s", req.AlertType)

	alertID := uuid.New().String()

	alert := &models.SecurityAlert{
		ID:              alertID,
		UserID:          req.UserId,
		AlertType:       req.AlertType,
		Severity:        req.Severity,
		Description:     req.Description,
		Metadata:        stringToPointer(req.Metadata),
		IPAddress:       stringToPointer(req.IpAddress),
		LocationCountry: stringToPointer(req.LocationCountry),
		LocationCity:    stringToPointer(req.LocationCity),
		IsResolved:      false,
		CreatedAt:       time.Now(),
	}

	err := h.repo.CreateSecurityAlert(alert)
	if err != nil {
		log.Printf("Failed to create security alert: %v", err)
		return &pb.CreateSecurityAlertResponse{
			Success: false,
			Message: "Failed to create security alert",
		}, nil
	}

	return &pb.CreateSecurityAlertResponse{
		Success: true,
		Message: "Security alert created successfully",
		AlertId: alertID,
	}, nil
}

// ResolveSecurityAlert marks a security alert as resolved
func (h *AuditHandler) ResolveSecurityAlert(ctx context.Context, req *pb.ResolveSecurityAlertRequest) (*pb.ResolveSecurityAlertResponse, error) {
	log.Printf("ResolveSecurityAlert request received: %s", req.AlertId)

	err := h.repo.ResolveSecurityAlert(req.AlertId)
	if err != nil {
		log.Printf("Failed to resolve security alert: %v", err)
		return &pb.ResolveSecurityAlertResponse{
			Success: false,
			Message: "Failed to resolve security alert",
		}, nil
	}

	return &pb.ResolveSecurityAlertResponse{
		Success: true,
		Message: "Security alert resolved successfully",
	}, nil
}

// GetComplianceReport generates a compliance report
func (h *AuditHandler) GetComplianceReport(ctx context.Context, req *pb.GetComplianceReportRequest) (*pb.GetComplianceReportResponse, error) {
	log.Printf("GetComplianceReport request received for user: %s", req.UserId)

	// Parse dates
	startDate := time.Now().AddDate(0, 0, -30) // Default: 30 days ago
	if req.StartDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.StartDate)
		if err == nil {
			startDate = parsed
		}
	}

	endDate := time.Now()
	if req.EndDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.EndDate)
		if err == nil {
			endDate = parsed
		}
	}

	report, err := h.repo.GetComplianceReport(req.UserId, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get compliance report: %v", err)
		return &pb.GetComplianceReportResponse{
			Success: false,
			Message: "Failed to generate compliance report",
		}, nil
	}

	// Convert event breakdown
	eventBreakdown := make([]*pb.EventCount, 0, len(report.EventBreakdown))
	for _, ec := range report.EventBreakdown {
		eventBreakdown = append(eventBreakdown, &pb.EventCount{
			EventType: ec.EventType,
			Count:     int32(ec.Count),
		})
	}

	return &pb.GetComplianceReportResponse{
		Success:            true,
		Message:            "Compliance report generated successfully",
		TotalEvents:        int32(report.TotalEvents),
		SuccessfulLogins:   int32(report.SuccessfulLogins),
		FailedLogins:       int32(report.FailedLogins),
		SessionRevocations: int32(report.SessionRevocations),
		SecurityAlerts:     int32(report.SecurityAlerts),
		MfaEvents:          int32(report.MFAEvents),
		EventBreakdown:     eventBreakdown,
		TopLocations:       report.TopLocations,
	}, nil
}

// GetActivitySummary retrieves activity summary
func (h *AuditHandler) GetActivitySummary(ctx context.Context, req *pb.GetActivitySummaryRequest) (*pb.GetActivitySummaryResponse, error) {
	log.Printf("GetActivitySummary request received for user: %s", req.UserId)

	days := int(req.Days)
	if days == 0 {
		days = 30 // Default: 30 days
	}

	summary, err := h.repo.GetActivitySummary(req.UserId, days)
	if err != nil {
		log.Printf("Failed to get activity summary: %v", err)
		return &pb.GetActivitySummaryResponse{
			Success: false,
			Message: "Failed to retrieve activity summary",
		}, nil
	}

	// Convert daily activity
	dailyActivity := make([]*pb.DailyActivity, 0, len(summary.DailyActivity))
	for _, da := range summary.DailyActivity {
		dailyActivity = append(dailyActivity, &pb.DailyActivity{
			Date:             da.Date,
			LoginCount:       int32(da.LoginCount),
			FailedLoginCount: int32(da.FailedLoginCount),
		})
	}

	return &pb.GetActivitySummaryResponse{
		Success:             true,
		Message:             "Activity summary retrieved successfully",
		TotalLogins:         int32(summary.TotalLogins),
		FailedLoginAttempts: int32(summary.FailedLoginAttempts),
		UniqueDevices:       int32(summary.UniqueDevices),
		UniqueLocations:     int32(summary.UniqueLocations),
		DailyActivity:       dailyActivity,
	}, nil
}

// Helper functions
func stringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func pointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}