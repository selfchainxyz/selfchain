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

	if params.NetworkType == "" {
		return errors.New("network type cannot be empty")
	}

	if params.ChainId == "" {
		return errors.New("chain ID cannot be empty")
	}

	if params.SigningAlgorithm == "" {
		return errors.New("signing algorithm cannot be empty")
	}

	if params.CurveType == "" {
		return errors.New("curve type cannot be empty")
	}

	if params.AddressPrefix == "" {
		return errors.New("address prefix cannot be empty")
	}

	if !isValidSigningAlgorithm(params.SigningAlgorithm) {
		return fmt.Errorf("invalid signing algorithm: %s", params.SigningAlgorithm)
	}

	return nil
}

// isValidSigningAlgorithm checks if the signing algorithm is supported
func isValidSigningAlgorithm(algo string) bool {
	validAlgos := []string{
		"ECDSA",
		"EdDSA",
		"BLS",
		"Schnorr",
	}

	algo = strings.ToUpper(algo)
	for _, validAlgo := range validAlgos {
		if algo == validAlgo {
			return true
		}
	}
	return false
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
func (h *NetworkHandler) ValidatePublicKey(pubKey string) error {
	if pubKey == "" {
		return errors.New("public key cannot be empty")
	}

	// Add network-specific public key validation based on signing algorithm
	switch strings.ToUpper(h.networkParams.SigningAlgorithm) {
	case "ECDSA":
		return validateECDSAPublicKey(pubKey)
	case "EDDSA":
		return validateEdDSAPublicKey(pubKey)
	case "BLS":
		return validateBLSPublicKey(pubKey)
	case "SCHNORR":
		return validateSchnorrPublicKey(pubKey)
	default:
		return fmt.Errorf("unsupported signing algorithm: %s", h.networkParams.SigningAlgorithm)
	}
}

func validateECDSAPublicKey(pubKey string) error {
	// TODO: Implement ECDSA public key validation
	return nil
}

func validateEdDSAPublicKey(pubKey string) error {
	// TODO: Implement EdDSA public key validation
	return nil
}

func validateBLSPublicKey(pubKey string) error {
	// TODO: Implement BLS public key validation
	return nil
}

func validateSchnorrPublicKey(pubKey string) error {
	// TODO: Implement Schnorr public key validation
	return nil
}
