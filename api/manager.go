package api

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/naren-m/panchangam/observability"
)

// DefaultPluginManager implements the PluginManager interface
type DefaultPluginManager struct {
	plugins         map[string]Plugin
	pluginConfigs   map[string]PluginConfig
	extensionPoints map[string]ExtensionPoint
	mu              sync.RWMutex
	logger          *observability.ErrorRecorder
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *DefaultPluginManager {
	return &DefaultPluginManager{
		plugins:         make(map[string]Plugin),
		pluginConfigs:   make(map[string]PluginConfig),
		extensionPoints: make(map[string]ExtensionPoint),
		logger:          observability.NewErrorRecorder(),
	}
}

// RegisterPlugin registers a new plugin
func (pm *DefaultPluginManager) RegisterPlugin(plugin Plugin) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	info := plugin.GetInfo()
	if _, exists := pm.plugins[info.Name]; exists {
		return NewPluginError(info.Name, "register", "plugin already registered", nil)
	}

	pm.plugins[info.Name] = plugin
	pm.pluginConfigs[info.Name] = PluginConfig{
		Name:    info.Name,
		Enabled: true,
		Config:  make(map[string]interface{}),
	}

	return nil
}

// UnregisterPlugin removes a plugin
func (pm *DefaultPluginManager) UnregisterPlugin(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return NewPluginError(name, "unregister", "plugin not found", nil)
	}

	// Shutdown plugin if it's running
	if plugin.IsEnabled() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := plugin.Shutdown(ctx); err != nil {
			return NewPluginError(name, "unregister", "failed to shutdown plugin", err)
		}
	}

	delete(pm.plugins, name)
	delete(pm.pluginConfigs, name)
	return nil
}

// GetPlugin retrieves a plugin by name
func (pm *DefaultPluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, NewPluginError(name, "get", "plugin not found", nil)
	}

	return plugin, nil
}

// GetPluginsByType retrieves all plugins of a specific type
func (pm *DefaultPluginManager) GetPluginsByType(pluginType string) []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var result []Plugin
	for _, plugin := range pm.plugins {
		info := plugin.GetInfo()
		for _, capability := range info.Capabilities {
			if capability == pluginType {
				result = append(result, plugin)
				break
			}
		}
	}

	return result
}

// ListPlugins returns all registered plugins
func (pm *DefaultPluginManager) ListPlugins() []PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var result []PluginInfo
	for _, plugin := range pm.plugins {
		result = append(result, plugin.GetInfo())
	}

	// Sort by name for consistent output
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// EnablePlugin enables a plugin
func (pm *DefaultPluginManager) EnablePlugin(ctx context.Context, name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return NewPluginError(name, "enable", "plugin not found", nil)
	}

	if plugin.IsEnabled() {
		return nil // Already enabled
	}

	config := pm.pluginConfigs[name]
	if err := plugin.Initialize(ctx, config.Config); err != nil {
		return NewPluginError(name, "enable", "failed to initialize plugin", err)
	}

	config.Enabled = true
	pm.pluginConfigs[name] = config

	return nil
}

// DisablePlugin disables a plugin
func (pm *DefaultPluginManager) DisablePlugin(ctx context.Context, name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return NewPluginError(name, "disable", "plugin not found", nil)
	}

	if !plugin.IsEnabled() {
		return nil // Already disabled
	}

	if err := plugin.Shutdown(ctx); err != nil {
		return NewPluginError(name, "disable", "failed to shutdown plugin", err)
	}

	config := pm.pluginConfigs[name]
	config.Enabled = false
	pm.pluginConfigs[name] = config

	return nil
}

// InitializeAll initializes all registered plugins
func (pm *DefaultPluginManager) InitializeAll(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var errors []error
	for name, plugin := range pm.plugins {
		config := pm.pluginConfigs[name]
		if config.Enabled {
			if err := plugin.Initialize(ctx, config.Config); err != nil {
				errors = append(errors, NewPluginError(name, "initialize", "failed to initialize", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to initialize %d plugins: %v", len(errors), errors)
	}

	return nil
}

// ShutdownAll shuts down all plugins
func (pm *DefaultPluginManager) ShutdownAll(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var errors []error
	for name, plugin := range pm.plugins {
		if plugin.IsEnabled() {
			if err := plugin.Shutdown(ctx); err != nil {
				errors = append(errors, NewPluginError(name, "shutdown", "failed to shutdown", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to shutdown %d plugins: %v", len(errors), errors)
	}

	return nil
}

// SetPluginConfig sets configuration for a plugin
func (pm *DefaultPluginManager) SetPluginConfig(name string, config PluginConfig) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.plugins[name]; !exists {
		return NewPluginError(name, "config", "plugin not found", nil)
	}

	pm.pluginConfigs[name] = config
	return nil
}

// GetPluginConfig gets configuration for a plugin
func (pm *DefaultPluginManager) GetPluginConfig(name string) (PluginConfig, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	config, exists := pm.pluginConfigs[name]
	if !exists {
		return PluginConfig{}, NewPluginError(name, "config", "plugin config not found", nil)
	}

	return config, nil
}

// RegisterExtensionPoint registers a new extension point
func (pm *DefaultPluginManager) RegisterExtensionPoint(extensionPoint ExtensionPoint) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := extensionPoint.GetName()
	if _, exists := pm.extensionPoints[name]; exists {
		return fmt.Errorf("extension point %s already registered", name)
	}

	pm.extensionPoints[name] = extensionPoint
	return nil
}

// GetExtensionPoint retrieves an extension point by name
func (pm *DefaultPluginManager) GetExtensionPoint(name string) (ExtensionPoint, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	extensionPoint, exists := pm.extensionPoints[name]
	if !exists {
		return nil, fmt.Errorf("extension point %s not found", name)
	}

	return extensionPoint, nil
}

// ExecuteExtensionPoint executes all plugins registered for an extension point
func (pm *DefaultPluginManager) ExecuteExtensionPoint(ctx context.Context, name string, data interface{}) (interface{}, error) {
	extensionPoint, err := pm.GetExtensionPoint(name)
	if err != nil {
		return nil, err
	}

	return extensionPoint.Execute(ctx, data)
}

// GetPluginsByCapability returns plugins that have a specific capability
func (pm *DefaultPluginManager) GetPluginsByCapability(capability PluginCapability) []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var result []Plugin
	for _, plugin := range pm.plugins {
		info := plugin.GetInfo()
		for _, cap := range info.Capabilities {
			if cap == string(capability) {
				result = append(result, plugin)
				break
			}
		}
	}

	return result
}

// GetEnabledPlugins returns all currently enabled plugins
func (pm *DefaultPluginManager) GetEnabledPlugins() []Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var result []Plugin
	for name, plugin := range pm.plugins {
		config := pm.pluginConfigs[name]
		if config.Enabled && plugin.IsEnabled() {
			result = append(result, plugin)
		}
	}

	return result
}

// GetPluginStats returns statistics about plugins
func (pm *DefaultPluginManager) GetPluginStats() map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	stats := map[string]interface{}{
		"total_plugins":   len(pm.plugins),
		"enabled_plugins": 0,
		"plugin_types":    make(map[string]int),
	}

	for name, plugin := range pm.plugins {
		config := pm.pluginConfigs[name]
		if config.Enabled && plugin.IsEnabled() {
			stats["enabled_plugins"] = stats["enabled_plugins"].(int) + 1
		}

		info := plugin.GetInfo()
		for _, capability := range info.Capabilities {
			if count, exists := stats["plugin_types"].(map[string]int)[capability]; exists {
				stats["plugin_types"].(map[string]int)[capability] = count + 1
			} else {
				stats["plugin_types"].(map[string]int)[capability] = 1
			}
		}
	}

	return stats
}

// ValidatePluginDependencies checks if all plugin dependencies are satisfied
func (pm *DefaultPluginManager) ValidatePluginDependencies() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, plugin := range pm.plugins {
		info := plugin.GetInfo()
		for _, dependency := range info.Dependencies {
			if _, exists := pm.plugins[dependency]; !exists {
				return NewPluginError(info.Name, "validate", fmt.Sprintf("missing dependency: %s", dependency), nil)
			}
		}
	}

	return nil
}

// LoadPluginsFromConfig loads and configures plugins from configuration
func (pm *DefaultPluginManager) LoadPluginsFromConfig(ctx context.Context, configs []PluginConfig) error {
	for _, config := range configs {
		if err := pm.SetPluginConfig(config.Name, config); err != nil {
			return fmt.Errorf("failed to set config for plugin %s: %w", config.Name, err)
		}

		if config.Enabled {
			if err := pm.EnablePlugin(ctx, config.Name); err != nil {
				return fmt.Errorf("failed to enable plugin %s: %w", config.Name, err)
			}
		}
	}

	return nil
}

// Health check for plugin manager
func (pm *DefaultPluginManager) Health(ctx context.Context) map[string]interface{} {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"plugins":   make(map[string]interface{}),
	}

	unhealthyCount := 0
	for name, plugin := range pm.plugins {
		config := pm.pluginConfigs[name]
		pluginHealth := map[string]interface{}{
			"enabled":    config.Enabled,
			"is_running": plugin.IsEnabled(),
			"info":       plugin.GetInfo(),
		}

		if config.Enabled && !plugin.IsEnabled() {
			pluginHealth["status"] = "unhealthy"
			unhealthyCount++
		} else {
			pluginHealth["status"] = "healthy"
		}

		health["plugins"].(map[string]interface{})[name] = pluginHealth
	}

	if unhealthyCount > 0 {
		health["status"] = "degraded"
		health["unhealthy_plugins"] = unhealthyCount
	}

	return health
}
