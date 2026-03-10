package utils

import (
	"cli-sorter/types"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func PrintStats(stats map[string]int) {
	color.New(color.FgHiGreen, color.Bold).Println("\nSorting statistics:")

	total := 0

	for category, count := range stats {
		color.New(color.FgHiCyan).Printf("%-10s : %d\n", category, count)
		total += count
	}

	color.New(color.Bold).Printf("\nTotal files: %d\n", total)
}

func GetUniqueFilePath(dst string) string {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return dst
	}

	ext := filepath.Ext(dst)
	name := strings.TrimSuffix(filepath.Base(dst), ext)
	dir := filepath.Dir(dst)

	i := 1
	for {
		newName := fmt.Sprintf("%s (%d)%s", name, i, ext)
		newPath := filepath.Join(dir, newName)

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}

		i++
	}
}

func WaitForEnter() {
	fmt.Println()
	color.New(color.FgHiWhite).Println("Press Enter to exit...")
	fmt.Scanln()
}

func ParseIgnore(ignore string) map[string]bool {
	result := map[string]bool{}

	if ignore == "" {
		return result
	}

	items := strings.SplitSeq(ignore, ",")

	for item := range items {
		item = strings.TrimSpace(item)
		result[item] = true
	}

	return result
}

func RemoveEmptyDirs(root string) error {
	var dirs []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() && path != root {
			dirs = append(dirs, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		dir := dirs[i]

		if types.IsCategoryDir(filepath.Base(dir)) {
			continue
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		if len(entries) == 0 {
			os.Remove(dir)
		}
	}

	return nil
}
