package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/schema"
	"gopkg.in/yaml.v3"
)

func main() {
	jsonOut := flag.Bool("json", false, "emit machine-readable JSON")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <system-description.yaml|json>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Classify an AI system under the EU AI Act. Runs fully offline.\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	path := flag.Arg(0)
	req, err := loadRequest(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load input: %v\n", err)
		os.Exit(1)
	}

	eng, err := engine.Default()
	if err != nil {
		fmt.Fprintf(os.Stderr, "engine: %v\n", err)
		os.Exit(1)
	}

	resp, err := eng.Classify(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "classify: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(resp); err != nil {
			fmt.Fprintf(os.Stderr, "encode: %v\n", err)
			os.Exit(1)
		}
		return
	}

	printHuman(resp)
}

func loadRequest(path string) (schema.ClassifyRequest, error) {
	// #nosec G304 -- CLI intentionally reads a user-provided local file path.
	data, err := os.ReadFile(path)
	if err != nil {
		return schema.ClassifyRequest{}, err
	}

	var req schema.ClassifyRequest
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".yaml"), strings.HasSuffix(lower, ".yml"):
		if err := yaml.Unmarshal(data, &req); err != nil {
			return schema.ClassifyRequest{}, err
		}
	default:
		if err := json.Unmarshal(data, &req); err != nil {
			// try YAML as a fallback for extension-less files
			if err2 := yaml.Unmarshal(data, &req); err2 != nil {
				return schema.ClassifyRequest{}, fmt.Errorf("json: %v; yaml: %v", err, err2)
			}
		}
	}
	return req, nil
}

func printHuman(resp schema.ClassifyResponse) {
	if resp.Name != "" {
		fmt.Printf("System:           %s\n", resp.Name)
	}
	fmt.Printf("Risk tier:        %s\n", resp.RiskTier)
	fmt.Printf("Ruleset:          %s (%s)\n", resp.RulesetVersion, resp.LastUpdated)
	fmt.Printf("\nRationale\n%s\n", resp.Rationale)

	if len(resp.MatchedRules) > 0 {
		fmt.Println("\nMatched rules")
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTIER\tREF")
		for _, r := range resp.MatchedRules {
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.ID, r.Tier, r.ArticleOrAnnex)
		}
		_ = w.Flush()
	}

	if len(resp.RecommendedControls) > 0 {
		fmt.Println("\nRecommended controls")
		for _, c := range resp.RecommendedControls {
			fmt.Printf("  - %s\n", c)
		}
	}

	if len(resp.JudgmentCalls) > 0 {
		fmt.Println("\nJudgment calls (review these)")
		for _, j := range resp.JudgmentCalls {
			fmt.Printf("  - %s\n", j)
		}
	}

	fmt.Printf("\nDisclaimer\n%s\n", resp.Disclaimer)
}
