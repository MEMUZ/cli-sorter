package cli

import "flag"

type Config struct {
	Dir        string
	DryRun     bool
	Quiet      bool
	Ignore     string
	Recursive  bool
	ConfigPath string
}

func ParseFlags() Config {
	cfg := Config{}

	flag.BoolVar(&cfg.DryRun, "dry-run", false, "preview sorting without moving files")
	flag.BoolVar(&cfg.Quiet, "quiet", false, "show only final statistics")
	flag.StringVar(&cfg.Ignore, "ignore", "", "comma separated list of files or extensions to ignore")
	flag.BoolVar(&cfg.Recursive, "recursive", false, "sort files recursively")
	flag.StringVar(&cfg.ConfigPath, "config", "", "path to your custom config file")

	flag.BoolVar(&cfg.DryRun, "d", false, "dry-run (short)")
	flag.BoolVar(&cfg.Quiet, "q", false, "quiet mode (short)")
	flag.StringVar(&cfg.Ignore, "i", "", "ignore files (short)")
	flag.BoolVar(&cfg.Recursive, "r", false, "recursive (short)")
	flag.StringVar(&cfg.ConfigPath, "c", "", "config (short)")

	flag.Parse()

	if flag.NArg() > 0 {
		cfg.Dir = flag.Arg(0)
	}

	return cfg
}
