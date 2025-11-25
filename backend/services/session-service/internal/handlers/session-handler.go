package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	pb "github.com/aashiq-04/session-management-system/backend/services/session-service/proto"
	"github.com/aashiq-04/session-management-system/backend/services/session-service/internal/models"
	"github.com/aashiq-04/session-management-system/backend/services/session-service/internal/repository"
)

// SessionHandler implements the SessionService gRPC service
type SessionHandler struct {
	pb.UnimplementedSessionServiceServer
	repo *repository.SessionRepository
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(db *sql.DB) *SessionHandler {
	return &SessionHandler{
		repo: repository.NewSessionRepository(db),
	}
}

// GetUserSessions retrieves all sessions for a user
// GetUserSessions retrieves all sessions for a user
func (h *SessionHandler) GetUserSessions(ctx context.Context, req *pb.GetUserSessionsRequest) (*pb.GetUserSessionsResponse, error) {
    log.Printf("GetUserSessions request received for user: %s", req.UserId)

    sessions, err := h.repo.GetUserSessions(req.UserId, req.IncludeInactive)
    if err != nil {
        log.Printf("Failed to get user sessions: %v", err)
        return &pb.GetUserSessionsResponse{
            Success: false,
            Message: "Failed to retrieve sessions",
        }, nil
    }

    pbSessions := make([]*pb.Session, 0, len(sessions))
    activeCount := 0

    for _, s := range sessions {

        // Determine active status
        isActive := s.IsActive && time.Now().Before(s.ExpiresAt)
        if isActive {
            activeCount++
        }

        // Device name fallback
        deviceName := "Unknown Device"
        if s.DeviceName != nil && *s.DeviceName != "" {
            deviceName = *s.DeviceName
        }

        // Device type fallback
        deviceType := "unknown"
        if s.DeviceType != nil && *s.DeviceType != "" {
            deviceType = *s.DeviceType
        }

        // -------- SAFE LOCATION MERGE ----------
        var finalLocationCountry string
        var finalLocationCity string

        if s.LocationCountry != nil {
            finalLocationCountry = *s.LocationCountry
        } else {
            finalLocationCountry = ""
        }

        if s.LocationCity != nil && *s.LocationCity != "" && s.LocationCountry != nil {
            finalLocationCity = fmt.Sprintf("%s, %s", *s.LocationCity, *s.LocationCountry)
        } else if s.LocationCountry != nil {
            finalLocationCity = *s.LocationCountry
        } else {
            finalLocationCity = ""
        }

        // Last seen logic
        lastSeen := s.CreatedAt
        if s.RevokedAt != nil {
            lastSeen = *s.RevokedAt
        }

        pbSessions = append(pbSessions, &pb.Session{
            Id:              s.ID,
            UserId:          s.UserID,
            DeviceId:        s.DeviceID,
            DeviceName:      deviceName,
            DeviceType:      deviceType,
            IpAddress:       s.IPAddress,
            UserAgent:       getStringValue(s.UserAgent),

            // Location safely mapped
            LocationCountry: finalLocationCountry,
            LocationCity:    finalLocationCity,
            Latitude:        getFloat64Value(s.Latitude),
            Longitude:       getFloat64Value(s.Longitude),

            IsActive:        isActive,
            CreatedAt:       s.CreatedAt.Format(time.RFC3339),
            LastSeenAt:      lastSeen.Format(time.RFC3339),
            ExpiresAt:       s.ExpiresAt.Format(time.RFC3339),
            IsCurrent:       false,
        })
    }

    return &pb.GetUserSessionsResponse{
        Success:     true,
        Message:     "Sessions retrieved successfully",
        Sessions:    pbSessions,
        TotalCount:  int32(len(sessions)),
        ActiveCount: int32(activeCount),
    }, nil
}

// GetSessionDetails retrieves details of a specific session
func (h *SessionHandler) GetSessionDetails(ctx context.Context, req *pb.GetSessionDetailsRequest) (*pb.GetSessionDetailsResponse, error) {
	log.Printf("GetSessionDetails request received for session: %s", req.SessionId)

	session, err := h.repo.GetSessionByID(req.SessionId)
	if err != nil {
		log.Printf("Failed to get session details: %v", err)
		return &pb.GetSessionDetailsResponse{
			Success: false,
			Message: "Session not found",
		}, nil
	}

	// Verify user owns this session
	if session.UserID != req.UserId {
		return &pb.GetSessionDetailsResponse{
			Success: false,
			Message: "Unauthorized",
		}, nil
	}

	// Check if session is active
	isActive := session.IsActive && time.Now().Before(session.ExpiresAt)

	deviceName := "Unknown Device"
	if session.DeviceName != nil && *session.DeviceName != "" {
		deviceName = *session.DeviceName
	}

	deviceType := "unknown"
	if session.DeviceType != nil {
		deviceType = *session.DeviceType
	}

	pbSession := &pb.Session{
		Id:              session.ID,
		UserId:          session.UserID,
		DeviceId:        session.DeviceID,
		DeviceName:      deviceName,
		DeviceType:      deviceType,
		IpAddress:       session.IPAddress,
		UserAgent:       getStringValue(session.UserAgent),
		LocationCountry: getStringValue(session.LocationCountry),
		LocationCity:    getStringValue(session.LocationCity),
		Latitude:        getFloat64Value(session.Latitude),
		Longitude:       getFloat64Value(session.Longitude),
		IsActive:        isActive,
		CreatedAt:       session.CreatedAt.Format(time.RFC3339),
		ExpiresAt:       session.ExpiresAt.Format(time.RFC3339),
	}

	return &pb.GetSessionDetailsResponse{
		Success: true,
		Message: "Session details retrieved",
		Session: pbSession,
	}, nil
}

// RevokeSession revokes a specific session
func (h *SessionHandler) RevokeSession(ctx context.Context, req *pb.RevokeSessionRequest) (*pb.RevokeSessionResponse, error) {
	log.Printf("RevokeSession request received for session: %s", req.SessionId)

	// Get session to verify ownership
	session, err := h.repo.GetSessionByID(req.SessionId)
	if err != nil {
		return &pb.RevokeSessionResponse{
			Success: false,
			Message: "Session not found",
		}, nil
	}

	// Verify user owns this session
	if session.UserID != req.UserId {
		return &pb.RevokeSessionResponse{
			Success: false,
			Message: "Unauthorized",
		}, nil
	}

	// Revoke the session
	err = h.repo.RevokeSession(req.SessionId)
	if err != nil {
		log.Printf("Failed to revoke session: %v", err)
		return &pb.RevokeSessionResponse{
			Success: false,
			Message: "Failed to revoke session",
		}, nil
	}

	// Create audit log
	h.createAuditLog(&models.AuditLog{
		ID:            uuid.New().String(),
		UserID:        &req.UserId,
		SessionID:     &req.SessionId,
		EventType:     "session_revoked",
		EventCategory: "session_management",
		Severity:      "info",
		IPAddress:     &req.RevokedByIp,
		Success:       true,
		CreatedAt:     time.Now(),
	})

	return &pb.RevokeSessionResponse{
		Success: true,
		Message: "Session revoked successfully",
	}, nil
}

// RevokeAllSessions revokes all sessions for a user
func (h *SessionHandler) RevokeAllSessions(ctx context.Context, req *pb.RevokeAllSessionsRequest) (*pb.RevokeAllSessionsResponse, error) {
	log.Printf("RevokeAllSessions request received for user: %s", req.UserId)

	count, err := h.repo.RevokeAllSessions(req.UserId, req.ExceptSessionId)
	if err != nil {
		log.Printf("Failed to revoke all sessions: %v", err)
		return &pb.RevokeAllSessionsResponse{
			Success: false,
			Message: "Failed to revoke sessions",
		}, nil
	}

	// Create audit log
	h.createAuditLog(&models.AuditLog{
		ID:            uuid.New().String(),
		UserID:        &req.UserId,
		EventType:     "all_sessions_revoked",
		EventCategory: "session_management",
		Severity:      "warning",
		IPAddress:     &req.RevokedByIp,
		Success:       true,
		CreatedAt:     time.Now(),
	})

	return &pb.RevokeAllSessionsResponse{
		Success:      true,
		Message:      fmt.Sprintf("Revoked %d session(s)", count),
		RevokedCount: int32(count),
	}, nil
}

// GetUserDevices retrieves all devices for a user
func (h *SessionHandler) GetUserDevices(ctx context.Context, req *pb.GetUserDevicesRequest) (*pb.GetUserDevicesResponse, error) {
	log.Printf("GetUserDevices request received for user: %s", req.UserId)

	devices, err := h.repo.GetUserDevices(req.UserId)
	if err != nil {
		log.Printf("Failed to get user devices: %v", err)
		return &pb.GetUserDevicesResponse{
			Success: false,
			Message: "Failed to retrieve devices",
		}, nil
	}

	// Convert to protobuf format
	pbDevices := make([]*pb.Device, 0, len(devices))
	trustedCount := 0

	for _, d := range devices {
		if d.IsTrusted {
			trustedCount++
		}

		deviceName := "Unknown Device"
		if d.DeviceName != nil && *d.DeviceName != "" {
			deviceName = *d.DeviceName
		}

		pbDevice := &pb.Device{
			Id:                d.ID,
			UserId:            d.UserID,
			DeviceFingerprint: d.DeviceFingerprint,
			DeviceName:        deviceName,
			DeviceType:        getStringValue(d.DeviceType),
			Os:                getStringValue(d.OS),
			Browser:           getStringValue(d.Browser),
			IsTrusted:         d.IsTrusted,
			FirstSeenAt:       d.FirstSeenAt.Format(time.RFC3339),
			LastSeenAt:        d.LastSeenAt.Format(time.RFC3339),
		}

		pbDevices = append(pbDevices, pbDevice)
	}

	return &pb.GetUserDevicesResponse{
		Success:      true,
		Message:      "Devices retrieved successfully",
		Devices:      pbDevices,
		TotalCount:   int32(len(devices)),
		TrustedCount: int32(trustedCount),
	}, nil
}

// TrustDevice marks a device as trusted
func (h *SessionHandler) TrustDevice(ctx context.Context, req *pb.TrustDeviceRequest) (*pb.TrustDeviceResponse, error) {
	log.Printf("TrustDevice request received for device: %s", req.DeviceId)

	err := h.repo.TrustDevice(req.DeviceId)
	if err != nil {
		log.Printf("Failed to trust device: %v", err)
		return &pb.TrustDeviceResponse{
			Success: false,
			Message: "Failed to trust device",
		}, nil
	}

	// Create audit log
	h.createAuditLog(&models.AuditLog{
		ID:            uuid.New().String(),
		UserID:        &req.UserId,
		DeviceID:      &req.DeviceId,
		EventType:     "device_trusted",
		EventCategory: "security",
		Severity:      "info",
		Success:       true,
		CreatedAt:     time.Now(),
	})

	return &pb.TrustDeviceResponse{
		Success: true,
		Message: "Device trusted successfully",
	}, nil
}

// GetSessionStats retrieves session statistics for a user
func (h *SessionHandler) GetSessionStats(ctx context.Context, req *pb.GetSessionStatsRequest) (*pb.GetSessionStatsResponse, error) {
	log.Printf("GetSessionStats request received for user: %s", req.UserId)

	stats, err := h.repo.GetSessionStats(req.UserId)
	if err != nil {
		log.Printf("Failed to get session stats: %v", err)
		return &pb.GetSessionStatsResponse{
			Success: false,
			Message: "Failed to retrieve statistics",
		}, nil
	}

	lastLogin := ""
	if stats.LastLogin != nil {
		lastLogin = stats.LastLogin.Format(time.RFC3339)
	}

	lastLoginLocation := ""
	if stats.LastLoginLocation != nil {
		lastLoginLocation = *stats.LastLoginLocation
	}

	return &pb.GetSessionStatsResponse{
		Success:           true,
		Message:           "Statistics retrieved successfully",
		TotalSessions:     int32(stats.TotalSessions),
		ActiveSessions:    int32(stats.ActiveSessions),
		TotalDevices:      int32(stats.TotalDevices),
		TrustedDevices:    int32(stats.TrustedDevices),
		LastLogin:         lastLogin,
		LastLoginLocation: lastLoginLocation,
		RecentLocations:   stats.RecentLocations,
	}, nil
}

// Helper function to get string value from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Helper function to get float64 value from pointer
func getFloat64Value(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

// Helper function to create audit log
func (h *SessionHandler) createAuditLog(log *models.AuditLog) {
	err := h.repo.CreateAuditLog(log)
	if err != nil {
		logMsg := fmt.Sprintf("Failed to create audit log: %v", err)
		fmt.Println(logMsg)
	}
}