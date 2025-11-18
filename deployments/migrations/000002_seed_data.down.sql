-- Rollback seed data

-- Remove seeded feature flags
DELETE FROM feature_flags WHERE name IN (
    'regional_variations',
    'sky_view_3d',
    'festival_notifications',
    'api_v2',
    'advanced_calculations',
    'user_preferences',
    'cache_optimization',
    'real_time_updates',
    'planetary_positions',
    'muhurta_recommendations'
);

-- Remove seeded locations
DELETE FROM locations WHERE country = 'India' AND is_popular = true;
