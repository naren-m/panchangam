package implementations

import "github.com/naren-m/panchangam/api"

// isActivitySuitable checks if a muhurta is suitable for given activities
func (m *MuhurtaPlugin) isActivitySuitable(muhurta api.Muhurta, activities []string) bool {
	if len(muhurta.Purpose) == 0 {
		return false
	}

	// Check if any of the requested activities match the muhurta's purposes
	for _, activity := range activities {
		for _, purpose := range muhurta.Purpose {
			if activity == purpose || purpose == "all_auspicious_activities" {
				return true
			}
		}
	}

	// Check if any activities are in the avoid list
	for _, activity := range activities {
		for _, avoid := range muhurta.Avoid {
			if activity == avoid || avoid == "all_auspicious_activities" {
				return false
			}
		}
	}

	return false
}
