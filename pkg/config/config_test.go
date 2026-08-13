package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/new-world-coder/riskline/pkg/ruleset"
)

func TestResolveDefault(t *testing.T) {
	t.Setenv(EnvRegimes, "")
	dir := t.TempDir()
	regs, err := ResolveRegimes(nil, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(regs) != 1 || regs[0] != ruleset.RegimeEUAIAct {
		t.Fatalf("got %v", regs)
	}
}

func TestResolveExplicitWins(t *testing.T) {
	t.Setenv(EnvRegimes, "eu-ai-act")
	regs, err := ResolveRegimes([]string{"eu-ai-act"}, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(regs) != 1 || regs[0] != ruleset.RegimeEUAIAct {
		t.Fatalf("got %v", regs)
	}
}

func TestResolveEnv(t *testing.T) {
	t.Setenv(EnvRegimes, " eu-ai-act , eu-ai-act ")
	regs, err := ResolveRegimes(nil, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(regs) != 1 || regs[0] != ruleset.RegimeEUAIAct {
		t.Fatalf("got %v", regs)
	}
}

func TestResolveFile(t *testing.T) {
	t.Setenv(EnvRegimes, "")
	dir := t.TempDir()
	path := filepath.Join(dir, FileName)
	if err := os.WriteFile(path, []byte("regimes:\n  - eu-ai-act\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	regs, err := ResolveRegimes(nil, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(regs) != 1 || regs[0] != ruleset.RegimeEUAIAct {
		t.Fatalf("got %v", regs)
	}
}

func TestParseRegimesFlag(t *testing.T) {
	got := ParseRegimesFlag("eu-ai-act, eu-ai-act")
	if len(got) != 1 || got[0] != ruleset.RegimeEUAIAct {
		t.Fatalf("got %v", got)
	}
	if ParseRegimesFlag("") != nil {
		t.Fatal("empty should be nil")
	}
}
