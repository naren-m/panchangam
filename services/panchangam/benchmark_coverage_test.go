package panchangam

import (
	"testing"

	"github.com/naren-m/panchangam/observability"
)

func BenchmarkFeatureCoverage(b *testing.B) {
	observability.NewLocalObserver()

	b.Run("Feature_Performance_All_Elements", func(b *testing.B) {
		benchmarkAllPanchangamElements(b)
	})

	b.Run("Feature_Performance_Service_Layer", func(b *testing.B) {
		benchmarkServiceLayer(b)
	})

	b.Run("Feature_Performance_Astronomy", func(b *testing.B) {
		benchmarkAstronomyCalculations(b)
	})

	b.Run("Feature_Performance_Observability", func(b *testing.B) {
		benchmarkObservability(b)
	})
}
