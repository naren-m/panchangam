# CORS and Region Configuration Updates

## CORS Configuration

The CORS (Cross-Origin Resource Sharing) configuration has been updated to be more flexible and secure:

### Environment Variable Configuration

Instead of hardcoding allowed origins, the system now uses the `CORS_ALLOWED_ORIGINS` environment variable:

```bash
# Example: Allow specific origins (comma-separated)
export CORS_ALLOWED_ORIGINS="https://panchangam.app,https://staging.panchangam.app,http://localhost:5173"

# If not set, defaults to development origins:
# - http://localhost:5173 (Vite dev server)
# - http://localhost:3000 (React dev server)  
# - http://localhost:8086 (Docker frontend container)
```

### Security Note

The previous `ALLOW_ALL_ORIGINS` flag that allowed wildcard (`*`) origins has been removed for security reasons. Always specify explicit allowed origins in production.

## Calendar System by Region

The calendar system determination has been refactored to use a maintainable map-based approach:

### Supported Regions

#### Amanta System (South India)
- Tamil Nadu
- Kerala
- Karnataka
- Andhra Pradesh
- Telangana

#### Purnimanta System (North India and default)
- Maharashtra
- Gujarat
- Rajasthan
- Uttar Pradesh
- Madhya Pradesh
- Bihar
- West Bengal
- Odisha
- Punjab
- Haryana
- Himachal Pradesh
- Uttarakhand
- Delhi

#### International Regions (default to Purnimanta)
- California
- New York
- Texas
- New Jersey

### Adding New Regions

To add support for a new region, simply update the `calendarSystemByRegion` map in `services/panchangam/service.go`:

```go
var calendarSystemByRegion = map[string]string{
    // Add your region here
    "New Region": "Amanta", // or "Purnimanta"
}
```

### Note about Tamil Nadu

As per the comment, Tamil Nadu uses the Amanta system, which has been correctly configured in the map.
