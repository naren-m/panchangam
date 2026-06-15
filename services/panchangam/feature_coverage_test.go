package panchangam

import (
	"context"
	"testing"
	"time"

	"github.com/naren-m/panchangam/observability"
)

func TestFeatureCoverage(t *testing.T) {
	observability.NewLocalObserver()

	t.Run("Feature_Coverage_All_Elements", func(t *testing.T) {
		ctx := context.Background()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		testFeatureTITHI_001(t, ctx, testDate)
		testFeatureNAKSHATRA_001(t, ctx, testDate)
		testFeatureYOGA_001(t, ctx, testDate)
		testFeatureKARANA_001(t, ctx, testDate)
		testFeatureVARA_001(t, ctx, testDate)
	})

	t.Run("Feature_Coverage_Service_Integration", func(t *testing.T) {
		testFeatureSERVICE_001(t)
		testFeatureASTRONOMY_001(t)
		testFeatureOBSERVABILITY_001(t)
	})

	t.Run("Feature_Coverage_Quality_Assurance", func(t *testing.T) {
		testFeatureQA_001(t)
		testFeatureQA_002(t)
	})
}
