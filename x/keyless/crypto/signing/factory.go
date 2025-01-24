package signing

import (
	"context"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"selfchain/x/keyless/networks"
	"selfchain/x/keyless/types"
	"selfchain/x/keyless/crypto/signing/ecdsa"
	"selfchain/x/keyless/crypto/signing/format"
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

// SignerFactory creates and manages signers for different networks
type SignerFactory struct {
	registry *networks.NetworkRegistry
	signer   *ecdsa.ECDSASigner
}

// NewSignerFactory creates a new signer factory
func NewSignerFactory(registry *networks.NetworkRegistry) *SignerFactory {
	return &SignerFactory{
		registry: registry,
		signer:   ecdsa.NewECDSASigner(nil, nil), // Default to ECDSA signer
	}
}

// Sign signs a message for a specific network
func (f *SignerFactory) Sign(ctx context.Context, networkID string, message []byte, metadata map[string]interface{}, signResult *format.SignatureResult) ([]byte, error) {
	// Parse network ID
	networkType, _, err := networks.ParseNetworkID(networkID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network ID: %w", err)
	}

	// Get network parameters
	params := networks.GetDefaultNetworkParams(networkID)
	if params == nil {
		return nil, fmt.Errorf("network parameters not found for: %s", networkID)
	}

	// Use network parameters for signing
	signingCtx := &SigningContext{
		NetworkParams: params,
		Message:       message,
		Metadata:      metadata,
	}

	// Sign the message using the ECDSA signer
	signResult, err = f.signer.Sign(ctx, signingCtx.Message, ecdsa.ECDSA)
	if err != nil {
		return nil, err
	}

	// Format signature based on network type
	var formattedSig []byte
	switch networkType {
	case networks.Bitcoin:
		formattedSig, err = format.FormatBitcoinSignature(signResult)
	case networks.Ethereum:
		formattedSig, err = format.FormatEthereumSignature(signResult)
	case networks.Cosmos:
		formattedSig, err = format.FormatCosmosSignature(signResult)
	default:
		return nil, fmt.Errorf("unsupported network type: %s", networkType)
	}
	if err != nil {
		return nil, err
	}

	// Store formatted signature in SignatureResult
	signResult.Bytes = formattedSig
	return formattedSig, nil
}

// Verify verifies a signature for a specific network
func (f *SignerFactory) Verify(networkID string, pubKeyBytes []byte, message []byte, signature []byte) (bool, error) {
	// Parse network ID
	networkType, _, err := networks.ParseNetworkID(networkID)
	if err != nil {
		return false, fmt.Errorf("failed to parse network ID: %w", err)
	}

	// Get network parameters
	params := networks.GetDefaultNetworkParams(networkID)
	if params == nil {
		return false, fmt.Errorf("network parameters not found for: %s", networkID)
	}

	// Parse signature based on network type
	var signatureResult *format.SignatureResult
	switch networkType {
	case networks.Bitcoin:
		signatureResult, err = format.ParseBitcoinSignature(signature)
	case networks.Ethereum:
		signatureResult, err = format.ParseEthereumSignature(signature)
		if err != nil {
			return false, err
		}

		// For Ethereum, we need the raw public key (64 bytes)
		if len(pubKeyBytes) == 33 {
			pubKey, err := btcec.ParsePubKey(pubKeyBytes)
			if err != nil {
				return false, fmt.Errorf("failed to parse public key: %w", err)
			}
			pubKeyBytes = pubKey.SerializeUncompressed()[1:] // Remove prefix and use raw 64 bytes
		} else if len(pubKeyBytes) == 65 {
			pubKeyBytes = pubKeyBytes[1:] // Remove prefix and use raw 64 bytes
		}

	case networks.Cosmos:
		signatureResult, err = format.ParseCosmosSignature(signature)
	default:
		return false, fmt.Errorf("unsupported network type: %s", networkType)
	}

	if err != nil {
		return false, fmt.Errorf("failed to parse signature: %w", err)
	}

	// Verify the signature using the ECDSA signer
	return f.signer.Verify(context.Background(), message, signatureResult, pubKeyBytes)
}
