import { gql } from '@apollo/client';

// Auth Mutations
export const REGISTER = gql`
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
      success
      message
      userId
      accessToken
      refreshToken
    }
  }
`;

export const LOGIN = gql`
  mutation Login($input: LoginInput!) {
    login(input: $input) {
      success
      message
      userId
      accessToken
      refreshToken
      mfaRequired
      sessionId
    }
  }
`;

export const REFRESH_TOKEN = gql`
  mutation RefreshToken($refreshToken: String!) {
    refreshToken(refreshToken: $refreshToken) {
      success
      message
      accessToken
      refreshToken
    }
  }
`;

export const ENABLE_MFA = gql`
  mutation EnableMFA {
    enableMFA {
      success
      message
      secret
      qrCodeUrl
      backupCodes
    }
  }
`;

export const VERIFY_MFA = gql`
  mutation VerifyMFA($code: String!) {
    verifyMFA(code: $code) {
      success
      message
    }
  }
`;

// Session Mutations
export const REVOKE_SESSION = gql`
  mutation RevokeSession($sessionId: ID!) {
    revokeSession(sessionId: $sessionId) {
      success
      message
    }
  }
`;

export const REVOKE_ALL_SESSIONS = gql`
  mutation RevokeAllSessions($exceptCurrent: Boolean) {
    revokeAllSessions(exceptCurrent: $exceptCurrent) {
      success
      message
    }
  }
`;

export const TRUST_DEVICE = gql`
  mutation TrustDevice($deviceId: ID!) {
    trustDevice(deviceId: $deviceId) {
      success
      message
    }
  }
`;

// Security Mutations
export const RESOLVE_SECURITY_ALERT = gql`
  mutation ResolveSecurityAlert($alertId: ID!) {
    resolveSecurityAlert(alertId: $alertId) {
      success
      message
    }
  }
`;