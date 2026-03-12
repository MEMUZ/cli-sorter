package sorter

import (
	"cli-sorter/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type Options struct {
	Dir        string
	DryRun     bool
	Quiet      bool
	Ignore     map[string]bool
	Recursive  bool
	Rules      map[string]string
	Categories map[string]bool
}

func Sort(opts Options) error {
	statsMap := map[string]int{}

	if !opts.Recursive {
		files, err := os.ReadDir(opts.Dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			processFile(opts, opts.Dir, opts.Dir, file.Name(), statsMap)
		}
	} else {
		filepath.WalkDir(opts.Dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if opts.Categories[d.Name()] {
					return filepath.SkipDir
				}

				return nil
			}

			processFile(opts, opts.Dir, filepath.Dir(path), d.Name(), statsMap)

			return nil
		})
	}

	utils.PrintStats(statsMap)

	if opts.Recursive && !opts.DryRun {
		utils.RemoveEmptyDirs(opts.Dir, opts.Categories)
	}

	return nil
}

func processFile(opts Options, rootDir, currentDir, name string, stats map[string]int) {
	ext := strings.ToLower(filepath.Ext(name))

	if opts.Ignore[name] || opts.Ignore[ext] {
		return
	}

	category, ok := opts.Rules[ext]
	if !ok {
		category = "other"
	}

	src := filepath.Join(currentDir, name)
	dstDir := filepath.Join(rootDir, category)
	dst := utils.GetUniqueFilePath(filepath.Join(dstDir, name))

	if opts.DryRun {
		if !opts.Quiet {
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

	if !opts.Quiet {
		color.New(color.FgHiBlue).Printf("Moved: %s -> %s\n", name, category)
	}
}
