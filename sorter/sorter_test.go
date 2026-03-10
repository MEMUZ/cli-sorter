package sorter

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

func checkFileMoved(t *testing.T, baseDir, fileName, category string) {
	t.Helper()
	src := filepath.Join(baseDir, fileName)
	dst := filepath.Join(baseDir, category, fileName)

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Errorf("File %s should be moved from source", fileName)
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Errorf("File %s should exist in %s directory", fileName, category)
	}
}

func TestSort_NonRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{"photo.jpg", "document.pdf", "song.mp3", "archive.zip", "unknown.xyz", ".gitignore"}
	for _, f := range files {
		createTestFile(t, tmpDir, f)
	}

	err := Sort(tmpDir, false, true, ".gitignore", false)
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "photo.jpg", "images")
	checkFileMoved(t, tmpDir, "document.pdf", "documents")
	checkFileMoved(t, tmpDir, "song.mp3", "audios")
	checkFileMoved(t, tmpDir, "archive.zip", "archives")
	checkFileMoved(t, tmpDir, "unknown.xyz", "other")

	if _, err := os.Stat(filepath.Join(tmpDir, ".gitignore")); os.IsNotExist(err) {
		t.Error(".gitignore should not be moved")
	}

	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)
	createTestFile(t, subDir, "nested.jpg")

	err = Sort(tmpDir, false, true, "", false)
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(subDir, "nested.jpg")); os.IsNotExist(err) {
		t.Error("File in subdir should not be moved in non-recursive mode")
	}

	runtime.GC()
}

func TestSort_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	createTestFile(t, tmpDir, "test.jpg")

	err := Sort(tmpDir, true, true, "", false)
	if err != nil {
		t.Fatalf("Sort dry-run failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "test.jpg")); os.IsNotExist(err) {
		t.Error("File should not be moved in dry-run mode")
	}

	imagesDir := filepath.Join(tmpDir, "images")
	if _, err := os.Stat(imagesDir); !os.IsNotExist(err) {
		t.Error("Directory should not be created in dry-run mode")
	}

	runtime.GC()
}

func TestSort_FileCollision(t *testing.T) {
	tmpDir := t.TempDir()

	imagesDir := filepath.Join(tmpDir, "images")
	os.MkdirAll(imagesDir, 0755)
	createTestFile(t, imagesDir, "photo.jpg")

	createTestFile(t, tmpDir, "photo.jpg")

	err := Sort(tmpDir, false, true, "", false)
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	original := filepath.Join(imagesDir, "photo.jpg")
	if _, err := os.Stat(original); os.IsNotExist(err) {
		t.Error("Original file should not be overwritten")
	}

	renamed := filepath.Join(imagesDir, "photo (1).jpg")
	if _, err := os.Stat(renamed); os.IsNotExist(err) {
		t.Error("Colliding file should be renamed with suffix")
	}

	runtime.GC()
}

func TestSort_Recursive_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	// folder structure:
	// tmpDir/
	// ├── root.jpg                    -> images/ (root)
	// ├── subdir1/
	// │   ├── photo.png               -> images/ (root, not in subdir1)
	// │   └── nested/
	// │       └── audio.mp3           -> audios/ (root)
	// ├── subdir2/
	// │   └── doc.pdf                 -> documents/ (root)
	// └── .gitignore                  <- ignore

	createTestFile(t, tmpDir, "root.jpg")

	sub1 := filepath.Join(tmpDir, "subdir1")
	os.MkdirAll(filepath.Join(sub1, "nested"), 0755)
	createTestFile(t, sub1, "photo.png")
	createTestFile(t, filepath.Join(sub1, "nested"), "audio.mp3")

	sub2 := filepath.Join(tmpDir, "subdir2")
	os.MkdirAll(sub2, 0755)
	createTestFile(t, sub2, "doc.pdf")

	createTestFile(t, tmpDir, ".gitignore")

	err := Sort(tmpDir, false, true, ".gitignore", true)
	if err != nil {
		t.Fatalf("Recursive sort failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "root.jpg", "images")
	checkFileMoved(t, tmpDir, "photo.png", "images")
	checkFileMoved(t, tmpDir, "audio.mp3", "audios")
	checkFileMoved(t, tmpDir, "doc.pdf", "documents")

	if _, err := os.Stat(filepath.Join(tmpDir, ".gitignore")); os.IsNotExist(err) {
		t.Error(".gitignore should not be moved")
	}

	if _, err := os.Stat(sub1); !os.IsNotExist(err) {
		entries, _ := os.ReadDir(sub1)
		if len(entries) == 0 {
			t.Error("Empty subdir1 should be removed")
		}
	}

	runtime.GC()
}

func TestSort_Recursive_DryRun(t *testing.T) {
	tmpDir := t.TempDir()

	sub := filepath.Join(tmpDir, "photos")
	os.MkdirAll(sub, 0755)
	createTestFile(t, tmpDir, "img1.jpg")
	createTestFile(t, sub, "img2.png")

	err := Sort(tmpDir, true, true, "", true)
	if err != nil {
		t.Fatalf("Recursive dry-run failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "img1.jpg")); os.IsNotExist(err) {
		t.Error("File should not be moved in dry-run")
	}
	if _, err := os.Stat(filepath.Join(sub, "img2.png")); os.IsNotExist(err) {
		t.Error("File should not be moved in dry-run")
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "images")); !os.IsNotExist(err) {
		t.Error("Directory should not be created in dry-run mode")
	}

	runtime.GC()
}

func TestSort_Recursive_WithIgnore(t *testing.T) {
	tmpDir := t.TempDir()

	sub := filepath.Join(tmpDir, "work")
	os.MkdirAll(sub, 0755)

	createTestFile(t, tmpDir, "report.pdf")
	createTestFile(t, tmpDir, "temp.tmp")
	createTestFile(t, sub, "notes.txt")
	createTestFile(t, sub, "debug.log")

	ignore := ".tmp,debug.log"
	err := Sort(tmpDir, false, true, ignore, true)
	if err != nil {
		t.Fatalf("Sort with ignore failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "report.pdf", "documents")
	checkFileMoved(t, tmpDir, "notes.txt", "documents")

	if _, err := os.Stat(filepath.Join(tmpDir, "temp.tmp")); os.IsNotExist(err) {
		t.Error("temp.tmp should be ignored")
	}
	if _, err := os.Stat(filepath.Join(sub, "debug.log")); os.IsNotExist(err) {
		t.Error("debug.log should be ignored")
	}

	runtime.GC()
}

func TestSort_Recursive_SkipCategoryDirs(t *testing.T) {
	tmpDir := t.TempDir()

	imagesDir := filepath.Join(tmpDir, "images")
	os.MkdirAll(imagesDir, 0755)
	createTestFile(t, imagesDir, "already_sorted.jpg")
	createTestFile(t, imagesDir, "another.png")

	createTestFile(t, tmpDir, "new_photo.webp")

	err := Sort(tmpDir, false, true, "", true)
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "new_photo.webp", "images")

	if _, err := os.Stat(filepath.Join(imagesDir, "already_sorted.jpg")); os.IsNotExist(err) {
		t.Error("Files inside category directory should not be processed")
	}
	if _, err := os.Stat(filepath.Join(imagesDir, "another.png")); os.IsNotExist(err) {
		t.Error("Files inside category directory should not be processed")
	}

	runtime.GC()
}

func TestSort_Recursive_RemoveEmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()

	deep := filepath.Join(tmpDir, "level1", "level2", "level3")
	os.MkdirAll(deep, 0755)
	createTestFile(t, deep, "file.jpg")

	err := Sort(tmpDir, false, true, "", true)
	if err != nil {
		t.Fatalf("Sort failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "file.jpg", "images")

	if _, err := os.Stat(deep); !os.IsNotExist(err) {
		t.Error("Empty nested directories should be removed")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "level1", "level2")); !os.IsNotExist(err) {
		t.Error("Empty nested directories should be removed")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "level1")); !os.IsNotExist(err) {
		t.Error("Empty nested directories should be removed")
	}

	runtime.GC()
}
