import { gql } from '@apollo/client';

// Auth Queries
export const GET_ME = gql`
  query GetMe {
    me {
      id
      email
      fullName
      isActive
      mfaEnabled
      createdAt
      updatedAt
    }
  }
`;

export const VALIDATE_TOKEN = gql`
  query ValidateToken($token: String!) {
    validateToken(token: $token) {
      valid
      userId
      email
      message
    }
  }
`;

// Session Queries
export const GET_SESSIONS = gql`
  query GetSessions($includeInactive: Boolean) {
    sessions(includeInactive: $includeInactive) {
      success
      message
      sessions {
        id
        userId
        deviceId
        deviceName
        deviceType
        ipAddress
        userAgent
        locationCountry
        locationCity
        latitude
        longitude
        isActive
        createdAt
        lastSeenAt
        expiresAt
        isCurrent
      }
      totalCount
      activeCount
    }
  }
`;

export const GET_SESSION_STATS = gql`
  query GetSessionStats {
    sessionStats {
      totalSessions
      activeSessions
      totalDevices
      trustedDevices
      lastLogin
      lastLoginLocation
      recentLocations
    }
  }
`;

// Device Queries
export const GET_DEVICES = gql`
  query GetDevices {
    devices {
      success
      message
      devices {
        id
        userId
        deviceFingerprint
        deviceName
        deviceType
        os
        browser
        isTrusted
        firstSeenAt
        lastSeenAt
        sessionCount
      }
      totalCount
      trustedCount
    }
  }
`;

// Audit Queries
export const GET_AUDIT_LOGS = gql`
  query GetAuditLogs(
    $limit: Int
    $offset: Int
    $eventCategory: String
    $severity: String
    $successOnly: Boolean
  ) {
    auditLogs(
      limit: $limit
      offset: $offset
      eventCategory: $eventCategory
      severity: $severity
      successOnly: $successOnly
    ) {
      success
      message
      logs {
        id
        userId
        sessionId
        deviceId
        eventType
        eventCategory
        severity
        ipAddress
        userAgent
        locationCountry
        locationCity
        metadata
        success
        failureReason
        createdAt
      }
      totalCount
    }
  }
`;

export const GET_SECURITY_ALERTS = gql`
  query GetSecurityAlerts($includeResolved: Boolean, $severity: String) {
    securityAlerts(includeResolved: $includeResolved, severity: $severity) {
      success
      message
      alerts {
        id
        userId
        alertType
        severity
        description
        metadata
        ipAddress
        locationCountry
        locationCity
        isResolved
        resolvedAt
        createdAt
      }
      totalCount
      unresolvedCount
    }
  }
`;

export const GET_ACTIVITY_SUMMARY = gql`
  query GetActivitySummary($days: Int) {
    activitySummary(days: $days) {
      totalLogins
      failedLoginAttempts
      uniqueDevices
      uniqueLocations
      dailyActivity {
        date
        loginCount
        failedLoginCount
      }
    }
  }
`;

export const GET_COMPLIANCE_REPORT = gql`
  query GetComplianceReport($startDate: String, $endDate: String) {
    complianceReport(startDate: $startDate, endDate: $endDate) {
      totalEvents
      successfulLogins
      failedLogins
      sessionRevocations
      securityAlerts
      mfaEvents
      eventBreakdown {
        eventType
        count
      }
      topLocations
    }
  }
`;