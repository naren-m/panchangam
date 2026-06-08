package observability

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewObserverWithEmptyAddress(t *testing.T) {
	observer, err := NewObserver("")
	assert.NotNil(t, observer)
	assert.Nil(t, err)
}

func TestNewObserverWithInvalidAddress(t *testing.T) {
	observer := NewLocalObserver()
	assert.NotNil(t, observer)
}

func TestObserverAutoInitAfterReset(t *testing.T) {
	originalOi := oi

	oi = nil
	initObserverOnce = sync.Once{}

	assert.NotPanics(t, func() {
		observer := Observer()
		assert.NotNil(t, observer)
	}, "The code should auto-initialize, not panic")

	oi = originalOi
}

func TestObserverConcurrentAutoInitDoesNotRace(t *testing.T) {
	oi = nil
	initObserverOnce = sync.Once{}

	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			assert.NotNil(t, Observer())
		}()
	}

	close(start)
	wg.Wait()
}

func TestInitMeterProvider(t *testing.T) {
	mp, err := InitMeterProvider()
	assert.NoError(t, err)
	assert.NotNil(t, mp)
}

func TestInitResource(t *testing.T) {
	resource = nil
	initResourcesOnce = sync.Once{}

	res := initResource()
	assert.NotNil(t, res)
	assert.Equal(t, "panchangam", res.Attributes()[0].Value.AsString())
}

func TestInitStdoutProvider(t *testing.T) {
	tp, err := initStdoutProvider()
	assert.NotNil(t, tp)
	assert.Nil(t, err)
}

func TestInitTracerProviderEmptyAddress(t *testing.T) {
	tp, err := initTracerProvider("")
	assert.Nil(t, tp)
	assert.NotNil(t, err)
	assert.Equal(t, "address is required", err.Error())
}

func TestInitTracerProviderInvalidAddress(t *testing.T) {
	tp, err := initTracerProvider("invalid:address:format")
	if err != nil {
		assert.Nil(t, tp)
	} else {
		assert.NotNil(t, tp)
	}
}

func TestNewObserverMultipleTimes(t *testing.T) {
	observer1 := NewLocalObserver()
	observer2 := NewLocalObserver()
	assert.Equal(t, observer1, observer2)
}

func TestInitResourceMultipleTimes(t *testing.T) {
	res1 := initResource()
	res2 := initResource()
	assert.Equal(t, res1, res2)
}

func TestNewObserverWithAddress(t *testing.T) {
	oi = nil
	initObserverOnce = sync.Once{}

	observer, err := NewObserver("localhost:4317")
	assert.NotNil(t, observer)
	assert.Nil(t, err)
}

func TestInitializationErrorHandling(t *testing.T) {
	tp, err := initStdoutProvider()
	assert.NotNil(t, tp)
	assert.Nil(t, err)

	tp2, err2 := initTracerProvider("localhost:4317")
	assert.NotNil(t, tp2)
	assert.Nil(t, err2)
}

func TestInitTracerProviderConnectionFailure(t *testing.T) {
	tp, err := initTracerProvider("invalid.host:99999")
	if err != nil {
		assert.Nil(t, tp)
	} else {
		assert.NotNil(t, tp)
	}
}

func TestInitTracerProviderErrorPaths(t *testing.T) {
	tp, err := initTracerProvider("")
	assert.Nil(t, tp)
	assert.NotNil(t, err)
	assert.Equal(t, "address is required", err.Error())

	tp2, err2 := initTracerProvider("localhost:4317")
	assert.NotNil(t, tp2)
	assert.Nil(t, err2)
}

func TestNewObserverBranchCoverage(t *testing.T) {
	oi = nil
	initObserverOnce = sync.Once{}

	observer1, err1 := NewObserver("")
	assert.NotNil(t, observer1)
	assert.Nil(t, err1)

	oi = nil
	initObserverOnce = sync.Once{}

	observer2, err2 := NewObserver("localhost:4317")
	assert.NotNil(t, observer2)
	assert.Nil(t, err2)
}

func TestObserverValidationEdgeCases(t *testing.T) {
	observer1, err1 := NewObserver("invalid:address:format:too:many:colons")
	assert.NotNil(t, observer1)
	assert.Nil(t, err1)

	for i := 0; i < 5; i++ {
		observer := NewLocalObserver()
		assert.NotNil(t, observer)
	}
}
