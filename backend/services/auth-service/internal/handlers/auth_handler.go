package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"net"
	"net/http"
	"github.com/google/uuid"
	pb "github.com/aashiq-04/session-management-system/backend/services/auth-service/proto"
	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/models"
	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/repository"
	"github.com/aashiq-04/session-management-system/backend/services/auth-service/internal/utils"
)


func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func floatPtr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}


// AuthHandler implements the AuthService gRPC service
type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	repo      *repository.UserRepository
	jwtSecret string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *sql.DB) *AuthHandler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	return &AuthHandler{
		repo:      repository.NewUserRepository(db),
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register request received for email: %s", req.Email)

	// Validate input
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Email, password, and full name are required",
		}, nil
	}

	// Check if user already exists
	existingUser, _ := h.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "User with this email already exists",
		}, nil
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to process password",
		}, nil
	}

	// Create user model
	userID := uuid.New().String()
	user := &models.User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		IsActive:     true,
		MFAEnabled:   false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user to database
	err = h.repo.CreateUser(user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to create user account",
		}, nil
	}

	// Create or get device
	deviceID, err := h.handleDevice(req.DeviceInfo, userID)
	if err != nil {
		log.Printf("Failed to handle device: %v", err)
		// Continue anyway - device tracking is not critical for registration
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(userID, req.Email, h.jwtSecret)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to generate authentication token",
		}, nil
	}

	refreshToken, err := utils.GenerateRefreshToken(userID, req.Email, h.jwtSecret)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to generate refresh token",
		}, nil
	}

	// Create session
	sessionID := uuid.New().String()
	session := &models.Session{
		ID:           sessionID,
		UserID:       userID,
		DeviceID:     deviceID,
		RefreshToken: refreshToken,
		IPAddress:    req.DeviceInfo.IpAddress,
		UserAgent:    &req.DeviceInfo.UserAgent,
		LocationCountry: strPtr(req.DeviceInfo.LocationCountry),
		LocationCity:    strPtr(req.DeviceInfo.LocationCity),
		Latitude:        floatPtr(req.DeviceInfo.Latitude),
		Longitude:       floatPtr(req.DeviceInfo.Longitude),
		IsActive:     true,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	err = h.repo.CreateSession(session)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		// Continue anyway - session tracking is not critical
	}

	// Create audit log
	h.createAuditLog(&models.AuditLog{
		ID:            uuid.New().String(),
		UserID:        &userID,
		SessionID:     &sessionID,
		DeviceID:      &deviceID,
		EventType:     "user_registered",
		EventCategory: "authentication",
		Severity:      "info",
		IPAddress:     strPtr(req.DeviceInfo.IpAddress),
		UserAgent:     &req.DeviceInfo.UserAgent,
		LocationCountry: strPtr(req.DeviceInfo.LocationCountry),
		LocationCity:    strPtr(req.DeviceInfo.LocationCity),
		Success:       true,
		CreatedAt:     time.Now(),
	})

	log.Printf("User registered successfully: %s", userID)

	return &pb.RegisterResponse{
		Success:      true,
		Message:      "User registered successfully",
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login handles user authentication
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Login request received for email: %s", req.Email)

	// Validate input
	if req.Email == "" || req.Password == "" {
		return &pb.LoginResponse{
			Success: false,
			Message: "Email and password are required",
		}, nil
	}

	// Get user from database
	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("User not found: %s", req.Email)
		h.createFailedLoginAuditLog(req.Email, req.DeviceInfo, "user_not_found")
		return &pb.LoginResponse{
			Success: false,
			Message: "Invalid email or password",
		}, nil
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("User account is inactive: %s", req.Email)
		h.createFailedLoginAuditLog(req.Email, req.DeviceInfo, "account_inactive")
		return &pb.LoginResponse{
			Success: false,
			Message: "Account is inactive",
		}, nil
	}

	// Verify password
	err = utils.ComparePassword(user.PasswordHash, req.Password)
	if err != nil {
		log.Printf("Invalid password for user: %s", req.Email)
		h.createFailedLoginAuditLog(req.Email, req.DeviceInfo, "invalid_password")
		return &pb.LoginResponse{
			Success: false,
			Message: "Invalid email or password",
		}, nil
	}

	// Check if MFA is enabled
	if user.MFAEnabled {
		// If MFA code is not provided, request it
		if req.MfaCode == "" {
			return &pb.LoginResponse{
				Success:     false,
				Message:     "MFA code required",
				MfaRequired: true,
			}, nil
		}

		// Validate MFA code
		if user.MFASecret == nil {
			log.Printf("MFA secret not found for user: %s", user.ID)
			return &pb.LoginResponse{
				Success: false,
				Message: "MFA configuration error",
			}, nil
		}

		valid := utils.ValidateMFACode(req.MfaCode, *user.MFASecret)
		if !valid {
			log.Printf("Invalid MFA code for user: %s", user.ID)
			h.createFailedLoginAuditLog(req.Email, req.DeviceInfo, "invalid_mfa_code")
			return &pb.LoginResponse{
				Success: false,
				Message: "Invalid MFA code",
			}, nil
		}
	}

	// Create or get device
	deviceID, err := h.handleDevice(req.DeviceInfo, user.ID)
	if err != nil {
		log.Printf("Failed to handle device: %v", err)
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		return &pb.LoginResponse{
			Success: false,
			Message: "Failed to generate authentication token",
		}, nil
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		return &pb.LoginResponse{
			Success: false,
			Message: "Failed to generate refresh token",
		}, nil
	}

	// Create session
	sessionID := uuid.New().String()
	session := &models.Session{
		ID:           sessionID,
		UserID:       user.ID,
		DeviceID:     deviceID,
		RefreshToken: refreshToken,
		IPAddress:    *strPtr(req.DeviceInfo.IpAddress),
		UserAgent:    &req.DeviceInfo.UserAgent,
		LocationCountry: strPtr(req.DeviceInfo.LocationCountry),
		LocationCity:    strPtr(req.DeviceInfo.LocationCity),
		Latitude:        floatPtr(req.DeviceInfo.Latitude),
		Longitude:       floatPtr(req.DeviceInfo.Longitude),
		IsActive:     true,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	err = h.repo.CreateSession(session)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
	}

	// Create audit log
	userIDCopy := user.ID
	h.createAuditLog(&models.AuditLog{
		ID:            uuid.New().String(),
		UserID:        &userIDCopy,
		SessionID:     &sessionID,
		DeviceID:      &deviceID,
		EventType:     "user_login",
		EventCategory: "authentication",
		Severity:      "info",
		IPAddress:     strPtr(req.DeviceInfo.IpAddress),
		UserAgent:     &req.DeviceInfo.UserAgent,
		LocationCountry: strPtr(req.DeviceInfo.LocationCountry),
		LocationCity:    strPtr(req.DeviceInfo.LocationCity),
		Success:       true,
		CreatedAt:     time.Now(),
	})

	log.Printf("User logged in successfully: %s", user.ID)

	return &pb.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionId:    sessionID,
	}, nil
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := utils.ValidateToken(req.Token, h.jwtSecret)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:   true,
		UserId:  claims.UserID,
		Email:   claims.Email,
		Message: "Token is valid",
	}, nil
}

// RefreshToken generates a new access token from a refresh token
func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Success: false,
			Message: "Invalid refresh token",
		}, nil
	}

	// Check if session exists and is active
	session, err := h.repo.GetSessionByRefreshToken(req.RefreshToken)
	if err != nil || !session.IsActive {
		return &pb.RefreshTokenResponse{
			Success: false,
			Message: "Session not found or inactive",
		}, nil
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return &pb.RefreshTokenResponse{
			Success: false,
			Message: "Session expired",
		}, nil
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Email, h.jwtSecret)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Success: false,
			Message: "Failed to generate new access token",
		}, nil
	}

	return &pb.RefreshTokenResponse{
		Success:      true,
		Message:      "Token refreshed successfully",
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken, // Keep the same refresh token
	}, nil
}

// EnableMFA enables multi-factor authentication for a user
func (h *AuthHandler) EnableMFA(ctx context.Context, req *pb.EnableMFARequest) (*pb.EnableMFAResponse, error) {
	// Get user
	user, err := h.repo.GetUserByID(req.UserId)
	if err != nil {
		return &pb.EnableMFAResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	// Generate MFA secret
	secret, qrCodeURL, err := utils.GenerateMFASecret(user.Email)
	if err != nil {
		return &pb.EnableMFAResponse{
			Success: false,
			Message: "Failed to generate MFA secret",
		}, nil
	}

	// Generate backup codes
	backupCodes, err := utils.GenerateBackupCodes()
	if err != nil {
		return &pb.EnableMFAResponse{
			Success: false,
			Message: "Failed to generate backup codes",
		}, nil
	}

	// Enable MFA for user
	err = h.repo.EnableMFA(user.ID, secret)
	if err != nil {
		return &pb.EnableMFAResponse{
			Success: false,
			Message: "Failed to enable MFA",
		}, nil
	}

	// Create audit log
	h.createAuditLog(&models.AuditLog{
		ID:            uuid.New().String(),
		UserID:        &user.ID,
		EventType:     "mfa_enabled",
		EventCategory: "security",
		Severity:      "info",
		Success:       true,
		CreatedAt:     time.Now(),
	})

	return &pb.EnableMFAResponse{
		Success:     true,
		Message:     "MFA enabled successfully",
		Secret:      secret,
		QrCodeUrl:   qrCodeURL,
		BackupCodes: backupCodes,
	}, nil
}

// VerifyMFA verifies an MFA code
func (h *AuthHandler) VerifyMFA(ctx context.Context, req *pb.VerifyMFARequest) (*pb.VerifyMFAResponse, error) {
	// Get user
	user, err := h.repo.GetUserByID(req.UserId)
	if err != nil {
		return &pb.VerifyMFAResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	if !user.MFAEnabled || user.MFASecret == nil {
		return &pb.VerifyMFAResponse{
			Success: false,
			Message: "MFA not enabled",
		}, nil
	}

	// Validate MFA code
	valid := utils.ValidateMFACode(req.Code, *user.MFASecret)
	if !valid {
		return &pb.VerifyMFAResponse{
			Success: false,
			Message: "Invalid MFA code",
		}, nil
	}

	return &pb.VerifyMFAResponse{
		Success: true,
		Message: "MFA code verified",
	}, nil
}

// GetUserProfile retrieves user profile information
func (h *AuthHandler) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	user, err := h.repo.GetUserByID(req.UserId)
	if err != nil {
		return &pb.GetUserProfileResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	return &pb.GetUserProfileResponse{
		Success: true,
		Message: "User profile retrieved",
		Profile: &pb.UserProfile{
			Id:         user.ID,
			Email:      user.Email,
			FullName:   user.FullName,
			IsActive:   user.IsActive,
			MfaEnabled: user.MFAEnabled,
			CreatedAt:  user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// Helper function to handle device creation/retrieval
func (h *AuthHandler) handleDevice(deviceInfo *pb.DeviceInfo, userID string) (string, error) {
	if deviceInfo == nil || deviceInfo.DeviceFingerprint == "" {
		return "", fmt.Errorf("device info is required")
	}

	// Check if device exists
	existingDevice, err := h.repo.GetDeviceByFingerprint(deviceInfo.DeviceFingerprint)
	if err != nil {
		return "", err
	}

	if existingDevice != nil {
		// Update last seen
		h.repo.UpdateDeviceLastSeen(existingDevice.ID)
		return existingDevice.ID, nil
	}

	// Create new device
	deviceID := uuid.New().String()
	device := &models.Device{
		ID:                deviceID,
		UserID:            userID,
		DeviceFingerprint: deviceInfo.DeviceFingerprint,
		DeviceName:        &deviceInfo.DeviceName,
		DeviceType:        &deviceInfo.DeviceType,
		OS:                &deviceInfo.Os,
		Browser:           &deviceInfo.Browser,
		IsTrusted:         false,
		FirstSeenAt:       time.Now(),
		LastSeenAt:        time.Now(),
		CreatedAt:         time.Now(),
	}

	err = h.repo.CreateDevice(device)
	if err != nil {
		return "", err
	}

	return deviceID, nil
}
func getIPFromContext(ctx context.Context) string {
	req, ok := ctx.Value("httpRequest").(*http.Request)
	if !ok || req == nil {
		return "0.0.0.0"
	}

	// 1) Check proxy header
	if forwarded := req.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}

	// 2) Check X-Real-IP (nginx)
	if realIP := req.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// 3) Fall back
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		return ip
	}

	return "0.0.0.0"
}

// Helper function to create audit log
func (h *AuthHandler) createAuditLog(log *models.AuditLog) {
	err := h.repo.CreateAuditLog(log)
	if err != nil {
		log := fmt.Sprintf("Failed to create audit log: %v", err)
		fmt.Println(log)
	}
}

// Helper function to create failed login audit log
func (h *AuthHandler) createFailedLoginAuditLog(email string, deviceInfo *pb.DeviceInfo, reason string) {
	h.createAuditLog(&models.AuditLog{
		ID:            uuid.New().String(),
		EventType:     "login_failed",
		EventCategory: "authentication",
		Severity:      "warning",
		IPAddress:     &deviceInfo.IpAddress,
		UserAgent:     &deviceInfo.UserAgent,
		LocationCountry: &deviceInfo.LocationCountry,
		LocationCity:  &deviceInfo.LocationCity,
		Success:       false,
		FailureReason: &reason,
		CreatedAt:     time.Now(),
	})
}