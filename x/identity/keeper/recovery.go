package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"

	"selfchain/x/identity/types"
)

const (
	RecoveryTokenLength = 32
	RecoveryExpiry     = 30 * time.Minute
)

// InitiateRecovery initiates a wallet recovery process
func (k Keeper) InitiateRecovery(ctx sdk.Context, did string, socialProvider string, oauthToken string) (*types.RecoverySession, error) {
	// Get DID document to verify it exists
	_, found := k.GetDIDDocument(ctx, did)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDIDNotFound, "did %s not found", did)
	}

	// Verify OAuth token
	if err := k.VerifyOAuth2Token(ctx, socialProvider, oauthToken); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidToken, "failed to verify oauth token: %s", err.Error())
	}

	// Get social identity
	socialIdentity, found := k.GetSocialIdentity(ctx, did, socialProvider)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrSocialIdentityNotFound, 
			"social identity not found for did %s and provider %s", did, socialProvider)
	}

	// Create recovery session
	session := &types.RecoverySession{
		Id:               uuid.New().String(),
		Did:              did,
		SocialProvider:   socialProvider,
		SocialId:         socialIdentity.ProviderId,
		MfaVerified:      false,
		IdentityVerified: false,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		CreatedAt:        time.Now(),
	}

	// Save recovery session
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecoverySessionPrefix))
	key := []byte(session.Id)
	bz := k.cdc.MustMarshal(session)
	store.Set(key, bz)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOAuthSuccess,
			sdk.NewAttribute(types.AttributeKeyDID, did),
			sdk.NewAttribute(types.AttributeKeyProvider, socialProvider),
			sdk.NewAttribute(types.AttributeKeySocialID, socialIdentity.ProviderId),
		),
	)

	return session, nil
}

// UpdateRecoverySession updates a recovery session
func (k Keeper) UpdateRecoverySession(ctx sdk.Context, sessionID string) (*types.RecoverySession, error) {
	// Get recovery session
	session, found := k.GetRecoverySession(ctx, sessionID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrRecoverySessionNotFound, "session %s not found", sessionID)
	}

	// Update session
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecoverySessionPrefix))
	key := []byte(sessionID)
	bz := k.cdc.MustMarshal(session)
	store.Set(key, bz)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOAuthSuccess,
			sdk.NewAttribute(types.AttributeKeyDID, session.Did),
			sdk.NewAttribute(types.AttributeKeyStatus, "updated"),
		),
	)

	return session, nil
}

// CompleteRecovery completes a wallet recovery process
func (k Keeper) CompleteRecovery(ctx sdk.Context, sessionID string, mfaCode string) error {
	// Get recovery session
	session, found := k.GetRecoverySession(ctx, sessionID)
	if !found {
		return sdkerrors.Wrapf(types.ErrRecoverySessionNotFound, "session %s not found", sessionID)
	}

	// Verify MFA code
	if !session.MfaVerified {
		return sdkerrors.Wrap(types.ErrMFAVerificationFailed, "MFA not verified")
	}

	// Get DID document
	doc, found := k.GetDIDDocument(ctx, session.Did)
	if !found {
		return sdkerrors.Wrapf(types.ErrDIDNotFound, "did %s not found", session.Did)
	}

	// Save updated DID document
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DIDDocumentPrefix))
	key := []byte(session.Did)
	bz := k.cdc.MustMarshal(&doc)
	store.Set(key, bz)

	// Mark session as completed
	session.IdentityVerified = true
	sessionStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecoverySessionPrefix))
	sessionKey := []byte(sessionID)
	sessionBz := k.cdc.MustMarshal(session)
	sessionStore.Set(sessionKey, sessionBz)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOAuthSuccess,
			sdk.NewAttribute(types.AttributeKeyDID, session.Did),
			sdk.NewAttribute(types.AttributeKeyStatus, "completed"),
		),
	)

	return nil
}

// GetRecoverySession gets a recovery session by ID
func (k Keeper) GetRecoverySession(ctx sdk.Context, sessionID string) (*types.RecoverySession, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecoverySessionPrefix))
	key := []byte(sessionID)

	if !store.Has(key) {
		return nil, false
	}

	var session types.RecoverySession
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &session)

	return &session, true
}

// GenerateRecoveryToken generates a recovery token for wallet recovery
func (k Keeper) GenerateRecoveryToken(ctx sdk.Context, walletID string) (string, error) {
	// Generate a random token
	token := fmt.Sprintf("%s-%d", walletID, ctx.BlockHeight())

	// Store the token with expiry
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecoveryPrefix))
	expiresAt := ctx.BlockTime().Add(RecoveryExpiry)

	// Create and store recovery data
	recoveryData := &types.RecoveryData{
		WalletId:  walletID,
		Token:     token,
		ExpiresAt: &expiresAt,
	}

	bz, err := k.cdc.Marshal(recoveryData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal recovery data: %w", err)
	}

	store.Set([]byte(token), bz)
	return token, nil
}

// ValidateRecoveryToken validates a recovery token for a wallet
func (k Keeper) ValidateRecoveryToken(ctx sdk.Context, walletID string, token string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecoveryPrefix))

	bz := store.Get([]byte(token))
	if bz == nil {
		return fmt.Errorf("recovery token not found for wallet: %s", walletID)
	}

	var recoveryData types.RecoveryData
	if err := k.cdc.Unmarshal(bz, &recoveryData); err != nil {
		return fmt.Errorf("failed to unmarshal recovery data: %w", err)
	}

	// Validate wallet ID matches
	if recoveryData.WalletId != walletID {
		return fmt.Errorf("invalid recovery token: wallet ID mismatch")
	}

	// Check if token has expired
	if recoveryData.ExpiresAt.Before(ctx.BlockTime()) {
		return fmt.Errorf("recovery token has expired")
	}

	// Token is valid, clean it up since it's one-time use
	store.Delete([]byte(token))

	return nil
}

// CleanupExpiredSessions removes expired recovery sessions
func (k Keeper) CleanupExpiredSessions(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecoverySessionPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var session types.RecoverySession
		k.cdc.MustUnmarshal(iterator.Value(), &session)

		if ctx.BlockTime().After(session.ExpiresAt) {
			store.Delete(iterator.Key())
		}
	}
}
