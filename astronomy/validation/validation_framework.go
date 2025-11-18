package validation

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/observability"
	"go.opentelemetry.io/otel/attribute"
)

// ValidationResult represents the result of a validation check
type ValidationResult struct {
	TestName        string    `json:"test_name"`
	Passed          bool      `json:"passed"`
	Expected        any       `json:"expected"`
	Actual          any       `json:"actual"`
	Tolerance       float64   `json:"tolerance"`
	Error           float64   `json:"error"`
	ErrorUnit       string    `json:"error_unit"`
	TestedAt        time.Time `json:"tested_at"`
	Duration        time.Duration `json:"duration"`
	Notes           string    `json:"notes,omitempty"`
	Source          string    `json:"source"`          // e.g., "Drik Panchang", "Swiss Ephemeris"
	CalculationType string    `json:"calculation_type"` // e.g., "Tithi", "Nakshatra", "Sunrise"
}

// ValidationSuite holds a collection of validation results
type ValidationSuite struct {
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Results         []ValidationResult `json:"results"`
	TotalTests      int                `json:"total_tests"`
	PassedTests     int                `json:"passed_tests"`
	FailedTests     int                `json:"failed_tests"`
	SuccessRate     float64            `json:"success_rate"`
	ExecutionTime   time.Duration      `json:"execution_time"`
	TestedAt        time.Time          `json:"tested_at"`
}

// ReferenceData represents reference data from established sources
type ReferenceData struct {
	Source       string                 `json:"source"`
	Date         time.Time              `json:"date"`
	Location     astronomy.Location     `json:"location"`
	Sunrise      time.Time              `json:"sunrise,omitempty"`
	Sunset       time.Time              `json:"sunset,omitempty"`
	TithiName    string                 `json:"tithi_name,omitempty"`
	NakshatraName string                `json:"nakshatra_name,omitempty"`
	YogaName     string                 `json:"yoga_name,omitempty"`
	Data         map[string]interface{} `json:"data"` // Additional custom data
}

// Validator validates astronomical calculations against reference sources
type Validator struct {
	tithiCalc        *astronomy.TithiCalculator
	nakshatraCalc    *astronomy.NakshatraCalculator
	yogaCalc         *astronomy.YogaCalculator
	karanaCalc       *astronomy.KaranaCalculator
	varaCalc         *astronomy.VaraCalculator
	sunriseCalc      *astronomy.SunriseCalculator
	ephemerisManager *ephemeris.Manager
	observer         observability.ObserverInterface
}

// NewValidator creates a new validator
func NewValidator(
	tithiCalc *astronomy.TithiCalculator,
	nakshatraCalc *astronomy.NakshatraCalculator,
	yogaCalc *astronomy.YogaCalculator,
	karanaCalc *astronomy.KaranaCalculator,
	varaCalc *astronomy.VaraCalculator,
	sunriseCalc *astronomy.SunriseCalculator,
	ephemerisManager *ephemeris.Manager,
) *Validator {
	return &Validator{
		tithiCalc:        tithiCalc,
		nakshatraCalc:    nakshatraCalc,
		yogaCalc:         yogaCalc,
		karanaCalc:       karanaCalc,
		varaCalc:         varaCalc,
		sunriseCalc:      sunriseCalc,
		ephemerisManager: ephemerisManager,
		observer:         observability.Observer(),
	}
}

// ValidateTithi validates Tithi calculation against reference data
func (v *Validator) ValidateTithi(ctx context.Context, ref ReferenceData, tolerance float64) ValidationResult {
	ctx, span := v.observer.CreateSpan(ctx, "validator.ValidateTithi")
	defer span.End()

	start := time.Now()
	result := ValidationResult{
		TestName:        "Tithi Validation",
		TestedAt:        time.Now(),
		Source:          ref.Source,
		CalculationType: "Tithi",
		Tolerance:       tolerance,
		ErrorUnit:       "minutes",
	}

	// Calculate Tithi
	tithi, err := v.tithiCalc.Calculate(ctx, ref.Date, ref.Location)
	if err != nil {
		result.Passed = false
		result.Notes = fmt.Sprintf("Calculation error: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Expected = ref.TithiName
	result.Actual = tithi.Name
	result.Passed = tithi.Name == ref.TithiName

	// Calculate time-based error if possible
	// This would require more detailed reference data
	result.Error = 0 // Placeholder
	result.Duration = time.Since(start)

	span.SetAttributes(
		attribute.Bool("passed", result.Passed),
		attribute.String("expected", ref.TithiName),
		attribute.String("actual", tithi.Name),
	)

	return result
}

// ValidateNakshatra validates Nakshatra calculation against reference data
func (v *Validator) ValidateNakshatra(ctx context.Context, ref ReferenceData, tolerance float64) ValidationResult {
	ctx, span := v.observer.CreateSpan(ctx, "validator.ValidateNakshatra")
	defer span.End()

	start := time.Now()
	result := ValidationResult{
		TestName:        "Nakshatra Validation",
		TestedAt:        time.Now(),
		Source:          ref.Source,
		CalculationType: "Nakshatra",
		Tolerance:       tolerance,
		ErrorUnit:       "minutes",
	}

	// Calculate Nakshatra
	nakshatra, err := v.nakshatraCalc.Calculate(ctx, ref.Date, ref.Location)
	if err != nil {
		result.Passed = false
		result.Notes = fmt.Sprintf("Calculation error: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	result.Expected = ref.NakshatraName
	result.Actual = nakshatra.Name
	result.Passed = nakshatra.Name == ref.NakshatraName
	result.Error = 0 // Placeholder
	result.Duration = time.Since(start)

	span.SetAttributes(
		attribute.Bool("passed", result.Passed),
		attribute.String("expected", ref.NakshatraName),
		attribute.String("actual", nakshatra.Name),
	)

	return result
}

// ValidateSunrise validates Sunrise calculation against reference data
func (v *Validator) ValidateSunrise(ctx context.Context, ref ReferenceData, toleranceMinutes float64) ValidationResult {
	ctx, span := v.observer.CreateSpan(ctx, "validator.ValidateSunrise")
	defer span.End()

	start := time.Now()
	result := ValidationResult{
		TestName:        "Sunrise Validation",
		TestedAt:        time.Now(),
		Source:          ref.Source,
		CalculationType: "Sunrise",
		Tolerance:       toleranceMinutes,
		ErrorUnit:       "minutes",
	}

	// Calculate Sunrise
	sunrise, err := v.sunriseCalc.CalculateSunrise(ctx, ref.Date, ref.Location)
	if err != nil {
		result.Passed = false
		result.Notes = fmt.Sprintf("Calculation error: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	// Calculate time difference in minutes
	diff := math.Abs(sunrise.Sub(ref.Sunrise).Minutes())

	result.Expected = ref.Sunrise.Format("15:04:05")
	result.Actual = sunrise.Format("15:04:05")
	result.Error = diff
	result.Passed = diff <= toleranceMinutes
	result.Duration = time.Since(start)

	if !result.Passed {
		result.Notes = fmt.Sprintf("Time difference: %.2f minutes", diff)
	}

	span.SetAttributes(
		attribute.Bool("passed", result.Passed),
		attribute.Float64("error_minutes", diff),
		attribute.Float64("tolerance", toleranceMinutes),
	)

	return result
}

// ValidatePlanetaryPosition validates planetary position against reference data
func (v *Validator) ValidatePlanetaryPosition(ctx context.Context, planet string, expectedLongitude float64, date time.Time, toleranceDegrees float64) ValidationResult {
	ctx, span := v.observer.CreateSpan(ctx, "validator.ValidatePlanetaryPosition")
	defer span.End()

	start := time.Now()
	result := ValidationResult{
		TestName:        fmt.Sprintf("%s Position Validation", planet),
		TestedAt:        time.Now(),
		Source:          "Reference Ephemeris",
		CalculationType: "Planetary Position",
		Tolerance:       toleranceDegrees,
		ErrorUnit:       "degrees",
	}

	// Get planetary position
	jd := ephemeris.TimeToJulianDay(date)
	positions, err := v.ephemerisManager.GetPlanetaryPositions(ctx, jd)
	if err != nil {
		result.Passed = false
		result.Notes = fmt.Sprintf("Calculation error: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	// Extract planet position
	var actualLongitude float64
	switch planet {
	case "sun":
		actualLongitude = positions.Sun.Longitude
	case "moon":
		actualLongitude = positions.Moon.Longitude
	case "mercury":
		actualLongitude = positions.Mercury.Longitude
	case "venus":
		actualLongitude = positions.Venus.Longitude
	case "mars":
		actualLongitude = positions.Mars.Longitude
	case "jupiter":
		actualLongitude = positions.Jupiter.Longitude
	case "saturn":
		actualLongitude = positions.Saturn.Longitude
	default:
		result.Passed = false
		result.Notes = fmt.Sprintf("Unknown planet: %s", planet)
		result.Duration = time.Since(start)
		return result
	}

	// Calculate angular difference
	diff := math.Abs(actualLongitude - expectedLongitude)
	if diff > 180 {
		diff = 360 - diff
	}

	result.Expected = fmt.Sprintf("%.6f°", expectedLongitude)
	result.Actual = fmt.Sprintf("%.6f°", actualLongitude)
	result.Error = diff
	result.Passed = diff <= toleranceDegrees
	result.Duration = time.Since(start)

	if !result.Passed {
		result.Notes = fmt.Sprintf("Angular difference: %.6f°", diff)
	}

	span.SetAttributes(
		attribute.Bool("passed", result.Passed),
		attribute.Float64("error_degrees", diff),
		attribute.String("planet", planet),
	)

	return result
}

// RunValidationSuite runs a complete validation suite
func (v *Validator) RunValidationSuite(ctx context.Context, name string, referenceData []ReferenceData) *ValidationSuite {
	ctx, span := v.observer.CreateSpan(ctx, "validator.RunValidationSuite")
	defer span.End()

	start := time.Now()
	suite := &ValidationSuite{
		Name:     name,
		TestedAt: time.Now(),
		Results:  make([]ValidationResult, 0),
	}

	// Run validations for each reference data point
	for _, ref := range referenceData {
		// Validate Tithi if available
		if ref.TithiName != "" {
			result := v.ValidateTithi(ctx, ref, 5.0) // 5 minutes tolerance
			suite.Results = append(suite.Results, result)
		}

		// Validate Nakshatra if available
		if ref.NakshatraName != "" {
			result := v.ValidateNakshatra(ctx, ref, 5.0) // 5 minutes tolerance
			suite.Results = append(suite.Results, result)
		}

		// Validate Sunrise if available
		if !ref.Sunrise.IsZero() {
			result := v.ValidateSunrise(ctx, ref, 3.0) // 3 minutes tolerance
			suite.Results = append(suite.Results, result)
		}
	}

	// Calculate statistics
	suite.TotalTests = len(suite.Results)
	for _, result := range suite.Results {
		if result.Passed {
			suite.PassedTests++
		} else {
			suite.FailedTests++
		}
	}

	if suite.TotalTests > 0 {
		suite.SuccessRate = float64(suite.PassedTests) / float64(suite.TotalTests) * 100
	}

	suite.ExecutionTime = time.Since(start)

	span.SetAttributes(
		attribute.Int("total_tests", suite.TotalTests),
		attribute.Int("passed_tests", suite.PassedTests),
		attribute.Int("failed_tests", suite.FailedTests),
		attribute.Float64("success_rate", suite.SuccessRate),
	)

	return suite
}

// ValidateAgainstDrikPanchang validates calculations against Drik Panchang data
func (v *Validator) ValidateAgainstDrikPanchang(ctx context.Context, referenceData []ReferenceData) *ValidationSuite {
	// Set source for all reference data
	for i := range referenceData {
		referenceData[i].Source = "Drik Panchang"
	}

	return v.RunValidationSuite(ctx, "Drik Panchang Validation", referenceData)
}

// ValidateRegionalCalculations validates calculations for different regions
func (v *Validator) ValidateRegionalCalculations(ctx context.Context, region string, referenceData []ReferenceData) *ValidationSuite {
	return v.RunValidationSuite(ctx, fmt.Sprintf("Regional Validation (%s)", region), referenceData)
}

// GenerateValidationReport generates a human-readable validation report
func (suite *ValidationSuite) GenerateReport() string {
	report := fmt.Sprintf("Validation Suite: %s\n", suite.Name)
	report += fmt.Sprintf("Executed at: %s\n", suite.TestedAt.Format("2006-01-02 15:04:05"))
	report += fmt.Sprintf("Total Tests: %d\n", suite.TotalTests)
	report += fmt.Sprintf("Passed: %d\n", suite.PassedTests)
	report += fmt.Sprintf("Failed: %d\n", suite.FailedTests)
	report += fmt.Sprintf("Success Rate: %.2f%%\n", suite.SuccessRate)
	report += fmt.Sprintf("Execution Time: %s\n\n", suite.ExecutionTime)

	// List failed tests
	if suite.FailedTests > 0 {
		report += "Failed Tests:\n"
		for _, result := range suite.Results {
			if !result.Passed {
				report += fmt.Sprintf("  - %s: Expected %v, Got %v (Error: %.2f %s)\n",
					result.TestName, result.Expected, result.Actual, result.Error, result.ErrorUnit)
				if result.Notes != "" {
					report += fmt.Sprintf("    Notes: %s\n", result.Notes)
				}
			}
		}
	}

	return report
}

// ExportToJSON exports validation suite results to JSON
func (suite *ValidationSuite) ExportToJSON() ([]byte, error) {
	return nil, fmt.Errorf("JSON export not yet implemented")
}

// CompareSuites compares two validation suites and returns differences
func CompareSuites(suite1, suite2 *ValidationSuite) map[string]interface{} {
	comparison := make(map[string]interface{})

	comparison["suite1_name"] = suite1.Name
	comparison["suite2_name"] = suite2.Name
	comparison["success_rate_diff"] = suite2.SuccessRate - suite1.SuccessRate
	comparison["passed_tests_diff"] = suite2.PassedTests - suite1.PassedTests
	comparison["failed_tests_diff"] = suite2.FailedTests - suite1.FailedTests

	return comparison
}
