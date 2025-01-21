package keeper

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"selfchain/x/identity/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RecoverySessionPrefix = "recovery_session/"
	RecoveryExpiry        = 30 * time.Minute
)

// InitiateRecovery starts a new wallet recovery session
func (k Keeper) InitiateRecovery(ctx sdk.Context, socialProvider string, oauthToken string) (*types.RecoverySession, error) {
	// Verify OAuth token and get social ID
	socialID, err := k.VerifyOAuthToken(ctx, socialProvider, oauthToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OAuth token: %v", err)
	}

	// Find social identity
	identity, found := k.GetSocialIdentityBySocialID(ctx, socialProvider, socialID)
	if !found {
		return nil, fmt.Errorf("no identity found for social ID: %s", socialID)
	}

	// Generate session ID
	sessionID := make([]byte, 16)
	if _, err := rand.Read(sessionID); err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %v", err)
	}

	session := &types.RecoverySession{
		Id:               base32.StdEncoding.EncodeToString(sessionID),
		Did:              identity.Did,
		SocialProvider:   socialProvider,
		SocialId:         socialID,
		MfaVerified:      false,
		IdentityVerified: true,
		RecoveryData:     nil, // Will be set by keyless module
		ExpiresAt:        ctx.BlockTime().Add(RecoveryExpiry),
		CreatedAt:        ctx.BlockTime(),
	}

	// Store session
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RecoverySessionPrefix))
	sessionBytes := k.cdc.MustMarshal(session)
	store.Set([]byte(session.Id), sessionBytes)

	return session, nil
}

// GetRecoverySession retrieves a recovery session
func (k Keeper) GetRecoverySession(ctx sdk.Context, sessionID string) (types.RecoverySession, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RecoverySessionPrefix))
	sessionBytes := store.Get([]byte(sessionID))
	if sessionBytes == nil {
		return types.RecoverySession{}, false
	}

	var session types.RecoverySession
	k.cdc.MustUnmarshal(sessionBytes, &session)
	return session, true
}

// CompleteRecovery completes a wallet recovery session
func (k Keeper) CompleteRecovery(ctx sdk.Context, sessionID string, mfaCode string) error {
	session, found := k.GetRecoverySession(ctx, sessionID)
	if !found {
		return fmt.Errorf("recovery session not found: %s", sessionID)
	}

	if ctx.BlockTime().After(session.ExpiresAt) {
		return fmt.Errorf("recovery session expired")
	}

	// Check if MFA is required
	mfaConfig, hasMFA := k.GetMFAConfig(ctx, session.Did)
	if hasMFA && len(mfaConfig.Methods) > 0 {
		// Create MFA challenge if not verified
		if !session.MfaVerified {
			_, err := k.CreateMFAChallenge(ctx, session.Did, "totp")
			if err != nil {
				return fmt.Errorf("failed to create MFA challenge: %v", err)
			}

			if err := k.VerifyMFAChallenge(ctx, session.Did, "totp", mfaCode); err != nil {
				return fmt.Errorf("MFA verification failed: %v", err)
			}

			session.MfaVerified = true
		}
	}

	// Get DID document
	didDoc, found := k.GetDIDDocument(ctx, session.Did)
	if !found {
		return fmt.Errorf("DID document not found: %s", session.Did)
	}

	// Call keyless module to reconstruct wallet
	recoveryData, err := k.keyless.ReconstructWallet(ctx, didDoc)
	if err != nil {
		return fmt.Errorf("failed to reconstruct wallet: %v", err)
	}

	session.RecoveryData = recoveryData

	// Update session
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RecoverySessionPrefix))
	sessionBytes := k.cdc.MustMarshal(&session)
	store.Set([]byte(session.Id), sessionBytes)

	return nil
}

// CleanupExpiredSessions removes expired recovery sessions
func (k Keeper) CleanupExpiredSessions(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(RecoverySessionPrefix))
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
