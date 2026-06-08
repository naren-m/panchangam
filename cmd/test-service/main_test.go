package main

import (
	"bytes"
	"testing"
)

func TestWriteServiceTestFailure(t *testing.T) {
	var out bytes.Buffer

	if err := writeServiceTestFailure(&out, 2); err != nil {
		t.Fatalf("writeServiceTestFailure returned error: %v", err)
	}

	want := "WARN: 2 tests failed. Check error messages above.\nService test failed\n"
	if got := out.String(); got != want {
		t.Fatalf("writeServiceTestFailure() output = %q, want %q", got, want)
	}
}
