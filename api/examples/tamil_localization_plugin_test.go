package examples

import (
	"context"
	"strings"
	"testing"

	"github.com/naren-m/panchangam/api"
	"github.com/stretchr/testify/assert"
)

func newEnabledTamilLocalizationPlugin(t *testing.T) *TamilLocalizationPlugin {
	t.Helper()

	plugin := NewTamilLocalizationPlugin()
	if err := plugin.Initialize(context.Background(), map[string]interface{}{}); err != nil {
		t.Fatalf("Initialize returned error: %v", err)
	}
	return plugin
}

func TestTamilLocalizationPluginLocalizesSupportedLocale(t *testing.T) {
	plugin := newEnabledTamilLocalizationPlugin(t)

	tithi := &api.Tithi{
		Name:    "Ekadashi",
		Lord:    "Vishnu",
		Quality: "auspicious",
	}
	err := plugin.LocalizeTithi(tithi, "ta-IN", api.RegionTamilNadu)

	assert.NoError(t, err)
	assert.Equal(t, "ஏகாதசி", tithi.NameLocal)
	assert.Equal(t, "விஷ்ணு", tithi.Lord)
	assert.Equal(t, "நன்மை", tithi.Quality)
}

func TestTamilLocalizationPluginSkipsUnsupportedLocale(t *testing.T) {
	plugin := newEnabledTamilLocalizationPlugin(t)

	event := &api.Event{
		Name:         "Pongal",
		Significance: "festival celebration",
	}
	err := plugin.LocalizeEvent(event, "en-US", api.RegionTamilNadu)

	assert.NoError(t, err)
	assert.Empty(t, event.NameLocal)
	assert.Equal(t, "festival celebration", event.Significance)
}

func TestTamilLocalizationPluginLifecycle(t *testing.T) {
	plugin := NewTamilLocalizationPlugin()
	assert.False(t, plugin.IsEnabled())

	err := plugin.Initialize(context.Background(), map[string]interface{}{"mode": "test"})
	assert.NoError(t, err)
	assert.True(t, plugin.IsEnabled())

	err = plugin.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.False(t, plugin.IsEnabled())
}

func TestTamilLocalizationPluginRejectsLocalizationWhenDisabled(t *testing.T) {
	plugin := NewTamilLocalizationPlugin()

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "tithi",
			call: func() error {
				return plugin.LocalizeTithi(&api.Tithi{Name: "Ekadashi"}, "ta-IN", api.RegionTamilNadu)
			},
		},
		{
			name: "nakshatra",
			call: func() error {
				return plugin.LocalizeNakshatra(&api.Nakshatra{Name: "Rohini"}, "ta-IN", api.RegionTamilNadu)
			},
		},
		{
			name: "yoga",
			call: func() error {
				return plugin.LocalizeYoga(&api.Yoga{Name: "Siddhi"}, "ta-IN", api.RegionTamilNadu)
			},
		},
		{
			name: "karana",
			call: func() error {
				return plugin.LocalizeKarana(&api.Karana{Name: "Bava"}, "ta-IN", api.RegionTamilNadu)
			},
		},
		{
			name: "event",
			call: func() error {
				return plugin.LocalizeEvent(&api.Event{Name: "Pongal"}, "ta-IN", api.RegionTamilNadu)
			},
		},
		{
			name: "muhurta",
			call: func() error {
				return plugin.LocalizeMuhurta(&api.Muhurta{Name: "Abhijit Muhurta"}, "ta-IN", api.RegionTamilNadu)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("Expected disabled plugin to return an error")
			}
			if !strings.Contains(err.Error(), "tamil localization plugin is not enabled") {
				t.Fatalf("Expected disabled plugin error, got %v", err)
			}
		})
	}
}
