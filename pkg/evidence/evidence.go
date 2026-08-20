package evidence

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/new-world-coder/riskline/pkg/schema"
)

const Algorithm = "ed25519"

// GenerateKeyPair returns a new Ed25519 keypair for local dev/test.
// Production users should bring their own keys.
func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate ed25519 keypair: %w", err)
	}
	return pub, priv, nil
}

// SignBundle signs an AssureResponse with Ed25519 and returns a portable bundle.
// The signature covers canonical JSON of the payload only. This is local signing
// for offline verification — not a public blockchain ledger.
func SignBundle(resp schema.AssureResponse, priv ed25519.PrivateKey) (schema.EvidenceBundle, error) {
	if len(priv) != ed25519.PrivateKeySize {
		return schema.EvidenceBundle{}, fmt.Errorf("invalid ed25519 private key length: %d", len(priv))
	}

	payload, err := canonicalPayload(resp)
	if err != nil {
		return schema.EvidenceBundle{}, err
	}

	sig := ed25519.Sign(priv, payload)
	pub := priv.Public().(ed25519.PublicKey)

	return schema.EvidenceBundle{
		Payload:   resp,
		Signature: hex.EncodeToString(sig),
		PublicKey: hex.EncodeToString(pub),
		Algorithm: Algorithm,
	}, nil
}

// VerifyBundle checks an Ed25519 signature over the bundle payload (local signing, not blockchain).
func VerifyBundle(bundle schema.EvidenceBundle) (bool, string) {
	if bundle.Algorithm != Algorithm {
		return false, fmt.Sprintf("unsupported algorithm: %s", bundle.Algorithm)
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

	payload, err := canonicalPayload(bundle.Payload)
	if err != nil {
		return false, fmt.Sprintf("canonical payload: %v", err)
	}

	if !ed25519.Verify(ed25519.PublicKey(pubBytes), payload, sigBytes) {
		return false, "signature verification failed"
	}
	return true, "signature valid"
}

func canonicalPayload(resp schema.AssureResponse) ([]byte, error) {
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	return b, nil
}

// PrivateKeyFromHex loads an Ed25519 private key from hex (64-byte seed+pub form).
func PrivateKeyFromHex(hexKey string) (ed25519.PrivateKey, error) {
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}
	if len(b) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("expected %d-byte private key, got %d", ed25519.PrivateKeySize, len(b))
	}
	return ed25519.PrivateKey(b), nil
}
