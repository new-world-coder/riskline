package ruleset

import (
	"fmt"
	"sort"
	"strings"
)

// Canonical regime identifiers. Pack JSON version strings stay independent
// (e.g. eu-ai-act → eu-ai-act-2024-v0.1.0) so releases can bump pack content
// without renaming the regime selector.
const (
	RegimeEUAIAct = "eu-ai-act"
)

// Character describes the legal nature of a regime pack's outcomes.
const (
	CharacterHardLaw    = "hard_law"
	CharacterPrinciples = "principles"
	CharacterGuidance   = "guidance"
	CharacterMapping    = "mapping"
)

// Pack is a named, embedded ruleset ready for classification.
type Pack struct {
	ID        string
	Character string
	Set       *Set
}

// Loader resolves regime IDs to embedded packs.
type Loader struct {
	packs map[string]*Pack
}

// DefaultLoader returns a loader with every shipped pack registered.
// P0 ships eu-ai-act only; additional regimes are additive later.
func DefaultLoader() (*Loader, error) {
	l := &Loader{packs: make(map[string]*Pack)}
	set, err := LoadDefault()
	if err != nil {
		return nil, err
	}
	if err := l.Register(&Pack{
		ID:        RegimeEUAIAct,
		Character: CharacterHardLaw,
		Set:       set,
	}); err != nil {
		return nil, err
	}
	return l, nil
}

// Register adds a pack. Duplicate IDs are rejected.
func (l *Loader) Register(p *Pack) error {
	if l == nil {
		return fmt.Errorf("ruleset: nil loader")
	}
	if p == nil || p.ID == "" || p.Set == nil {
		return fmt.Errorf("ruleset: pack requires id and set")
	}
	id := NormalizeRegime(p.ID)
	if _, exists := l.packs[id]; exists {
		return fmt.Errorf("ruleset: duplicate regime %q", id)
	}
	cp := *p
	cp.ID = id
	l.packs[id] = &cp
	return nil
}

// Load returns the pack for a regime ID.
func (l *Loader) Load(regime string) (*Pack, error) {
	if l == nil {
		return nil, fmt.Errorf("ruleset: nil loader")
	}
	id := NormalizeRegime(regime)
	if id == "" {
		return nil, fmt.Errorf("ruleset: empty regime")
	}
	p, ok := l.packs[id]
	if !ok {
		return nil, fmt.Errorf("ruleset: unknown regime %q (available: %s)", id, strings.Join(l.List(), ", "))
	}
	return p, nil
}

// List returns sorted regime IDs.
func (l *Loader) List() []string {
	if l == nil {
		return nil
	}
	out := make([]string, 0, len(l.packs))
	for id := range l.packs {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// NormalizeRegime trims and lowercases a regime selector.
func NormalizeRegime(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// NormalizeRegimes deduplicates and normalizes a regime list, preserving order.
func NormalizeRegimes(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, r := range in {
		id := NormalizeRegime(r)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}
