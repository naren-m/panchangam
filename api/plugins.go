package api

import (
	"context"
	"fmt"
	"time"
)

// PluginInfo provides metadata about a plugin
type PluginInfo struct {
	Name         string                 `json:"name"`
	Version      Version                `json:"version"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author,omitempty"`
	Website      string                 `json:"website,omitempty"`
	License      string                 `json:"license,omitempty"`
	Capabilities []string               `json:"capabilities"`
	Dependencies []string               `json:"dependencies,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// Plugin represents the base plugin interface
type Plugin interface {
	// GetInfo returns plugin metadata
	GetInfo() PluginInfo

	// Initialize sets up the plugin with configuration
	Initialize(ctx context.Context, config map[string]interface{}) error

	// IsEnabled returns whether the plugin is currently enabled
	IsEnabled() bool

	// Shutdown cleans up plugin resources
	Shutdown(ctx context.Context) error
}

// CalculationPlugin defines the interface for custom calculation plugins
type CalculationPlugin interface {
	Plugin

	// GetSupportedMethods returns calculation methods this plugin supports
	GetSupportedMethods() []CalculationMethod

	// GetSupportedRegions returns regions this plugin supports
	GetSupportedRegions() []Region

	// CalculateTithi calculates tithi for given parameters
	CalculateTithi(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Tithi, error)

	// CalculateNakshatra calculates nakshatra for given parameters
	CalculateNakshatra(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Nakshatra, error)

	// CalculateYoga calculates yoga for given parameters
	CalculateYoga(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Yoga, error)

	// CalculateKarana calculates karana for given parameters
	CalculateKarana(ctx context.Context, date time.Time, location Location, method CalculationMethod) (*Karana, error)

	// CalculateSunMoonTimes calculates sun and moon timings
	CalculateSunMoonTimes(ctx context.Context, date time.Time, location Location) (*SunMoonTimes, error)
}

// EventPlugin defines the interface for event and festival plugins
type EventPlugin interface {
	Plugin

	// GetSupportedRegions returns regions this plugin supports
	GetSupportedRegions() []Region

	// GetSupportedEventTypes returns event types this plugin can generate
	GetSupportedEventTypes() []EventType

	// GetEvents returns events for a specific date and location
	GetEvents(ctx context.Context, date time.Time, location Location, region Region) ([]Event, error)

	// GetEventsInRange returns events for a date range
	GetEventsInRange(ctx context.Context, start, end time.Time, location Location, region Region) ([]Event, error)
}

// MuhurtaPlugin defines the interface for muhurta calculation plugins
type MuhurtaPlugin interface {
	Plugin

	// GetSupportedRegions returns regions this plugin supports
	GetSupportedRegions() []Region

	// GetMuhurtas returns muhurtas for a specific date and location
	GetMuhurtas(ctx context.Context, date time.Time, location Location, region Region) ([]Muhurta, error)

	// FindAuspiciousTimes finds auspicious times for specific activities
	FindAuspiciousTimes(ctx context.Context, date time.Time, location Location, activities []string) ([]Muhurta, error)

	// IsTimeAuspicious checks if a specific time is auspicious for given activities
	IsTimeAuspicious(ctx context.Context, dateTime time.Time, location Location, activities []string) (bool, string, error)
}

// LocalizationPlugin defines the interface for localization plugins
type LocalizationPlugin interface {
	Plugin

	// GetSupportedLocales returns supported locale codes
	GetSupportedLocales() []string

	// GetSupportedRegions returns regions this plugin supports
	GetSupportedRegions() []Region

	// LocalizeTithi returns localized tithi information
	LocalizeTithi(tithi *Tithi, locale string, region Region) error

	// LocalizeNakshatra returns localized nakshatra information
	LocalizeNakshatra(nakshatra *Nakshatra, locale string, region Region) error

	// LocalizeYoga returns localized yoga information
	LocalizeYoga(yoga *Yoga, locale string, region Region) error

	// LocalizeKarana returns localized karana information
	LocalizeKarana(karana *Karana, locale string, region Region) error

	// LocalizeEvent returns localized event information
	LocalizeEvent(event *Event, locale string, region Region) error

	// LocalizeMuhurta returns localized muhurta information
	LocalizeMuhurta(muhurta *Muhurta, locale string, region Region) error
}

// ValidationPlugin defines the interface for data validation plugins
type ValidationPlugin interface {
	Plugin

	// ValidateRequest validates a Panchangam request
	ValidateRequest(ctx context.Context, req PanchangamRequest) error

	// ValidateLocation validates location data
	ValidateLocation(ctx context.Context, location Location) error

	// ValidatePanchangamData validates calculated Panchangam data
	ValidatePanchangamData(ctx context.Context, data *PanchangamData) error

	// GetValidationRules returns validation rules for a specific region
	GetValidationRules(region Region) map[string]interface{}
}

// CachePlugin defines the interface for caching plugins
type CachePlugin interface {
	Plugin

	// Get retrieves cached data by key
	Get(ctx context.Context, key string) (interface{}, error)

	// Set stores data in cache with optional TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete removes data from cache
	Delete(ctx context.Context, key string) error

	// Clear clears all cached data
	Clear(ctx context.Context) error

	// GetStats returns cache statistics
	GetStats(ctx context.Context) map[string]interface{}
}

// PluginManager manages plugin lifecycle and orchestration
type PluginManager interface {
	// RegisterPlugin registers a new plugin
	RegisterPlugin(plugin Plugin) error

	// UnregisterPlugin removes a plugin
	UnregisterPlugin(name string) error

	// GetPlugin retrieves a plugin by name
	GetPlugin(name string) (Plugin, error)

	// GetPluginsByType retrieves all plugins of a specific type
	GetPluginsByType(pluginType string) []Plugin

	// GetPluginsByCapability returns plugins that have a specific capability
	GetPluginsByCapability(capability PluginCapability) []Plugin

	// ListPlugins returns all registered plugins
	ListPlugins() []PluginInfo

	// EnablePlugin enables a plugin
	EnablePlugin(ctx context.Context, name string) error

	// DisablePlugin disables a plugin
	DisablePlugin(ctx context.Context, name string) error

	// InitializeAll initializes all registered plugins
	InitializeAll(ctx context.Context) error

	// ShutdownAll shuts down all plugins
	ShutdownAll(ctx context.Context) error
}

// PluginRegistry provides plugin discovery and management
type PluginRegistry interface {
	// DiscoverPlugins discovers plugins from specified directories
	DiscoverPlugins(ctx context.Context, directories []string) error

	// LoadPlugin loads a plugin from a file or directory
	LoadPlugin(ctx context.Context, path string) (Plugin, error)

	// GetAvailablePlugins returns all available plugins
	GetAvailablePlugins() []PluginInfo

	// InstallPlugin installs a plugin from a package
	InstallPlugin(ctx context.Context, packagePath string) error

	// UninstallPlugin removes an installed plugin
	UninstallPlugin(ctx context.Context, name string) error

	// UpdatePlugin updates a plugin to a newer version
	UpdatePlugin(ctx context.Context, name string) error
}

// ExtensionPoint defines points where plugins can extend functionality
type ExtensionPoint interface {
	// GetName returns the extension point name
	GetName() string

	// GetDescription returns a description of what this extension point does
	GetDescription() string

	// GetRequiredInterface returns the interface that plugins must implement
	GetRequiredInterface() string

	// Execute executes all registered plugins for this extension point
	Execute(ctx context.Context, data interface{}) (interface{}, error)

	// RegisterPlugin registers a plugin for this extension point
	RegisterPlugin(plugin Plugin) error

	// UnregisterPlugin removes a plugin from this extension point
	UnregisterPlugin(name string) error
}

// RegionalExtension defines region-specific calculation extensions
type RegionalExtension interface {
	Plugin

	// GetRegion returns the region this extension supports
	GetRegion() Region

	// GetCalendarSystem returns the calendar system used
	GetCalendarSystem() CalendarSystem

	// ApplyRegionalRules applies region-specific rules to Panchangam data
	ApplyRegionalRules(ctx context.Context, data *PanchangamData) error

	// GetRegionalEvents returns region-specific events
	GetRegionalEvents(ctx context.Context, date time.Time, location Location) ([]Event, error)

	// GetRegionalMuhurtas returns region-specific muhurtas
	GetRegionalMuhurtas(ctx context.Context, date time.Time, location Location) ([]Muhurta, error)

	// GetRegionalNames returns localized names for Panchangam elements
	GetRegionalNames(locale string) map[string]string
}

// PluginConfig represents plugin configuration
type PluginConfig struct {
	Name     string                 `json:"name"`
	Enabled  bool                   `json:"enabled"`
	Priority int                    `json:"priority,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
}

// PluginError represents plugin-specific errors
type PluginError struct {
	PluginName string
	Operation  string
	Message    string
	Cause      error
}

func (e *PluginError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("plugin %s error in %s: %s (caused by: %v)", e.PluginName, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("plugin %s error in %s: %s", e.PluginName, e.Operation, e.Message)
}

func (e *PluginError) Unwrap() error {
	return e.Cause
}

// NewPluginError creates a new plugin error
func NewPluginError(pluginName, operation, message string, cause error) *PluginError {
	return &PluginError{
		PluginName: pluginName,
		Operation:  operation,
		Message:    message,
		Cause:      cause,
	}
}

// PluginCapability represents different plugin capabilities
type PluginCapability string

const (
	CapabilityCalculation  PluginCapability = "calculation"
	CapabilityEvent        PluginCapability = "event"
	CapabilityMuhurta      PluginCapability = "muhurta"
	CapabilityLocalization PluginCapability = "localization"
	CapabilityValidation   PluginCapability = "validation"
	CapabilityCache        PluginCapability = "cache"
	CapabilityRegional     PluginCapability = "regional"
	CapabilityData         PluginCapability = "data"
	CapabilityNotification PluginCapability = "notification"
	CapabilityExport       PluginCapability = "export"
)

// PluginMetadata provides additional plugin metadata
type PluginMetadata struct {
	Tags          []string `json:"tags,omitempty"`
	Categories    []string `json:"categories,omitempty"`
	SupportURL    string   `json:"support_url,omitempty"`
	DocsURL       string   `json:"docs_url,omitempty"`
	SourceURL     string   `json:"source_url,omitempty"`
	MinAPIVersion Version  `json:"min_api_version"`
	MaxAPIVersion Version  `json:"max_api_version,omitempty"`
	Experimental  bool     `json:"experimental,omitempty"`
	Deprecated    bool     `json:"deprecated,omitempty"`
	ReleaseNotes  string   `json:"release_notes,omitempty"`
}
