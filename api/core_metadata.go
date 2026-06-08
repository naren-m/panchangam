package api

// GetVersion returns the API version
func (api *CorePanchangamAPI) GetVersion() Version {
	return api.version
}

// GetSupportedRegions returns all supported regions
func (api *CorePanchangamAPI) GetSupportedRegions() []Region {
	return []Region{
		RegionGlobal,
		RegionNorthIndia,
		RegionSouthIndia,
		RegionTamilNadu,
		RegionKerala,
		RegionBengal,
		RegionGujarat,
		RegionMaha,
	}
}

// GetSupportedMethods returns all supported calculation methods
func (api *CorePanchangamAPI) GetSupportedMethods() []CalculationMethod {
	return []CalculationMethod{
		MethodDrik,
		MethodVakya,
		MethodAuto,
	}
}

// GetSupportedCalendars returns all supported calendar systems
func (api *CorePanchangamAPI) GetSupportedCalendars() []CalendarSystem {
	return []CalendarSystem{
		CalendarPurnimanta,
		CalendarAmanta,
		CalendarLunar,
		CalendarSolar,
	}
}

// GetPluginManager returns the plugin manager
func (api *CorePanchangamAPI) GetPluginManager() PluginManager {
	return api.pluginManager
}

// RegisterPlugin registers a plugin with the API
func (api *CorePanchangamAPI) RegisterPlugin(plugin Plugin) error {
	return api.pluginManager.RegisterPlugin(plugin)
}
