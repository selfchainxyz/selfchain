package signing

import (
	"context"
	"crypto/ecdsa"

	"fmt"
	"math/big"
	"sync"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/btcsuite/btcd/btcec/v2"

	"selfchain/x/keyless/networks"
	"selfchain/x/keyless/types"
)

// SigningContext contains the context for signing operations
type SigningContext struct {
	NetworkParams *types.NetworkParams
	Message       []byte
	Party1Data    interface{}
	Party2Data    interface{}
	Metadata      map[string]interface{} // Network-specific metadata
}

// SignRequest contains the parameters for a signing request
type SignRequest struct {
	NetworkID string
	Params    *types.NetworkParams
	Message   []byte
	MetaData  map[string]interface{}
}

// SignerFactory manages signer instances
type SignerFactory struct {
	mu       sync.RWMutex
	signer   *UniversalSigner
	registry *networks.NetworkRegistry
}

// NewSignerFactory creates a new signer factory
func NewSignerFactory(registry *networks.NetworkRegistry) *SignerFactory {
	return &SignerFactory{
		signer:   NewUniversalSigner(nil, nil),
		registry: registry,
	}
}

// Sign signs a message for the specified network
func (f *SignerFactory) Sign(ctx context.Context, networkID string, message []byte, metadata map[string]interface{}, signResult *SignatureResult) ([]byte, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Get network configuration
	networkType, _, err := networks.ParseNetworkID(networkID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network ID: %w", err)
	}

	// Format signature based on network type
	var signature []byte
	switch networkType {
	case networks.Bitcoin:
		signature, err = formatBitcoinSignature(signResult)
	case networks.Ethereum:
		signature, err = formatEthereumSignature(signResult, "")
	case networks.Cosmos:
		signature, err = formatCosmosSignature(signResult)
	default:
		return nil, fmt.Errorf("unsupported network type: %s", networkType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to format signature: %w", err)
	}

	return signature, nil
}

// Verify verifies a signature for the specified network
func (f *SignerFactory) Verify(networkID string, pubKey []byte, message []byte, signature []byte) (bool, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Get network configuration
	networkType, _, err := networks.ParseNetworkID(networkID)
	if err != nil {
		return false, fmt.Errorf("failed to parse network ID: %w", err)
	}

	// Parse the compressed public key
	pubKeyObj, err := btcec.ParsePubKey(pubKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Convert to ECDSA public key
	ecdsaPubKey := &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     pubKeyObj.X(),
		Y:     pubKeyObj.Y(),
	}

	// Handle network-specific signature verification
	switch networkType {
	case networks.Bitcoin:
		r, s, err := UnmarshalDERSignature(signature)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal DER signature: %w", err)
		}
		return ecdsa.Verify(ecdsaPubKey, message, r, s), nil

	case networks.Ethereum:
		if len(signature) != 65 {
			return false, fmt.Errorf("invalid ethereum signature length")
		}
		r := new(big.Int).SetBytes(signature[:32])
		s := new(big.Int).SetBytes(signature[32:64])
		return ecdsa.Verify(ecdsaPubKey, message, r, s), nil

	case networks.Cosmos:
		if len(signature) != 64 {
			return false, fmt.Errorf("invalid cosmos signature length")
		}
		r := new(big.Int).SetBytes(signature[:32])
		s := new(big.Int).SetBytes(signature[32:])
		return ecdsa.Verify(ecdsaPubKey, message, r, s), nil

	default:
		return false, fmt.Errorf("unsupported network type: %s", networkType)
	}
}

// UnmarshalDERSignature parses a DER encoded signature
func UnmarshalDERSignature(sig []byte) (r, s *big.Int, err error) {
	// DER format: 0x30 [total-length] 0x02 [R-length] [R] 0x02 [S-length] [S]
	if len(sig) < 8 || sig[0] != 0x30 {
		return nil, nil, fmt.Errorf("invalid DER signature format")
	}

	// Get R component
	rLen := int(sig[3])
	if len(sig) < 4+rLen {
		return nil, nil, fmt.Errorf("invalid R length")
	}
	r = new(big.Int).SetBytes(sig[4 : 4+rLen])

	// Get S component
	sOffset := 4 + rLen + 2
	if len(sig) < sOffset {
		return nil, nil, fmt.Errorf("invalid S offset")
	}
	sLen := int(sig[4+rLen+1])
	if len(sig) < sOffset+sLen {
		return nil, nil, fmt.Errorf("invalid S length")
	}
	s = new(big.Int).SetBytes(sig[sOffset : sOffset+sLen])

	return r, s, nil
}

// Signer interface defines the methods required for network-specific signing
type Signer interface {
	Sign(ctx context.Context, message []byte, algorithm SigningAlgorithm) (*SignatureResult, error)
	Verify(ctx context.Context, message []byte, signature *SignatureResult, pubKey *btcec.PublicKey) (bool, error)
}

// UniversalSigner implements Signer for all networks
type UniversalSigner struct {
	networkParams *types.NetworkParams
	ecdsaSigner   *ECDSASigner
	eddsaSigner   *EdDSASigner
}

// NewUniversalSigner creates a new universal signer
func NewUniversalSigner(party1Data, party2Data *keygen.LocalPartySaveData) *UniversalSigner {
	return &UniversalSigner{
		ecdsaSigner: NewECDSASigner(party1Data, party2Data),
		eddsaSigner: NewEdDSASigner(party1Data, party2Data),
	}
}

// Sign signs a message using the appropriate algorithm
func (s *UniversalSigner) Sign(ctx context.Context, signingCtx *SigningContext) (*SignatureResult, error) {
	switch signingCtx.NetworkParams.SigningAlgorithm {
	case string(networks.ECDSA):
		return s.ecdsaSigner.Sign(ctx, signingCtx.Message, ECDSA)
	case string(networks.EdDSA):
		return s.eddsaSigner.Sign(ctx, signingCtx.Message, EdDSA)
	default:
		return nil, fmt.Errorf("unsupported signing algorithm: %s", signingCtx.NetworkParams.SigningAlgorithm)
	}
}

// Verify verifies a signature using the appropriate algorithm
func (s *UniversalSigner) Verify(ctx context.Context, pubKey *btcec.PublicKey, msg []byte, signature *SignatureResult) (bool, error) {
	// Try ECDSA first
	if valid, err := s.ecdsaSigner.Verify(ctx, msg, signature, pubKey); err == nil {
		return valid, nil
	}

	// Try EdDSA if ECDSA fails
	return s.eddsaSigner.Verify(ctx, msg, signature, pubKey)
}

// Helper functions for signature formatting
func formatBitcoinSignature(result *SignatureResult) ([]byte, error) {
	// Convert R and S to bytes
	rBytes := padTo32(result.R.Bytes())
	sBytes := padTo32(result.S.Bytes())

	// DER format: 0x30 [total-length] 0x02 [R-length] [R] 0x02 [S-length] [S]
	rLen := len(rBytes)
	sLen := len(sBytes)
	totalLen := 2 + rLen + 2 + sLen

	signature := make([]byte, 2+totalLen)
	signature[0] = 0x30 // Sequence tag
	signature[1] = byte(totalLen)
	signature[2] = 0x02 // Integer tag for R
	signature[3] = byte(rLen)
	copy(signature[4:], rBytes)
	signature[4+rLen] = 0x02 // Integer tag for S
	signature[4+rLen+1] = byte(sLen)
	copy(signature[4+rLen+2:], sBytes)

	return signature, nil
}

func formatEthereumSignature(result *SignatureResult, chainID string) ([]byte, error) {
	// Convert R and S to 32-byte arrays
	rBytes := padTo32(result.R.Bytes())
	sBytes := padTo32(result.S.Bytes())

	// Calculate V value (27 or 28 for legacy transactions)
	v := byte(27)

	// Combine R, S, and V
	signature := make([]byte, 65)
	copy(signature[:32], rBytes)
	copy(signature[32:64], sBytes)
	signature[64] = v

	return signature, nil
}

func formatCosmosSignature(result *SignatureResult) ([]byte, error) {
	// Convert R and S to 32-byte arrays
	rBytes := padTo32(result.R.Bytes())
	sBytes := padTo32(result.S.Bytes())

	// Combine R and S
	signature := append(rBytes, sBytes...)
	return signature, nil
}

// Helper function to pad a byte slice to 32 bytes
func padTo32(b []byte) []byte {
	if len(b) > 32 {
		return b[:32]
	}
	if len(b) == 32 {
		return b
	}
	result := make([]byte, 32)
	copy(result[32-len(b):], b)
	return result
}
