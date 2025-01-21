package eddsa

import (
	"math/big"
	"selfchain/x/keyless/crypto/types"
)

// EdDSAConverter implements the CryptoService interface for EdDSA operations
type EdDSAConverter struct{}

var _ types.CryptoService = (*EdDSAConverter)(nil)

// ConvertSignatureToEdDSA converts an ECDSA signature (r,s) to EdDSA format
func (e *EdDSAConverter) ConvertSignatureToEdDSA(r, s *big.Int) ([]byte, error) {
	return ConvertSignatureToEdDSA(r, s)
}
