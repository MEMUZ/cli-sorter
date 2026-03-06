package sorter

import (
	"cli-sorter/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Sort(dir string, dryRun bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := strings.ToLower(filepath.Ext(name))

		category, ok := types.Rules[ext]
		if !ok {
			category = "other"
		}

		src := filepath.Join(dir, name)
		dstDir := filepath.Join(dir, category)
		dst := filepath.Join(dstDir, name)

		if dryRun {
			color.New(color.FgHiYellow).Printf("[DRY] %s -> %s", name, category)
		}

		os.MkdirAll(dstDir, os.ModePerm)

		err := os.Rename(src, dst)
		if err != nil {
			color.New(color.FgHiRed).Println("Failed to move", name)
			continue
		}

		color.New(color.FgHiBlue).Printf("Moved: %s -> %s", name, category)
	}

	return nil
}
