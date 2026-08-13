package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/new-world-coder/riskline/pkg/ruleset"
	"gopkg.in/yaml.v3"
)

// FileName is the project-level config searched from the working directory.
const FileName = ".riskline.yaml"

// EnvRegimes is the comma-separated override for default regimes.
const EnvRegimes = "RISKLINE_REGIMES"

// File is the on-disk project config shape.
type File struct {
	Regimes []string `yaml:"regimes" json:"regimes"`
}

// ResolveRegimes applies precedence:
//  1. explicit (request body or -regimes flag)
//  2. RISKLINE_REGIMES
//  3. .riskline.yaml in dir (usually ".")
//  4. default [eu-ai-act]
func ResolveRegimes(explicit []string, configDir string) ([]string, error) {
	if regs := ruleset.NormalizeRegimes(explicit); len(regs) > 0 {
		return regs, nil
	}
	if env := strings.TrimSpace(os.Getenv(EnvRegimes)); env != "" {
		parts := strings.Split(env, ",")
		if regs := ruleset.NormalizeRegimes(parts); len(regs) > 0 {
			return regs, nil
		}
	}
	if configDir == "" {
		configDir = "."
	}
	path := configDir
	if !strings.HasSuffix(path, FileName) {
		path = strings.TrimRight(configDir, "/\\") + string(os.PathSeparator) + FileName
	}
	data, err := os.ReadFile(path)
	if err == nil {
		var f File
		if err := yaml.Unmarshal(data, &f); err != nil {
			return nil, fmt.Errorf("config: parse %s: %w", path, err)
		}
		if regs := ruleset.NormalizeRegimes(f.Regimes); len(regs) > 0 {
			return regs, nil
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}
	return []string{ruleset.RegimeEUAIAct}, nil
}

// ParseRegimesFlag splits a comma-separated -regimes flag value.
func ParseRegimesFlag(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return ruleset.NormalizeRegimes(strings.Split(s, ","))
}
