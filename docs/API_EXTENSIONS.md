# API Extensions and Versioning

## Overview

The Panchangam API provides extensibility through a plugin architecture and maintains backward compatibility through semantic versioning.

## API Versioning

### Semantic Versioning

We follow [Semantic Versioning 2.0.0](https://semver.org/):

```
MAJOR.MINOR.PATCH

Example: 2.1.3
```

- **MAJOR**: Breaking changes (incompatible API changes)
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Version Structure

```go
package api

type Version struct {
    Major int `json:"major"`
    Minor int `json:"minor"`
    Patch int `json:"patch"`
}

func (v Version) String() string {
    return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) GreaterThan(other Version) bool {
    if v.Major != other.Major {
        return v.Major > other.Major
    }
    if v.Minor != other.Minor {
        return v.Minor > other.Minor
    }
    return v.Patch > other.Patch
}
```

### Current API Versions

| Version | Status | Release Date | End of Support |
|---------|--------|--------------|----------------|
| v1.0.x  | Deprecated | 2023-01-15 | 2024-12-31 |
| v2.0.x  | Stable | 2024-01-15 | TBD |
| v2.1.x  | Current | 2024-06-15 | TBD |
| v3.0.x  | Beta | 2024-11-15 | TBD |

### Version Header

All API requests should include version header:

```
X-API-Version: 2.1.0
```

### Version Negotiation

```go
func NegotiateVersion(requestedVersion string, supportedVersions []string) (string, error) {
    requested, err := ParseVersion(requestedVersion)
    if err != nil {
        return "", err
    }

    // Find best compatible version
    for _, sv := range supportedVersions {
        supported, _ := ParseVersion(sv)
        if supported.Major == requested.Major && supported.Minor >= requested.Minor {
            return sv, nil
        }
    }

    return "", fmt.Errorf("unsupported API version: %s", requestedVersion)
}
```

## Plugin Architecture

### Plugin Interface

```go
package api

import "context"

// PanchangamPlugin is the base interface for all plugins
type PanchangamPlugin interface {
    // GetName returns the plugin name
    GetName() string

    // GetVersion returns the plugin version
    GetVersion() Version

    // ProcessPanchangamData processes panchangam data
    ProcessPanchangamData(ctx context.Context, data *PanchangamData) (*ProcessedData, error)
}

// PanchangamData holds the raw panchangam data
type PanchangamData struct {
    Date     time.Time
    Location Location
    Tithi    *Tithi
    Nakshatra *Nakshatra
    Yoga     *Yoga
    Karana   *Karana
    Vara     *Vara
    Sunrise  time.Time
    Sunset   time.Time
}

// ProcessedData holds the processed data with plugin modifications
type ProcessedData struct {
    OriginalData  *PanchangamData
    Modifications map[string]interface{}
    Metadata      map[string]string
}
```

### Plugin Types

#### 1. Calculation Method Plugin

Modifies how astronomical calculations are performed.

```go
type CalculationMethodPlugin interface {
    PanchangamPlugin

    // GetCalculationMethod returns the method name (e.g., "drik", "vakya")
    GetCalculationMethod() string

    // ModifyTithiCalculation modifies tithi calculation
    ModifyTithiCalculation(ctx context.Context, tithi *Tithi) (*Tithi, error)

    // ModifyNakshatraCalculation modifies nakshatra calculation
    ModifyNakshatraCalculation(ctx context.Context, nakshatra *Nakshatra) (*Nakshatra, error)
}
```

**Example Implementation:**

```go
type VakyaMethodPlugin struct {
    name    string
    version Version
}

func (p *VakyaMethodPlugin) GetName() string {
    return p.name
}

func (p *VakyaMethodPlugin) GetVersion() Version {
    return p.version
}

func (p *VakyaMethodPlugin) GetCalculationMethod() string {
    return "vakya"
}

func (p *VakyaMethodPlugin) ProcessPanchangamData(
    ctx context.Context,
    data *PanchangamData,
) (*ProcessedData, error) {
    // Apply Vakya method adjustments
    modifications := make(map[string]interface{})

    // Adjust for mean vs true positions
    adjustedTithi, err := p.ModifyTithiCalculation(ctx, data.Tithi)
    if err != nil {
        return nil, err
    }

    modifications["tithi"] = adjustedTithi
    modifications["method"] = "vakya"

    return &ProcessedData{
        OriginalData:  data,
        Modifications: modifications,
    }, nil
}
```

#### 2. Regional Plugin

Handles regional variations and customizations.

```go
type RegionalPlugin interface {
    PanchangamPlugin

    // GetRegion returns the region name
    GetRegion() string

    // GetCalendarSystem returns the calendar system (amanta/purnimanta)
    GetCalendarSystem() string

    // GetFestivals returns regional festivals for a year
    GetFestivals(year int) ([]Festival, error)

    // AdjustForRegion adjusts calculations for regional variations
    AdjustForRegion(ctx context.Context, data *PanchangamData) (*PanchangamData, error)
}
```

**Example Implementation:**

```go
type TamilNaduPlugin struct {
    name    string
    version Version
}

func (p *TamilNaduPlugin) GetRegion() string {
    return "tamil_nadu"
}

func (p *TamilNaduPlugin) GetCalendarSystem() string {
    return "amanta"
}

func (p *TamilNaduPlugin) GetFestivals(year int) ([]Festival, error) {
    festivals := []Festival{
        {
            Name:  "Pongal",
            Date:  calculatePongalDate(year),
            Type:  "solar",
            Region: "tamil_nadu",
        },
        {
            Name:  "Tamil New Year",
            Date:  time.Date(year, 4, 14, 0, 0, 0, 0, time.UTC),
            Type:  "solar",
            Region: "tamil_nadu",
        },
    }

    return festivals, nil
}

func (p *TamilNaduPlugin) ProcessPanchangamData(
    ctx context.Context,
    data *PanchangamData,
) (*ProcessedData, error) {
    // Apply Tamil Nadu specific adjustments
    adjusted, err := p.AdjustForRegion(ctx, data)
    if err != nil {
        return nil, err
    }

    return &ProcessedData{
        OriginalData:  data,
        Modifications: map[string]interface{}{
            "adjusted_data": adjusted,
            "region":        "tamil_nadu",
        },
    }, nil
}
```

#### 3. Muhurta Plugin

Calculates auspicious timings.

```go
type MuhurtaPlugin interface {
    PanchangamPlugin

    // CalculateMuhurta calculates auspicious timings
    CalculateMuhurta(ctx context.Context, date time.Time, location Location, purpose string) ([]Muhurta, error)

    // GetAuspiciousHours returns auspicious hours for a day
    GetAuspiciousHours(ctx context.Context, date time.Time) ([]TimeRange, error)

    // GetInauspiciousPeriods returns inauspicious periods (Rahu Kalam, etc.)
    GetInauspiciousPeriods(ctx context.Context, date time.Time) ([]InauspiciousPeriod, error)
}
```

**Example Implementation:**

```go
type StandardMuhurtaPlugin struct {
    name    string
    version Version
}

func (p *StandardMuhurtaPlugin) CalculateMuhurta(
    ctx context.Context,
    date time.Time,
    location Location,
    purpose string,
) ([]Muhurta, error) {
    // Calculate based on purpose
    switch purpose {
    case "marriage":
        return p.calculateMarriageMuhurta(ctx, date, location)
    case "griha_pravesh":
        return p.calculateGrihaPraveshMuhurta(ctx, date, location)
    default:
        return p.calculateGeneralMuhurta(ctx, date, location)
    }
}

func (p *StandardMuhurtaPlugin) GetInauspiciousPeriods(
    ctx context.Context,
    date time.Time,
) ([]InauspiciousPeriod, error) {
    periods := []InauspiciousPeriod{
        {
            Name:  "Rahu Kalam",
            Start: calculateRahuKalamStart(date),
            End:   calculateRahuKalamEnd(date),
        },
        {
            Name:  "Yamaganda Kalam",
            Start: calculateYamagandaStart(date),
            End:   calculateYamagandaEnd(date),
        },
    }

    return periods, nil
}
```

#### 4. Festival Plugin

Manages festival calculations.

```go
type AdvancedFestivalPlugin interface {
    PanchangamPlugin

    // GetFestivalDates calculates festival dates for a year
    GetFestivalDates(year int, region string, locale string) ([]Festival, error)

    // ValidateFestivalDate validates if a date is a festival
    ValidateFestivalDate(date time.Time, festival string, region string) (bool, error)

    // CalculateFestivalMuhurta calculates auspicious timing for festival
    CalculateFestivalMuhurta(festival Festival, location Location) (*Muhurta, error)
}
```

### Plugin Manager

```go
type PluginManager struct {
    plugins map[string]PanchangamPlugin
    order   []string
}

func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]PanchangamPlugin),
        order:   make([]string, 0),
    }
}

func (pm *PluginManager) RegisterPlugin(name string, plugin PanchangamPlugin) error {
    if _, exists := pm.plugins[name]; exists {
        return fmt.Errorf("plugin %s already registered", name)
    }

    pm.plugins[name] = plugin
    pm.order = append(pm.order, name)

    return nil
}

func (pm *PluginManager) UnregisterPlugin(name string) {
    delete(pm.plugins, name)

    // Remove from order
    for i, n := range pm.order {
        if n == name {
            pm.order = append(pm.order[:i], pm.order[i+1:]...)
            break
        }
    }
}

func (pm *PluginManager) ProcessPanchangamData(
    ctx context.Context,
    data *PanchangamData,
) (*ProcessedData, error) {
    result := &ProcessedData{
        OriginalData:  data,
        Modifications: make(map[string]interface{}),
        Metadata:      make(map[string]string),
    }

    // Process through plugins in order
    for _, name := range pm.order {
        plugin := pm.plugins[name]

        processed, err := plugin.ProcessPanchangamData(ctx, data)
        if err != nil {
            return nil, fmt.Errorf("plugin %s failed: %w", name, err)
        }

        // Merge modifications
        for k, v := range processed.Modifications {
            result.Modifications[k] = v
        }

        // Add plugin metadata
        result.Metadata[name] = plugin.GetVersion().String()
    }

    return result, nil
}

func (pm *PluginManager) GetPlugin(name string) (PanchangamPlugin, bool) {
    plugin, exists := pm.plugins[name]
    return plugin, exists
}

func (pm *PluginManager) ListPlugins() []string {
    return append([]string{}, pm.order...)
}
```

## Extension Points

### 1. Custom Ephemeris Provider

```go
type CustomEphemerisProvider struct {
    dataSource string
}

func (p *CustomEphemerisProvider) GetPlanetaryPositions(
    ctx context.Context,
    jd ephemeris.JulianDay,
) (*ephemeris.PlanetaryPositions, error) {
    // Custom implementation
    return &ephemeris.PlanetaryPositions{}, nil
}

func (p *CustomEphemerisProvider) IsAvailable(ctx context.Context) bool {
    return true
}

func (p *CustomEphemerisProvider) GetProviderName() string {
    return "custom-provider"
}
```

### 2. Custom Ayanamsa Calculation

```go
type CustomAyanamasaPlugin struct {
    name string
}

func (p *CustomAyanamasaPlugin) CalculateAyanamsa(jd float64) float64 {
    // Custom ayanamsa calculation
    // Example: Lahiri, Raman, KP, etc.
    return customValue
}

func (p *CustomAyanamasaPlugin) ProcessPanchangamData(
    ctx context.Context,
    data *PanchangamData,
) (*ProcessedData, error) {
    jd := ephemeris.TimeToJulianDay(data.Date)
    ayanamsa := p.CalculateAyanamsa(float64(jd))

    return &ProcessedData{
        OriginalData: data,
        Modifications: map[string]interface{}{
            "ayanamsa": ayanamsa,
        },
    }, nil
}
```

### 3. Custom Localization

```go
type LocalizationPlugin struct {
    locale string
}

func (p *LocalizationPlugin) TranslateTithi(tithi string) string {
    translations := map[string]map[string]string{
        "hi": {
            "Pratipada":   "प्रतिपदा",
            "Dwitiya":     "द्वितीया",
            // ... more translations
        },
        "ta": {
            "Pratipada":   "பிரதமை",
            "Dwitiya":     "துவிதியை",
            // ... more translations
        },
    }

    if trans, ok := translations[p.locale][tithi]; ok {
        return trans
    }

    return tithi
}

func (p *LocalizationPlugin) ProcessPanchangamData(
    ctx context.Context,
    data *PanchangamData,
) (*ProcessedData, error) {
    return &ProcessedData{
        OriginalData: data,
        Modifications: map[string]interface{}{
            "tithi_localized":     p.TranslateTithi(data.Tithi.Name),
            "nakshatra_localized": p.TranslateNakshatra(data.Nakshatra.Name),
        },
    }, nil
}
```

## Usage Examples

### Basic Plugin Usage

```go
// Create manager
manager := api.NewPluginManager()

// Register plugins
manager.RegisterPlugin("tamil", &TamilNaduPlugin{
    name:    "tamil-nadu-plugin",
    version: api.Version{Major: 1, Minor: 0, Patch: 0},
})

manager.RegisterPlugin("muhurta", &StandardMuhurtaPlugin{
    name:    "standard-muhurta",
    version: api.Version{Major: 1, Minor: 0, Patch: 0},
})

// Process data
result, err := manager.ProcessPanchangamData(ctx, data)
if err != nil {
    log.Fatal(err)
}

// Access modifications
fmt.Printf("Processed data: %+v\n", result.Modifications)
```

### Version-Aware Client

```go
type VersionedClient struct {
    apiVersion Version
    client     *Client
}

func (vc *VersionedClient) GetPanchangam(
    ctx context.Context,
    req *Request,
) (*Response, error) {
    // Add version header
    ctx = metadata.AppendToOutgoingContext(
        ctx,
        "x-api-version", vc.apiVersion.String(),
    )

    return vc.client.Get(ctx, req)
}
```

## Migration Guide

### Migrating from v1 to v2

#### Breaking Changes

1. **Request/Response Structure**
   ```go
   // v1
   type Request struct {
       Date string
       Lat  float64
       Lon  float64
   }

   // v2
   type Request struct {
       Date     string
       Location Location  // Changed to struct
       Options  Options   // New field
   }
   ```

2. **Plugin Interface**
   ```go
   // v1
   type Plugin interface {
       Process(data Data) Data
   }

   // v2
   type Plugin interface {
       ProcessPanchangamData(ctx context.Context, data *PanchangamData) (*ProcessedData, error)
   }
   ```

#### Migration Steps

1. Update imports
2. Modify struct definitions
3. Add context parameters
4. Update error handling
5. Test thoroughly

## Best Practices

### Plugin Development

1. **Versioning**: Always version your plugins
2. **Error Handling**: Return descriptive errors
3. **Context**: Respect context cancellation
4. **Testing**: Write comprehensive tests
5. **Documentation**: Document plugin behavior

### API Compatibility

1. **Deprecation**: Announce deprecations early
2. **Support Window**: Maintain old versions for 12 months
3. **Migration Path**: Provide clear migration guides
4. **Version Detection**: Support version negotiation

### Performance

1. **Caching**: Cache expensive calculations
2. **Concurrency**: Use goroutines appropriately
3. **Resource Cleanup**: Always clean up resources
4. **Monitoring**: Add observability

## Security Considerations

1. **Input Validation**: Validate all plugin inputs
2. **Resource Limits**: Set limits on plugin execution
3. **Sandboxing**: Consider plugin sandboxing for untrusted code
4. **Audit Logging**: Log plugin activities

---

*Last Updated: 2025-11-18*
*Maintainer: Panchangam Development Team*
