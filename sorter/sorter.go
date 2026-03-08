package sorter

import (
	"cli-sorter/types"
	"cli-sorter/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Sort(dir string, dryRun bool, quiet bool, ignore string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	statsMap := map[string]int{}
	ignoreMap := utils.ParseIgnore(ignore)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := strings.ToLower(filepath.Ext(name))

		if ignoreMap[name] || ignoreMap[ext] {
			continue
		}

		category, ok := types.Rules[ext]
		if !ok {
			category = "other"
		}

		src := filepath.Join(dir, name)
		dstDir := filepath.Join(dir, category)
		dst := utils.GetUniqueFilePath(filepath.Join(dstDir, name))

		if dryRun {
			if !quiet {
				color.New(color.FgHiYellow).Printf("[DRY] %s -> %s\n", name, category)
			}
			statsMap[category]++
			continue
		}

		os.MkdirAll(dstDir, os.ModePerm)

		err := os.Rename(src, dst)
		if err != nil {
			color.New(color.FgHiRed).Println("Failed to move", name)
			continue
		}

		statsMap[category]++

		if !quiet {
			color.New(color.FgHiBlue).Printf("Moved: %s -> %s\n", name, category)
		}
	}

	utils.PrintStats(statsMap)

	return nil
}
