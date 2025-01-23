package keeper

import (
	"context"
	"fmt"

	"selfchain/x/keyless/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RecoverWallet handles the wallet recovery process
func (k Keeper) RecoverWallet(ctx context.Context, msg *types.MsgRecoverWallet) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 1. Verify DID ownership through identity module
	didDoc, found := k.identityKeeper.GetDIDDocument(sdkCtx, msg.Creator)
	if !found {
		return fmt.Errorf("DID document not found: %s", msg.Creator)
	}

	// 2. Verify recovery proof through identity module
	if err := k.identityKeeper.VerifyRecoveryToken(sdkCtx, msg.Creator, msg.RecoveryProof); err != nil {
		return fmt.Errorf("failed to verify recovery proof: %w", err)
	}

	// 3. Get recovery share from identity module and use it for reconstruction
	if _, found := k.identityKeeper.GetKeyShare(sdkCtx, msg.Creator); !found {
		return fmt.Errorf("key share not found for DID: %s", msg.Creator)
	}

	// 4. Reconstruct wallet using recovery share
	wallet, err := k.identityKeeper.ReconstructWallet(sdkCtx, didDoc)
	if err != nil {
		return fmt.Errorf("failed to reconstruct wallet: %w", err)
	}

	// 5. Store reconstructed wallet data in keeper's state
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	blockTime := sdkCtx.BlockTime()
	newWallet := &types.Wallet{
		Id:            msg.WalletAddress,
		Address:       msg.WalletAddress,
		Status:        types.WalletStatusActive,
		SecurityLevel: uint32(SecurityLevelHigh), // Recovered wallets use high security
		Threshold:     2,                         // Default threshold for recovered wallets
		Parties:       3,                         // Default number of parties for recovered wallets
		CreatedAt:     &blockTime,
		UpdatedAt:     &blockTime,
		Metadata: map[string]string{
			"recovery_method": "did_recovery",
			"recovered_by":    msg.Creator,
		},
	}

	// Copy reconstructed wallet data
	if w, ok := wallet.(*types.Wallet); ok {
		newWallet.PublicKey = w.PublicKey
		newWallet.KeyVersion = w.KeyVersion
		newWallet.ChainId = w.ChainId
	}

	walletBytes := k.cdc.MustMarshal(newWallet)
	store.Set([]byte(msg.WalletAddress), walletBytes)

	// 6. Emit recovery event
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWalletRecovered,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyStatus, types.WalletStatusActive),
		),
	)

	return nil
}
