package types

import (
	"math/big"
)

// CryptoService defines the interface for cryptographic operations
type CryptoService interface {
	ConvertSignatureToEdDSA(r, s *big.Int) ([]byte, error)
}
