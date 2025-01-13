package signing

import (
	"context"
	"fmt"
	"sync"

	"selfchain/x/keyless/networks"
	"selfchain/x/keyless/tss"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
)

// SigningContext contains the context for signing operations
type SigningContext struct {
	NetworkInfo *networks.NetworkInfo
	Message     []byte
	Party1Data  interface{}
	Party2Data  interface{}
	Metadata    map[string]interface{} // Network-specific metadata
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

// CreateSigner creates an appropriate signer for the given network
func (f *SignerFactory) CreateSigner(networkType networks.NetworkType, chainID string) (Signer, error) {
	networkInfo, err := f.registry.GetNetwork(networkType, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network info: %w", err)
	}

	return &UniversalSigner{networkInfo: networkInfo}, nil
}

// Sign signs a message for the specified network
func (f *SignerFactory) Sign(ctx context.Context, networkID string, message []byte, metadata map[string]interface{}) ([]byte, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	params := networks.DefaultNetworkParams(networkID)
	if params == nil {
		return nil, fmt.Errorf("unsupported network: %s", networkID)
	}

	req := &SignRequest{
		NetworkID: networkID,
		Params:    params,
		Message:   message,
		MetaData:  metadata,
	}

	resp, err := f.signer.Sign(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("signing failed: %w", err)
	}

	return resp.Signature, nil
}

// Verify verifies a signature for the specified network
func (f *SignerFactory) Verify(networkID string, pubKey, message, signature []byte) (bool, error) {
	params := networks.DefaultNetworkParams(networkID)
	if params == nil {
		return false, fmt.Errorf("unsupported network: %s", networkID)
	}

	// Implement verification logic based on network parameters
	return false, fmt.Errorf("verification not implemented")
}

// Signer interface defines the methods required for network-specific signing
type Signer interface {
	Sign(ctx context.Context, signingCtx *SigningContext) ([]byte, error)
	Verify(pubKey []byte, msg []byte, signature []byte) (bool, error)
}

// UniversalSigner implements Signer for all networks
type UniversalSigner struct {
	networkInfo *networks.NetworkInfo
}

func (s *UniversalSigner) Sign(ctx context.Context, signingCtx *SigningContext) ([]byte, error) {
	// Convert generic interface{} to specific TSS party data
	party1Data, ok := signingCtx.Party1Data.(*keygen.LocalPartySaveData)
	if !ok {
		return nil, fmt.Errorf("invalid party1 data type")
	}
	party2Data, ok := signingCtx.Party2Data.(*keygen.LocalPartySaveData)
	if !ok {
		return nil, fmt.Errorf("invalid party2 data type")
	}

	// Use TSS to sign
	result, err := tss.SignMessage(ctx, signingCtx.Message, party1Data, party2Data)
	if err != nil {
		return nil, fmt.Errorf("tss signing failed: %w", err)
	}

	// Format signature according to network requirements
	switch s.networkInfo.NetworkType {
	case networks.Bitcoin:
		return formatBitcoinSignature(result)
	case networks.Ethereum:
		return formatEthereumSignature(result, s.networkInfo.SigningConfig.ChainID)
	case networks.Cosmos:
		return formatCosmosSignature(result)
	case networks.Solana:
		// Implement Solana Ed25519 signing
		return nil, fmt.Errorf("not implemented")
	case networks.Cardano:
		// Implement Cardano signing
		return nil, fmt.Errorf("not implemented")
	case networks.Aptos:
		// Implement Aptos signing
		return nil, fmt.Errorf("not implemented")
	case networks.Sui:
		// Implement Sui signing
		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("unsupported network type: %s", s.networkInfo.NetworkType)
	}
}

func (s *UniversalSigner) Verify(pubKey []byte, msg []byte, signature []byte) (bool, error) {
	// Implement verification logic based on network parameters
	return false, fmt.Errorf("verification not implemented")
}

// Helper functions for signature formatting
func formatBitcoinSignature(result *tss.SignResult) ([]byte, error) {
	// TODO: Implement DER encoding for Bitcoin signatures
	return nil, fmt.Errorf("not implemented")
}

func formatEthereumSignature(result *tss.SignResult, chainID string) ([]byte, error) {
	// TODO: Implement Ethereum signature formatting with chainID
	return nil, fmt.Errorf("not implemented")
}

func formatCosmosSignature(result *tss.SignResult) ([]byte, error) {
	// TODO: Implement Cosmos signature formatting
	return nil, fmt.Errorf("not implemented")
}
