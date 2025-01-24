package signing

import (
	"context"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

// ECDSASigner implements the SigningService interface for ECDSA signatures
type ECDSASigner struct {
	party1Data interface{}
	party2Data interface{}
	signer     *ecdsaSigner
}

// NewECDSASigner creates a new ECDSA signer with party data
func NewECDSASigner(party1Data, party2Data interface{}) *ECDSASigner {
	return &ECDSASigner{
		party1Data: party1Data,
		party2Data: party2Data,
		signer:     &ecdsaSigner{},
	}
}

// Sign creates an ECDSA signature for the given message
func (s *ECDSASigner) Sign(ctx context.Context, message []byte, algorithm SigningAlgorithm) (*SignatureResult, error) {
	if algorithm != ECDSA {
		return nil, errors.New("unsupported algorithm")
	}

	return s.signer.Sign(ctx, message)
}

// Verify verifies an ECDSA signature
func (s *ECDSASigner) Verify(ctx context.Context, message []byte, signature *SignatureResult, pubKeyBytes []byte) (bool, error) {
	if signature == nil {
		return false, errors.New("signature is nil")
	}

	return s.signer.Verify(ctx, message, signature, pubKeyBytes)
}

// GetPublicKey returns the ECDSA public key
func (s *ECDSASigner) GetPublicKey(ctx context.Context, algorithm SigningAlgorithm) ([]byte, error) {
	if algorithm != ECDSA {
		return nil, errors.New("unsupported algorithm")
	}

	return s.signer.GetPublicKey(ctx)
}

// FormatBitcoinSignature formats a signature in Bitcoin's DER format
func FormatBitcoinSignature(sig *SignatureResult) ([]byte, error) {
	if sig == nil {
		return nil, errors.New("signature is nil")
	}

	// Create DER signature
	der, err := FormatDERSignature(sig)
	if err != nil {
		return nil, err
	}

	// Append SIGHASH_ALL
	return append(der, 0x01), nil
}

// FormatDERSignature formats a signature in DER format
func FormatDERSignature(sig *SignatureResult) ([]byte, error) {
	return ecdsa.FormatDERSignature(sig)
}

// UnmarshalDERSignature parses a DER signature into R and S components
func UnmarshalDERSignature(sig []byte) (*big.Int, *big.Int, error) {
	return ecdsa.UnmarshalDERSignature(sig)
}
