package signing

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"
	"sync"
	"bytes"

	"github.com/btcsuite/btcd/btcec/v2"

	"selfchain/x/keyless/networks"
	"selfchain/x/keyless/types"
	"selfchain/x/keyless/tss"
)

// SigningContext contains the context for signing operations
type SigningContext struct {
	NetworkParams *types.NetworkParams
	Message     []byte
	Party1Data  interface{}
	Party2Data  interface{}
	Metadata    map[string]interface{} // Network-specific metadata
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
		signer:   &UniversalSigner{},
		registry: registry,
	}
}

// Sign signs a message for the specified network
func (f *SignerFactory) Sign(ctx context.Context, networkID string, message []byte, metadata map[string]interface{}, signResult *tss.SignResult) ([]byte, error) {
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

	// Convert public key to ECDSA
	var curve elliptic.Curve
	switch networkType {
	case networks.Bitcoin:
		curve = btcec.S256()
	case networks.Ethereum:
		curve = btcec.S256()
	case networks.Cosmos:
		curve = btcec.S256()
	default:
		return false, fmt.Errorf("unsupported network type: %s", networkType)
	}

	x, y := elliptic.Unmarshal(curve, pubKey)
	if x == nil {
		return false, fmt.Errorf("failed to unmarshal public key")
	}
	ecdsaPubKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	// Handle network-specific verification
	switch networkType {
	case networks.Bitcoin:
		// For Bitcoin, signature is in DER format
		r, s, err := UnmarshalDERSignature(signature)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal DER signature: %w", err)
		}
		return ecdsa.Verify(ecdsaPubKey, message, r, s), nil

	case networks.Ethereum:
		// For Ethereum, signature is R || S || V
		if len(signature) != 65 {
			return false, fmt.Errorf("invalid ethereum signature length")
		}
		r := new(big.Int).SetBytes(signature[:32])
		s := new(big.Int).SetBytes(signature[32:64])
		return ecdsa.Verify(ecdsaPubKey, message, r, s), nil

	case networks.Cosmos:
		// For Cosmos, signature is R || S
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

	// Skip header byte and length byte
	pos := 2
	if pos+2 > len(sig) {
		return nil, nil, fmt.Errorf("invalid DER signature length")
	}

	// Check for integer marker
	if sig[pos] != 0x02 {
		return nil, nil, fmt.Errorf("invalid R marker in DER signature")
	}
	pos++

	// Get R length
	rLen := int(sig[pos])
	pos++
	if pos+rLen > len(sig) {
		return nil, nil, fmt.Errorf("invalid R length in DER signature")
	}

	// Extract R value
	r = new(big.Int).SetBytes(sig[pos : pos+rLen])
	pos += rLen

	// Check for integer marker for S
	if pos+2 > len(sig) || sig[pos] != 0x02 {
		return nil, nil, fmt.Errorf("invalid S marker in DER signature")
	}
	pos++

	// Get S length
	sLen := int(sig[pos])
	pos++
	if pos+sLen > len(sig) {
		return nil, nil, fmt.Errorf("invalid S length in DER signature")
	}

	// Extract S value
	s = new(big.Int).SetBytes(sig[pos : pos+sLen])

	return r, s, nil
}

// Signer interface defines the methods required for network-specific signing
type Signer interface {
	Sign(ctx context.Context, signingCtx *SigningContext) ([]byte, error)
	Verify(pubKey []byte, msg []byte, signature []byte) (bool, error)
}

// UniversalSigner implements Signer for all networks
type UniversalSigner struct {
	networkParams *types.NetworkParams
}

func (s *UniversalSigner) Sign(ctx context.Context, signingCtx *SigningContext) ([]byte, error) {
	// Implementation for signing
	networkType := networks.NetworkType(s.networkParams.NetworkType)
	switch networkType {
	case networks.Bitcoin:
		// Bitcoin signing logic
		return nil, fmt.Errorf("bitcoin signing not implemented")
	case networks.Ethereum:
		// Ethereum signing logic
		return nil, fmt.Errorf("ethereum signing not implemented with chain ID %s", s.networkParams.SigningConfig.ChainId)
	case networks.Cosmos:
		// Cosmos signing logic
		return nil, fmt.Errorf("cosmos signing not implemented")
	case networks.Solana:
		// Solana signing logic
		return nil, fmt.Errorf("solana signing not implemented")
	case networks.Cardano:
		// Cardano signing logic
		return nil, fmt.Errorf("cardano signing not implemented")
	case networks.Aptos:
		// Aptos signing logic
		return nil, fmt.Errorf("aptos signing not implemented")
	case networks.Sui:
		// Sui signing logic
		return nil, fmt.Errorf("sui signing not implemented")
	default:
		return nil, fmt.Errorf("unsupported network type: %s", networkType)
	}
}

func (s *UniversalSigner) Verify(pubKey []byte, msg []byte, signature []byte) (bool, error) {
	// Implement verification logic based on network parameters
	return false, fmt.Errorf("verification not implemented")
}

// Helper functions for signature formatting
func formatBitcoinSignature(result *tss.SignResult) ([]byte, error) {
	if result == nil {
		return nil, fmt.Errorf("sign result is nil")
	}

	// Convert R and S to bytes with 32-byte padding
	rBytes := make([]byte, 32)
	sBytes := make([]byte, 32)
	result.R.FillBytes(rBytes)
	result.S.FillBytes(sBytes)

	// DER encoding:
	// 0x30 [total-length] 0x02 [R-length] [R] 0x02 [S-length] [S]
	
	// Remove leading zeros from R and S
	rLen := len(bytes.TrimLeft(rBytes, "\x00"))
	sLen := len(bytes.TrimLeft(sBytes, "\x00"))
	
	// Total length is:
	// 2 bytes for type and length of R
	// rLen bytes for R
	// 2 bytes for type and length of S
	// sLen bytes for S
	totalLen := 2 + rLen + 2 + sLen
	
	// Create the DER signature
	der := make([]byte, 2+totalLen)
	der[0] = 0x30 // Sequence tag
	der[1] = byte(totalLen)
	
	// Encode R
	der[2] = 0x02 // Integer tag
	der[3] = byte(rLen)
	copy(der[4:], bytes.TrimLeft(rBytes, "\x00"))
	
	// Encode S
	der[4+rLen] = 0x02 // Integer tag
	der[5+rLen] = byte(sLen)
	copy(der[6+rLen:], bytes.TrimLeft(sBytes, "\x00"))
	
	return der, nil
}

func formatEthereumSignature(result *tss.SignResult, chainID string) ([]byte, error) {
	if result == nil {
		return nil, fmt.Errorf("sign result is nil")
	}

	// Ethereum signatures are 65 bytes: R (32 bytes) + S (32 bytes) + V (1 byte)
	signature := make([]byte, 65)

	// Convert R and S to bytes with 32-byte padding
	rBytes := make([]byte, 32)
	sBytes := make([]byte, 32)
	result.R.FillBytes(rBytes)
	result.S.FillBytes(sBytes)

	// Copy R and S
	copy(signature[:32], rBytes)
	copy(signature[32:64], sBytes)

	// Calculate V based on chainID
	// For EIP-155, V = 27 + chainID * 2 + 35
	v := byte(27) // Default V value
	if chainID != "" {
		chainIDInt := new(big.Int)
		chainIDInt, ok := chainIDInt.SetString(chainID, 10)
		if !ok {
			return nil, fmt.Errorf("invalid chainID: %s", chainID)
		}
		
		vBig := new(big.Int).Add(big.NewInt(27), new(big.Int).Mul(chainIDInt, big.NewInt(2)))
		vBig.Add(vBig, big.NewInt(35))
		
		if !vBig.IsUint64() {
			return nil, fmt.Errorf("V value overflow")
		}
		v = byte(vBig.Uint64())
	}
	
	signature[64] = v

	return signature, nil
}

func formatCosmosSignature(result *tss.SignResult) ([]byte, error) {
	if result == nil {
		return nil, fmt.Errorf("sign result is nil")
	}

	// Cosmos signatures are just R || S concatenated (64 bytes)
	signature := make([]byte, 64)

	// Convert R and S to bytes with 32-byte padding
	rBytes := make([]byte, 32)
	sBytes := make([]byte, 32)
	result.R.FillBytes(rBytes)
	result.S.FillBytes(sBytes)

	// Copy R and S
	copy(signature[:32], rBytes)
	copy(signature[32:], sBytes)

	return signature, nil
}
