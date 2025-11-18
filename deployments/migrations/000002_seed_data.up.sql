-- Seed initial data

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
    ('Tirupati', 13.6288, 79.4192, 'Asia/Kolkata', 'India', 'Andhra Pradesh', true),
    ('Jaipur', 26.9124, 75.7873, 'Asia/Kolkata', 'India', 'Rajasthan', true),
    ('Lucknow', 26.8467, 80.9462, 'Asia/Kolkata', 'India', 'Uttar Pradesh', true),
    ('Chandigarh', 30.7333, 76.7794, 'Asia/Kolkata', 'India', 'Chandigarh', true),
    ('Kochi', 9.9312, 76.2673, 'Asia/Kolkata', 'India', 'Kerala', true),
    ('Thiruvananthapuram', 8.5241, 76.9366, 'Asia/Kolkata', 'India', 'Kerala', true)
ON CONFLICT (name, country) DO NOTHING;

-- Insert default feature flags
INSERT INTO feature_flags (name, description, enabled, rollout_percentage) VALUES
    ('regional_variations', 'Enable regional calendar variations', true, 100),
    ('sky_view_3d', 'Enable 3D sky visualization', true, 100),
    ('festival_notifications', 'Enable festival notifications', false, 0),
    ('api_v2', 'Enable API version 2', false, 50),
    ('advanced_calculations', 'Enable advanced astronomical calculations', true, 100),
    ('user_preferences', 'Enable user preference storage', true, 100),
    ('cache_optimization', 'Enable advanced caching optimizations', true, 100),
    ('real_time_updates', 'Enable real-time calculation updates', false, 25),
    ('planetary_positions', 'Enable detailed planetary position calculations', true, 100),
    ('muhurta_recommendations', 'Enable auspicious time recommendations', false, 0)
ON CONFLICT (name) DO NOTHING;
