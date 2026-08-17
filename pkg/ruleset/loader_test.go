package ruleset

import (
	"testing"
)

func TestDefaultLoader(t *testing.T) {
	l, err := DefaultLoader()
	if err != nil {
		t.Fatal(err)
	}
	list := l.List()
	if len(list) != 2 {
		t.Fatalf("list = %v, want eu-ai-act and nist-ai-rmf", list)
	}
	p, err := l.Load("EU-AI-Act")
	if err != nil {
		t.Fatal(err)
	}
	if p.Character != CharacterHardLaw {
		t.Fatalf("character = %q", p.Character)
	}
	if p.Set.Version != "eu-ai-act-2024-v0.1.0" {
		t.Fatalf("version = %q", p.Set.Version)
	}
	nist, err := l.Load(RegimeNISTAIRMF)
	if err != nil {
		t.Fatal(err)
	}
	if nist.Character != CharacterMapping {
		t.Fatalf("nist character = %q", nist.Character)
	}
	if nist.Set.Version != "nist-ai-rmf-2023-v0.1.0" {
		t.Fatalf("nist version = %q", nist.Set.Version)
	}
}

func TestLoadUnknownRegime(t *testing.T) {
	l, err := DefaultLoader()
	if err != nil {
		t.Fatal(err)
	}
	_, err = l.Load("mas-feat")
	if err == nil {
		t.Fatal("expected unknown regime error")
	}
}

func TestNormalizeRegimes(t *testing.T) {
	got := NormalizeRegimes([]string{" eu-ai-act ", "EU-AI-ACT", "", "eu-ai-act"})
	if len(got) != 1 || got[0] != RegimeEUAIAct {
		t.Fatalf("got %v", got)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	l, err := DefaultLoader()
	if err != nil {
		t.Fatal(err)
	}
	set, err := LoadDefault()
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Register(&Pack{ID: RegimeEUAIAct, Character: CharacterHardLaw, Set: set}); err == nil {
		t.Fatal("expected duplicate error")
	}
}
