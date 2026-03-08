package main

import (
	"cli-sorter/cli"
	"cli-sorter/sorter"
	"cli-sorter/utils"

	"github.com/fatih/color"
)

func main() {
	cfg := cli.ParseFlags()

	if cfg.Dir == "" {
		color.New(color.FgHiRed).Println("Usage: sorter [flags] <directory>")
		return
	}

	color.New(color.Bold).Println("Sorting folder:", cfg.Dir)

	err := sorter.Sort(cfg.Dir, cfg.DryRun, cfg.Quiet, cfg.Ignore)
	if err != nil {
		color.New(color.FgHiRed).Println("Error:", err)
		return
	}

	color.New(color.FgHiGreen, color.Bold).Println("Sorting complete")
	utils.WaitForEnter()
}
