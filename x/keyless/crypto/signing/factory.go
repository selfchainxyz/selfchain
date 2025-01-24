package signing

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	ecdsa_signer "selfchain/x/keyless/crypto/signing/ecdsa"
	"selfchain/x/keyless/crypto/signing/types"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

// SignerFactory is responsible for creating signing service instances
type SignerFactory struct{}

// NewSignerFactory creates a new signer factory
func NewSignerFactory() *SignerFactory {
	return &SignerFactory{}
}

// CreateSigner creates a new signing service based on the algorithm and key pair
func (f *SignerFactory) CreateSigner(ctx context.Context, algorithm types.SigningAlgorithm, privKeyBytes []byte, pubKeyBytes []byte) (types.SigningService, error) {
	switch algorithm {
	case types.ECDSA:
		return f.createECDSASigner(privKeyBytes, pubKeyBytes)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %v", algorithm)
	}
}

// createECDSASigner creates a new ECDSA signing service
func (f *SignerFactory) createECDSASigner(privKeyBytes []byte, pubKeyBytes []byte) (types.SigningService, error) {
	// Parse private key if provided
	var privKey *btcec.PrivateKey
	var pubKey *btcec.PublicKey
	var err error
	if len(privKeyBytes) > 0 {
		privKey, _ = btcec.PrivKeyFromBytes(privKeyBytes)
		if privKey == nil {
			return nil, fmt.Errorf("failed to parse private key")
		}
	}

	// Parse public key if provided
	if len(pubKeyBytes) > 0 {
		pubKey, err = btcec.ParsePubKey(pubKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}
	} else if privKey != nil {
		// If public key is not provided but private key is, derive public key from private key
		pubKey = privKey.PubKey()
	}

	if privKey == nil && pubKey == nil {
		return nil, errors.New("either private key or public key must be provided")
	}

	return ecdsa_signer.NewECDSASigner(privKey, pubKey), nil
}

// Sign signs a message using the specified algorithm
func (f *SignerFactory) Sign(ctx context.Context, message []byte, algorithm types.SigningAlgorithm, signer types.SigningService) (*types.SignatureResult, error) {
	if signer == nil {
		return nil, errors.New("signer is nil")
	}

	// Sign the message
	return signer.Sign(ctx, message, algorithm)
}

// Verify verifies a signature using the specified algorithm
func (f *SignerFactory) Verify(ctx context.Context, message []byte, signature *types.SignatureResult, pubKey []byte, algorithm types.SigningAlgorithm, signer types.SigningService) (bool, error) {
	if signer == nil {
		return false, errors.New("signer is nil")
	}

	// Verify the signature
	return signer.Verify(ctx, message, signature, pubKey)
}

// FormatSignature formats a signature according to the specified algorithm
func (f *SignerFactory) FormatSignature(ctx context.Context, sig *types.SignatureResult, algorithm types.SigningAlgorithm) ([]byte, error) {
	if sig == nil {
		return nil, errors.New("signature is nil")
	}

	// Return the DER encoded signature
	return sig.Bytes, nil
}

// UnformatSignature unformats a signature according to the specified algorithm
func (f *SignerFactory) UnformatSignature(ctx context.Context, sigBytes []byte, algorithm types.SigningAlgorithm) (*types.SignatureResult, error) {
	if len(sigBytes) == 0 {
		return nil, errors.New("signature bytes are empty")
	}

	// Parse the DER signature
	parsedSig, err := ecdsa.ParseDERSignature(sigBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature: %w", err)
	}

	// Get the serialized signature
	serialized := parsedSig.Serialize()

	// Convert to SignatureResult
	return &types.SignatureResult{
		R:     new(big.Int).SetBytes(serialized[:32]),
		S:     new(big.Int).SetBytes(serialized[32:]),
		Bytes: sigBytes,
	}, nil
}
