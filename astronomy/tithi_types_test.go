package astronomy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTithiTypeDescription(t *testing.T) {
	tests := []struct {
		tithiType    TithiType
		expectedDesc string
	}{
		{TithiTypeNanda, "Joyful, good for celebrations and new beginnings"},
		{TithiTypeBhadra, "Auspicious, good for all activities"},
		{TithiTypeJaya, "Victorious, good for achieving success"},
		{TithiTypeRikta, "Empty, avoid starting new ventures"},
		{TithiTypePurna, "Complete, excellent for completion of tasks"},
		{TithiType("Invalid"), "Unknown Tithi type"},
	}

	for _, test := range tests {
		t.Run(string(test.tithiType), func(t *testing.T) {
			desc := GetTithiTypeDescription(test.tithiType)
			assert.Equal(t, test.expectedDesc, desc)
		})
	}
}

func TestGetTithiType(t *testing.T) {
	tests := []struct {
		tithiNumber  int
		expectedType TithiType
	}{
		{1, TithiTypeNanda},
		{2, TithiTypeBhadra},
		{3, TithiTypeJaya},
		{4, TithiTypeRikta},
		{5, TithiTypePurna},
		{6, TithiTypeNanda},
		{7, TithiTypeBhadra},
		{8, TithiTypeJaya},
		{9, TithiTypeRikta},
		{10, TithiTypePurna},
		{11, TithiTypeNanda},
		{12, TithiTypeBhadra},
		{13, TithiTypeJaya},
		{14, TithiTypeRikta},
		{15, TithiTypePurna},
		{16, TithiTypeNanda},
		{17, TithiTypeBhadra},
		{18, TithiTypeJaya},
		{19, TithiTypeRikta},
		{20, TithiTypePurna},
		{25, TithiTypePurna},
		{30, TithiTypePurna},
	}

	for _, test := range tests {
		t.Run(TithiNames[test.tithiNumber], func(t *testing.T) {
			tithiType := getTithiType(test.tithiNumber)
			assert.Equal(t, test.expectedType, tithiType)
		})
	}
}

func TestTithiNames(t *testing.T) {
	for i := 1; i <= 30; i++ {
		name, exists := TithiNames[i]
		assert.True(t, exists, "Tithi number %d should have a name", i)
		assert.NotEmpty(t, name, "Tithi name for %d should not be empty", i)
	}

	assert.Equal(t, "Pratipada", TithiNames[1])
	assert.Equal(t, "Purnima", TithiNames[15])
	assert.Equal(t, "Pratipada", TithiNames[16])
	assert.Equal(t, "Amavasya", TithiNames[30])
}

func TestNormalizeTithiCalendarSystem(t *testing.T) {
	tests := []struct {
		name           string
		calendarSystem string
		expected       string
	}{
		{
			name:           "empty defaults to purnimanta",
			calendarSystem: "",
			expected:       "Purnimanta",
		},
		{
			name:           "traditional defaults to purnimanta",
			calendarSystem: "traditional",
			expected:       "Purnimanta",
		},
		{
			name:           "amanta is normalized",
			calendarSystem: " AmAnTa ",
			expected:       "Amanta",
		},
		{
			name:           "unknown value is preserved",
			calendarSystem: "CustomCalendar",
			expected:       "CustomCalendar",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, normalizeTithiCalendarSystem(test.calendarSystem))
		})
	}
}

func BenchmarkGetTithiType(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getTithiType((i % 30) + 1)
	}
}
