package signing

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"selfchain/x/keyless/crypto/signing/types"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

// ECDSASigner implements the SigningService interface for ECDSA signatures
type ECDSASigner struct {
	privKey *btcec.PrivateKey
	pubKey  *btcec.PublicKey
}

// NewECDSASigner creates a new ECDSA signer
func NewECDSASigner(privKey *btcec.PrivateKey, pubKey *btcec.PublicKey) *ECDSASigner {
	return &ECDSASigner{
		privKey: privKey,
		pubKey:  pubKey,
	}
}

// Sign creates an ECDSA signature for the given message
func (s *ECDSASigner) Sign(ctx context.Context, message []byte, algorithm types.SigningAlgorithm) (*types.SignatureResult, error) {
	if algorithm != types.ECDSA {
		return nil, errors.New("unsupported algorithm")
	}

	if s.privKey == nil {
		return nil, errors.New("private key not initialized")
	}

	// Hash the message
	hash := sha256.Sum256(message)

	// Sign the hash
	signature := ecdsa.Sign(s.privKey, hash[:])

	// Get the serialized signature
	der := signature.Serialize()

	// Get the signature components
	rBytes := signature.Serialize()[:32]  // First 32 bytes are R
	sBytes := signature.Serialize()[32:64] // Next 32 bytes are S
	
	rBig := new(big.Int).SetBytes(rBytes)
	sBig := new(big.Int).SetBytes(sBytes)

	return &types.SignatureResult{
		R:     rBig,
		S:     sBig,
		Bytes: der,
	}, nil
}

// Verify verifies an ECDSA signature
func (e *ECDSASigner) Verify(ctx context.Context, message []byte, signature *types.SignatureResult, pubKeyBytes []byte) (bool, error) {
	// Parse public key
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Hash the message
	hash := sha256.Sum256(message)

	// Create signature from R and S
	rScalar := new(btcec.ModNScalar)
	sScalar := new(btcec.ModNScalar)
	rScalar.SetByteSlice(signature.R.Bytes())
	sScalar.SetByteSlice(signature.S.Bytes())
	sig := ecdsa.NewSignature(rScalar, sScalar)

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

	pubKey := s.privKey.PubKey()
	return pubKey.SerializeCompressed(), nil
}

// FormatDERSignature formats a signature in DER format
func FormatDERSignature(sig *types.SignatureResult) ([]byte, error) {
	if sig == nil {
		return nil, errors.New("signature is nil")
	}

	// Create signature from R and S
	rScalar := new(btcec.ModNScalar)
	sScalar := new(btcec.ModNScalar)
	rScalar.SetByteSlice(sig.R.Bytes())
	sScalar.SetByteSlice(sig.S.Bytes())
	signature := ecdsa.NewSignature(rScalar, sScalar)

	// Format signature in DER format
	return signature.Serialize(), nil
}

// UnmarshalDERSignature parses a DER signature into R and S components
func UnmarshalDERSignature(sigBytes []byte) (*big.Int, *big.Int, error) {
	sig, err := ecdsa.ParseDERSignature(sigBytes)
	if err != nil {
		return nil, nil, err
	}

	// Get the signature components
	sigData := sig.Serialize()
	rBytes := sigData[:32]  // First 32 bytes are R
	sBytes := sigData[32:64] // Next 32 bytes are S
	
	rBig := new(big.Int).SetBytes(rBytes)
	sBig := new(big.Int).SetBytes(sBytes)

	return rBig, sBig, nil
}
