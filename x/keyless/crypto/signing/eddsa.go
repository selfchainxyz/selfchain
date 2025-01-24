package signing

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
)

// EdDSASigner implements the SigningService interface for EdDSA signatures
type EdDSASigner struct {
	party1Data interface{}
	party2Data interface{}
	privKey    ed25519.PrivateKey // For testing only
	pubKey     ed25519.PublicKey  // For testing only
}

// NewEdDSASigner creates a new EdDSA signer with party data
func NewEdDSASigner(party1Data, party2Data interface{}) *EdDSASigner {
	// For testing, generate a dummy key pair
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	return &EdDSASigner{
		party1Data: party1Data,
		party2Data: party2Data,
		privKey:    priv,
		pubKey:     pub,
	}
}

// Sign creates an EdDSA signature for the given message
func (s *EdDSASigner) Sign(ctx context.Context, message []byte, algorithm SigningAlgorithm) (*SignatureResult, error) {
	if algorithm != EdDSA {
		return nil, errors.New("unsupported algorithm")
	}

	// Sign the message (for testing only)
	sig := ed25519.Sign(s.privKey, message)

	// For EdDSA, we store the raw signature bytes
	return &SignatureResult{
		Bytes: sig,
	}, nil
}

// Verify verifies an EdDSA signature
func (s *EdDSASigner) Verify(ctx context.Context, message []byte, signature *SignatureResult, pubKeyBytes []byte) (bool, error) {
	if signature == nil {
		return false, errors.New("signature is nil")
	}

	if len(pubKeyBytes) != ed25519.PublicKeySize {
		return false, fmt.Errorf("invalid public key size: expected %d, got %d", ed25519.PublicKeySize, len(pubKeyBytes))
	}

	// Convert bytes to ed25519.PublicKey
	pubKey := ed25519.PublicKey(pubKeyBytes)

	// Verify the signature
	return ed25519.Verify(pubKey, message, signature.Bytes), nil
}

// GetPublicKey returns the EdDSA public key
func (s *EdDSASigner) GetPublicKey(ctx context.Context, algorithm SigningAlgorithm) ([]byte, error) {
	if algorithm != EdDSA {
		return nil, errors.New("unsupported algorithm")
	}

	if s.pubKey == nil {
		return nil, errors.New("public key not initialized")
	}

	// Return the raw public key bytes
	return []byte(s.pubKey), nil
}
