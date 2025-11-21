// User Types
export interface User {
    id: string;
    email: string;
    fullName: string;
    isActive: boolean;
    mfaEnabled: boolean;
    createdAt: string;
    updatedAt: string;
  }
  
  // Session Types
  export interface Session {
    id: string;
    userId: string;
    deviceId: string;
    deviceName: string;
    deviceType: string;
    ipAddress: string;
    userAgent?: string;
    locationCountry?: string;
    locationCity?: string;
    latitude?: number;
    longitude?: number;
    isActive: boolean;
    createdAt: string;
    lastSeenAt: string;
    expiresAt: string;
    isCurrent: boolean;
  }
  
  export interface SessionStats {
    totalSessions: number;
    activeSessions: number;
    totalDevices: number;
    trustedDevices: number;
    lastLogin?: string;
    lastLoginLocation?: string;
    recentLocations: string[];
  }
  
  // Device Types
  export interface Device {
    id: string;
    userId: string;
    deviceFingerprint: string;
    deviceName: string;
    deviceType: string;
    os: string;
    browser: string;
    isTrusted: boolean;
    firstSeenAt: string;
    lastSeenAt: string;
    sessionCount: number;
  }
  
  // Audit Types
  export interface AuditLog {
    id: string;
    userId?: string;
    sessionId?: string;
    deviceId?: string;
    eventType: string;
    eventCategory: string;
    severity: string;
    ipAddress?: string;
    userAgent?: string;
    locationCountry?: string;
    locationCity?: string;
    metadata?: string;
    success: boolean;
    failureReason?: string;
    createdAt: string;
  }
  
  // Security Alert Types
  export interface SecurityAlert {
    id: string;
    userId: string;
    alertType: string;
    severity: string;
    description: string;
    metadata?: string;
    ipAddress?: string;
    locationCountry?: string;
    locationCity?: string;
    isResolved: boolean;
    resolvedAt?: string;
    createdAt: string;
  }
  
  // Activity Types
  export interface DailyActivity {
    date: string;
    loginCount: number;
    failedLoginCount: number;
  }
  
  export interface ActivitySummary {
    totalLogins: number;
    failedLoginAttempts: number;
    uniqueDevices: number;
    uniqueLocations: number;
    dailyActivity: DailyActivity[];
  }
  
  // Compliance Types
  export interface EventCount {
    eventType: string;
    count: number;
  }
  
  export interface ComplianceReport {
    totalEvents: number;
    successfulLogins: number;
    failedLogins: number;
    sessionRevocations: number;
    securityAlerts: number;
    mfaEvents: number;
    eventBreakdown: EventCount[];
    topLocations: string[];
  }
  
  // Auth Types
  export interface AuthPayload {
    success: boolean;
    message: string;
    userId?: string;
    accessToken?: string;
    refreshToken?: string;
    mfaRequired?: boolean;
    sessionId?: string;
  }
  
  export interface MFASetup {
    success: boolean;
    message: string;
    secret?: string;
    qrCodeUrl?: string;
    backupCodes?: string[];
  }
  
  export interface GenericResponse {
    success: boolean;
    message: string;
  }