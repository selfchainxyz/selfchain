package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// SetWallet sets a wallet in the store
func (k Keeper) SetWallet(ctx sdk.Context, wallet *types.Wallet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	b := k.cdc.MustMarshal(wallet)
	store.Set([]byte(wallet.Id), b)
}

// GetWalletFromStore returns a wallet from the store
func (k Keeper) GetWalletFromStore(ctx sdk.Context, id string) (types.Wallet, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	b := store.Get([]byte(id))
	if b == nil {
		return types.Wallet{}, false
	}

	var wallet types.Wallet
	k.cdc.MustUnmarshal(b, &wallet)
	return wallet, true
}

// DeleteWallet removes a wallet from the store
func (k Keeper) DeleteWallet(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	store.Delete([]byte(id))
}

// ValidateWalletAccess validates if the permission exists in the wallet's permissions list
func (k Keeper) ValidateWalletAccess(ctx sdk.Context, id string, permission string) error {
	wallet, found := k.GetWalletFromStore(ctx, id)
	if !found {
		return fmt.Errorf("wallet not found: %s", id)
	}

	for _, p := range wallet.Permissions {
		if p == permission {
			return nil
		}
	}

	return fmt.Errorf("permission denied: %s", permission)
}

// ValidateRecoveryProof validates a recovery proof for a wallet
func (k Keeper) ValidateRecoveryProof(ctx sdk.Context, walletId, recoveryProof string) error {
	// TODO: Implement recovery proof validation
	// This should integrate with the identity module for verification
	return nil
}

// SignWithTSS signs a transaction using TSS
func (k Keeper) SignWithTSS(ctx sdk.Context, wallet *types.Wallet, unsignedTx string) ([]byte, error) {
	// TODO: Implement TSS signing
	// This should:
	// 1. Validate security level and threshold requirements
	// 2. Coordinate with parties to perform distributed signing
	// 3. Verify the signature before returning
	return nil, nil
}
