package sorter

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSort(t *testing.T) {
	tmpDir := t.TempDir()

	files := []string{
		"photo.jpg",
		"document.pdf",
		"song.mp3",
		"archive.zip",
		"unknown.xyz",
		".gitignore",
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		file, err := os.Create(path)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", f, err)
		}
		file.Close()
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

	runtime.GC()
}

func TestSortDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	filePath := filepath.Join(tmpDir, "test.jpg")
	file, error := os.Create(filePath)
	if error != nil {
		t.Fatalf("Failed to create test.jpg: %v", error)
	}
	file.Close()

	err := Sort(tmpDir, true, true, "", false)
	if err != nil {
		t.Fatalf("Sort dry-run failed: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("File should not be moved in dry-run mode")
	}

	imagesDir := filepath.Join(tmpDir, "images")
	if _, err := os.Stat(imagesDir); !os.IsNotExist(err) {
		t.Error("Directory should not be created in dry-run mode")
	}

	runtime.GC()
}

func TestSortRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	// Folder structure:
	// tmpDir/
	// ├── root.jpg                    -> images/
	// ├── subdir1/
	// │   ├── photo.png               -> subdir1/images/
	// │   └── nested/
	// │       └── audio.mp3           -> subdir1/nested/audios/
	// ├── subdir2/
	// │   ├── doc.pdf                 -> subdir2/documents/
	// │   └── images/                 <- this should be skipped
	// │       └── should_skip.jpg     <- no moves
	// └── .gitignore                  <- ignores

	createTestFile(t, tmpDir, "root.jpg")

	sub1 := filepath.Join(tmpDir, "subdir1")
	os.MkdirAll(filepath.Join(sub1, "nested"), 0755)
	createTestFile(t, sub1, "photo.png")
	createTestFile(t, filepath.Join(sub1, "nested"), "audio.mp3")

	sub2 := filepath.Join(tmpDir, "subdir2")
	os.MkdirAll(filepath.Join(sub2, "images"), 0755)
	createTestFile(t, sub2, "doc.pdf")
	createTestFile(t, filepath.Join(sub2, "images"), "should_skip.jpg")

	createTestFile(t, tmpDir, ".gitignore")

	err := Sort(tmpDir, false, true, ".gitignore", true)
	if err != nil {
		t.Fatalf("Recursive sort failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "root.jpg", "images")

	checkFileMoved(t, sub1, "photo.png", "images")
	checkFileMoved(t, filepath.Join(sub1, "nested"), "audio.mp3", "audios")

	checkFileMoved(t, sub2, "doc.pdf", "documents")

	skippedFile := filepath.Join(tmpDir, "subdir2", "images", "should_skip.jpg")
	if _, err := os.Stat(skippedFile); os.IsNotExist(err) {
		t.Error("File inside category directory should NOT be moved")
	}

	if _, err := os.Stat(filepath.Join(tmpDir, ".gitignore")); os.IsNotExist(err) {
		t.Error(".gitignore should not be moved")
	}

	runtime.GC()
}

func TestSortRecursiveDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	sub := filepath.Join(tmpDir, "photos")
	os.MkdirAll(sub, 0755)
	createTestFile(t, tmpDir, "img1.jpg")
	createTestFile(t, sub, "img2.png")

	err := Sort(tmpDir, true, true, "", true)
	if err != nil {
		t.Fatalf("Recursive dry-run failed: %v", err)
	}

	roots := []string{
		filepath.Join(tmpDir, "img1.jpg"),
		filepath.Join(sub, "img2.png"),
	}
	for _, p := range roots {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("File should not be moved in dry-run: %s", p)
		}
	}

	imagesRoot := filepath.Join(tmpDir, "images")
	imagesSub := filepath.Join(sub, "images")
	if _, err := os.Stat(imagesRoot); !os.IsNotExist(err) {
		t.Error("Directory should not be created in dry-run mode")
	}
	if _, err := os.Stat(imagesSub); !os.IsNotExist(err) {
		t.Error("Directory should not be created in dry-run mode (subdir)")
	}

	runtime.GC()
}

func TestSortRecursiveWithIgnore(t *testing.T) {
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
	checkFileMoved(t, sub, "notes.txt", "documents")

	if _, err := os.Stat(filepath.Join(tmpDir, "temp.tmp")); os.IsNotExist(err) {
		t.Error("temp.tmp should be ignored and stay in place")
	}
	if _, err := os.Stat(filepath.Join(sub, "debug.log")); os.IsNotExist(err) {
		t.Error("debug.log should be ignored and stay in place")
	}

	runtime.GC()
}

func TestSortRecursiveSkipCategoryDirs(t *testing.T) {
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

	original1 := filepath.Join(imagesDir, "already_sorted.jpg")
	original2 := filepath.Join(imagesDir, "another.png")
	if _, err := os.Stat(original1); os.IsNotExist(err) {
		t.Error("Files inside category directory should not be processed")
	}
	if _, err := os.Stat(original2); os.IsNotExist(err) {
		t.Error("Files inside category directory should not be processed")
	}

	runtime.GC()
}

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
	src := filepath.Join(baseDir, fileName)
	dst := filepath.Join(baseDir, category, fileName)

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Errorf("File %s should be moved from source", fileName)
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Errorf("File %s should exist in %s directory", fileName, category)
	}
}
