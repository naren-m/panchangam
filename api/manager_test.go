package api

import (
	"context"
	"errors"
	"testing"
)

type managerInfoPlugin struct {
	info    PluginInfo
	enabled bool
}

func (p managerInfoPlugin) GetInfo() PluginInfo {
	return p.info
}

func (p managerInfoPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	return nil
}

func (p managerInfoPlugin) IsEnabled() bool {
	return p.enabled
}

func (p managerInfoPlugin) Shutdown(ctx context.Context) error {
	return nil
}

type failingLifecyclePlugin struct {
	name        string
	enabled     bool
	initErr     error
	shutdownErr error
}

func (p failingLifecyclePlugin) GetInfo() PluginInfo {
	return PluginInfo{
		Name:         p.name,
		Version:      Version{Major: 1},
		Capabilities: []string{string(CapabilityCalculation)},
	}
}

func (p failingLifecyclePlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	return p.initErr
}

func (p failingLifecyclePlugin) IsEnabled() bool {
	return p.enabled
}

func (p failingLifecyclePlugin) Shutdown(ctx context.Context) error {
	return p.shutdownErr
}

func TestPluginManagerInitializeAllPreservesErrorCause(t *testing.T) {
	initErr := errors.New("init failed")
	manager := NewPluginManager()
	if err := manager.RegisterPlugin(failingLifecyclePlugin{name: "bad_init", initErr: initErr}); err != nil {
		t.Fatalf("failed to register plugin: %v", err)
	}

	err := manager.InitializeAll(context.Background())
	if err == nil {
		t.Fatal("expected initialize error")
	}
	if !errors.Is(err, initErr) {
		t.Fatalf("expected initialize error to preserve cause, got %v", err)
	}
}

func TestPluginManagerShutdownAllPreservesErrorCause(t *testing.T) {
	shutdownErr := errors.New("shutdown failed")
	manager := NewPluginManager()
	if err := manager.RegisterPlugin(failingLifecyclePlugin{
		name:        "bad_shutdown",
		enabled:     true,
		shutdownErr: shutdownErr,
	}); err != nil {
		t.Fatalf("failed to register plugin: %v", err)
	}

	err := manager.ShutdownAll(context.Background())
	if err == nil {
		t.Fatal("expected shutdown error")
	}
	if !errors.Is(err, shutdownErr) {
		t.Fatalf("expected shutdown error to preserve cause, got %v", err)
	}
}

func TestPluginManagerFindsPluginsByCapability(t *testing.T) {
	manager := NewPluginManager()

	eventPlugin := managerInfoPlugin{
		info: PluginInfo{
			Name:         "event_plugin",
			Capabilities: []string{string(CapabilityEvent), string(CapabilityRegional)},
		},
	}
	calculationPlugin := managerInfoPlugin{
		info: PluginInfo{
			Name:         "calculation_plugin",
			Capabilities: []string{string(CapabilityCalculation)},
		},
	}

	for _, plugin := range []Plugin{eventPlugin, calculationPlugin} {
		if err := manager.RegisterPlugin(plugin); err != nil {
			t.Fatalf("failed to register plugin: %v", err)
		}
	}

	regionalPlugins := manager.GetPluginsByType(string(CapabilityRegional))
	if len(regionalPlugins) != 1 || regionalPlugins[0].GetInfo().Name != "event_plugin" {
		t.Fatalf("expected regional plugin to be event_plugin, got %#v", regionalPlugins)
	}

	calculationPlugins := manager.GetPluginsByCapability(CapabilityCalculation)
	if len(calculationPlugins) != 1 || calculationPlugins[0].GetInfo().Name != "calculation_plugin" {
		t.Fatalf("expected calculation plugin to be calculation_plugin, got %#v", calculationPlugins)
	}

	if validationPlugins := manager.GetPluginsByCapability(CapabilityValidation); len(validationPlugins) != 0 {
		t.Fatalf("expected no validation plugins, got %#v", validationPlugins)
	}
}

func TestPluginManagerStatsCountsEnabledPluginsAndCapabilities(t *testing.T) {
	manager := NewPluginManager()

	plugins := []Plugin{
		managerInfoPlugin{
			info: PluginInfo{
				Name:         "event_plugin",
				Capabilities: []string{string(CapabilityEvent), string(CapabilityRegional)},
			},
			enabled: true,
		},
		managerInfoPlugin{
			info: PluginInfo{
				Name:         "second_event_plugin",
				Capabilities: []string{string(CapabilityEvent)},
			},
		},
	}

	for _, plugin := range plugins {
		if err := manager.RegisterPlugin(plugin); err != nil {
			t.Fatalf("failed to register plugin: %v", err)
		}
	}

	stats := manager.GetPluginStats()
	if got := stats["total_plugins"]; got != 2 {
		t.Fatalf("expected 2 total plugins, got %v", got)
	}
	if got := stats["enabled_plugins"]; got != 1 {
		t.Fatalf("expected 1 enabled plugin, got %v", got)
	}

	pluginTypes, ok := stats["plugin_types"].(map[string]int)
	if !ok {
		t.Fatalf("expected plugin_types map, got %T", stats["plugin_types"])
	}
	if got := pluginTypes[string(CapabilityEvent)]; got != 2 {
		t.Fatalf("expected 2 event plugins, got %d", got)
	}
	if got := pluginTypes[string(CapabilityRegional)]; got != 1 {
		t.Fatalf("expected 1 regional plugin, got %d", got)
	}
}

func TestPluginManagerHealthReportsDegradedStoppedPlugins(t *testing.T) {
	manager := NewPluginManager()
	stoppedPlugin := managerInfoPlugin{
		info: PluginInfo{
			Name:         "stopped_plugin",
			Capabilities: []string{string(CapabilityEvent)},
		},
	}

	if err := manager.RegisterPlugin(stoppedPlugin); err != nil {
		t.Fatalf("failed to register plugin: %v", err)
	}

	health := manager.Health(context.Background())
	if got := health["status"]; got != "degraded" {
		t.Fatalf("expected degraded status, got %v", got)
	}
	if got := health["unhealthy_plugins"]; got != 1 {
		t.Fatalf("expected 1 unhealthy plugin, got %v", got)
	}

	plugins, ok := health["plugins"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected plugins map, got %T", health["plugins"])
	}
	pluginHealth, ok := plugins["stopped_plugin"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected stopped plugin health map, got %T", plugins["stopped_plugin"])
	}
	if got := pluginHealth["status"]; got != "unhealthy" {
		t.Fatalf("expected stopped plugin to be unhealthy, got %v", got)
	}
}
