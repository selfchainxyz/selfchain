package eddsa

import (
	"crypto/ed25519"
	"crypto/sha512"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
)

// ConversionParameters stores parameters for ECDSA to EdDSA conversion
type ConversionParameters struct {
	ECDSAPublicKey  *btcec.PublicKey
	EdDSAPublicKey  ed25519.PublicKey
	ConversionProof []byte
}

// ConvertECDSAToEdDSA converts an ECDSA public key to EdDSA format
func ConvertECDSAToEdDSA(ecdsaPub *btcec.PublicKey) (*ConversionParameters, error) {
	if ecdsaPub == nil {
		return nil, fmt.Errorf("ECDSA public key is nil")
	}

	// Generate Ed25519 key from ECDSA public key components
	h := sha512.New()
	h.Write(ecdsaPub.SerializeCompressed())
	seed := h.Sum(nil)[:32]

	// Generate Ed25519 keypair from seed
	edPub := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)

	// Create conversion proof
	proof := generateConversionProof(ecdsaPub, edPub)

	return &ConversionParameters{
		ECDSAPublicKey:  ecdsaPub,
		EdDSAPublicKey:  edPub,
		ConversionProof: proof,
	}, nil
}

// ConvertSignatureToEdDSA converts an ECDSA signature to EdDSA format
func ConvertSignatureToEdDSA(r, s *big.Int) ([]byte, error) {
	if r == nil || s == nil {
		return nil, fmt.Errorf("signature components cannot be nil")
	}

	// EdDSA signature is 64 bytes: R (32 bytes) || S (32 bytes)
	signature := make([]byte, ed25519.SignatureSize)

	// Convert R to fixed 32 bytes
	rBytes := r.Bytes()
	if len(rBytes) > 32 {
		return nil, fmt.Errorf("R component too large")
	}
	copy(signature[32-len(rBytes):32], rBytes)

	// Convert S to fixed 32 bytes
	sBytes := s.Bytes()
	if len(sBytes) > 32 {
		return nil, fmt.Errorf("S component too large")
	}
	copy(signature[64-len(sBytes):], sBytes)

	return signature, nil
}

// VerifyEdDSASignature verifies an EdDSA signature
func VerifyEdDSASignature(pubKey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pubKey, message, signature)
}

// generateConversionProof creates a proof of correct conversion
// This is a placeholder - in production, implement a zero-knowledge proof
func generateConversionProof(ecdsaPub *btcec.PublicKey, eddsaPub ed25519.PublicKey) []byte {
	h := sha512.New()
	h.Write(ecdsaPub.SerializeCompressed())
	h.Write(eddsaPub)
	return h.Sum(nil)
}
