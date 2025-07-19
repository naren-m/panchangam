# Panchangam API Extension Examples

This directory contains example plugins that demonstrate the extensibility of the Panchangam API architecture. These examples showcase how to implement different types of plugins to extend the core functionality.

## Example Plugins

### 1. Tamil Nadu Regional Extension (`tamil_nadu_plugin.go`)

A comprehensive regional extension plugin that provides Tamil Nadu specific calculations and traditions.

**Features:**
- Tamil calendar adjustments (Amanta system)
- Regional festivals (Pongal, Chithirai Festival, Thai Pusam)
- Tamil-specific muhurtas (Abhijit, Brahma)
- Tamil language names for Panchangam elements
- Regional calculation rules and interpretations

**Usage:**
```go
// Create and register the plugin
tamilPlugin := examples.NewTamilNaduExtension()
api.RegisterPlugin(tamilPlugin)

// The plugin will automatically apply Tamil Nadu specific rules
// when the region is set to RegionTamilNadu
req := api.PanchangamRequest{
    Date:     time.Now(),
    Location: chennaiLocation,
    Region:   api.RegionTamilNadu,
}
```

**Key Methods:**
- `ApplyRegionalRules()` - Applies Tamil calendar adjustments
- `GetRegionalEvents()` - Returns Tamil festivals and observances
- `GetRegionalMuhurtas()` - Provides Tamil-specific auspicious times
- `GetRegionalNames()` - Tamil language translations

### 2. Tamil Localization Plugin (`tamil_localization_plugin.go`)

A specialized localization plugin that provides comprehensive Tamil language support for all Panchangam elements.

**Features:**
- Complete Tamil translations for Tithis, Nakshatras, Yogas, Karanas
- Localized divine names and symbols
- Cultural adaptation of significance descriptions
- Support for multiple Tamil locale codes (ta, ta-IN, tamil)
- Quality and event type translations

**Usage:**
```go
// Create and register the localization plugin
tamilLocPlugin := examples.NewTamilLocalizationPlugin()
api.RegisterPlugin(tamilLocPlugin)

// Request data with Tamil locale
req := api.PanchangamRequest{
    Date:     time.Now(),
    Location: location,
    Locale:   "ta",  // Tamil locale
    Region:   api.RegionTamilNadu,
}
```

**Translation Examples:**
- Ashwini → அசுவினி
- Pratipada → பிரதமை
- Vishkambha → விஷ்கம்பா
- Bava → பவ
- Sunday → ஞாயிறு

### 3. Hindu Festival Plugin (`hindu_festival_plugin.go`)

A comprehensive event plugin that provides detailed Hindu festival calculations across different regions and traditions.

**Features:**
- Major festivals (Diwali, Holi, Navaratri, Janmashtami)
- Solar festivals (Makar Sankranti, Ram Navami)
- Ekadashi observances with names and significance
- Monthly observances (Amavasya, Purnima)
- Regional festival variations
- Detailed metadata including rituals, deities, and traditions

**Usage:**
```go
// Create and register the festival plugin
festivalPlugin := examples.NewHinduFestivalPlugin()
api.RegisterPlugin(festivalPlugin)

// Request data with events included
req := api.PanchangamRequest{
    Date:          time.Now(),
    Location:      location,
    IncludeEvents: true,
    Region:        api.RegionNorthIndia,
}
```

**Festival Categories:**
- **Lunar Festivals**: Based on lunar calendar positions
- **Solar Festivals**: Based on solar transitions
- **Ekadashi Events**: Monthly Vishnu fasting days
- **Regional Festivals**: State-specific celebrations

## Plugin Architecture Benefits

### 1. Extensibility
- Easy to add new regional calculations
- Support for additional languages and locales
- Custom festival and event definitions
- Flexible validation and business rules

### 2. Modularity
- Plugins can be enabled/disabled independently
- No impact on core API functionality
- Clean separation of concerns
- Easy testing and maintenance

### 3. Configurability
- Runtime configuration of plugin behavior
- Priority-based plugin execution
- Regional and locale-specific activation
- Feature flags and toggles

## Implementation Patterns

### Plugin Interface Implementation
```go
type MyPlugin struct {
    enabled bool
    config  map[string]interface{}
}

func (p *MyPlugin) GetInfo() api.PluginInfo { /* ... */ }
func (p *MyPlugin) Initialize(ctx context.Context, config map[string]interface{}) error { /* ... */ }
func (p *MyPlugin) IsEnabled() bool { /* ... */ }
func (p *MyPlugin) Shutdown(ctx context.Context) error { /* ... */ }
```

### Capability-Specific Interfaces
```go
// For regional extensions
func (p *MyPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error

// For event plugins
func (p *MyPlugin) GetEvents(ctx context.Context, date time.Time, location api.Location, region api.Region) ([]api.Event, error)

// For localization plugins
func (p *MyPlugin) LocalizeTithi(tithi *api.Tithi, locale string, region api.Region) error
```

### Error Handling
```go
if err := plugin.SomeMethod(); err != nil {
    return api.NewPluginError("plugin_name", "operation", "description", err)
}
```

## Best Practices

### 1. Configuration Management
- Validate configuration during initialization
- Provide sensible defaults
- Support runtime configuration updates
- Document all configuration options

### 2. Error Handling
- Use plugin-specific error types
- Provide meaningful error messages
- Don't fail the entire request for optional features
- Log errors appropriately

### 3. Performance
- Cache frequently used data
- Minimize external dependencies
- Use efficient algorithms for date calculations
- Consider memory usage for large datasets

### 4. Localization
- Support multiple locale formats
- Provide fallback mechanisms
- Consider cultural variations
- Test with native speakers

### 5. Regional Variations
- Research authentic regional practices
- Consult with subject matter experts
- Provide configuration for local variations
- Document sources and reasoning

## Testing

### Unit Tests
```go
func TestTamilNaduExtension(t *testing.T) {
    plugin := examples.NewTamilNaduExtension()
    
    // Test initialization
    err := plugin.Initialize(context.Background(), map[string]interface{}{})
    assert.NoError(t, err)
    assert.True(t, plugin.IsEnabled())
    
    // Test regional rules
    data := &api.PanchangamData{Region: api.RegionTamilNadu}
    err = plugin.ApplyRegionalRules(context.Background(), data)
    assert.NoError(t, err)
}
```

### Integration Tests
```go
func TestAPIWithPlugins(t *testing.T) {
    api := api.NewCorePanchangamAPI(observer)
    
    // Register plugins
    api.RegisterPlugin(examples.NewTamilNaduExtension())
    api.RegisterPlugin(examples.NewTamilLocalizationPlugin())
    
    // Test with Tamil Nadu region and locale
    req := api.PanchangamRequest{
        Date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
        Location: chennaiLocation,
        Region:   api.RegionTamilNadu,
        Locale:   "ta",
    }
    
    result, err := api.GetPanchangam(context.Background(), req)
    assert.NoError(t, err)
    assert.Equal(t, "பிரதமை", result.Tithi.NameLocal)
}
```

## Contributing

When creating new plugins:

1. Follow the established interface patterns
2. Provide comprehensive documentation
3. Include unit and integration tests
4. Consider cultural authenticity and accuracy
5. Optimize for performance and memory usage
6. Handle errors gracefully
7. Support configuration and customization

## Future Extensions

Potential areas for additional plugins:

- **Jyotish Calculations**: Planetary positions and astrological analysis
- **Vastu Integration**: Directional and timing recommendations
- **Audio Pronunciation**: Voice synthesis for Sanskrit/regional names
- **Calendar Integration**: Export to modern calendar applications
- **Notification System**: Alerts for upcoming festivals and muhurtas
- **Historical Data**: Archaeological and historical festival information
- **Medical Astrology**: Ayurvedic timing recommendations
- **Agricultural Calendar**: Farming and seasonal guidance