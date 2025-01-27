package keeper

import (
	"fmt"
	"time"

	"selfchain/x/identity/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt/v4"
)

const (
	OAuthSessionPrefix = "oauth_session/"
	MFASessionPrefix   = "mfa_session/"
	KeySharePrefix     = "key_share/"
	AuditEventPrefix   = "audit_event/"
)

// VerifyDIDOwnership verifies that an address owns a DID
func (k Keeper) VerifyDIDOwnership(ctx sdk.Context, did string, owner sdk.AccAddress) error {
	doc, found := k.GetDIDDocument(ctx, did)
	if !found {
		return fmt.Errorf("DID document not found: %s", did)
	}

	if !doc.HasController(owner.String()) {
		return fmt.Errorf("address %s is not the controller of DID %s", owner.String(), did)
	}

	return nil
}

// VerifyOAuth2Token verifies an OAuth2 token for a DID
func (k Keeper) VerifyOAuth2Token(ctx sdk.Context, did string, token string) error {
	claims := &types.OAuthClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return k.GetOAuthPublicKey(ctx, did)
	})
	if err != nil {
		return fmt.Errorf("failed to verify OAuth2 token: %v", err)
	}

	// Verify claims
	if claims.Subject == "" || claims.Issuer == "" {
		return fmt.Errorf("invalid OAuth2 token claims")
	}

	return nil
}

// VerifyMFA verifies MFA for a DID
func (k Keeper) VerifyMFA(ctx sdk.Context, did string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(MFASessionPrefix))
	sessionBytes := store.Get([]byte(did))
	if sessionBytes == nil {
		return fmt.Errorf("MFA session not found for DID: %s", did)
	}

	var session types.MFASession
	k.cdc.MustUnmarshal(sessionBytes, &session)

	if session.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("MFA session expired")
	}

	if session.Status != types.MFASessionStatus_MFA_SESSION_STATUS_VERIFIED {
		return fmt.Errorf("MFA not verified")
	}

	return nil
}

// VerifyRecoveryToken verifies a recovery token for a DID
func (k Keeper) VerifyRecoveryToken(ctx sdk.Context, did string, token string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RecoveryPrefix))
	sessionBytes := store.Get([]byte(did))
	if sessionBytes == nil {
		return fmt.Errorf("recovery session not found for DID: %s", did)
	}

	var session types.RecoverySession
	k.cdc.MustUnmarshal(sessionBytes, &session)

	// Verify session state
	if !session.IdentityVerified {
		return fmt.Errorf("identity not verified for recovery session")
	}

	if !session.MfaVerified {
		return fmt.Errorf("MFA not verified for recovery session")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("recovery session expired")
	}

	return nil
}

// GetKeyShare returns a key share for a DID
func (k Keeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(KeySharePrefix))
	share := store.Get([]byte(did))
	return share, share != nil
}

// LogAuditEvent logs an audit event
func (k Keeper) LogAuditEvent(ctx sdk.Context, event *types.AuditEvent) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(AuditEventPrefix))
	key := []byte(fmt.Sprintf("%s/%d", event.Did, ctx.BlockHeight()))
	store.Set(key, k.cdc.MustMarshal(event))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAuditLog,
			sdk.NewAttribute(types.AttributeKeyDID, event.Did),
			sdk.NewAttribute(types.AttributeKeyEventType, event.EventType),
			sdk.NewAttribute(types.AttributeKeySuccess, fmt.Sprintf("%t", event.Success)),
		),
	)

	return nil
}

// ReconstructWallet reconstructs a wallet from a DID document using TSS
func (k Keeper) ReconstructWallet(ctx sdk.Context, didDoc types.DIDDocument) (interface{}, error) {
	// Get key share
	share, found := k.GetKeyShare(ctx, didDoc.Id)
	if !found {
		return nil, fmt.Errorf("key share not found for DID: %s", didDoc.Id)
	}

	// Log audit event
	k.LogAuditEvent(ctx, &types.AuditEvent{
		Did:       didDoc.Id,
		EventType: "wallet_reconstruction",
		Success:   true,
		Details:   "Wallet reconstruction initiated",
	})

	// TODO: Implement TSS-based wallet reconstruction
	// 1. Coordinate with other parties
	// 2. Run TSS protocol to reconstruct the key
	// 3. Generate new key shares
	// 4. Store new shares and update metadata

	// For now, return the share as a placeholder
	return share, nil
}
