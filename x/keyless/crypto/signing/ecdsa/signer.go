package ecdsa

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

// ECDSASigner implements the types.SigningService interface for ECDSA signatures
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

	// Create a deterministic nonce
	nonce := sha256.Sum256(append(hash[:], s.privKey.Serialize()...))
	k := new(btcec.ModNScalar)
	k.SetBytes(&nonce)

	// Get the curve parameters
	curve := btcec.S256()
	n := curve.N

	// Calculate R = k*G
	kBytes := k.Bytes()
	x, _ := curve.ScalarBaseMult(kBytes[:])
	r := new(big.Int).Set(x)
	r.Mod(r, n)

	// Calculate s = k^-1(hash + r*privKey) mod n
	rInt := new(big.Int).Set(r)
	privInt := new(big.Int).SetBytes(s.privKey.Serialize())
	hashInt := new(big.Int).SetBytes(hash[:])

	// k^-1
	kBytes = k.Bytes()
	kInt := new(big.Int).SetBytes(kBytes[:])
	kInv := new(big.Int).ModInverse(kInt, n)
	if kInv == nil {
		return nil, errors.New("failed to calculate k inverse")
	}

	// r*privKey
	rPriv := new(big.Int).Mul(rInt, privInt)
	rPriv.Mod(rPriv, n)

	// hash + r*privKey
	sInt := new(big.Int).Add(hashInt, rPriv)
	sInt.Mod(sInt, n)

	// k^-1 * (hash + r*privKey)
	sInt.Mul(sInt, kInv)
	sInt.Mod(sInt, n)

	// Convert to ModNScalar for signature creation
	rScalar := new(btcec.ModNScalar)
	sScalar := new(btcec.ModNScalar)
	rScalar.SetByteSlice(r.Bytes())
	sScalar.SetByteSlice(sInt.Bytes())

	// Create ECDSA signature
	signature := ecdsa.NewSignature(rScalar, sScalar)

	// Get the serialized signature
	der := signature.Serialize()

	// Return signature result
	return &types.SignatureResult{
		R:     r,
		S:     sInt,
		Bytes: der,
	}, nil
}

// Verify verifies an ECDSA signature
func (s *ECDSASigner) Verify(ctx context.Context, message []byte, signature *types.SignatureResult, pubKeyBytes []byte) (bool, error) {
	if signature == nil {
		return false, errors.New("signature is nil")
	}

	// Parse public key
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Hash the message
	hash := sha256.Sum256(message)

	// Convert big.Int to ModNScalar
	rScalar := new(btcec.ModNScalar)
	sScalar := new(btcec.ModNScalar)
	rScalar.SetByteSlice(signature.R.Bytes())
	sScalar.SetByteSlice(signature.S.Bytes())

	// Create ECDSA signature
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

	return s.pubKey.SerializeCompressed(), nil
}
