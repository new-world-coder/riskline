package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/new-world-coder/riskline/pkg/assure"
	"github.com/new-world-coder/riskline/pkg/config"
	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/ruleset"
	"github.com/new-world-coder/riskline/pkg/schema"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "diff":
			runDiff(os.Args[2:])
			return
		case "assure":
			runAssure(os.Args[2:])
			return
		case "help", "-h", "--help":
			printUsage()
			return
		}
	}
	runClassify(os.Args[1:])
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  %s [flags] <system-description.yaml|json>     Classify (default)\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s diff <baseline> <current>                  Material change detection\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s assure <classification.json> --probes f   Verify controls + conformity state\n", os.Args[0])
}

func runClassify(args []string) {
	fs := flag.NewFlagSet("classify", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	regimesFlag := fs.String("regimes", "", "comma-separated regime packs")
	listRegimes := fs.Bool("list-regimes", false, "print shipped regime pack IDs and exit")
	_ = fs.Parse(args)

	loader, err := ruleset.DefaultLoader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loader: %v\n", err)
		os.Exit(1)
	}

	if *listRegimes {
		for _, id := range loader.List() {
			fmt.Println(id)
		}
		return
	}

	if fs.NArg() != 1 {
		printUsage()
		os.Exit(2)
	}

	path := fs.Arg(0)
	req, err := loadRequest(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load input: %v\n", err)
		os.Exit(1)
	}

	explicit := config.ParseRegimesFlag(*regimesFlag)
	if len(explicit) == 0 {
		explicit = req.Regimes
	}
	regs, err := config.ResolveRegimes(explicit, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	req.Regimes = regs

	eng, err := engine.NewWithLoader(loader, regs)
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
		emitJSON(resp)
		return
	}
	printHuman(resp)
}

func runDiff(args []string) {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	_ = fs.Parse(args)

	if fs.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s diff <baseline.yaml> <current.yaml>\n", os.Args[0])
		os.Exit(2)
	}

	baseline, err := loadRequest(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "baseline: %v\n", err)
		os.Exit(1)
	}
	current, err := loadRequest(fs.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "current: %v\n", err)
		os.Exit(1)
	}

	resp := engine.DetectMaterialChange(baseline, current)
	if *jsonOut {
		emitJSON(resp)
		return
	}

	fmt.Printf("Material:         %v\n", resp.Material)
	fmt.Printf("Impact:           %s\n", resp.ConformityImpact)
	fmt.Printf("Changed fields:   %s\n", strings.Join(resp.ChangedFields, ", "))
	fmt.Printf("Prior fingerprint:  %s\n", resp.PriorFingerprint)
	fmt.Printf("Current fingerprint: %s\n", resp.CurrentFingerprint)
	fmt.Printf("\n%s\n", resp.Summary)
	fmt.Printf("\nDisclaimer\n%s\n", resp.Disclaimer)
}

func runAssure(args []string) {
	fs := flag.NewFlagSet("assure", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	probesPath := fs.String("probes", "", "JSON file mapping technical_hook to pass/fail")
	prevHash := fs.String("previous-hash", "", "optional prior evidence chain tail hash")
	_ = fs.Parse(args)

	if fs.NArg() != 1 || *probesPath == "" {
		fmt.Fprintf(os.Stderr, "usage: %s assure --probes probes.json [--json] <classification.json>\n", os.Args[0])
		os.Exit(2)
	}

	var classification schema.ClassifyResponse
	if err := loadJSONFile(fs.Arg(0), &classification); err != nil {
		fmt.Fprintf(os.Stderr, "classification: %v\n", err)
		os.Exit(1)
	}

	probes := map[string]bool{}
	if err := loadJSONFile(*probesPath, &probes); err != nil {
		fmt.Fprintf(os.Stderr, "probes: %v\n", err)
		os.Exit(1)
	}

	resp := assure.Evaluate(schema.AssureRequest{
		Classification: classification,
		Probes:         probes,
		PreviousHash:   *prevHash,
	})

	if *jsonOut {
		emitJSON(resp)
		return
	}

	fmt.Printf("Conformity state: %s\n", resp.ConformityState)
	fmt.Printf("Failed:           %d\n", resp.FailedCount)
	fmt.Printf("Unverified:       %d\n", resp.UnverifiedCount)
	fmt.Printf("\n%s\n", resp.Summary)
	if len(resp.EvidenceRecords) > 0 {
		fmt.Println("\nEvidence records (SHA-256 chain)")
		for _, rec := range resp.EvidenceRecords {
			fmt.Printf("  %s  %s  hash=%s\n", rec.ControlID, rec.Result, rec.ContentHash[:16]+"...")
		}
	}
	fmt.Printf("\nDisclaimer\n%s\n", resp.Disclaimer)
}

func loadJSONFile(path string, v any) error {
	data, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func emitJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}
}

func loadRequest(path string) (schema.ClassifyRequest, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- intentional offline classify input
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
	if resp.Regime != "" {
		fmt.Printf("Primary regime:   %s\n", resp.Regime)
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

	if len(resp.TechnicalControls) > 0 {
		fmt.Println("\nTechnical controls (Assure probes)")
		for _, tc := range resp.TechnicalControls {
			fmt.Printf("  - %s (%s) hook=%s\n", tc.ID, tc.PaperRef, tc.TechnicalHook)
		}
	}

	if len(resp.Classifications) > 1 {
		fmt.Println("\nAll regimes")
		for _, c := range resp.Classifications {
			fmt.Printf("  - %s (%s): %s [%s]\n", c.Regime, c.Character, c.RiskTier, c.RulesetVersion)
		}
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
