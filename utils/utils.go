package utils

import (
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
