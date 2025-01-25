package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// GetBatchSignStatus retrieves the batch sign status for a wallet
func (k Keeper) GetBatchSignStatus(ctx sdk.Context, walletAddress string) (*types.BatchSignStatusInfo, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.BatchSignStatusStoreKey(walletAddress)
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("batch sign status not found")
	}

	var status types.BatchSignStatusInfo
	err := k.cdc.Unmarshal(bz, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// SetBatchSignStatus stores the batch sign status for a wallet
func (k Keeper) SetBatchSignStatus(ctx sdk.Context, status *types.BatchSignStatusInfo) error {
	store := ctx.KVStore(k.storeKey)
	key := types.BatchSignStatusStoreKey(status.WalletAddress)
	bz, err := k.cdc.Marshal(status)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

// DeleteBatchSignStatus deletes the batch sign status for a wallet
func (k Keeper) DeleteBatchSignStatus(ctx sdk.Context, walletAddress string) {
	store := ctx.KVStore(k.storeKey)
	key := types.BatchSignStatusStoreKey(walletAddress)
	store.Delete(key)
}

// BatchSign handles the batch signing request
func (k Keeper) BatchSign(ctx sdk.Context, msg *types.MsgBatchSignRequest) error {
	// Get wallet
	_, err := k.GetWallet(ctx, msg.WalletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %s", msg.WalletAddress)
	}

	// Create batch sign status
	status := &types.BatchSignStatusInfo{
		WalletAddress: msg.WalletAddress,
		Messages:      msg.Messages,
		Signatures:    make([]string, len(msg.Messages)),
		Status:        types.BatchSignStatus_BATCH_SIGN_STATUS_IN_PROGRESS,
	}

	// Store status
	if err := k.SetBatchSignStatus(ctx, status); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBatchSignRequested,
			sdk.NewAttribute(types.AttributeKeyWalletAddress, msg.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyBatchStatus, status.Status.String()),
		),
	)

	return nil
}
