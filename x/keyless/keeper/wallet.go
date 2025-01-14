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
	store.Set([]byte(wallet.Address), b)
}

// GetWalletFromStore returns a wallet from the store
func (k Keeper) GetWalletFromStore(ctx sdk.Context, address string) (types.Wallet, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	b := store.Get([]byte(address))
	if b == nil {
		return types.Wallet{}, false
	}

	var wallet types.Wallet
	k.cdc.MustUnmarshal(b, &wallet)
	return wallet, true
}

// DeleteWallet removes a wallet from the store
func (k Keeper) DeleteWallet(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	store.Delete([]byte(address))
}

// ValidateWalletOwner validates if the given creator is the owner of the wallet
func (k Keeper) ValidateWalletOwner(ctx sdk.Context, address string, creator string) error {
	wallet, found := k.GetWalletFromStore(ctx, address)
	if !found {
		return fmt.Errorf("wallet not found: %s", address)
	}

	if wallet.Creator != creator {
		return fmt.Errorf("not the wallet owner")
	}

	return nil
}

// ValidateRecoveryProof validates a recovery proof for a wallet
func (k Keeper) ValidateRecoveryProof(ctx sdk.Context, walletAddress, recoveryProof, signature string) error {
	// TODO: Implement recovery proof validation
	return nil
}

// SignWithTSS signs a transaction using TSS
func (k Keeper) SignWithTSS(ctx sdk.Context, wallet *types.Wallet, unsignedTx string) (string, error) {
	// TODO: Implement TSS signing
	return "", nil
}
