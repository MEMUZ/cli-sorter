package types

import (
	"strings"
	"testing"
)

func TestRulesFormat(t *testing.T) {
	for ext, category := range Rules {
		if ext != strings.TrimSpace(ext) {
			t.Errorf("Extension '%s' contains leading/trailing spaces", ext)
		}
		if category != strings.TrimSpace(category) {
			t.Errorf("Category '%s' for extension '%s' contains spaces", category, ext)
		}
		if !strings.HasPrefix(ext, ".") {
			t.Errorf("Extension '%s' should start with a dot", ext)
		}
	}
}
