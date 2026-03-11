package sorter

import (
	"cli-sorter/config"
	"cli-sorter/types"
	"cli-sorter/utils"
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

	ignoreMap := utils.ParseIgnore(".gitignore")

	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  false,
		Rules:      types.Rules,
		Categories: categories,
	}

	err := Sort(opts)
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

	ignoreMap1 := utils.ParseIgnore("")

	opts1 := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     ignoreMap1,
		Recursive:  false,
		Rules:      types.Rules,
		Categories: categories,
	}

	err = Sort(opts1)
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

	ignoreMap := utils.ParseIgnore("")
	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     true,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  false,
		Rules:      types.Rules,
		Categories: categories,
	}

	err := Sort(opts)
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

	ignoreMap := utils.ParseIgnore("")
	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  false,
		Rules:      types.Rules,
		Categories: categories,
	}

	err := Sort(opts)
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

	ignoreMap := utils.ParseIgnore(".gitignore")
	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  true,
		Rules:      types.Rules,
		Categories: categories,
	}

	err := Sort(opts)
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

	ignoreMap := utils.ParseIgnore("")
	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     true,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  true,
		Rules:      types.Rules,
		Categories: categories,
	}

	err := Sort(opts)
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

	ignoreMap := utils.ParseIgnore(ignore)
	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  true,
		Rules:      types.Rules,
		Categories: categories,
	}
	err := Sort(opts)
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

	ignoreMap := utils.ParseIgnore("")
	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  true,
		Rules:      types.Rules,
		Categories: categories,
	}

	err := Sort(opts)
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

	ignoreMap := utils.ParseIgnore("")
	categories := types.BuildCategorySet(types.Rules)

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     ignoreMap,
		Recursive:  true,
		Rules:      types.Rules,
		Categories: categories,
	}

	err := Sort(opts)
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

func TestSort_WithCustomConfig(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "custom.json")
	configContent := `{
		"rules": {
			"my_images": [".jpg", ".png"],
			"my_docs": [".pdf", ".txt"],
			"my_code": [".go", ".py"]
		}
	}`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	rules := cfg.BuildRuleMap()
	categories := types.BuildCategorySet(rules)

	createTestFile(t, tmpDir, "photo.jpg")
	createTestFile(t, tmpDir, "document.pdf")
	createTestFile(t, tmpDir, "script.go")
	createTestFile(t, tmpDir, "unknown.xyz")

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     map[string]bool{},
		Recursive:  false,
		Rules:      rules,
		Categories: categories,
	}

	err = Sort(opts)
	if err != nil {
		t.Fatalf("Sort with custom config failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "photo.jpg", "my_images")
	checkFileMoved(t, tmpDir, "document.pdf", "my_docs")
	checkFileMoved(t, tmpDir, "script.go", "my_code")
	checkFileMoved(t, tmpDir, "unknown.xyz", "other")

	runtime.GC()
}

func TestSort_CustomConfigRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	configContent := `{
		"rules": {
			"custom_media": [".mp4", ".mp3"],
			"custom_docs": [".pdf"]
		}
	}`
	os.WriteFile(configPath, []byte(configContent), 0644)

	cfg, _ := config.LoadConfig(configPath)
	rules := cfg.BuildRuleMap()
	categories := types.BuildCategorySet(rules)

	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)
	createTestFile(t, tmpDir, "video.mp4")
	createTestFile(t, subDir, "audio.mp3")
	createTestFile(t, subDir, "doc.pdf")

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     map[string]bool{},
		Recursive:  true,
		Rules:      rules,
		Categories: categories,
	}

	err := Sort(opts)
	if err != nil {
		t.Fatalf("Recursive sort with config failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "video.mp4", "custom_media")
	checkFileMoved(t, tmpDir, "audio.mp3", "custom_media")
	checkFileMoved(t, tmpDir, "doc.pdf", "custom_docs")

	runtime.GC()
}

func TestSort_CustomConfigDryRun(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	configContent := `{"rules": {"custom": [".jpg"]}}`
	os.WriteFile(configPath, []byte(configContent), 0644)

	cfg, _ := config.LoadConfig(configPath)
	rules := cfg.BuildRuleMap()
	categories := types.BuildCategorySet(rules)

	createTestFile(t, tmpDir, "test.jpg")

	opts := Options{
		Dir:        tmpDir,
		DryRun:     true,
		Quiet:      true,
		Ignore:     map[string]bool{},
		Recursive:  false,
		Rules:      rules,
		Categories: categories,
	}

	err := Sort(opts)
	if err != nil {
		t.Fatalf("Dry-run with config failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "test.jpg")); os.IsNotExist(err) {
		t.Error("File should not be moved in dry-run mode")
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "custom")); !os.IsNotExist(err) {
		t.Error("Category directory should not be created in dry-run")
	}

	runtime.GC()
}

func TestSort_CustomConfigWithIgnore(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.json")
	configContent := `{"rules": {"my_docs": [".pdf", ".txt"]}}`
	os.WriteFile(configPath, []byte(configContent), 0644)

	cfg, _ := config.LoadConfig(configPath)
	rules := cfg.BuildRuleMap()
	categories := types.BuildCategorySet(rules)

	createTestFile(t, tmpDir, "report.pdf")
	createTestFile(t, tmpDir, "temp.txt")
	createTestFile(t, tmpDir, "ignore.txt")

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     map[string]bool{"ignore.txt": true},
		Recursive:  false,
		Rules:      rules,
		Categories: categories,
	}

	err := Sort(opts)
	if err != nil {
		t.Fatalf("Sort with config and ignore failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "report.pdf", "my_docs")
	checkFileMoved(t, tmpDir, "temp.txt", "my_docs")

	if _, err := os.Stat(filepath.Join(tmpDir, "ignore.txt")); os.IsNotExist(err) {
		t.Error("Ignored file should stay in place")
	}

	runtime.GC()
}

func TestSort_DefaultRulesFallback(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "empty.json")
	os.WriteFile(configPath, []byte(`{"rules": {}}`), 0644)

	cfg, _ := config.LoadConfig(configPath)
	rules := cfg.BuildRuleMap()
	categories := types.BuildCategorySet(rules)

	createTestFile(t, tmpDir, "photo.jpg")
	createTestFile(t, tmpDir, "unknown.xyz")

	opts := Options{
		Dir:        tmpDir,
		DryRun:     false,
		Quiet:      true,
		Ignore:     map[string]bool{},
		Recursive:  false,
		Rules:      rules,
		Categories: categories,
	}

	err := Sort(opts)
	if err != nil {
		t.Fatalf("Sort with empty config failed: %v", err)
	}

	checkFileMoved(t, tmpDir, "photo.jpg", "other")
	checkFileMoved(t, tmpDir, "unknown.xyz", "other")

	runtime.GC()
}
