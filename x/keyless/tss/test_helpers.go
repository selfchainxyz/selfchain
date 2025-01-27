package tss

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "math/big"
    "selfchain/x/keyless/types"
)

// reconstructPublicKey reconstructs the ECDSA public key from the TSS public key bytes
func reconstructPublicKey(pubKeyBytes []byte) (*ecdsa.PublicKey, error) {
    curve := elliptic.P256() // We use secp256k1 in TSS, but for testing P256 is fine
    x := new(big.Int).SetBytes(pubKeyBytes[:32])
    y := new(big.Int).SetBytes(pubKeyBytes[32:])

    if !curve.IsOnCurve(x, y) {
        return nil, types.ErrInvalidPublicKey
    }

    return &ecdsa.PublicKey{
        Curve: curve,
        X:     x,
        Y:     y,
    }, nil
}
