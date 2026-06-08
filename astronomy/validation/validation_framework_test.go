package validation

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
)

func TestValidationResult(t *testing.T) {
	result := ValidationResult{
		TestName:        "Test Validation",
		Passed:          true,
		Expected:        "Expected Value",
		Actual:          "Actual Value",
		Tolerance:       1.0,
		Error:           0.5,
		ErrorUnit:       "degrees",
		TestedAt:        time.Now(),
		Duration:        time.Millisecond * 100,
		Notes:           "Test notes",
		Source:          "Test Source",
		CalculationType: "Test Type",
	}

	if result.TestName != "Test Validation" {
		t.Errorf("Expected test name 'Test Validation', got %s", result.TestName)
	}

	if !result.Passed {
		t.Error("Expected result to pass")
	}
}

func TestValidationSuiteReport(t *testing.T) {
	suite := &ValidationSuite{
		Name:          "Test Suite",
		Description:   "Test description",
		TotalTests:    10,
		PassedTests:   8,
		FailedTests:   2,
		SuccessRate:   80.0,
		ExecutionTime: time.Second,
		TestedAt:      time.Now(),
		Results: []ValidationResult{
			{
				TestName:  "Test 1",
				Passed:    true,
				Expected:  "A",
				Actual:    "A",
				ErrorUnit: "degrees",
			},
			{
				TestName:  "Test 2",
				Passed:    false,
				Expected:  "B",
				Actual:    "C",
				Error:     5.0,
				ErrorUnit: "minutes",
				Notes:     "Failed due to timing",
			},
		},
	}

	report := suite.GenerateReport()

	if report == "" {
		t.Error("Expected non-empty report")
	}

	// Check if report contains key information
	if len(report) < 100 {
		t.Errorf("Report seems too short: %d characters", len(report))
	}

	t.Logf("Generated report:\n%s", report)
}

func TestValidationSuiteExportToJSON(t *testing.T) {
	suite := &ValidationSuite{
		Name:        "JSON Export Suite",
		TotalTests:  1,
		PassedTests: 1,
		SuccessRate: 100,
		TestedAt:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Results: []ValidationResult{
			{
				TestName:        "Tithi Validation",
				Passed:          true,
				Expected:        "Panchami",
				Actual:          "Panchami",
				ErrorUnit:       "name mismatch",
				CalculationType: "Tithi",
			},
		},
	}

	data, err := suite.ExportToJSON()
	if err != nil {
		t.Fatalf("expected JSON export to succeed, got error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty JSON export")
	}

	var exported map[string]interface{}
	if err := json.Unmarshal(data, &exported); err != nil {
		t.Fatalf("expected parseable JSON export, got error: %v", err)
	}
	if exported["name"] != suite.Name {
		t.Fatalf("expected exported suite name %q, got %v", suite.Name, exported["name"])
	}
}

func TestReferenceData(t *testing.T) {
	location := astronomy.Location{
		Latitude:  13.0827,
		Longitude: 80.2707,
	}

	ref := ReferenceData{
		Source:        "Test Source",
		Date:          time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Location:      location,
		Sunrise:       time.Date(2024, 1, 1, 6, 30, 0, 0, time.UTC),
		Sunset:        time.Date(2024, 1, 1, 18, 15, 0, 0, time.UTC),
		TithiName:     "Pratipada",
		NakshatraName: "Ashwini",
		YogaName:      "Vishkambha",
		Data:          make(map[string]interface{}),
	}

	if ref.Source != "Test Source" {
		t.Errorf("Expected source 'Test Source', got %s", ref.Source)
	}

	if ref.TithiName != "Pratipada" {
		t.Errorf("Expected tithi 'Pratipada', got %s", ref.TithiName)
	}

	if ref.Location.Latitude != location.Latitude || ref.Location.Longitude != location.Longitude {
		t.Errorf("Expected location %+v, got %+v", location, ref.Location)
	}
}

func TestCompareSuites(t *testing.T) {
	suite1 := &ValidationSuite{
		Name:        "Suite 1",
		TotalTests:  10,
		PassedTests: 8,
		FailedTests: 2,
		SuccessRate: 80.0,
	}

	suite2 := &ValidationSuite{
		Name:        "Suite 2",
		TotalTests:  10,
		PassedTests: 9,
		FailedTests: 1,
		SuccessRate: 90.0,
	}

	comparison := CompareSuites(suite1, suite2)

	if comparison["suite1_name"] != "Suite 1" {
		t.Errorf("Expected suite1_name 'Suite 1', got %v", comparison["suite1_name"])
	}

	if comparison["suite2_name"] != "Suite 2" {
		t.Errorf("Expected suite2_name 'Suite 2', got %v", comparison["suite2_name"])
	}

	successRateDiff := comparison["success_rate_diff"].(float64)
	if successRateDiff != 10.0 {
		t.Errorf("Expected success rate diff 10.0, got %f", successRateDiff)
	}

	passedDiff := comparison["passed_tests_diff"].(int)
	if passedDiff != 1 {
		t.Errorf("Expected passed tests diff 1, got %d", passedDiff)
	}
}

func TestValidatorCreation(t *testing.T) {
	// Create test components
	manager := createTestEphemerisManager(t)

	tithiCalc := astronomy.NewTithiCalculator(manager)
	nakshatraCalc := astronomy.NewNakshatraCalculator(manager)
	yogaCalc := astronomy.NewYogaCalculator(manager)
	karanaCalc := astronomy.NewKaranaCalculator(manager)
	varaCalc := astronomy.NewVaraCalculator()

	validator := NewValidator(
		tithiCalc,
		nakshatraCalc,
		yogaCalc,
		karanaCalc,
		varaCalc,
		manager,
	)

	if validator == nil {
		t.Fatal("Expected non-nil validator")
	}

	if validator.tithiCalc == nil {
		t.Error("Expected non-nil tithiCalc")
	}

	if validator.ephemerisManager == nil {
		t.Error("Expected non-nil ephemerisManager")
	}
}

func TestValidatePlanetaryPosition(t *testing.T) {
	manager := createTestEphemerisManager(t)

	validator := NewValidator(
		nil, nil, nil, nil, nil,
		manager,
	)

	// Test with known planetary position (approximate)
	date := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	// Sun's approximate longitude on Jan 1, 2024 is around 280 degrees
	result := validator.ValidatePlanetaryPosition(
		context.Background(),
		"sun",
		280.0,
		date,
		5.0, // 5 degree tolerance
	)

	if result.TestName == "" {
		t.Error("Expected non-empty test name")
	}

	if result.ErrorUnit != "degrees" {
		t.Errorf("Expected error unit 'degrees', got %s", result.ErrorUnit)
	}

	t.Logf("Planetary position validation result:")
	t.Logf("  Test: %s", result.TestName)
	t.Logf("  Passed: %v", result.Passed)
	t.Logf("  Expected: %v", result.Expected)
	t.Logf("  Actual: %v", result.Actual)
	t.Logf("  Error: %.6f %s", result.Error, result.ErrorUnit)
}

func TestNameValidationReportsMismatchError(t *testing.T) {
	manager := createTestEphemerisManager(t)
	validator := NewValidator(
		astronomy.NewTithiCalculator(manager),
		astronomy.NewNakshatraCalculator(manager),
		nil,
		nil,
		nil,
		manager,
	)

	ref := ReferenceData{
		Source:        "Test Source",
		Date:          time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		TithiName:     "Reference Tithi Not Produced",
		NakshatraName: "Reference Nakshatra Not Produced",
	}

	tithiResult := validator.ValidateTithi(context.Background(), ref, 0)
	if tithiResult.Passed {
		t.Fatal("expected mismatched tithi name to fail validation")
	}
	if tithiResult.Error <= 0 {
		t.Fatalf("expected mismatched tithi name to report non-zero error, got %.2f", tithiResult.Error)
	}
	if tithiResult.ErrorUnit != "name mismatch" {
		t.Fatalf("expected tithi error unit to describe a name mismatch, got %q", tithiResult.ErrorUnit)
	}

	nakshatraResult := validator.ValidateNakshatra(context.Background(), ref, 0)
	if nakshatraResult.Passed {
		t.Fatal("expected mismatched nakshatra name to fail validation")
	}
	if nakshatraResult.Error <= 0 {
		t.Fatalf("expected mismatched nakshatra name to report non-zero error, got %.2f", nakshatraResult.Error)
	}
	if nakshatraResult.ErrorUnit != "name mismatch" {
		t.Fatalf("expected nakshatra error unit to describe a name mismatch, got %q", nakshatraResult.ErrorUnit)
	}
}

func TestValidationSuiteStatistics(t *testing.T) {
	suite := &ValidationSuite{
		Name: "Statistics Test",
		Results: []ValidationResult{
			{Passed: true},
			{Passed: true},
			{Passed: false},
			{Passed: true},
			{Passed: false},
		},
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

	if suite.TotalTests != 5 {
		t.Errorf("Expected 5 total tests, got %d", suite.TotalTests)
	}

	if suite.PassedTests != 3 {
		t.Errorf("Expected 3 passed tests, got %d", suite.PassedTests)
	}

	if suite.FailedTests != 2 {
		t.Errorf("Expected 2 failed tests, got %d", suite.FailedTests)
	}

	expectedRate := 60.0
	if suite.SuccessRate != expectedRate {
		t.Errorf("Expected success rate %.2f%%, got %.2f%%", expectedRate, suite.SuccessRate)
	}
}

func TestInvalidPlanetValidation(t *testing.T) {
	manager := createTestEphemerisManager(t)

	validator := NewValidator(
		nil, nil, nil, nil, nil,
		manager,
	)

	// Test with invalid planet
	date := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	result := validator.ValidatePlanetaryPosition(
		context.Background(),
		"invalid_planet",
		0.0,
		date,
		1.0,
	)

	if result.Passed {
		t.Error("Expected validation to fail for invalid planet")
	}

	if result.Notes == "" {
		t.Error("Expected error notes for invalid planet")
	}
}

func TestValidationResultDuration(t *testing.T) {
	result := ValidationResult{
		Duration: time.Millisecond * 150,
	}

	if result.Duration < time.Millisecond {
		t.Errorf("Duration too short: %s", result.Duration)
	}

	if result.Duration > time.Second {
		t.Errorf("Duration too long: %s", result.Duration)
	}
}

func BenchmarkValidatePlanetaryPosition(b *testing.B) {
	manager := createTestEphemerisManager(b)

	validator := NewValidator(
		nil, nil, nil, nil, nil,
		manager,
	)

	date := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.ValidatePlanetaryPosition(
			ctx,
			"sun",
			280.0,
			date,
			5.0,
		)
	}
}

// Helper function to create test ephemeris manager
func createTestEphemerisManager(tb testing.TB) *ephemeris.Manager {
	tb.Helper()

	primary := ephemeris.NewJPLProvider()
	fallback := ephemeris.NewSwissProvider()

	cache := ephemeris.NewMemoryCache(100, 1*time.Hour)

	return ephemeris.NewManager(primary, fallback, cache)
}
