package networks

import (
	"errors"
	"fmt"
	"strings"

	"selfchain/x/keyless/types"
)

// NetworkHandler handles network-specific operations
type NetworkHandler struct {
	networkParams *types.NetworkParams
}

// NewNetworkHandler creates a new network handler
func NewNetworkHandler(networkParams *types.NetworkParams) (*NetworkHandler, error) {
	if networkParams == nil {
		return nil, errors.New("network params cannot be nil")
	}
	return &NetworkHandler{networkParams: networkParams}, nil
}

// GetNetworkParams returns the network params
func (h *NetworkHandler) GetNetworkParams() *types.NetworkParams {
	return h.networkParams
}

// ValidateNetworkParams validates network parameters
func ValidateNetworkParams(params *types.NetworkParams) error {
	if params == nil {
		return errors.New("network params cannot be nil")
	}

	if params.ChainId == "" {
		return errors.New("chain ID cannot be empty")
	}

	if params.SigningAlgorithm == "" {
		return errors.New("signing algorithm cannot be empty")
	}

	if !isValidSigningAlgorithm(params.SigningAlgorithm) {
		return fmt.Errorf("unsupported signing algorithm: %s", params.SigningAlgorithm)
	}

	if params.KeygenThreshold <= 0 {
		return errors.New("keygen threshold must be greater than 0")
	}

	if params.SigningThreshold <= 0 {
		return errors.New("signing threshold must be greater than 0")
	}

	if params.SigningThreshold > params.KeygenThreshold {
		return errors.New("signing threshold cannot be greater than keygen threshold")
	}

	return nil
}

// isValidSigningAlgorithm checks if the signing algorithm is supported
func isValidSigningAlgorithm(algo string) bool {
	switch algo {
	case "ecdsa":
		return true
	case "eddsa":
		return true
	case "bls":
		return true
	case "schnorr":
		return true
	default:
		return false
	}
}

// ValidateAddress validates a network-specific address
func (h *NetworkHandler) ValidateAddress(address string) error {
	if address == "" {
		return errors.New("address cannot be empty")
	}

	if !strings.HasPrefix(address, h.networkParams.AddressPrefix) {
		return fmt.Errorf("invalid address prefix, expected %s", h.networkParams.AddressPrefix)
	}

	return nil
}

// ValidatePublicKey validates a network-specific public key
func (h *NetworkHandler) ValidatePublicKey(pubKey []byte, keyType types.KeyType) error {
	switch keyType {
	case types.KeyType_KEY_TYPE_ECDSA:
		return h.validateECDSAPublicKey(pubKey)
	case types.KeyType_KEY_TYPE_EDDSA:
		return h.validateEdDSAPublicKey(pubKey)
	case types.KeyType_KEY_TYPE_BLS:
		return h.validateBLSPublicKey(pubKey)
	case types.KeyType_KEY_TYPE_SCHNORR:
		return h.validateSchnorrPublicKey(pubKey)
	default:
		return fmt.Errorf("unsupported key type: %v", keyType)
	}
}

func (h *NetworkHandler) validateECDSAPublicKey(pubKey []byte) error {
	if len(pubKey) != 65 && len(pubKey) != 33 {
		return fmt.Errorf("invalid ECDSA public key length: %d", len(pubKey))
	}
	
	// For uncompressed keys (65 bytes)
	if len(pubKey) == 65 && pubKey[0] != 0x04 {
		return fmt.Errorf("invalid ECDSA public key format")
	}
	
	// For compressed keys (33 bytes)
	if len(pubKey) == 33 && pubKey[0] != 0x02 && pubKey[0] != 0x03 {
		return fmt.Errorf("invalid ECDSA public key format")
	}
	
	return nil
}

func (h *NetworkHandler) validateEdDSAPublicKey(pubKey []byte) error {
	// Ed25519 public keys are always 32 bytes
	if len(pubKey) != 32 {
		return fmt.Errorf("invalid EdDSA public key length: %d", len(pubKey))
	}
	return nil
}

func (h *NetworkHandler) validateBLSPublicKey(pubKey []byte) error {
	// BLS12-381 G1 public keys are 48 bytes
	if len(pubKey) != 48 {
		return fmt.Errorf("invalid BLS public key length: %d", len(pubKey))
	}
	return nil
}

func (h *NetworkHandler) validateSchnorrPublicKey(pubKey []byte) error {
	// Schnorr public keys are 32 bytes (same as secp256k1)
	if len(pubKey) != 32 {
		return fmt.Errorf("invalid Schnorr public key length: %d", len(pubKey))
	}
	return nil
}
