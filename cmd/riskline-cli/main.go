package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/new-world-coder/riskline/pkg/assure"
	"github.com/new-world-coder/riskline/pkg/config"
	"github.com/new-world-coder/riskline/pkg/engine"
	"github.com/new-world-coder/riskline/pkg/evidence"
	"github.com/new-world-coder/riskline/pkg/runtime"
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
		case "verify":
			runVerify(os.Args[2:])
			return
		case "runtime":
			runRuntime(os.Args[2:])
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
	fmt.Fprintf(os.Stderr, "  %s assure --probes f <classification.json>    Verify controls + conformity state\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s verify <bundle.evidence.json>              Verify Ed25519-signed evidence or receipt bundle\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s runtime verify <request.json> [--sign]     Verify runtime observation + optional signed receipt\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s runtime register <request.json>            Register runtime baseline and policy\n", os.Args[0])
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
	signOut := fs.Bool("sign", false, "emit Ed25519-signed .evidence.json bundle (local signing, not blockchain)")
	probesPath := fs.String("probes", "", "JSON file mapping technical_hook to pass/fail")
	prevHash := fs.String("previous-hash", "", "optional prior evidence chain tail hash")
	_ = fs.Parse(args)

	if fs.NArg() != 1 || *probesPath == "" {
		fmt.Fprintf(os.Stderr, "usage: %s assure --probes probes.json [--json] [--sign] <classification.json>\n", os.Args[0])
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

	if *signOut {
		_, priv, err := evidence.GenerateKeyPair()
		if err != nil {
			fmt.Fprintf(os.Stderr, "sign: %v\n", err)
			os.Exit(1)
		}
		bundle, err := evidence.SignBundle(resp, priv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sign bundle: %v\n", err)
			os.Exit(1)
		}
		emitJSON(bundle)
		return
	}

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

func runVerify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	_ = fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s verify [--json] <bundle.evidence.json|bundle.receipt.json>\n", os.Args[0])
		os.Exit(2)
	}

	data, err := os.ReadFile(fs.Arg(0)) // #nosec G304
	if err != nil {
		fmt.Fprintf(os.Stderr, "bundle: %v\n", err)
		os.Exit(1)
	}

	var peek struct {
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		fmt.Fprintf(os.Stderr, "bundle: %v\n", err)
		os.Exit(1)
	}

	var idProbe struct {
		VerificationID string `json:"verification_id"`
	}
	_ = json.Unmarshal(peek.Payload, &idProbe)

	if idProbe.VerificationID != "" {
		runVerifyReceiptBundle(data, *jsonOut)
		return
	}
	runVerifyEvidenceBundle(data, *jsonOut)
}

func runVerifyEvidenceBundle(data []byte, jsonOut bool) {
	var bundle schema.EvidenceBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		fmt.Fprintf(os.Stderr, "bundle: %v\n", err)
		os.Exit(1)
	}

	ok, msg := evidence.VerifyBundle(bundle)
	emitVerifyResult(ok, msg, bundle.Algorithm, bundle.PublicKey, jsonOut)
}

func runVerifyReceiptBundle(data []byte, jsonOut bool) {
	var bundle schema.VerificationReceiptBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		fmt.Fprintf(os.Stderr, "bundle: %v\n", err)
		os.Exit(1)
	}

	ok, msg := evidence.VerifyReceiptBundle(bundle)
	emitVerifyResult(ok, msg, bundle.Algorithm, bundle.PublicKey, jsonOut)
}

func emitVerifyResult(ok bool, msg, algorithm, publicKey string, jsonOut bool) {
	result := map[string]any{
		"valid":      ok,
		"message":    msg,
		"algorithm":  algorithm,
		"public_key": publicKey,
		"disclaimer": schema.Disclaimer,
	}

	if jsonOut {
		emitJSON(result)
	} else {
		fmt.Printf("Valid:     %v\n", ok)
		fmt.Printf("Message:   %s\n", msg)
		fmt.Printf("Algorithm: %s\n", algorithm)
		fmt.Printf("\nDisclaimer\n%s\n", schema.Disclaimer)
	}

	if !ok {
		os.Exit(1)
	}
}

func runRuntime(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s runtime <register|verify> ...\n", os.Args[0])
		os.Exit(2)
	}
	switch args[0] {
	case "register":
		runRuntimeRegister(args[1:])
	case "verify":
		runRuntimeVerify(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown runtime command: %s\n", args[0])
		os.Exit(2)
	}
}

func runRuntimeRegister(args []string) {
	fs := flag.NewFlagSet("runtime register", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	_ = fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s runtime register [--json] <register-request.json>\n", os.Args[0])
		os.Exit(2)
	}

	var req schema.RegisterRuntimeRequest
	if err := loadJSONFile(fs.Arg(0), &req); err != nil {
		fmt.Fprintf(os.Stderr, "request: %v\n", err)
		os.Exit(1)
	}

	resp, err := runtime.Register(req, time.Now().UTC())
	if err != nil {
		fmt.Fprintf(os.Stderr, "register: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		emitJSON(resp)
		return
	}

	fmt.Printf("System:      %s\n", resp.Baseline.SystemID)
	fmt.Printf("Deployment:  %s\n", resp.Baseline.DeploymentID)
	fmt.Printf("Fingerprint: %s\n", resp.Baseline.Fingerprint)
	fmt.Printf("Risk tier:   %s\n", resp.Policy.RiskTier)
	fmt.Printf("\nDisclaimer\n%s\n", resp.Disclaimer)
}

func runRuntimeVerify(args []string) {
	fs := flag.NewFlagSet("runtime verify", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit machine-readable JSON")
	signOut := fs.Bool("sign", false, "emit Ed25519-signed .receipt.json bundle")
	_ = fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s runtime verify [--json] [--sign] <verify-request.json>\n", os.Args[0])
		os.Exit(2)
	}

	var req schema.VerifyRuntimeRequest
	if err := loadJSONFile(fs.Arg(0), &req); err != nil {
		fmt.Fprintf(os.Stderr, "request: %v\n", err)
		os.Exit(1)
	}

	resp, err := runtime.Verify(req, time.Now().UTC())
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify: %v\n", err)
		os.Exit(1)
	}

	if *signOut {
		_, priv, err := evidence.GenerateKeyPair()
		if err != nil {
			fmt.Fprintf(os.Stderr, "sign: %v\n", err)
			os.Exit(1)
		}
		bundle, err := evidence.SignReceiptBundle(resp.Receipt, priv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "sign receipt: %v\n", err)
			os.Exit(1)
		}
		emitJSON(bundle)
		return
	}

	if *jsonOut {
		emitJSON(resp)
		return
	}

	fmt.Printf("Conformity state: %s\n", resp.Result.ConformityState)
	fmt.Printf("Violations:       %d\n", len(resp.Result.Violations))
	fmt.Printf("Receipt hash:     %s\n", resp.Receipt.ReceiptHash)
	fmt.Printf("\n%s\n", resp.Result.Summary)
	fmt.Printf("\nDisclaimer\n%s\n", schema.Disclaimer)
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
