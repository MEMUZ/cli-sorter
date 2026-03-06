package main

import (
	"cli-sorter/sorter"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func waitForEnter() {
	fmt.Println()
	color.New(color.FgHiWhite).Println("Press Enter to exit...")
	fmt.Scanln()
}

func main() {
	if len(os.Args) < 2 {
		color.New(color.FgHiRed).Println("Usage: sorter <directory> [--dry-run]")
		return
	}

	dir := os.Args[1]
	dryRun := false

	if len(os.Args) > 2 && os.Args[2] == "--dry-run" {
		dryRun = true
	}

	color.New(color.Bold).Println("Sorting folder:", dir)

	err := sorter.Sort(dir, dryRun)
	if err != nil {
		color.New(color.FgHiRed).Println("Error:", err)
		return
	}

	color.New(color.FgHiGreen, color.Bold).Println("Sorting complete")
	waitForEnter()
}
