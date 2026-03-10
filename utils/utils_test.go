package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseIgnore(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]bool
	}{
		{"Empty string", "", map[string]bool{}},
		{"Single item", ".git", map[string]bool{".git": true}},
		{"Multiple items", ".git,.DS_Store,temp.txt", map[string]bool{".git": true, ".DS_Store": true, "temp.txt": true}},
		{"With spaces", ".jpg , .png", map[string]bool{".jpg": true, ".png": true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseIgnore(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d\n", len(tt.expected), len(result))
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("Expected %s to be %v, got %v", k, v, result[k])
				}
			}
		})
	}
}

func TestGetUniqueFilePath(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("File does not exist", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test.txt")
		result := GetUniqueFilePath(path)
		if result != path {
			t.Errorf("Expected %s, got %s", path, result)
		}
	})

	t.Run("File exists - should add postfix", func(t *testing.T) {
		basePath := filepath.Join(tmpDir, "duplicate.txt")
		f, err := os.Create(basePath)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		f.Close()

		result := GetUniqueFilePath(basePath)
		expected := filepath.Join(tmpDir, "duplicate (1).txt")
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}
