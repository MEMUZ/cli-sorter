package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func createTestFile(t *testing.T, dir, name string) {
	t.Helper()
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create %s: %v", name, err)
	}
	f.Close()
}

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

	t.Run("Multiple collisions", func(t *testing.T) {
		basePath := filepath.Join(tmpDir, "test.txt")

		for i := 0; i < 3; i++ {
			var name string
			if i == 0 {
				name = "test.txt"
			} else {
				name = "test (" + string(rune('0'+i)) + ").txt"
			}
			f, err := os.Create(filepath.Join(tmpDir, name))
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
			f.Close()
		}

		result := GetUniqueFilePath(basePath)
		expected := filepath.Join(tmpDir, "test (3).txt")
		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}

func TestRemoveEmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()

	empty1 := filepath.Join(tmpDir, "empty1")
	empty2 := filepath.Join(tmpDir, "level1", "level2")
	notEmpty := filepath.Join(tmpDir, "notEmpty")

	os.MkdirAll(empty1, 0755)
	os.MkdirAll(empty2, 0755)
	os.MkdirAll(notEmpty, 0755)

	createTestFile(t, notEmpty, "file.txt")

	categoryDir := filepath.Join(tmpDir, "images")
	os.MkdirAll(categoryDir, 0755)

	err := RemoveEmptyDirs(tmpDir)
	if err != nil {
		t.Fatalf("RemoveEmptyDirs failed: %v", err)
	}

	if _, err := os.Stat(empty1); !os.IsNotExist(err) {
		t.Error("Empty directory should be removed")
	}
	if _, err := os.Stat(empty2); !os.IsNotExist(err) {
		t.Error("Empty nested directory should be removed")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "level1")); !os.IsNotExist(err) {
		t.Error("Empty parent directory should be removed")
	}

	if _, err := os.Stat(notEmpty); os.IsNotExist(err) {
		t.Error("Non-empty directory should not be removed")
	}

	if _, err := os.Stat(categoryDir); os.IsNotExist(err) {
		t.Error("Category directory should not be removed")
	}

	runtime.GC()
}

func TestRemoveEmptyDirs_RootProtection(t *testing.T) {
	tmpDir := t.TempDir()

	err := RemoveEmptyDirs(tmpDir)
	if err != nil {
		t.Fatalf("RemoveEmptyDirs failed: %v", err)
	}

	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Root directory should not be removed")
	}

	runtime.GC()
}
