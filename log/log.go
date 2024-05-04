package log

import (
	"os"
"log/slog"
"github.com/vincentfree/opentelemetry/otelslog"
)



var Logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

var AddTracingContext = otelslog.AddTracingContext