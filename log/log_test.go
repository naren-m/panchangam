package log

import (
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	h := NewHandler(slog.LevelDebug, slog.NewTextHandler(os.Stdout, nil))

	if h.Handler() != h.handler {
		t.Errorf("Handler() = %v, want %v", h.Handler(), h.handler)
	}
}

func TestWithAttrs(t *testing.T) {
	h := NewHandler(slog.LevelDebug, slog.NewTextHandler(os.Stdout, nil))
	attrs := []slog.Attr{slog.String("key", "value")}

	newHandler := h.WithAttrs(attrs)

	if _, ok := newHandler.(*Handler); !ok {
		t.Errorf("WithAttrs() should return a Handler")
	}
}

func TestWithGroup(t *testing.T) {
	h := NewHandler(slog.LevelDebug, slog.NewTextHandler(os.Stdout, nil))
	groupName := "testGroup"

	newHandler := h.WithGroup(groupName)

	if _, ok := newHandler.(*Handler); !ok {
		t.Errorf("WithGroup() should return a Handler")
	}
}

func TestConvertSlogAttrToSpanAttr(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		attr    slog.Attr
		want    attribute.KeyValue
		wantErr bool
	}{
		{"string", "key", slog.String("key", "value"), attribute.String("key", "value"), false},
		{"bool", "key", slog.Bool("key", true), attribute.Bool("key", true), false},
		{"int64", "key", slog.Int64("key", 123), attribute.Int64("key", 123), false},
		{"uint64", "key", slog.Uint64("key", 123), attribute.Int64("key", 123), false},
		{"float64", "key", slog.Float64("key", 1.23), attribute.Float64("key", 1.23), false},
		{"duration", "key", slog.Duration("key", time.Second), attribute.String("key", "1s"), false},
		// {"time", "key", slog.Time("key", time.Unix(0, 0)), attribute.String("key", "1970-01-01 00:00:00 +0000 UTC"), false},
		// {"unsupported", "key", slog.Bytes("key", []byte("value")), attribute.String("key", "[118 97 108 117 101]"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertSlogAttrToSpanAttr(tt.key, tt.attr.Value)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertSlogAttrToSpanAttr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Value != tt.want.Value {
				t.Errorf("convertSlogAttrToSpanAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
