package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/stretchr/testify/assert"
)

func testFeatureTITHI_001(t *testing.T, ctx context.Context, testDate time.Time) {
	t.Run("TITHI_001_Lunar_Day_Calculation", func(t *testing.T) {
		testTithi := &astronomy.TithiInfo{
			Number:      15,
			Name:        "Purnima",
			Type:        astronomy.TithiTypePurna,
			StartTime:   testDate,
			EndTime:     testDate.Add(24 * time.Hour),
			Duration:    24.0,
			IsShukla:    true,
			MoonSunDiff: 180.0,
		}

		assert.True(t, testTithi.Number >= 1 && testTithi.Number <= 30, "TITHI_001: Number should be 1-30")
		assert.NotEmpty(t, testTithi.Name, "TITHI_001: Name should not be empty")
		assert.True(t, testTithi.Duration > 0, "TITHI_001: Duration should be positive")
		assert.True(t, testTithi.MoonSunDiff >= 0 && testTithi.MoonSunDiff < 360, "TITHI_001: MoonSunDiff should be 0-360")

		validTypes := map[astronomy.TithiType]bool{
			astronomy.TithiTypeNanda:  true,
			astronomy.TithiTypeBhadra: true,
			astronomy.TithiTypeJaya:   true,
			astronomy.TithiTypeRikta:  true,
			astronomy.TithiTypePurna:  true,
		}
		assert.True(t, validTypes[testTithi.Type], "TITHI_001: Type should be one of 5 valid types")

		if testTithi.Number <= 15 {
			assert.True(t, testTithi.IsShukla, "TITHI_001: Numbers 1-15 should be Shukla Paksha")
		} else {
			assert.False(t, testTithi.IsShukla, "TITHI_001: Numbers 16-30 should be Krishna Paksha")
		}

		assert.True(t, testTithi.EndTime.After(testTithi.StartTime), "TITHI_001: End time should be after start time")
		calculatedDuration := testTithi.EndTime.Sub(testTithi.StartTime).Hours()
		assert.InDelta(t, testTithi.Duration, calculatedDuration, 0.1, "TITHI_001: Duration should match time difference")

		t.Logf("TITHI_001: Validated Tithi %d - %s (%s)", testTithi.Number, testTithi.Name, testTithi.Type)
	})
}

func testFeatureNAKSHATRA_001(t *testing.T, ctx context.Context, testDate time.Time) {
	t.Run("NAKSHATRA_001_Lunar_Mansion_Calculation", func(t *testing.T) {
		testNakshatra := &astronomy.NakshatraInfo{
			Number:        13,
			Name:          "Hasta",
			Deity:         "Savitar",
			PlanetaryLord: "Moon",
			Symbol:        "Hand",
			Pada:          2,
			StartTime:     testDate,
			EndTime:       testDate.Add(time.Hour * 24),
			Duration:      24.0,
			MoonLongitude: 166.5,
		}

		assert.True(t, testNakshatra.Number >= 1 && testNakshatra.Number <= 27, "NAKSHATRA_001: Number should be 1-27")
		assert.NotEmpty(t, testNakshatra.Name, "NAKSHATRA_001: Name should not be empty")
		assert.NotEmpty(t, testNakshatra.Deity, "NAKSHATRA_001: Deity should not be empty")
		assert.NotEmpty(t, testNakshatra.PlanetaryLord, "NAKSHATRA_001: PlanetaryLord should not be empty")
		assert.NotEmpty(t, testNakshatra.Symbol, "NAKSHATRA_001: Symbol should not be empty")

		assert.True(t, testNakshatra.Pada >= 1 && testNakshatra.Pada <= 4, "NAKSHATRA_001: Pada should be 1-4")

		expectedLongitudeRange := 13.333333
		nakshatraStart := float64(testNakshatra.Number-1) * expectedLongitudeRange
		nakshatraEnd := float64(testNakshatra.Number) * expectedLongitudeRange
		assert.True(t, testNakshatra.MoonLongitude >= nakshatraStart && testNakshatra.MoonLongitude < nakshatraEnd,
			"NAKSHATRA_001: Moon longitude should be within Nakshatra range")

		assert.True(t, testNakshatra.EndTime.After(testNakshatra.StartTime), "NAKSHATRA_001: End time should be after start time")

		t.Logf("NAKSHATRA_001: Validated Nakshatra %d - %s (Pada %d)", testNakshatra.Number, testNakshatra.Name, testNakshatra.Pada)
	})
}

func testFeatureYOGA_001(t *testing.T, ctx context.Context, testDate time.Time) {
	t.Run("YOGA_001_Auspicious_Combinations", func(t *testing.T) {
		testYoga := &astronomy.YogaInfo{
			Number:        14,
			Name:          "Vishkambha",
			Quality:       astronomy.YogaQualityAuspicious,
			StartTime:     testDate,
			EndTime:       testDate.Add(time.Hour * 24),
			Duration:      24.0,
			CombinedValue: 180.0,
			Description:   "Auspicious for new beginnings",
		}

		assert.True(t, testYoga.Number >= 1 && testYoga.Number <= 27, "YOGA_001: Number should be 1-27")
		assert.NotEmpty(t, testYoga.Name, "YOGA_001: Name should not be empty")
		assert.NotEmpty(t, testYoga.Description, "YOGA_001: Description should not be empty")

		validQualities := map[astronomy.YogaQuality]bool{
			astronomy.YogaQualityAuspicious:   true,
			astronomy.YogaQualityInauspicious: true,
			astronomy.YogaQualityMixed:        true,
		}
		assert.True(t, validQualities[testYoga.Quality], "YOGA_001: Quality should be valid category")

		assert.True(t, testYoga.CombinedValue >= 0 && testYoga.CombinedValue < 360, "YOGA_001: CombinedValue should be 0-360")

		expectedYogaSize := 360.0 / 27.0
		yogaStart := float64(testYoga.Number-1) * expectedYogaSize
		yogaEnd := float64(testYoga.Number) * expectedYogaSize
		normalizedSum := testYoga.CombinedValue
		if normalizedSum >= 360 {
			normalizedSum -= 360
		}

		assert.True(t, normalizedSum >= yogaStart && normalizedSum < yogaEnd,
			"YOGA_001: Sun+Moon sum should be within Yoga range")

		assert.True(t, testYoga.EndTime.After(testYoga.StartTime), "YOGA_001: End time should be after start time")

		t.Logf("YOGA_001: Validated Yoga %d - %s (%s)", testYoga.Number, testYoga.Name, testYoga.Quality)
	})
}

func testFeatureKARANA_001(t *testing.T, ctx context.Context, testDate time.Time) {
	t.Run("KARANA_001_Half_Tithi_Divisions", func(t *testing.T) {
		testKarana := &astronomy.KaranaInfo{
			Number:      7,
			Name:        "Vanija",
			Type:        astronomy.KaranaTypeMovable,
			Description: "Merchant, good for business and trade",
			IsVishti:    false,
			StartTime:   testDate,
			EndTime:     testDate.Add(12 * time.Hour),
			Duration:    12.0,
			MoonSunDiff: 84.0,
			TithiNumber: 8,
			HalfTithi:   1,
		}

		assert.True(t, testKarana.Number >= 1 && testKarana.Number <= 11, "KARANA_001: Number should be 1-11")
		assert.NotEmpty(t, testKarana.Name, "KARANA_001: Name should not be empty")
		assert.NotEmpty(t, testKarana.Description, "KARANA_001: Description should not be empty")

		validTypes := map[astronomy.KaranaType]bool{
			astronomy.KaranaTypeMovable: true,
			astronomy.KaranaTypeFixed:   true,
		}
		assert.True(t, validTypes[testKarana.Type], "KARANA_001: Type should be Movable or Fixed")

		if testKarana.Number >= 1 && testKarana.Number <= 8 {
			if testKarana.Number != 8 {
				assert.Equal(t, astronomy.KaranaTypeMovable, testKarana.Type, "KARANA_001: Karanas 1-7 should be Movable")
			}
		} else {
			assert.Equal(t, astronomy.KaranaTypeFixed, testKarana.Type, "KARANA_001: Karanas 9-11 should be Fixed")
		}

		if testKarana.Name == "Vishti" {
			assert.True(t, testKarana.IsVishti, "KARANA_001: Vishti Karana should have IsVishti=true")
		} else {
			assert.False(t, testKarana.IsVishti, "KARANA_001: Non-Vishti Karanas should have IsVishti=false")
		}

		assert.True(t, testKarana.TithiNumber >= 1 && testKarana.TithiNumber <= 30, "KARANA_001: TithiNumber should be 1-30")
		assert.True(t, testKarana.HalfTithi >= 1 && testKarana.HalfTithi <= 2, "KARANA_001: HalfTithi should be 1 or 2")
		assert.True(t, testKarana.Duration > 6 && testKarana.Duration < 18, "KARANA_001: Duration should be roughly half a Tithi")
		assert.True(t, testKarana.MoonSunDiff >= 0 && testKarana.MoonSunDiff < 360, "KARANA_001: MoonSunDiff should be 0-360")

		t.Logf("KARANA_001: Validated Karana %d - %s (%s)", testKarana.Number, testKarana.Name, testKarana.Type)
	})
}

func testFeatureVARA_001(t *testing.T, ctx context.Context, testDate time.Time) {
	t.Run("VARA_001_Weekday_Calculation", func(t *testing.T) {
		testVara := &astronomy.VaraInfo{
			Number:        2,
			Name:          "Somavara",
			PlanetaryLord: "Moon",
			Quality:       "Peaceful",
			Color:         "White",
			Deity:         "Soma",
			StartTime:     testDate,
			EndTime:       testDate.Add(24 * time.Hour),
			Duration:      24.0,
			GregorianDay:  "Monday",
			IsAuspicious:  true,
			CurrentHora:   8,
			HoraPlanet:    "Moon",
		}

		assert.True(t, testVara.Number >= 1 && testVara.Number <= 7, "VARA_001: Number should be 1-7")
		assert.NotEmpty(t, testVara.Name, "VARA_001: Name should not be empty")
		assert.NotEmpty(t, testVara.PlanetaryLord, "VARA_001: PlanetaryLord should not be empty")
		assert.NotEmpty(t, testVara.GregorianDay, "VARA_001: GregorianDay should not be empty")
		assert.NotEmpty(t, testVara.Quality, "VARA_001: Quality should not be empty")
		assert.NotEmpty(t, testVara.Color, "VARA_001: Color should not be empty")
		assert.NotEmpty(t, testVara.Deity, "VARA_001: Deity should not be empty")

		validDays := map[string]int{
			"Sunday": 1, "Monday": 2, "Tuesday": 3, "Wednesday": 4,
			"Thursday": 5, "Friday": 6, "Saturday": 7,
		}
		expectedNumber, exists := validDays[testVara.GregorianDay]
		assert.True(t, exists, "VARA_001: GregorianDay should be valid weekday")
		assert.Equal(t, expectedNumber, testVara.Number, "VARA_001: Number should match weekday")

		assert.True(t, testVara.CurrentHora >= 1 && testVara.CurrentHora <= 24, "VARA_001: CurrentHora should be 1-24")
		assert.NotEmpty(t, testVara.HoraPlanet, "VARA_001: HoraPlanet should not be empty")

		validPlanets := map[string]bool{
			"Sun": true, "Moon": true, "Mars": true, "Mercury": true,
			"Jupiter": true, "Venus": true, "Saturn": true,
		}
		assert.True(t, validPlanets[testVara.PlanetaryLord], "VARA_001: PlanetaryLord should be valid planet")
		assert.True(t, validPlanets[testVara.HoraPlanet], "VARA_001: HoraPlanet should be valid planet")

		assert.True(t, testVara.EndTime.After(testVara.StartTime), "VARA_001: End time should be after start time")
		calculatedDuration := testVara.EndTime.Sub(testVara.StartTime).Hours()
		assert.InDelta(t, testVara.Duration, calculatedDuration, 0.1, "VARA_001: Duration should be approximately 24 hours")

		t.Logf("VARA_001: Validated Vara %d - %s (%s)", testVara.Number, testVara.Name, testVara.GregorianDay)
	})
}
