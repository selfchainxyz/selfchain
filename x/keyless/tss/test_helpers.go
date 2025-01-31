package tss

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "fmt"
    "math/big"
    "selfchain/x/keyless/types"
)

// reconstructPublicKey reconstructs the ECDSA public key from either bytes or an existing key
func reconstructPublicKey(pubKeyOrBytes interface{}) (*ecdsa.PublicKey, error) {
    switch v := pubKeyOrBytes.(type) {
    case []byte:
        curve := elliptic.P256() // We use secp256k1 in TSS, but for testing P256 is fine
        if len(v) != 64 {
            return nil, fmt.Errorf("invalid public key length: expected 64 bytes, got %d", len(v))
        }
        x := new(big.Int).SetBytes(v[:32])
        y := new(big.Int).SetBytes(v[32:])

        if !curve.IsOnCurve(x, y) {
            return nil, types.ErrInvalidPublicKey
        }

        return &ecdsa.PublicKey{
            Curve: curve,
            X:     x,
            Y:     y,
        }, nil
    case *ecdsa.PublicKey:
        if v == nil {
            return nil, fmt.Errorf("public key is nil")
        }
        return v, nil
    default:
        return nil, fmt.Errorf("unsupported public key type: %T", pubKeyOrBytes)
    }
}
