package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/sign/schnorr"
	"go.dedis.ch/kyber/v3/util/random"
)

// ValidateBasic performs basic validation of the ZK proof
func (p *ZKProof) ValidateBasic() error {
	if p == nil {
		return sdkerrors.Register(ModuleName, 1600, "proof cannot be nil")
	}

	if p.Type == "" {
		return sdkerrors.Register(ModuleName, 1601, "invalid proof type")
	}
	if len(p.ProofData) == 0 {
		return sdkerrors.Register(ModuleName, 1602, "proof data cannot be empty")
	}
	if p.ClaimsHash == "" {
		return sdkerrors.Register(ModuleName, 1603, "claims hash cannot be empty")
	}
	if len(p.DisclosedIndices) == 0 {
		return sdkerrors.Register(ModuleName, 1604, "must disclose at least one claim")
	}
	if p.Created == 0 {
		return sdkerrors.Register(ModuleName, 1605, "proof creation time cannot be zero")
	}
	if p.VerificationKey == "" {
		return sdkerrors.Register(ModuleName, 1606, "verification key cannot be empty")
	}
	return nil
}

// ValidateBasic performs basic validation of the presentation proof
func (p *ZKPresentationProof) ValidateBasic() error {
	if p == nil {
		return sdkerrors.Register(ModuleName, 1700, "presentation proof cannot be nil")
	}

	if p.Created == 0 {
		return sdkerrors.Register(ModuleName, 1701, "presentation proof creation time cannot be zero")
	}

	if p.ZkProof == nil {
		return sdkerrors.Register(ModuleName, 1702, "zero-knowledge proof cannot be nil")
	}

	if err := p.ZkProof.ValidateBasic(); err != nil {
		return sdkerrors.Wrapf(err, "invalid zero-knowledge proof")
	}

	return nil
}

// ComputeClaimsHash computes a hash of the credential claims
func ComputeClaimsHash(claims map[string]string) string {
	// Sort claims by key to ensure consistent hashing
	keys := make([]string, 0, len(claims))
	for k := range claims {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var claimStr string
	for _, k := range keys {
		claimStr += fmt.Sprintf("%s:%s;", k, claims[k])
	}

	hash := sha256.Sum256([]byte(claimStr))
	return hex.EncodeToString(hash[:])
}

// GenerateZKProof generates a zero-knowledge proof for selective disclosure
func GenerateZKProof(claims map[string]string, disclosedClaims []string, verificationKey string) (*ZKProof, error) {
	// Use Edwards25519 curve
	suite := edwards25519.NewBlakeSHA256Ed25519()

	// Create a new random generator
	rng := random.New()

	// Sort claims for consistent ordering
	keys := make([]string, 0, len(claims))
	for k := range claims {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create a map to track which claims are disclosed
	disclosedIndicesMap := make(map[string]uint32)
	for _, dc := range disclosedClaims {
		for i, k := range keys {
			if dc == k {
				disclosedIndicesMap[k] = uint32(i)
			}
		}
	}

	// Generate private/public key pair
	private := suite.Scalar().Pick(rng)
	public := suite.Point().Mul(private, nil)

	// Create commitment for each claim
	commitments := make([][]byte, 0)
	for _, k := range keys {
		// Convert claim value to bytes
		value := []byte(claims[k])

		// Create Schnorr signature as a commitment
		sig, err := schnorr.Sign(suite, private, value)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to create commitment for claim %s", k)
		}
		commitments = append(commitments, sig)
	}

	// Combine all commitments
	var proofData []byte
	for i, commitment := range commitments {
		if _, disclosed := disclosedIndicesMap[keys[i]]; !disclosed {
			proofData = append(proofData, commitment...)
		}
	}

	// Create disclosed indices array
	disclosedIndices := make([]uint32, 0)
	for _, idx := range disclosedIndicesMap {
		disclosedIndices = append(disclosedIndices, idx)
	}
	sort.Slice(disclosedIndices, func(i, j int) bool {
		return disclosedIndices[i] < disclosedIndices[j]
	})

	// Create metadata
	metadata := make(map[string]string)
	metadata["version"] = "1.0"
	metadata["curve"] = "edwards25519"
	metadata["hash"] = "blake2b"
	metadata["commitment_count"] = fmt.Sprintf("%d", len(commitments))

	// Marshal public key
	pubBytes, err := public.MarshalBinary()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to marshal public key")
	}
	metadata["public_key"] = hex.EncodeToString(pubBytes)

	now := time.Now().Unix()
	return &ZKProof{
		Type:             "Schnorr",
		ProofData:        proofData,
		ClaimsHash:       ComputeClaimsHash(claims),
		DisclosedIndices: disclosedIndices,
		Created:          now,
		VerificationKey:  verificationKey,
		Metadata:         metadata,
	}, nil
}

// VerifyZKProof verifies a zero-knowledge proof
func VerifyZKProof(proof *ZKProof, disclosedClaims map[string]string) (bool, error) {
	if err := proof.ValidateBasic(); err != nil {
		return false, err
	}

	// Verify the proof is not expired (24 hours validity)
	now := time.Now().Unix()
	if now-proof.Created > 24*60*60 {
		return false, sdkerrors.Register(ModuleName, 1620, "proof has expired")
	}

	// Use Edwards25519 curve
	suite := edwards25519.NewBlakeSHA256Ed25519()

	// Get public key from metadata
	publicKeyBytes, err := hex.DecodeString(proof.Metadata["public_key"])
	if err != nil {
		return false, sdkerrors.Wrapf(err, "failed to decode public key")
	}

	public := suite.Point()
	if err := public.UnmarshalBinary(publicKeyBytes); err != nil {
		return false, sdkerrors.Wrapf(err, "failed to unmarshal public key")
	}

	// Verify each commitment
	sigSize := 64 // Size of each Schnorr signature
	for i := 0; i < len(proof.ProofData); i += sigSize {
		sig := proof.ProofData[i : i+sigSize]
		msg := []byte(fmt.Sprintf("claim_%d", i/sigSize))
		err := schnorr.Verify(suite, public, msg, sig)
		if err != nil {
			return false, sdkerrors.Wrapf(err, "failed to verify commitment %d", i/sigSize)
		}
	}

	// Verify disclosed claims hash
	if disclosedClaimsHash := ComputeClaimsHash(disclosedClaims); disclosedClaimsHash != proof.ClaimsHash {
		return false, sdkerrors.Register(ModuleName, 1622, "disclosed claims hash mismatch")
	}

	return true, nil
}
