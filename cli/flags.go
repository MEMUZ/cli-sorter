package cli

import "flag"

type Config struct {
	Dir       string
	DryRun    bool
	Quiet     bool
	Ignore    string
	Recursive bool
}

func ParseFlags() Config {
	dryRun := flag.Bool("dry-run", false, "preview sorting without moving files")
	quiet := flag.Bool("quiet", false, "show only final statistics")
	ignore := flag.String("ignore", "", "comma separated list of files or extensions to ignore")
	recursive := flag.Bool("recursive", false, "sort files recursively")

	dryRunShort := flag.Bool("d", false, "dry-run (short)")
	quietShort := flag.Bool("q", false, "quiet mode (short)")
	ignoreShort := flag.String("i", "", "ignore files (short)")
	recursiveShort := flag.Bool("r", false, "recursive (short)")

	flag.Parse()

	cfg := Config{}

	if flag.NArg() > 0 {
		cfg.Dir = flag.Arg(0)
	}

	cfg.DryRun = *dryRun || *dryRunShort
	cfg.Quiet = *quiet || *quietShort
	cfg.Recursive = *recursive || *recursiveShort
	if *ignoreShort != "" {
		cfg.Ignore = *ignoreShort
	} else {
		cfg.Ignore = *ignore
	}

	return cfg
}
