package evidence

import (
	"encoding/hex"
	"testing"

	"github.com/new-world-coder/riskline/pkg/schema"
)

func TestSignAndVerifyBundle(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	resp := schema.AssureResponse{
		ConformityState: schema.ConformityGreen,
		Summary:         "All controls verified.",
		Disclaimer:      schema.Disclaimer,
	}

	bundle, err := SignBundle(resp, priv)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Algorithm != Algorithm {
		t.Fatalf("expected algorithm %s, got %s", Algorithm, bundle.Algorithm)
	}

	ok, msg := VerifyBundle(bundle)
	if !ok {
		t.Fatalf("expected valid signature, got: %s", msg)
	}
}

func TestVerifyBundleTamperedPayload(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	resp := schema.AssureResponse{
		ConformityState: schema.ConformityGreen,
		Summary:         "All controls verified.",
		Disclaimer:      schema.Disclaimer,
	}

	bundle, err := SignBundle(resp, priv)
	if err != nil {
		t.Fatal(err)
	}

	bundle.Payload.ConformityState = schema.ConformityRed
	ok, _ := VerifyBundle(bundle)
	if ok {
		t.Fatal("expected tampered payload to fail verification")
	}
}

func TestVerifyBundleWrongAlgorithm(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	bundle, err := SignBundle(schema.AssureResponse{Disclaimer: schema.Disclaimer}, priv)
	if err != nil {
		t.Fatal(err)
	}
	bundle.Algorithm = "rsa"

	ok, msg := VerifyBundle(bundle)
	if ok {
		t.Fatal("expected unsupported algorithm to fail")
	}
	if msg == "" {
		t.Fatal("expected error message")
	}
}

func TestPrivateKeyFromHex(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	loaded, err := PrivateKeyFromHex(hex.EncodeToString(priv))
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != len(priv) {
		t.Fatalf("loaded key length mismatch: %d vs %d", len(loaded), len(priv))
	}
	for i := range priv {
		if loaded[i] != priv[i] {
			t.Fatal("loaded private key does not match original")
		}
	}
}
