-- Database initialization script for Session Management System

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table: stores user credentials and basic info
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    mfa_enabled BOOLEAN DEFAULT false,
    mfa_secret VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Devices table: tracks unique devices per user
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_fingerprint VARCHAR(255) UNIQUE NOT NULL,
    device_name VARCHAR(255),
    device_type VARCHAR(50), -- mobile, desktop, tablet
    os VARCHAR(100),
    browser VARCHAR(100),
    is_trusted BOOLEAN DEFAULT false,
    first_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table: tracks active user sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    refresh_token VARCHAR(500) UNIQUE NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT,
    location_country VARCHAR(100),
    location_city VARCHAR(100),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE
);

-- Audit logs table: comprehensive security event logging
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    session_id UUID REFERENCES sessions(id) ON DELETE SET NULL,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    event_type VARCHAR(50) NOT NULL, -- login, logout, mfa_enable, session_revoke, etc.
    event_category VARCHAR(50) NOT NULL, -- authentication, authorization, session_management, security
    severity VARCHAR(20) NOT NULL, -- info, warning, critical
    ip_address INET,
    user_agent TEXT,
    location_country VARCHAR(100),
    location_city VARCHAR(100),
    metadata JSONB, -- flexible field for additional event data
    success BOOLEAN NOT NULL,
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Security alerts table: tracks detected anomalies
CREATE TABLE security_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    alert_type VARCHAR(50) NOT NULL, -- impossible_travel, new_device, suspicious_location, etc.
    severity VARCHAR(20) NOT NULL, -- low, medium, high, critical
    description TEXT NOT NULL,
    metadata JSONB,
    ip_address INET,
    location_country VARCHAR(100),
    location_city VARCHAR(100),
    is_resolved BOOLEAN DEFAULT false,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- MFA backup codes table: stores one-time backup codes
CREATE TABLE mfa_backup_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash VARCHAR(255) NOT NULL,
    is_used BOOLEAN DEFAULT false,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_devices_user_id ON devices(user_id);
CREATE INDEX idx_devices_fingerprint ON devices(device_fingerprint);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_is_active ON sessions(is_active);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_security_alerts_user_id ON security_alerts(user_id);
CREATE INDEX idx_security_alerts_is_resolved ON security_alerts(is_resolved);
CREATE INDEX idx_mfa_backup_codes_user_id ON mfa_backup_codes(user_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert a test user (password is "password123" - hashed with bcrypt)
-- This is just for development/testing
INSERT INTO users (email, password_hash, full_name, is_active) VALUES
('demo@example.com', '$2a$10$rKJ8qV.1Q9pZ9K0xZqYZFe5vYX9K0zZQZ9ZQZ9ZQZ9ZQZ9ZQZ9ZQZ', 'Demo User', true);