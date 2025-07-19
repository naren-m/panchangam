package astronomy

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAllPanchangamElements(t *testing.T) {
	calculator, mockProvider, _ := createTestTithiCalculator()
	manager := calculator.ephemerisManager

	// Set up mock planetary positions for comprehensive test
	mockPositions := &ephemeris.PlanetaryPositions{
		Sun: ephemeris.Position{
			Longitude: 285.5, // Capricorn
		},
		Moon: ephemeris.Position{
			Longitude: 95.7, // Cancer, should be around Pushya nakshatra
		},
	}
	mockProvider.On("GetPlanetaryPositions", mock.Anything, mock.Anything).Return(mockPositions, nil)

	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	t.Run("Tithi Calculation", func(t *testing.T) {
		tithi, err := calculator.GetTithiFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		
		assert.NoError(t, err)
		assert.NotNil(t, tithi)
		err = ValidateTithiCalculation(tithi)
		assert.NoError(t, err)
		
		t.Logf("Tithi: %s (#%d), Type: %s, Duration: %.2f hours", 
			tithi.Name, tithi.Number, tithi.Type, tithi.Duration)
	})

	t.Run("Nakshatra Calculation", func(t *testing.T) {
		calculator := NewNakshatraCalculator(manager)
		nakshatra, err := calculator.GetNakshatraFromLongitude(ctx, mockPositions.Moon.Longitude, testDate)
		
		assert.NoError(t, err)
		assert.NotNil(t, nakshatra)
		err = ValidateNakshatraCalculation(nakshatra)
		assert.NoError(t, err)
		
		t.Logf("Nakshatra: %s (#%d), Pada: %d, Planet: %s", 
			nakshatra.Name, nakshatra.Number, nakshatra.Pada, nakshatra.PlanetaryLord)
	})

	t.Run("Yoga Calculation", func(t *testing.T) {
		calculator := NewYogaCalculator(manager)
		yoga, err := calculator.GetYogaFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		
		assert.NoError(t, err)
		assert.NotNil(t, yoga)
		err = ValidateYogaCalculation(yoga)
		assert.NoError(t, err)
		
		t.Logf("Yoga: %s (#%d), Quality: %s, Combined: %.2f°", 
			yoga.Name, yoga.Number, yoga.Quality, yoga.CombinedValue)
	})

	t.Run("Karana Calculation", func(t *testing.T) {
		calculator := NewKaranaCalculator(manager)
		karana, err := calculator.GetKaranaFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		
		assert.NoError(t, err)
		assert.NotNil(t, karana)
		err = ValidateKaranaCalculation(karana)
		assert.NoError(t, err)
		
		t.Logf("Karana: %s (#%d), Type: %s, Tithi: %d/%d, Vishti: %v", 
			karana.Name, karana.Number, karana.Type, karana.TithiNumber, karana.HalfTithi, karana.IsVishti)
	})

	t.Run("Vara Calculation", func(t *testing.T) {
		varaCalculator := NewVaraCalculator()
		
		// For this test, we'll use the simple gregorian day approach since we don't need ephemeris
		gregorianDay := testDate.Weekday() // Monday
		sunrise := time.Date(2024, 1, 15, 6, 30, 0, 0, time.UTC)
		nextSunrise := time.Date(2024, 1, 16, 6, 31, 0, 0, time.UTC)
		
		vara, err := varaCalculator.GetVaraFromGregorianDay(ctx, gregorianDay, sunrise, nextSunrise, testDate)
		
		assert.NoError(t, err)
		assert.NotNil(t, vara)
		err = ValidateVaraCalculation(vara)
		assert.NoError(t, err)
		
		t.Logf("Vara: %s (#%d), Planet: %s, Hora: %d (%s), Auspicious: %v", 
			vara.Name, vara.Number, vara.PlanetaryLord, vara.CurrentHora, vara.HoraPlanet, vara.IsAuspicious)
	})
}

func TestPanchangamDataIntegrity(t *testing.T) {
	t.Run("Tithi Names", func(t *testing.T) {
		assert.Equal(t, 30, len(TithiNames), "Should have exactly 30 Tithi names")
		for i := 1; i <= 30; i++ {
			assert.NotEmpty(t, TithiNames[i], "Tithi %d name should not be empty", i)
		}
	})

	t.Run("Nakshatra Data", func(t *testing.T) {
		assert.Equal(t, 27, len(NakshatraData), "Should have exactly 27 Nakshatras")
		for i := 1; i <= 27; i++ {
			data := NakshatraData[i]
			assert.NotEmpty(t, data.Name, "Nakshatra %d name should not be empty", i)
			assert.NotEmpty(t, data.Deity, "Nakshatra %d deity should not be empty", i)
			assert.NotEmpty(t, data.PlanetaryLord, "Nakshatra %d planetary lord should not be empty", i)
		}
	})

	t.Run("Yoga Data", func(t *testing.T) {
		assert.Equal(t, 27, len(YogaData), "Should have exactly 27 Yogas")
		for i := 1; i <= 27; i++ {
			data := YogaData[i]
			assert.NotEmpty(t, data.Name, "Yoga %d name should not be empty", i)
			assert.NotEmpty(t, data.Description, "Yoga %d description should not be empty", i)
		}
	})

	t.Run("Karana Data", func(t *testing.T) {
		assert.Equal(t, 11, len(KaranaData), "Should have exactly 11 Karanas")
		for i := 1; i <= 11; i++ {
			data := KaranaData[i]
			assert.NotEmpty(t, data.Name, "Karana %d name should not be empty", i)
			assert.NotEmpty(t, data.Description, "Karana %d description should not be empty", i)
		}
		// Check that Vishti is properly marked
		assert.True(t, KaranaData[8].IsVishti, "Karana 8 should be Vishti")
	})

	t.Run("Vara Data", func(t *testing.T) {
		assert.Equal(t, 7, len(VaraData), "Should have exactly 7 Varas")
		for i := 1; i <= 7; i++ {
			data := VaraData[i]
			assert.NotEmpty(t, data.Name, "Vara %d name should not be empty", i)
			assert.NotEmpty(t, data.PlanetaryLord, "Vara %d planetary lord should not be empty", i)
			assert.NotEmpty(t, data.GregorianDay, "Vara %d gregorian day should not be empty", i)
		}
	})
}

func BenchmarkAllPanchangamCalculations(b *testing.B) {
	calculator, mockProvider, _ := createTestTithiCalculator()
	manager := calculator.ephemerisManager

	mockPositions := &ephemeris.PlanetaryPositions{
		Sun: ephemeris.Position{
			Longitude: 285.5,
		},
		Moon: ephemeris.Position{
			Longitude: 95.7,
		},
	}
	mockProvider.On("GetPlanetaryPositions", mock.Anything, mock.Anything).Return(mockPositions, nil)

	ctx := context.Background()
	testDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tithiCalc := calculator
	nakshatraCalc := NewNakshatraCalculator(manager)
	yogaCalc := NewYogaCalculator(manager)
	karanaCalc := NewKaranaCalculator(manager)
	varaCalc := NewVaraCalculator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Benchmark all calculations together
		_, err := tithiCalc.GetTithiFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		if err != nil {
			b.Fatal(err)
		}

		_, err = nakshatraCalc.GetNakshatraFromLongitude(ctx, mockPositions.Moon.Longitude, testDate)
		if err != nil {
			b.Fatal(err)
		}

		_, err = yogaCalc.GetYogaFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		if err != nil {
			b.Fatal(err)
		}

		_, err = karanaCalc.GetKaranaFromLongitudes(ctx, mockPositions.Sun.Longitude, mockPositions.Moon.Longitude, testDate)
		if err != nil {
			b.Fatal(err)
		}

		// For Vara, use simple calculation without ephemeris
		sunrise := time.Date(2024, 1, 15, 6, 30, 0, 0, time.UTC)
		nextSunrise := time.Date(2024, 1, 16, 6, 31, 0, 0, time.UTC)
		_, err = varaCalc.GetVaraFromGregorianDay(ctx, testDate.Weekday(), sunrise, nextSunrise, testDate)
		if err != nil {
			b.Fatal(err)
		}
	}
}