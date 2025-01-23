package signing

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
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
	ecdsaSigner   *ECDSASigner
	eddsaSigner   *EdDSASigner
}

func NewUniversalSigner(party1Data, party2Data *keygen.LocalPartySaveData) *UniversalSigner {
	return &UniversalSigner{
		ecdsaSigner: NewECDSASigner(party1Data, party2Data),
		eddsaSigner: NewEdDSASigner(party1Data, party2Data),
	}
}

func (s *UniversalSigner) Sign(ctx context.Context, signingCtx *SigningContext) ([]byte, error) {
	// Implementation for signing
	networkType := networks.NetworkType(s.networkParams.NetworkType)

	var signResult *SignatureResult
	var err error

	switch networkType {
	case networks.Bitcoin, networks.Ethereum:
		signResult, err = s.ecdsaSigner.Sign(ctx, signingCtx.Message, ECDSA)
		if err != nil {
			return nil, fmt.Errorf("ECDSA signing failed: %w", err)
		}

		if networkType == networks.Bitcoin {
			return formatBitcoinSignature(signResult)
		}
		return formatEthereumSignature(signResult, s.networkParams.SigningConfig.ChainId)

	case networks.Cosmos:
		signResult, err = s.ecdsaSigner.Sign(ctx, signingCtx.Message, ECDSA)
		if err != nil {
			return nil, fmt.Errorf("ECDSA signing failed: %w", err)
		}
		return formatCosmosSignature(signResult)

	case networks.Solana, networks.Cardano:
		signResult, err = s.eddsaSigner.Sign(ctx, signingCtx.Message, EdDSA)
		if err != nil {
			return nil, fmt.Errorf("EdDSA signing failed: %w", err)
		}
		return signResult.Bytes, nil

	default:
		return nil, fmt.Errorf("unsupported network type: %s", networkType)
	}
}

func (s *UniversalSigner) Verify(pubKey []byte, msg []byte, signature []byte) (bool, error) {
	networkType := networks.NetworkType(s.networkParams.NetworkType)

	switch networkType {
	case networks.Bitcoin, networks.Ethereum, networks.Cosmos:
		// Parse DER signature for ECDSA networks
		rVal, sVal, err := UnmarshalDERSignature(signature)
		if err != nil {
			return false, fmt.Errorf("failed to parse DER signature: %w", err)
		}

		// Parse public key
		pubKeyObj, err := btcec.ParsePubKey(pubKey)
		if err != nil {
			return false, fmt.Errorf("failed to parse public key: %w", err)
		}

		sig := &SignatureResult{R: rVal, S: sVal}
		return s.ecdsaSigner.Verify(context.Background(), msg, sig, pubKeyObj)

	case networks.Solana, networks.Cardano:
		// For EdDSA networks, use raw signature bytes
		sig := &SignatureResult{Bytes: signature}

		// Parse public key
		pubKeyObj, err := btcec.ParsePubKey(pubKey)
		if err != nil {
			return false, fmt.Errorf("failed to parse public key: %w", err)
		}

		return s.eddsaSigner.Verify(context.Background(), msg, sig, pubKeyObj)

	default:
		return false, fmt.Errorf("unsupported network type: %s", networkType)
	}
}

// Helper functions for signature formatting
func formatBitcoinSignature(result *SignatureResult) ([]byte, error) {
	// For Bitcoin, we need DER encoding
	var b bytes.Buffer

	// Write sequence marker
	b.WriteByte(0x30)

	// Leave a byte for the total length
	b.WriteByte(0x00)

	// Write R value
	b.WriteByte(0x02)
	rb := result.R.Bytes()
	if rb[0]&0x80 == 0x80 {
		b.WriteByte(byte(len(rb) + 1))
		b.WriteByte(0x00)
	} else {
		b.WriteByte(byte(len(rb)))
	}
	b.Write(rb)

	// Write S value
	b.WriteByte(0x02)
	sb := result.S.Bytes()
	if sb[0]&0x80 == 0x80 {
		b.WriteByte(byte(len(sb) + 1))
		b.WriteByte(0x00)
	} else {
		b.WriteByte(byte(len(sb)))
	}
	b.Write(sb)

	// Fill in total length
	der := b.Bytes()
	der[1] = byte(len(der) - 2)

	return der, nil
}

func formatEthereumSignature(result *SignatureResult, chainID string) ([]byte, error) {
	// For Ethereum, we need [R || S || V] format
	// R and S are padded to 32 bytes
	rBytes := padTo32(result.R.Bytes())
	sBytes := padTo32(result.S.Bytes())

	// V is recovery ID + 27 + (chainID * 2 + 35) if EIP-155 is used
	var v byte = result.Recovery + 27

	// If chainID is provided, apply EIP-155
	if chainID != "" {
		chainIDInt := new(big.Int)
		chainIDInt.SetString(chainID, 10)
		if chainIDInt.Sign() > 0 {
			v = result.Recovery + 35 + byte(chainIDInt.Uint64()*2)
		}
	}

	// Combine [R || S || V]
	signature := make([]byte, 65)
	copy(signature[0:32], rBytes)
	copy(signature[32:64], sBytes)
	signature[64] = v

	return signature, nil
}

func formatCosmosSignature(result *SignatureResult) ([]byte, error) {
	// For Cosmos, we use DER encoding similar to Bitcoin
	return formatBitcoinSignature(result)
}

// Helper function to pad a byte slice to 32 bytes
func padTo32(b []byte) []byte {
	if len(b) > 32 {
		return b[:32]
	}
	result := make([]byte, 32)
	copy(result[32-len(b):], b)
	return result
}
