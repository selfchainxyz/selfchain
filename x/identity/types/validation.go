package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateCredentialProof validates a credential proof
func ValidateCredentialProof(proof *CredentialProof) error {
	if proof == nil {
		return sdkerrors.Wrap(ErrInvalidProof, "proof cannot be nil")
	}

	if proof.Type == "" {
		return sdkerrors.Wrap(ErrInvalidProof, "proof type cannot be empty")
	}

	if proof.VerificationMethod == "" {
		return sdkerrors.Wrap(ErrInvalidProof, "verification method cannot be empty")
	}

	if proof.ProofPurpose == "" {
		return sdkerrors.Wrap(ErrInvalidProof, "proof purpose cannot be empty")
	}

	if proof.Created == 0 {
		return sdkerrors.Wrap(ErrInvalidProof, "proof creation time cannot be empty")
	}

	return nil
}
