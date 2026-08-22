package evidence

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/new-world-coder/riskline/pkg/schema"
)

// ReceiptContentHash returns SHA-256 over canonical receipt JSON excluding
// receipt_hash and signature metadata fields.
func ReceiptContentHash(receipt schema.VerificationReceipt) string {
	stub := receipt
	stub.ReceiptHash = ""
	stub.Signature = ""
	stub.PublicKey = ""
	stub.Algorithm = ""
	data, err := json.Marshal(stub)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func canonicalReceiptPayload(receipt schema.VerificationReceipt) ([]byte, error) {
	if receipt.ReceiptHash == "" {
		receipt.ReceiptHash = ReceiptContentHash(receipt)
	}
	stub := receipt
	stub.Signature = ""
	stub.PublicKey = ""
	stub.Algorithm = ""
	b, err := json.Marshal(stub)
	if err != nil {
		return nil, fmt.Errorf("marshal receipt payload: %w", err)
	}
	return b, nil
}

// SignReceiptBundle signs a VerificationReceipt with Ed25519.
// The signature covers canonical JSON of the receipt including receipt_hash.
func SignReceiptBundle(receipt schema.VerificationReceipt, priv ed25519.PrivateKey) (schema.VerificationReceiptBundle, error) {
	if len(priv) != ed25519.PrivateKeySize {
		return schema.VerificationReceiptBundle{}, fmt.Errorf("invalid ed25519 private key length: %d", len(priv))
	}

	if receipt.ReceiptHash == "" {
		receipt.ReceiptHash = ReceiptContentHash(receipt)
	}
	if ok, msg := VerifyReceiptHash(receipt); !ok {
		return schema.VerificationReceiptBundle{}, fmt.Errorf("receipt hash invalid: %s", msg)
	}

	payload, err := canonicalReceiptPayload(receipt)
	if err != nil {
		return schema.VerificationReceiptBundle{}, err
	}

	sig := ed25519.Sign(priv, payload)
	pub := priv.Public().(ed25519.PublicKey)

	signed := receipt
	signed.Signature = hex.EncodeToString(sig)
	signed.PublicKey = hex.EncodeToString(pub)
	signed.Algorithm = Algorithm

	return schema.VerificationReceiptBundle{
		Payload:    signed,
		Signature:  signed.Signature,
		PublicKey:  signed.PublicKey,
		Algorithm:  Algorithm,
		Disclaimer: schema.Disclaimer,
	}, nil
}

// VerifyReceiptBundle checks Ed25519 signature and receipt_hash integrity.
func VerifyReceiptBundle(bundle schema.VerificationReceiptBundle) (bool, string) {
	if bundle.Algorithm != Algorithm {
		return false, fmt.Sprintf("unsupported algorithm: %s", bundle.Algorithm)
	}

	if ok, msg := VerifyReceiptHash(bundle.Payload); !ok {
		return false, msg
	}

	pubBytes, err := hex.DecodeString(bundle.PublicKey)
	if err != nil {
		return false, fmt.Sprintf("invalid public_key hex: %v", err)
	}
	if len(pubBytes) != ed25519.PublicKeySize {
		return false, fmt.Sprintf("invalid public key length: %d", len(pubBytes))
	}

	sigBytes, err := hex.DecodeString(bundle.Signature)
	if err != nil {
		return false, fmt.Sprintf("invalid signature hex: %v", err)
	}

	payload, err := canonicalReceiptPayload(bundle.Payload)
	if err != nil {
		return false, err.Error()
	}

	if !ed25519.Verify(ed25519.PublicKey(pubBytes), payload, sigBytes) {
		return false, "signature verification failed"
	}
	return true, "signature valid"
}

// VerifyReceiptHash checks receipt_hash matches receipt content.
func VerifyReceiptHash(receipt schema.VerificationReceipt) (bool, string) {
	if receipt.ReceiptHash == "" {
		return false, "receipt_hash missing"
	}
	expected := ReceiptContentHash(receipt)
	if expected == "" {
		return false, "failed to compute receipt hash"
	}
	if receipt.ReceiptHash != expected {
		return false, fmt.Sprintf("receipt_hash mismatch: expected %s, got %s", expected, receipt.ReceiptHash)
	}
	return true, "receipt hash valid"
}

// VerifyReceiptChain checks that each receipt's previous_receipt_hash matches
// the prior receipt's receipt_hash (local integrity chain, not blockchain).
func VerifyReceiptChain(receipts []schema.VerificationReceipt) (bool, string) {
	var prev string
	for i, rec := range receipts {
		if ok, msg := VerifyReceiptHash(rec); !ok {
			return false, fmt.Sprintf("receipt %d: %s", i, msg)
		}
		if i > 0 && rec.PreviousReceiptHash != prev {
			return false, fmt.Sprintf("chain break at receipt %d: expected previous_receipt_hash %s, got %s", i, prev, rec.PreviousReceiptHash)
		}
		prev = rec.ReceiptHash
	}
	return true, "receipt chain intact"
}
