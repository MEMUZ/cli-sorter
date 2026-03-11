package main

import (
	"cli-sorter/cli"
	"cli-sorter/config"
	"cli-sorter/sorter"
	"cli-sorter/types"
	"cli-sorter/utils"

	"github.com/fatih/color"
)

func main() {
	cfg := cli.ParseFlags()

	if cfg.Dir == "" {
		color.New(color.FgHiRed).Println("Usage: sorter [flags] <directory>")
		return
	}

	var rules = types.Rules

	if cfg.ConfigPath != "" {
		conf, err := config.LoadConfig(cfg.ConfigPath)
		if err != nil {
			color.New(color.FgHiRed, color.Bold).Println("Failed to load config:", err)
			return
		}
		rules = conf.BuildRuleMap()
	}

	color.New(color.Bold).Println("Sorting folder:", cfg.Dir)

	ignoreMap := utils.ParseIgnore(cfg.Ignore)
	categories := types.BuildCategorySet(rules)

	opts := sorter.Options{
		Dir:        cfg.Dir,
		DryRun:     cfg.DryRun,
		Quiet:      cfg.Quiet,
		Ignore:     ignoreMap,
		Recursive:  cfg.Recursive,
		Rules:      rules,
		Categories: categories,
	}

	err := sorter.Sort(opts)
	if err != nil {
		color.New(color.FgHiRed).Println("Error:", err)
		return
	}

	color.New(color.FgHiGreen, color.Bold).Println("Sorting complete")
	utils.WaitForEnter()
}
