package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configContent := `{
		"rules": {
			"custom_images": [".jpg", ".png", ".webp"],
			"custom_docs": [".pdf", ".docx"],
			"custom_code": [".go", ".py", ".js"]
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Rules) != 3 {
		t.Errorf("Expected 3 rule categories, got %d", len(cfg.Rules))
	}

	if len(cfg.Rules["custom_images"]) != 3 {
		t.Errorf("Expected 3 extensions for custom_images, got %d", len(cfg.Rules["custom_images"]))
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	// Невалидный JSON
	err := os.WriteFile(configPath, []byte(`{invalid json}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadConfig_EmptyRules(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty.json")

	err := os.WriteFile(configPath, []byte(`{"rules": {}}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Rules) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(cfg.Rules))
	}
}

func TestBuildRuleMap(t *testing.T) {
	cfg := &Config{
		Rules: map[string][]string{
			"images": {".jpg", ".png", ".gif"},
			"docs":   {".pdf", ".docx"},
			"code":   {".go", ".py"},
		},
	}

	ruleMap := cfg.BuildRuleMap()

	expectedCount := 7 // 3 + 2 + 2
	if len(ruleMap) != expectedCount {
		t.Errorf("Expected %d rules, got %d", expectedCount, len(ruleMap))
	}

	tests := []struct {
		ext      string
		expected string
	}{
		{".jpg", "images"},
		{".png", "images"},
		{".gif", "images"},
		{".pdf", "docs"},
		{".docx", "docs"},
		{".go", "code"},
		{".py", "code"},
	}

	for _, tt := range tests {
		if ruleMap[tt.ext] != tt.expected {
			t.Errorf("Extension %s: expected category %s, got %s", tt.ext, tt.expected, ruleMap[tt.ext])
		}
	}
}

func TestBuildRuleMap_EmptyConfig(t *testing.T) {
	cfg := &Config{
		Rules: map[string][]string{},
	}

	ruleMap := cfg.BuildRuleMap()
	if len(ruleMap) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(ruleMap))
	}
}

func TestBuildRuleMap_DuplicateExtensions(t *testing.T) {
	cfg := &Config{
		Rules: map[string][]string{
			"images": {".jpg", ".png"},
			"backup": {".jpg", ".jpeg"},
		},
	}

	ruleMap := cfg.BuildRuleMap()

	if ruleMap[".jpg"] != "backup" {
		t.Errorf("Expected .jpg to map to 'backup' (last wins), got '%s'", ruleMap[".jpg"])
	}
}

func TestBuildRuleMap_WithSpaces(t *testing.T) {
	cfg := &Config{
		Rules: map[string][]string{
			"images": {".jpg ", " .png", " .gif "},
		},
	}

	ruleMap := cfg.BuildRuleMap()

	if _, ok := ruleMap[".jpg "]; !ok {
		t.Error("Extension with trailing space should be preserved")
	}
}

func TestLoadConfig_ComplexConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "complex.json")

	configContent := `{
		"rules": {
			"projects": [".go", ".py", ".js", ".ts", ".java"],
			"media": [".mp4", ".avi", ".mkv", ".mp3", ".wav"],
			"archives": [".zip", ".rar", ".7z", ".tar", ".gz"],
			"documents": [".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"],
			"images": [".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp"]
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Rules) != 5 {
		t.Errorf("Expected 5 categories, got %d", len(cfg.Rules))
	}

	ruleMap := cfg.BuildRuleMap()
	if len(ruleMap) != 29 { // 5 + 5 + 5 + 7 + 7
		t.Errorf("Expected 29 rules, got %d", len(ruleMap))
	}
}

func TestLoadConfig_UnicodeCategories(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "unicode.json")

	configContent := `{
		"rules": {
			"изображения": [".jpg", ".png"],
			"документы": [".pdf", ".docx"],
			"コード": [".go", ".py"]
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if _, ok := cfg.Rules["изображения"]; !ok {
		t.Error("Unicode category name should be supported")
	}
	if _, ok := cfg.Rules["コード"]; !ok {
		t.Error("Japanese category name should be supported")
	}
}

func TestConfig_ValidateExtensions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "validate.json")

	configContent := `{
		"rules": {
			"images": ["jpg", "png"],
			"docs": [".pdf", ".docx"]
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	ruleMap := cfg.BuildRuleMap()

	if _, ok := ruleMap["jpg"]; !ok {
		t.Error("Extension without dot should be loaded (validation is app's responsibility)")
	}
}

func TestConfig_CategoryNames(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "categories.json")

	configContent := `{
		"rules": {
			"my-custom-category": [".jpg"],
			"Category_With_Underscore": [".png"],
			"123numeric": [".gif"]
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	expectedCategories := []string{"my-custom-category", "Category_With_Underscore", "123numeric"}
	for _, cat := range expectedCategories {
		if _, ok := cfg.Rules[cat]; !ok {
			t.Errorf("Category %s should be loaded", cat)
		}
	}
}

func TestConfig_EmptyExtensionList(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty_ext.json")

	configContent := `{
		"rules": {
			"empty_category": [],
			"normal_category": [".jpg"]
		}
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if len(cfg.Rules["empty_category"]) != 0 {
		t.Error("Empty extension list should remain empty")
	}

	ruleMap := cfg.BuildRuleMap()
	if len(ruleMap) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(ruleMap))
	}
}

func TestConfig_NullValues(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "null.json")

	configContent := `{
		"rules": null
	}`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Rules != nil {
		t.Error("Rules should be nil when JSON has null")
	}
}

func TestConfig_LargeConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "large.json")

	var sb strings.Builder
	sb.WriteString(`{"rules": {`)
	for i := 0; i < 50; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(`"category_` + string(rune('0'+i/10)) + string(rune('0'+i%10)) + `": [".ext` + string(rune('0'+i/10)) + string(rune('0'+i%10)) + `"]`)
	}
	sb.WriteString(`}}`)

	err := os.WriteFile(configPath, []byte(sb.String()), 0644)
	if err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed for large config: %v", err)
	}

	if len(cfg.Rules) != 50 {
		t.Errorf("Expected 50 categories, got %d", len(cfg.Rules))
	}
}
