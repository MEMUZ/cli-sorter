package cli

import (
	"flag"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	origArgs := os.Args
	origFlag := flag.CommandLine
	defer func() {
		os.Args = origArgs
		flag.CommandLine = origFlag
	}()

	tests := []struct {
		name       string
		args       []string
		wantDir    string
		wantDry    bool
		wantQuiet  bool
		wantIgnore string
	}{
		{
			name:       "Directory only",
			args:       []string{"cmd", "/tmp/test"},
			wantDir:    "/tmp/test",
			wantDry:    false,
			wantQuiet:  false,
			wantIgnore: "",
		},
		{
			name:       "With dry-run",
			args:       []string{"cmd", "-dry-run", "/tmp/test"},
			wantDir:    "/tmp/test",
			wantDry:    true,
			wantQuiet:  false,
			wantIgnore: "",
		},
		{
			name:       "With quiet",
			args:       []string{"cmd", "-quiet", "/tmp/test"},
			wantDir:    "/tmp/test",
			wantDry:    false,
			wantQuiet:  true,
			wantIgnore: "",
		},
		{
			name:       "With ignore",
			args:       []string{"cmd", "-ignore", ".git,.DS_Store", "/tmp/test"},
			wantDir:    "/tmp/test",
			wantDry:    false,
			wantQuiet:  false,
			wantIgnore: ".git,.DS_Store",
		},
		{
			name:       "Short flags",
			args:       []string{"cmd", "-d", "-q", "-i", ".tmp", "/tmp/test"},
			wantDir:    "/tmp/test",
			wantDry:    true,
			wantQuiet:  true,
			wantIgnore: ".tmp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			cfg := ParseFlags()

			if cfg.Dir != tt.wantDir {
				t.Errorf("Dir = %v, want %v", cfg.Dir, tt.wantDir)
			}
			if cfg.DryRun != tt.wantDry {
				t.Errorf("DryRun = %v, want %v", cfg.DryRun, tt.wantDry)
			}
			if cfg.Quiet != tt.wantQuiet {
				t.Errorf("Quiet = %v, want %v", cfg.Quiet, tt.wantQuiet)
			}
			if cfg.Ignore != tt.wantIgnore {
				t.Errorf("Ignore = %v, want %v", cfg.Ignore, tt.wantIgnore)
			}
		})
	}
}
