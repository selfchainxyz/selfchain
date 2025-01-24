package format

import (
	"errors"
	"math/big"
)

// SignatureResult represents a signature result with R, S components
type SignatureResult struct {
	R        *big.Int
	S        *big.Int
	V        uint8
	Recovery byte
	Bytes    []byte
}

// FormatEthereumSignature formats a signature in Ethereum's format (R || S || V)
func FormatEthereumSignature(sig *SignatureResult) ([]byte, error) {
	if sig == nil {
		return nil, errors.New("signature is nil")
	}

	result := make([]byte, 65)

	// Ensure R and S are padded to 32 bytes
	rBytes := make([]byte, 32)
	sBytes := make([]byte, 32)
	sig.R.FillBytes(rBytes)  // Use FillBytes to properly pad
	sig.S.FillBytes(sBytes)  // Use FillBytes to properly pad

	// Copy R and S into result
	copy(result[0:32], rBytes)
	copy(result[32:64], sBytes)

	// Add recovery ID (V)
	result[64] = sig.Recovery
	if result[64] == 0 {
		result[64] = 27 // Default V value for Ethereum
	}

	return result, nil
}

// FormatCosmosSignature formats a signature in Cosmos' format (R || S)
func FormatCosmosSignature(sig *SignatureResult) ([]byte, error) {
	if sig == nil {
		return nil, errors.New("signature is nil")
	}

	result := make([]byte, 64)

	// Ensure R and S are padded to 32 bytes
	rBytes := make([]byte, 32)
	sBytes := make([]byte, 32)
	sig.R.FillBytes(rBytes)  // Use FillBytes to properly pad
	sig.S.FillBytes(sBytes)  // Use FillBytes to properly pad

	// Copy R and S into result
	copy(result[0:32], rBytes)
	copy(result[32:64], sBytes)

	return result, nil
}

// ParseEthereumSignature parses an Ethereum signature into R, S, V components
func ParseEthereumSignature(sig []byte) (*SignatureResult, error) {
	if len(sig) != 65 {
		return nil, errors.New("invalid ethereum signature length")
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := uint8(sig[64])

	return &SignatureResult{
		R:        r,
		S:        s,
		V:        v,
		Recovery: v,
		Bytes:    sig,
	}, nil
}

// ParseCosmosSignature parses a Cosmos signature into R, S components
func ParseCosmosSignature(sig []byte) (*SignatureResult, error) {
	if len(sig) != 64 {
		return nil, errors.New("invalid cosmos signature length")
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])

	return &SignatureResult{
		R:     r,
		S:     s,
		Bytes: sig,
	}, nil
}

// ParseBitcoinSignature parses a Bitcoin signature in DER format
func ParseBitcoinSignature(sig []byte) (*SignatureResult, error) {
	if len(sig) < 8 { // Minimum DER signature length
		return nil, errors.New("malformed bitcoin signature")
	}

	// Remove SIGHASH_ALL byte if present
	sigBytes := sig
	if len(sig) > 0 && sig[len(sig)-1] == 0x01 {
		sigBytes = sig[:len(sig)-1]
	}

	r, s, err := UnmarshalDERSignature(sigBytes)
	if err != nil {
		return nil, err
	}

	return &SignatureResult{
		R:     r,
		S:     s,
		Bytes: sigBytes,
	}, nil
}

// UnmarshalDERSignature parses a DER signature into R and S components
func UnmarshalDERSignature(sig []byte) (*big.Int, *big.Int, error) {
	// ASN.1 DER format:
	// 0x30 [total-length] 0x02 [R-length] [R] 0x02 [S-length] [S]
	if len(sig) < 8 {
		return nil, nil, errors.New("signature too short")
	}

	// Check sequence marker
	if sig[0] != 0x30 {
		return nil, nil, errors.New("signature not a sequence")
	}

	// Get R component
	rStart := 4 // Skip sequence header (0x30) + length byte + integer marker (0x02) + length byte
	rLen := int(sig[3])
	if rLen+rStart > len(sig) {
		return nil, nil, errors.New("R component length invalid")
	}
	r := new(big.Int).SetBytes(sig[rStart : rStart+rLen])

	// Get S component
	sStart := rStart + rLen + 2 // Skip R + integer marker (0x02) + length byte
	if sStart+1 > len(sig) {
		return nil, nil, errors.New("S component missing")
	}
	sLen := int(sig[sStart-1])
	if sLen+sStart > len(sig) {
		return nil, nil, errors.New("S component length invalid")
	}
	s := new(big.Int).SetBytes(sig[sStart : sStart+sLen])

	return r, s, nil
}

// FormatDERSignature formats a signature in DER format
func FormatDERSignature(sig *SignatureResult) ([]byte, error) {
	if sig == nil {
		return nil, errors.New("signature is nil")
	}
	if sig.R == nil {
		return nil, errors.New("R value is nil")
	}
	if sig.S == nil {
		return nil, errors.New("S value is nil")
	}

	// Validate R and S values
	if sig.R.Sign() <= 0 || sig.S.Sign() <= 0 {
		return nil, errors.New("invalid signature: R and S must be positive")
	}

	// Check if values are too large (> 32 bytes when padded)
	maxValue := new(big.Int).Lsh(big.NewInt(1), 256)
	if sig.R.Cmp(maxValue) >= 0 || sig.S.Cmp(maxValue) >= 0 {
		return nil, errors.New("invalid signature: R or S too large")
	}

	// Calculate lengths
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	rLen := len(rBytes)
	sLen := len(sBytes)

	// Add extra byte to maintain positive numbers
	if rBytes[0]&0x80 == 0x80 {
		rLen++
	}
	if sBytes[0]&0x80 == 0x80 {
		sLen++
	}

	// Total length is:
	// 0x30 + length byte + 0x02 + rLength + rValue + 0x02 + sLength + sValue
	totalLen := 2 + 2 + rLen + 2 + sLen

	// Create DER signature
	der := make([]byte, totalLen)
	i := 0

	// Sequence marker
	der[i] = 0x30
	i++
	der[i] = byte(2 + rLen + 2 + sLen)
	i++

	// R value
	der[i] = 0x02
	i++
	der[i] = byte(rLen)
	i++
	if rBytes[0]&0x80 == 0x80 {
		der[i] = 0x00
		i++
	}
	copy(der[i:], rBytes)
	i += len(rBytes)

	// S value
	der[i] = 0x02
	i++
	der[i] = byte(sLen)
	i++
	if sBytes[0]&0x80 == 0x80 {
		der[i] = 0x00
		i++
	}
	copy(der[i:], sBytes)

	return der, nil
}
