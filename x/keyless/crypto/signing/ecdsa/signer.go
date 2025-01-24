package ecdsa

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"selfchain/x/keyless/crypto/signing/types"
)

// SigningService defines the interface for signing operations
type SigningService interface {
	Sign(ctx context.Context, message []byte, algorithm types.SigningAlgorithm) (*types.SignatureResult, error)
	Verify(ctx context.Context, message []byte, signature *types.SignatureResult, pubKeyBytes []byte) (bool, error)
	GetPublicKey(ctx context.Context, algorithm types.SigningAlgorithm) ([]byte, error)
}

// ECDSASigner implements the SigningService interface for ECDSA signatures
type ECDSASigner struct {
	party1Data interface{}
	party2Data interface{}
	privKey    *btcec.PrivateKey // For testing only
	pubKey     *btcec.PublicKey  // For testing only
}

// NewECDSASigner creates a new ECDSA signer with party data
func NewECDSASigner(party1Data, party2Data interface{}) *ECDSASigner {
	// For testing, generate a dummy key pair
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil
	}
	return &ECDSASigner{
		party1Data: party1Data,
		party2Data: party2Data,
		privKey:    privKey,
		pubKey:     privKey.PubKey(),
	}
}

// Sign creates an ECDSA signature for the given message
func (s *ECDSASigner) Sign(ctx context.Context, message []byte, algorithm types.SigningAlgorithm) (*types.SignatureResult, error) {
	if algorithm != types.ECDSA {
		return nil, errors.New("unsupported algorithm")
	}

	// Hash the message
	hash := sha256.Sum256(message)

	// Sign the hash (for testing only)
	signature := ecdsa.Sign(s.privKey, hash[:])

	// Get signature components
	rBytes := signature.Serialize()[4:36]  // Skip DER header and length
	sBytes := signature.Serialize()[38:70] // Skip DER header and length

	rInt := new(big.Int).SetBytes(rBytes)
	sInt := new(big.Int).SetBytes(sBytes)

	return &types.SignatureResult{
		R:     rInt,
		S:     sInt,
		Bytes: signature.Serialize(),
	}, nil
}

// Verify verifies an ECDSA signature
func (s *ECDSASigner) Verify(ctx context.Context, message []byte, signature *types.SignatureResult, pubKeyBytes []byte) (bool, error) {
	if signature == nil {
		return false, errors.New("signature is nil")
	}

	// Hash the message
	hash := sha256.Sum256(message)

	// Parse public key based on format
	var pubKey *btcec.PublicKey
	var err error

	switch len(pubKeyBytes) {
	case 33: // Compressed public key
		pubKey, err = btcec.ParsePubKey(pubKeyBytes)
	case 65: // Uncompressed public key
		pubKey, err = btcec.ParsePubKey(pubKeyBytes)
	case 64: // Raw public key (x||y)
		// Add the uncompressed point marker
		fullPubKey := make([]byte, 65)
		fullPubKey[0] = 0x04
		copy(fullPubKey[1:], pubKeyBytes)
		pubKey, err = btcec.ParsePubKey(fullPubKey)
	default:
		return false, fmt.Errorf("invalid public key length: %d", len(pubKeyBytes))
	}

	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Parse signature
	sig, err := ecdsa.ParseDERSignature(signature.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse signature: %w", err)
	}

	// Verify the signature
	return sig.Verify(hash[:], pubKey), nil
}

// GetPublicKey returns the ECDSA public key
func (s *ECDSASigner) GetPublicKey(ctx context.Context, algorithm types.SigningAlgorithm) ([]byte, error) {
	if algorithm != types.ECDSA {
		return nil, errors.New("unsupported algorithm")
	}

	if s.pubKey == nil {
		return nil, errors.New("public key not initialized")
	}

	// Return the serialized public key
	return s.pubKey.SerializeCompressed(), nil
}
