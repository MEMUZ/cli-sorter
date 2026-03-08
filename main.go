package main

import (
	"cli-sorter/sorter"
	"cli-sorter/utils"
	"flag"

	"github.com/fatih/color"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "preview sorting without moving files")
	quiet := flag.Bool("quiet", false, "show only final statistics")

	// short
	dryRunShort := flag.Bool("d", false, "dry-run (short)")
	quietShort := flag.Bool("q", false, "quiet mode (short)")

	flag.Parse()

	if flag.NArg() < 1 {
		color.New(color.FgHiRed).Println("Usage: sorter [--dry-run | -d] [--quiet | -q] <directory>")
		return
	}

	dir := flag.Arg(0)

	if *dryRunShort {
		*dryRun = true
	}

	if *quietShort {
		*quiet = true
	}

	color.New(color.Bold).Println("Sorting folder:", dir)

	err := sorter.Sort(dir, *dryRun, *quiet)
	if err != nil {
		color.New(color.FgHiRed).Println("Error:", err)
		return
	}

	color.New(color.FgHiGreen, color.Bold).Println("Sorting complete")
	utils.WaitForEnter()
}
