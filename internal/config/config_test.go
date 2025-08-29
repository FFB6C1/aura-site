package config

import (
	"testing"
)

func TestSetPath(t *testing.T) {
	config.SetImportPath("test")
	if config.importPath != "test" {
		t.Fatalf("Expected: test, actual: %s", config.importPath)
	}
}
