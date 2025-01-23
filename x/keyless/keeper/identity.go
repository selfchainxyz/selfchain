package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"selfchain/x/identity/types"
)

// SecurityLevel represents the security level for wallet operations
type SecurityLevel uint32

const (
	SecurityLevelLow SecurityLevel = iota + 1
	SecurityLevelMedium
	SecurityLevelHigh
)

// VerifyIdentity verifies a user's identity using DID and OAuth2
func (k Keeper) VerifyIdentity(ctx context.Context, did string, oauthToken string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 1. Verify DID existence and ownership
	_, found := k.identityKeeper.GetDIDDocument(sdkCtx, did)
	if !found {
		return fmt.Errorf("DID document not found: %s", did)
	}

	// 2. Verify OAuth2 token
	if err := k.identityKeeper.VerifyOAuth2Token(sdkCtx, did, oauthToken); err != nil {
		return fmt.Errorf("failed to verify OAuth2 token: %w", err)
	}

	// 3. Check if MFA is required
	if k.IsMFARequired(sdkCtx, did) {
		if err := k.identityKeeper.VerifyMFA(sdkCtx, did); err != nil {
			return fmt.Errorf("MFA verification failed: %w", err)
		}
	}

	// 4. Check rate limits
	if err := k.identityKeeper.CheckRateLimit(sdkCtx, did, "identity_verification"); err != nil {
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	// 5. Log audit event
	if err := k.identityKeeper.LogAuditEvent(sdkCtx, &types.AuditEvent{
		Did:       did,
		EventType: "identity_verification",
		Success:   true,
		Details:   fmt.Sprintf("OAuth2 verification successful, MFA required: %v", k.IsMFARequired(sdkCtx, did)),
	}); err != nil {
		return fmt.Errorf("failed to log audit event: %w", err)
	}

	return nil
}

// IsMFARequired checks if MFA is required for a given DID
func (k Keeper) IsMFARequired(ctx sdk.Context, did string) bool {
	// Get security level from params
	params := k.GetParams(ctx)
	return params.MaxSecurityLevel >= uint32(SecurityLevelHigh)
}
