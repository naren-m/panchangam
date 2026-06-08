package panchangam

import (
	"time"

	"github.com/naren-m/panchangam/astronomy"
	"github.com/naren-m/panchangam/astronomy/ephemeris"
	"github.com/naren-m/panchangam/log"
	"github.com/naren-m/panchangam/observability"
	ppb "github.com/naren-m/panchangam/proto"
)

var logger = log.Logger()

type PanchangamServer struct {
	config           Config
	observer         observability.ObserverInterface
	ephemerisManager *ephemeris.Manager
	tithiCalc        *astronomy.TithiCalculator
	nakshatraCalc    *astronomy.NakshatraCalculator
	yogaCalc         *astronomy.YogaCalculator
	karanaCalc       *astronomy.KaranaCalculator
	varaCalc         *astronomy.VaraCalculator
	ppb.UnimplementedPanchangamServer
}

// NewPanchangamServer creates a server with the default ephemeris providers.
func NewPanchangamServer() *PanchangamServer {
	return NewPanchangamServerWithDependencies(defaultEphemerisManager(), DefaultConfig())
}

// NewPanchangamServerWithDependencies creates a server with explicit dependencies.
func NewPanchangamServerWithDependencies(manager *ephemeris.Manager, config Config) *PanchangamServer {
	if manager == nil {
		manager = defaultEphemerisManager()
	}

	tithiCalc := astronomy.NewTithiCalculator(manager)
	nakshatraCalc := astronomy.NewNakshatraCalculator(manager)
	yogaCalc := astronomy.NewYogaCalculator(manager)
	karanaCalc := astronomy.NewKaranaCalculator(manager)
	varaCalc := astronomy.NewVaraCalculator()

	return &PanchangamServer{
		config:           config,
		observer:         observability.Observer(),
		ephemerisManager: manager,
		tithiCalc:        tithiCalc,
		nakshatraCalc:    nakshatraCalc,
		yogaCalc:         yogaCalc,
		karanaCalc:       karanaCalc,
		varaCalc:         varaCalc,
	}
}

func defaultEphemerisManager() *ephemeris.Manager {
	jplProvider := ephemeris.NewJPLProvider()
	swissProvider := ephemeris.NewSwissProvider()
	memoryCache := ephemeris.NewMemoryCache(1000, 1*time.Hour)
	return ephemeris.NewManager(swissProvider, jplProvider, memoryCache)
}
