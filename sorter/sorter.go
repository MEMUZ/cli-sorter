package sorter

import (
	"cli-sorter/types"
	"cli-sorter/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Sort(dir string, dryRun bool, quiet bool, ignore string, recursive bool) error {
	statsMap := map[string]int{}
	ignoreMap := utils.ParseIgnore(ignore)

	if !recursive {
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			processFile(dir, dir, file.Name(), dryRun, quiet, ignoreMap, statsMap)
		}
	} else {
		filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if types.IsCategoryDir(d.Name()) {
					return filepath.SkipDir
				}

				return nil
			}

			processFile(dir, filepath.Dir(path), d.Name(), dryRun, quiet, ignoreMap, statsMap)

			return nil
		})
	}

	utils.PrintStats(statsMap)

	if recursive && !dryRun {
		utils.RemoveEmptyDirs(dir)
	}

	return nil
}

func processFile(rootDir, currentDir, name string, dryRun, quiet bool, ignoreMap map[string]bool, stats map[string]int) {
	ext := strings.ToLower(filepath.Ext(name))

	if ignoreMap[name] || ignoreMap[ext] {
		return
	}

	category, ok := types.Rules[ext]
	if !ok {
		category = "other"
	}

	src := filepath.Join(currentDir, name)
	dstDir := filepath.Join(rootDir, category)
	dst := utils.GetUniqueFilePath(filepath.Join(dstDir, name))

	if dryRun {
		if !quiet {
			color.New(color.FgHiYellow).Printf("[DRY] %s -> %s\n", name, category)
		}
		stats[category]++
		return
	}

	os.MkdirAll(dstDir, os.ModePerm)

	err := os.Rename(src, dst)
	if err != nil {
		color.New(color.FgHiRed).Println("Failed to move", name)
		return
	}

	stats[category]++

	if !quiet {
		color.New(color.FgHiBlue).Printf("Moved: %s -> %s\n", name, category)
	}
}
