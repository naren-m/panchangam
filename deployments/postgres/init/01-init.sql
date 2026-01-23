-- Panchangam Database Initialization Script
-- This script sets up the initial database schema for the Panchangam application

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search

-- Create custom types
CREATE TYPE calculation_system AS ENUM ('purnimanta', 'amanta');
CREATE TYPE regional_variation AS ENUM ('north_indian', 'south_indian', 'telugu', 'tamil', 'malayalam', 'kannada');

-- ==================================
-- Users and Authentication
-- ==================================

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    preferences JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_active ON users(is_active);

-- ==================================
-- User Preferences
-- ==================================

CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    timezone VARCHAR(100) DEFAULT 'UTC',
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    location_name VARCHAR(255),
    calculation_system calculation_system DEFAULT 'purnimanta',
    regional_variation regional_variation DEFAULT 'north_indian',
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);

-- ==================================
-- Panchangam Calculations Cache
-- ==================================

CREATE TABLE IF NOT EXISTS panchangam_cache (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cache_key VARCHAR(512) UNIQUE NOT NULL,
    date DATE NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    timezone VARCHAR(100) NOT NULL,
    calculation_system calculation_system,
    regional_variation regional_variation,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_cache_entry UNIQUE(date, latitude, longitude, timezone, calculation_system, regional_variation)
);

CREATE INDEX idx_panchangam_cache_key ON panchangam_cache(cache_key);
CREATE INDEX idx_panchangam_cache_date ON panchangam_cache(date);
CREATE INDEX idx_panchangam_cache_expires ON panchangam_cache(expires_at);
CREATE INDEX idx_panchangam_cache_location ON panchangam_cache(latitude, longitude);

-- ==================================
-- User Queries Log (for analytics)
-- ==================================

CREATE TABLE IF NOT EXISTS query_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    query_type VARCHAR(100) NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    date DATE,
    timezone VARCHAR(100),
    calculation_system calculation_system,
    response_time_ms INTEGER,
    cache_hit BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_query_log_user_id ON query_log(user_id);
CREATE INDEX idx_query_log_type ON query_log(query_type);
CREATE INDEX idx_query_log_created_at ON query_log(created_at);
CREATE INDEX idx_query_log_cache_hit ON query_log(cache_hit);

-- ==================================
-- Locations (predefined)
-- ==================================

CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    timezone VARCHAR(100) NOT NULL,
    country VARCHAR(100),
    state_province VARCHAR(100),
    is_popular BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_location UNIQUE(name, country)
);

CREATE INDEX idx_locations_name ON locations USING gin(name gin_trgm_ops);
CREATE INDEX idx_locations_country ON locations(country);
CREATE INDEX idx_locations_popular ON locations(is_popular);
CREATE INDEX idx_locations_coords ON locations(latitude, longitude);

-- ==================================
-- Feature Flags
-- ==================================

CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT false,
    rollout_percentage INTEGER DEFAULT 0 CHECK (rollout_percentage >= 0 AND rollout_percentage <= 100),
    conditions JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255)
);

CREATE INDEX idx_feature_flags_name ON feature_flags(name);
CREATE INDEX idx_feature_flags_enabled ON feature_flags(enabled);

-- ==================================
-- API Keys (for external integrations)
-- ==================================

CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    scopes JSONB DEFAULT '[]'::jsonb,
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    rate_limit INTEGER DEFAULT 1000 -- requests per hour
);

CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_active ON api_keys(is_active);

-- ==================================
-- Audit Log
-- ==================================

CREATE TABLE IF NOT EXISTS audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_action ON audit_log(action);
CREATE INDEX idx_audit_log_entity ON audit_log(entity_type, entity_id);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at);

-- ==================================
-- Insert default data
-- ==================================

-- Insert popular Indian locations
INSERT INTO locations (name, latitude, longitude, timezone, country, state_province, is_popular) VALUES
    ('Delhi', 28.6139, 77.2090, 'Asia/Kolkata', 'India', 'Delhi', true),
    ('Mumbai', 19.0760, 72.8777, 'Asia/Kolkata', 'India', 'Maharashtra', true),
    ('Bangalore', 12.9716, 77.5946, 'Asia/Kolkata', 'India', 'Karnataka', true),
    ('Chennai', 13.0827, 80.2707, 'Asia/Kolkata', 'India', 'Tamil Nadu', true),
    ('Kolkata', 22.5726, 88.3639, 'Asia/Kolkata', 'India', 'West Bengal', true),
    ('Hyderabad', 17.3850, 78.4867, 'Asia/Kolkata', 'India', 'Telangana', true),
    ('Pune', 18.5204, 73.8567, 'Asia/Kolkata', 'India', 'Maharashtra', true),
    ('Ahmedabad', 23.0225, 72.5714, 'Asia/Kolkata', 'India', 'Gujarat', true),
    ('Varanasi', 25.3176, 82.9739, 'Asia/Kolkata', 'India', 'Uttar Pradesh', true),
    ('Tirupati', 13.6288, 79.4192, 'Asia/Kolkata', 'India', 'Andhra Pradesh', true)
ON CONFLICT (name, country) DO NOTHING;

-- Insert default feature flags
INSERT INTO feature_flags (name, description, enabled, rollout_percentage) VALUES
    ('regional_variations', 'Enable regional calendar variations', true, 100),
    ('sky_view_3d', 'Enable 3D sky visualization', true, 100),
    ('festival_notifications', 'Enable festival notifications', false, 0),
    ('api_v2', 'Enable API version 2', false, 50),
    ('advanced_calculations', 'Enable advanced astronomical calculations', true, 100),
    ('user_preferences', 'Enable user preference storage', true, 100),
    ('cache_optimization', 'Enable advanced caching optimizations', true, 100)
ON CONFLICT (name) DO NOTHING;

-- ==================================
-- Functions and Triggers
-- ==================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_preferences_updated_at BEFORE UPDATE ON user_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_feature_flags_updated_at BEFORE UPDATE ON feature_flags
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to clean expired cache entries
CREATE OR REPLACE FUNCTION clean_expired_cache()
RETURNS void AS $$
BEGIN
    DELETE FROM panchangam_cache WHERE expires_at < CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO panchangam;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO panchangam;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO panchangam;
