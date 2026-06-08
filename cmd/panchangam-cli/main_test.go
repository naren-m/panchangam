package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriteCommandErrorWritesErrorLine(t *testing.T) {
	var output bytes.Buffer

	err := writeCommandError(&output, errors.New("command failed"))

	if err != nil {
		t.Fatalf("expected no write error, got %v", err)
	}
	if output.String() != "command failed\n" {
		t.Fatalf("expected command error on writer, got %q", output.String())
	}
}

func TestWriteCommandErrorIgnoresNilError(t *testing.T) {
	var output bytes.Buffer

	err := writeCommandError(&output, nil)

	if err != nil {
		t.Fatalf("expected no write error, got %v", err)
	}
	if output.String() != "" {
		t.Fatalf("expected no output, got %q", output.String())
	}
}
