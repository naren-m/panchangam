package examples

import "time"

func (h *HinduFestivalPlugin) isDiwali(date time.Time) bool {
	return date.Month() == time.October && date.Day() >= 20 ||
		date.Month() == time.November && date.Day() <= 15
}

func (h *HinduFestivalPlugin) isHoli(date time.Time) bool {
	return date.Month() == time.March && date.Day() >= 10 && date.Day() <= 25
}

func (h *HinduFestivalPlugin) isNavaratri(date time.Time) bool {
	return date.Month() == time.September && date.Day() >= 20 ||
		date.Month() == time.October && date.Day() <= 10
}

func (h *HinduFestivalPlugin) isJanmashtami(date time.Time) bool {
	return date.Month() == time.August && date.Day() >= 15 ||
		date.Month() == time.September && date.Day() <= 5
}

func (h *HinduFestivalPlugin) isMakarSankranti(date time.Time) bool {
	return date.Month() == time.January && (date.Day() == 14 || date.Day() == 15)
}

func (h *HinduFestivalPlugin) isRamNavami(date time.Time) bool {
	return date.Month() == time.March && date.Day() >= 20 ||
		date.Month() == time.April && date.Day() <= 15
}

func (h *HinduFestivalPlugin) isEkadashi(date time.Time) bool {
	return date.Day()%15 == 11 || date.Day()%15 == 26
}

func (h *HinduFestivalPlugin) isAmavasya(date time.Time) bool {
	return date.Day() == 1 || date.Day() == 15 || date.Day() == 30
}

func (h *HinduFestivalPlugin) isPurnima(date time.Time) bool {
	return date.Day() == 15
}

func (h *HinduFestivalPlugin) isPongal(date time.Time) bool {
	return date.Month() == time.January && date.Day() >= 14 && date.Day() <= 17
}

func (h *HinduFestivalPlugin) isOnam(date time.Time) bool {
	return date.Month() == time.August && date.Day() >= 20 ||
		date.Month() == time.September && date.Day() <= 10
}

func (h *HinduFestivalPlugin) isDurgaPuja(date time.Time) bool {
	return date.Month() == time.September && date.Day() >= 25 ||
		date.Month() == time.October && date.Day() <= 15
}

func (h *HinduFestivalPlugin) getNavaratriDay(date time.Time) int {
	return 1
}

func (h *HinduFestivalPlugin) getNavaratriGoddess(day int) string {
	goddesses := []string{
		"Shailaputri", "Brahmacharini", "Chandraghanta", "Kushmanda",
		"Skandamata", "Katyayani", "Kaalratri", "Mahagauri", "Siddhidatri",
	}
	if day >= 1 && day <= 9 {
		return goddesses[day-1]
	}
	return ""
}

func (h *HinduFestivalPlugin) getEkadashiName(date time.Time) string {
	ekadashiNames := []string{
		"Utpanna", "Mokshada", "Saphala", "Putrada", "Shattila", "Jaya",
		"Vijaya", "Amalaki", "Papamochani", "Kamada", "Varuthini", "Mohini",
		"Apara", "Nirjala", "Yogini", "Devshayani", "Kamika", "Shravana",
		"Aja", "Parsva", "Indira", "Papankusha", "Rama", "Haribodhini",
	}
	monthIndex := int(date.Month()) % len(ekadashiNames)
	return ekadashiNames[monthIndex]
}

func (h *HinduFestivalPlugin) getPongalDay(date time.Time) string {
	switch date.Day() {
	case 14:
		return "Bhogi Pongal"
	case 15:
		return "Thai Pongal"
	case 16:
		return "Mattu Pongal"
	case 17:
		return "Kaanum Pongal"
	default:
		return "Pongal"
	}
}

func (h *HinduFestivalPlugin) getPongalTamilName(day string) string {
	names := map[string]string{
		"Bhogi Pongal":  "போகி பொங்கல்",
		"Thai Pongal":   "தைப் பொங்கல்",
		"Mattu Pongal":  "மாட்டுப் பொங்கல்",
		"Kaanum Pongal": "காணும் பொங்கல்",
		"Pongal":        "பொங்கல்",
	}
	return names[day]
}

func (h *HinduFestivalPlugin) getPongalRituals(day string) []string {
	rituals := map[string][]string{
		"Bhogi Pongal":  {"discarding_old_items", "bonfire", "cleaning"},
		"Thai Pongal":   {"cooking_pongal", "surya_worship", "sugarcane_offering"},
		"Mattu Pongal":  {"cattle_worship", "decoration", "thanksgiving"},
		"Kaanum Pongal": {"family_gathering", "outing", "cultural_programs"},
	}
	return rituals[day]
}

func (h *HinduFestivalPlugin) getOnamDay(date time.Time) string {
	days := []string{
		"Atham", "Chithira", "Chodhi", "Vishakam", "Anizham",
		"Thriketa", "Moolam", "Pooradam", "Uthradom", "Thiruvonam",
	}
	dayIndex := (date.Day() - 1) % len(days)
	return days[dayIndex]
}
