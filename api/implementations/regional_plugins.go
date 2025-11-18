package implementations

import (
	"context"
	"time"

	"github.com/naren-m/panchangam/api"
)

// TamilNaduRegionalPlugin provides Tamil Nadu specific calculations
type TamilNaduRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewTamilNaduRegionalPlugin creates a new Tamil Nadu regional plugin
func NewTamilNaduRegionalPlugin() *TamilNaduRegionalPlugin {
	return &TamilNaduRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (t *TamilNaduRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "tamil_nadu_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Tamil Nadu regional extensions with Amanta calendar and Naazhikai time system",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
			string(api.CapabilityEvent),
			string(api.CapabilityMuhurta),
		},
		Dependencies: []string{"astronomy"},
		Metadata: map[string]interface{}{
			"calendar_system": "amanta",
			"time_system":     "naazhikai",
			"language":        "tamil",
			"festivals":       []string{"Pongal", "Tamil New Year", "Aadi Perukku"},
		},
	}
}

// Initialize sets up the plugin
func (t *TamilNaduRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	t.config = config
	t.enabled = true
	return nil
}

// IsEnabled returns whether the plugin is enabled
func (t *TamilNaduRegionalPlugin) IsEnabled() bool {
	return t.enabled
}

// Shutdown cleans up resources
func (t *TamilNaduRegionalPlugin) Shutdown(ctx context.Context) error {
	t.enabled = false
	return nil
}

// GetRegion returns the region this plugin supports
func (t *TamilNaduRegionalPlugin) GetRegion() api.Region {
	return api.RegionTamilNadu
}

// GetCalendarSystem returns the calendar system used
func (t *TamilNaduRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarAmanta
}

// ApplyRegionalRules applies Tamil Nadu specific rules
func (t *TamilNaduRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	// Set calendar system to Amanta
	data.CalendarSystem = api.CalendarAmanta

	// Add regional metadata via events
	regionalInfo := api.Event{
		Name:      "Tamil Nadu Regional Info",
		NameLocal: "தமிழ்நாடு பிராந்திய தகவல்",
		Type:      api.EventTypeLunar,
		StartTime: data.Date,
		EndTime:   data.Date.Add(24 * time.Hour),
		Significance: "Tamil Nadu follows Amanta calendar system and uses Naazhikai time units",
		Region:    api.RegionTamilNadu,
		Metadata: map[string]interface{}{
			"type":            "regional_info",
			"calendar_system": "amanta",
			"time_system":     "naazhikai",
			"language":        "tamil",
		},
	}

	data.Events = append(data.Events, regionalInfo)
	return nil
}

// GetRegionalEvents returns Tamil Nadu specific events
func (t *TamilNaduRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	var events []api.Event

	// Tamil New Year (Puthandu) - April 14th
	if date.Month() == time.April && date.Day() == 14 {
		events = append(events, api.Event{
			Name:         "Tamil New Year",
			NameLocal:    "தமிழ் புத்தாண்டு",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Tamil New Year celebrated on the first day of Chithirai month",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"importance":       "highest",
				"month":            "Chithirai",
				"solar_event":      true,
				"pan_tamil_festival": true,
			},
		})
	}

	// Pongal - January 14-17
	if date.Month() == time.January && date.Day() >= 14 && date.Day() <= 17 {
		pongalDay := ""
		switch date.Day() {
		case 14:
			pongalDay = "Bhogi Pongal"
		case 15:
			pongalDay = "Thai Pongal"
		case 16:
			pongalDay = "Maattu Pongal"
		case 17:
			pongalDay = "Kaanum Pongal"
		}

		events = append(events, api.Event{
			Name:         pongalDay,
			NameLocal:    "பொங்கல்",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Tamil harvest festival celebrating the Sun God",
			Region:       api.RegionTamilNadu,
			Metadata: map[string]interface{}{
				"importance":    "highest",
				"pongal_day":    pongalDay,
				"harvest_festival": true,
				"duration_days": 4,
			},
		})
	}

	return events, nil
}

// GetRegionalMuhurtas returns Tamil Nadu specific muhurtas
func (t *TamilNaduRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	// Tamil muhurtas use Naazhikai system
	// Will be calculated in Naazhikai converter
	return []api.Muhurta{}, nil
}

// GetRegionalNames returns localized Tamil names
func (t *TamilNaduRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	tamilNames := map[string]string{
		"Sunday":    "ஞாயிறு",
		"Monday":    "திங்கள்",
		"Tuesday":   "செவ்வாய்",
		"Wednesday": "புதன்",
		"Thursday":  "வியாழன்",
		"Friday":    "வெள்ளி",
		"Saturday":  "சனி",
	}
	return tamilNames
}

// KeralaRegionalPlugin provides Kerala specific calculations
type KeralaRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

// NewKeralaRegionalPlugin creates a new Kerala regional plugin
func NewKeralaRegionalPlugin() *KeralaRegionalPlugin {
	return &KeralaRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

// GetInfo returns plugin metadata
func (k *KeralaRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "kerala_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Kerala regional extensions with Amanta calendar and Malayalam calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
			string(api.CapabilityEvent),
		},
		Dependencies: []string{"astronomy"},
		Metadata: map[string]interface{}{
			"calendar_system": "amanta",
			"language":        "malayalam",
			"festivals":       []string{"Onam", "Vishu", "Thiruvathira"},
		},
	}
}

func (k *KeralaRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	k.config = config
	k.enabled = true
	return nil
}

func (k *KeralaRegionalPlugin) IsEnabled() bool {
	return k.enabled
}

func (k *KeralaRegionalPlugin) Shutdown(ctx context.Context) error {
	k.enabled = false
	return nil
}

func (k *KeralaRegionalPlugin) GetRegion() api.Region {
	return api.RegionKerala
}

func (k *KeralaRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarAmanta
}

func (k *KeralaRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	data.CalendarSystem = api.CalendarAmanta
	return nil
}

func (k *KeralaRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	var events []api.Event

	// Vishu (Kerala New Year) - April 14/15
	if date.Month() == time.April && (date.Day() == 14 || date.Day() == 15) {
		events = append(events, api.Event{
			Name:         "Vishu",
			NameLocal:    "വിഷു",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Malayalam New Year celebrated with Vishukkani",
			Region:       api.RegionKerala,
			Metadata: map[string]interface{}{
				"importance":    "highest",
				"new_year":      true,
				"solar_festival": true,
			},
		})
	}

	return events, nil
}

func (k *KeralaRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	return []api.Muhurta{}, nil
}

func (k *KeralaRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	malayalamNames := map[string]string{
		"Sunday":    "ഞായർ",
		"Monday":    "തിങ്കൾ",
		"Tuesday":   "ചൊവ്വ",
		"Wednesday": "ബുധൻ",
		"Thursday":  "വ്യാഴം",
		"Friday":    "വെള്ളി",
		"Saturday":  "ശനി",
	}
	return malayalamNames
}

// BengalRegionalPlugin provides Bengal specific calculations
type BengalRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

func NewBengalRegionalPlugin() *BengalRegionalPlugin {
	return &BengalRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

func (b *BengalRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "bengal_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Bengal regional extensions with Amanta calendar and Bengali calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
			string(api.CapabilityEvent),
		},
		Metadata: map[string]interface{}{
			"calendar_system": "amanta",
			"language":        "bengali",
			"festivals":       []string{"Durga Puja", "Pohela Boishakh", "Kali Puja"},
		},
	}
}

func (b *BengalRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	b.config = config
	b.enabled = true
	return nil
}

func (b *BengalRegionalPlugin) IsEnabled() bool {
	return b.enabled
}

func (b *BengalRegionalPlugin) Shutdown(ctx context.Context) error {
	b.enabled = false
	return nil
}

func (b *BengalRegionalPlugin) GetRegion() api.Region {
	return api.RegionBengal
}

func (b *BengalRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarAmanta
}

func (b *BengalRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	data.CalendarSystem = api.CalendarAmanta
	return nil
}

func (b *BengalRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	var events []api.Event

	// Pohela Boishakh (Bengali New Year) - April 14/15
	if date.Month() == time.April && (date.Day() == 14 || date.Day() == 15) {
		events = append(events, api.Event{
			Name:         "Pohela Boishakh",
			NameLocal:    "পহেলা বৈশাখ",
			Type:         api.EventTypeFestival,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Bengali New Year celebration",
			Region:       api.RegionBengal,
			Metadata: map[string]interface{}{
				"importance": "highest",
				"new_year":   true,
			},
		})
	}

	return events, nil
}

func (b *BengalRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	return []api.Muhurta{}, nil
}

func (b *BengalRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	bengaliNames := map[string]string{
		"Sunday":    "রবিবার",
		"Monday":    "সোমবার",
		"Tuesday":   "মঙ্গলবার",
		"Wednesday": "বুধবার",
		"Thursday":  "বৃহস্পতিবার",
		"Friday":    "শুক্রবার",
		"Saturday":  "শনিবার",
	}
	return bengaliNames
}

// GujaratRegionalPlugin provides Gujarat specific calculations
type GujaratRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

func NewGujaratRegionalPlugin() *GujaratRegionalPlugin {
	return &GujaratRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

func (g *GujaratRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "gujarat_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Gujarat regional extensions with Purnimanta calendar and Gujarati calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
			string(api.CapabilityEvent),
		},
		Metadata: map[string]interface{}{
			"calendar_system": "purnimanta",
			"language":        "gujarati",
			"festivals":       []string{"Uttarayan", "Navratri", "Bestu Varas"},
		},
	}
}

func (g *GujaratRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	g.config = config
	g.enabled = true
	return nil
}

func (g *GujaratRegionalPlugin) IsEnabled() bool {
	return g.enabled
}

func (g *GujaratRegionalPlugin) Shutdown(ctx context.Context) error {
	g.enabled = false
	return nil
}

func (g *GujaratRegionalPlugin) GetRegion() api.Region {
	return api.RegionGujarat
}

func (g *GujaratRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarPurnimanta
}

func (g *GujaratRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	data.CalendarSystem = api.CalendarPurnimanta
	return nil
}

func (g *GujaratRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	var events []api.Event

	// Uttarayan - January 14
	if date.Month() == time.January && date.Day() == 14 {
		events = append(events, api.Event{
			Name:         "Uttarayan",
			NameLocal:    "ઉત્તરાયણ",
			Type:         api.EventTypeSolar,
			StartTime:    date,
			EndTime:      date.Add(24 * time.Hour),
			Significance: "Kite festival celebrating the sun's northward journey",
			Region:       api.RegionGujarat,
			Metadata: map[string]interface{}{
				"importance": "highest",
				"solar_event": true,
				"kite_festival": true,
			},
		})
	}

	return events, nil
}

func (g *GujaratRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	return []api.Muhurta{}, nil
}

func (g *GujaratRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	gujaratiNames := map[string]string{
		"Sunday":    "રવિવાર",
		"Monday":    "સોમવાર",
		"Tuesday":   "મંગળવાર",
		"Wednesday": "બુધવાર",
		"Thursday":  "ગુરુવાર",
		"Friday":    "શુક્રવાર",
		"Saturday":  "શનિવાર",
	}
	return gujaratiNames
}

// MaharashtraRegionalPlugin provides Maharashtra specific calculations
type MaharashtraRegionalPlugin struct {
	enabled bool
	config  map[string]interface{}
}

func NewMaharashtraRegionalPlugin() *MaharashtraRegionalPlugin {
	return &MaharashtraRegionalPlugin{
		enabled: false,
		config:  make(map[string]interface{}),
	}
}

func (m *MaharashtraRegionalPlugin) GetInfo() api.PluginInfo {
	return api.PluginInfo{
		Name:        "maharashtra_regional_plugin",
		Version:     api.Version{Major: 1, Minor: 0, Patch: 0},
		Description: "Maharashtra regional extensions with Purnimanta calendar and Marathi calendar",
		Author:      "Panchangam Team",
		Capabilities: []string{
			string(api.CapabilityRegional),
			string(api.CapabilityEvent),
		},
		Metadata: map[string]interface{}{
			"calendar_system": "purnimanta",
			"language":        "marathi",
			"festivals":       []string{"Gudi Padwa", "Ganesh Chaturthi", "Gokul Ashtami"},
		},
	}
}

func (m *MaharashtraRegionalPlugin) Initialize(ctx context.Context, config map[string]interface{}) error {
	m.config = config
	m.enabled = true
	return nil
}

func (m *MaharashtraRegionalPlugin) IsEnabled() bool {
	return m.enabled
}

func (m *MaharashtraRegionalPlugin) Shutdown(ctx context.Context) error {
	m.enabled = false
	return nil
}

func (m *MaharashtraRegionalPlugin) GetRegion() api.Region {
	return api.RegionMaha
}

func (m *MaharashtraRegionalPlugin) GetCalendarSystem() api.CalendarSystem {
	return api.CalendarPurnimanta
}

func (m *MaharashtraRegionalPlugin) ApplyRegionalRules(ctx context.Context, data *api.PanchangamData) error {
	data.CalendarSystem = api.CalendarPurnimanta
	return nil
}

func (m *MaharashtraRegionalPlugin) GetRegionalEvents(ctx context.Context, date time.Time, location api.Location) ([]api.Event, error) {
	var events []api.Event

	// Gudi Padwa (Marathi New Year) - typically in March/April
	// This is a lunar date, would need precise calculation
	// For now, we'll note it as a regional festival

	return events, nil
}

func (m *MaharashtraRegionalPlugin) GetRegionalMuhurtas(ctx context.Context, date time.Time, location api.Location) ([]api.Muhurta, error) {
	return []api.Muhurta{}, nil
}

func (m *MaharashtraRegionalPlugin) GetRegionalNames(locale string) map[string]string {
	marathiNames := map[string]string{
		"Sunday":    "रविवार",
		"Monday":    "सोमवार",
		"Tuesday":   "मंगळवार",
		"Wednesday": "बुधवार",
		"Thursday":  "गुरुवार",
		"Friday":    "शुक्रवार",
		"Saturday":  "शनिवार",
	}
	return marathiNames
}
