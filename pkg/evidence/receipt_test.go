package evidence

import (
	"testing"
	"time"

	"github.com/new-world-coder/riskline/pkg/schema"
)

func sampleReceipt() schema.VerificationReceipt {
	return schema.VerificationReceipt{
		VerificationID:      "abc123",
		SystemID:            "hiring-assist-prod",
		VerifiedAt:          time.Date(2026, 8, 22, 12, 5, 0, 0, time.UTC),
		ObservedAt:          time.Date(2026, 8, 22, 12, 5, 0, 0, time.UTC),
		RuntimeFingerprint:  "runtime-fp",
		BaselineFingerprint: "baseline-fp",
		PolicyVersion:       "v1",
		ConformityState:     schema.ConformityGreen,
	}
}

func TestSignAndVerifyReceiptBundle(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	receipt := sampleReceipt()
	receipt.ReceiptHash = ReceiptContentHash(receipt)

	bundle, err := SignReceiptBundle(receipt, priv)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Algorithm != Algorithm {
		t.Fatalf("expected algorithm %s, got %s", Algorithm, bundle.Algorithm)
	}
	if bundle.Payload.Signature == "" {
		t.Fatal("expected signed payload")
	}

	ok, msg := VerifyReceiptBundle(bundle)
	if !ok {
		t.Fatalf("expected valid signature, got: %s", msg)
	}
}

func TestVerifyReceiptBundleTamperedPayload(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	receipt := sampleReceipt()
	receipt.ReceiptHash = ReceiptContentHash(receipt)

	bundle, err := SignReceiptBundle(receipt, priv)
	if err != nil {
		t.Fatal(err)
	}

	bundle.Payload.ConformityState = schema.ConformityRed
	ok, _ := VerifyReceiptBundle(bundle)
	if ok {
		t.Fatal("expected tampered payload to fail verification")
	}
}

func TestVerifyReceiptHash(t *testing.T) {
	receipt := sampleReceipt()
	receipt.ReceiptHash = ReceiptContentHash(receipt)

	ok, msg := VerifyReceiptHash(receipt)
	if !ok {
		t.Fatalf("expected valid hash, got: %s", msg)
	}

	receipt.ConformityState = schema.ConformityRed
	ok, _ = VerifyReceiptHash(receipt)
	if ok {
		t.Fatal("expected hash mismatch after tamper")
	}
}

func TestVerifyReceiptChain(t *testing.T) {
	r1 := sampleReceipt()
	r1.ReceiptHash = ReceiptContentHash(r1)

	r2 := sampleReceipt()
	r2.VerificationID = "def456"
	r2.ConformityState = schema.ConformityAmber
	r2.PreviousReceiptHash = r1.ReceiptHash
	r2.ReceiptHash = ReceiptContentHash(r2)

	ok, msg := VerifyReceiptChain([]schema.VerificationReceipt{r1, r2})
	if !ok {
		t.Fatalf("expected intact chain, got: %s", msg)
	}

	r2.PreviousReceiptHash = "broken"
	ok, _ = VerifyReceiptChain([]schema.VerificationReceipt{r1, r2})
	if ok {
		t.Fatal("expected broken chain to fail")
	}
}

func TestSignReceiptBundleRejectsBadHash(t *testing.T) {
	_, priv, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	receipt := sampleReceipt()
	receipt.ReceiptHash = "deadbeef"

	_, err = SignReceiptBundle(receipt, priv)
	if err == nil {
		t.Fatal("expected error for invalid receipt hash")
	}
}
