package format

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTestSignature() *SignatureResult {
	r := new(big.Int).SetBytes([]byte{
		0x1b, 0x84, 0xc5, 0x56, 0x7b, 0x23, 0x4c, 0x8a,
		0x87, 0x32, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc,
		0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc,
		0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc,
	})
	s := new(big.Int).SetBytes([]byte{
		0x2c, 0x95, 0xd6, 0x67, 0x8c, 0x34, 0x5d, 0x9b,
		0x98, 0x43, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd,
		0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd,
		0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd,
	})
	return &SignatureResult{
		R:        r,
		S:        s,
		V:        27,
		Recovery: 27,
		Bytes:    append(r.Bytes(), append(s.Bytes(), byte(27))...),
	}
}

func Test_FormatEthereumSignature(t *testing.T) {
	tests := []struct {
		name      string
		sig       *SignatureResult
		wantError bool
	}{
		{
			name:      "valid signature",
			sig:       setupTestSignature(),
			wantError: false,
		},
		{
			name:      "nil signature",
			sig:       nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := FormatEthereumSignature(tt.sig)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, formatted)
			} else {
				require.NoError(t, err)
				require.Equal(t, 65, len(formatted))
				require.Equal(t, tt.sig.R.FillBytes(make([]byte, 32)), formatted[:32])
				require.Equal(t, tt.sig.S.FillBytes(make([]byte, 32)), formatted[32:64])
				require.Equal(t, tt.sig.Recovery, formatted[64])
			}
		})
	}
}

func Test_FormatCosmosSignature(t *testing.T) {
	tests := []struct {
		name      string
		sig       *SignatureResult
		wantError bool
	}{
		{
			name:      "valid signature",
			sig:       setupTestSignature(),
			wantError: false,
		},
		{
			name:      "nil signature",
			sig:       nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := FormatCosmosSignature(tt.sig)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, formatted)
			} else {
				require.NoError(t, err)
				require.Equal(t, 64, len(formatted))
				require.Equal(t, tt.sig.R.FillBytes(make([]byte, 32)), formatted[:32])
				require.Equal(t, tt.sig.S.FillBytes(make([]byte, 32)), formatted[32:64])
			}
		})
	}
}

func Test_ParseEthereumSignature(t *testing.T) {
	validSig := setupTestSignature()
	formattedSig, err := FormatEthereumSignature(validSig)
	require.NoError(t, err)

	tests := []struct {
		name      string
		sig       []byte
		wantError bool
	}{
		{
			name:      "valid signature",
			sig:       formattedSig,
			wantError: false,
		},
		{
			name:      "invalid length",
			sig:       []byte{1, 2, 3},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseEthereumSignature(tt.sig)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, parsed)
			} else {
				require.NoError(t, err)
				require.Equal(t, validSig.R.Bytes(), parsed.R.Bytes())
				require.Equal(t, validSig.S.Bytes(), parsed.S.Bytes())
				require.Equal(t, validSig.V, parsed.V)
			}
		})
	}
}

func Test_ParseCosmosSignature(t *testing.T) {
	validSig := setupTestSignature()
	formattedSig, err := FormatCosmosSignature(validSig)
	require.NoError(t, err)

	tests := []struct {
		name      string
		sig       []byte
		wantError bool
	}{
		{
			name:      "valid signature",
			sig:       formattedSig,
			wantError: false,
		},
		{
			name:      "invalid length",
			sig:       []byte{1, 2, 3},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseCosmosSignature(tt.sig)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, parsed)
			} else {
				require.NoError(t, err)
				require.Equal(t, validSig.R.Bytes(), parsed.R.Bytes())
				require.Equal(t, validSig.S.Bytes(), parsed.S.Bytes())
			}
		})
	}
}

func Test_ParseBitcoinSignature(t *testing.T) {
	validSig := setupTestSignature()
	derSig, err := FormatDERSignature(validSig)
	require.NoError(t, err)

	tests := []struct {
		name      string
		sig       []byte
		wantError bool
	}{
		{
			name:      "valid DER signature",
			sig:       derSig,
			wantError: false,
		},
		{
			name:      "valid DER signature with SIGHASH_ALL",
			sig:       append(derSig, 0x01),
			wantError: false,
		},
		{
			name:      "invalid length",
			sig:       []byte{1, 2, 3},
			wantError: true,
		},
		{
			name:      "invalid DER format",
			sig:       []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseBitcoinSignature(tt.sig)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, parsed)
			} else {
				require.NoError(t, err)
				require.NotNil(t, parsed.R)
				require.NotNil(t, parsed.S)
			}
		})
	}
}

func Test_UnmarshalDERSignature(t *testing.T) {
	validSig := setupTestSignature()
	derSig, err := FormatDERSignature(validSig)
	require.NoError(t, err)

	tests := []struct {
		name      string
		sig       []byte
		wantError bool
	}{
		{
			name:      "valid DER signature",
			sig:       derSig,
			wantError: false,
		},
		{
			name:      "too short",
			sig:       []byte{1, 2, 3},
			wantError: true,
		},
		{
			name:      "invalid sequence marker",
			sig:       []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			wantError: true,
		},
		{
			name:      "invalid R length",
			sig:       append([]byte{0x30, 0xff, 0x02, 0xff}, make([]byte, 4)...),
			wantError: true,
		},
		{
			name:      "invalid S length",
			sig:       append([]byte{0x30, 0x06, 0x02, 0x01, 0x01, 0x02, 0xff}, make([]byte, 1)...),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, s, err := UnmarshalDERSignature(tt.sig)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, r)
				require.Nil(t, s)
			} else {
				require.NoError(t, err)
				require.Equal(t, validSig.R.Bytes(), r.Bytes())
				require.Equal(t, validSig.S.Bytes(), s.Bytes())
			}
		})
	}
}

func Test_FormatDERSignature(t *testing.T) {
	tests := []struct {
		name      string
		sig       *SignatureResult
		wantError bool
	}{
		{
			name:      "valid signature",
			sig:       setupTestSignature(),
			wantError: false,
		},
		{
			name:      "nil signature",
			sig:       nil,
			wantError: true,
		},
		{
			name: "high bit set in R",
			sig: &SignatureResult{
				R: new(big.Int).SetBytes([]byte{0x80, 0x01, 0x02, 0x03}),
				S: new(big.Int).SetBytes([]byte{0x04, 0x05, 0x06, 0x07}),
			},
			wantError: false,
		},
		{
			name: "high bit set in S",
			sig: &SignatureResult{
				R: new(big.Int).SetBytes([]byte{0x01, 0x02, 0x03, 0x04}),
				S: new(big.Int).SetBytes([]byte{0x80, 0x05, 0x06, 0x07}),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := FormatDERSignature(tt.sig)
			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, formatted)
			} else {
				require.NoError(t, err)
				require.NotNil(t, formatted)
				require.Equal(t, byte(0x30), formatted[0]) // Check sequence marker

				// Parse back and verify
				r, s, err := UnmarshalDERSignature(formatted)
				require.NoError(t, err)
				require.Equal(t, tt.sig.R.Bytes(), r.Bytes())
				require.Equal(t, tt.sig.S.Bytes(), s.Bytes())
			}
		})
	}
}
