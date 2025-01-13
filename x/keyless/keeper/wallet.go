package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// setWallet stores a wallet
func (k Keeper) setWallet(ctx sdk.Context, wallet types.Wallet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	b := k.cdc.MustMarshal(&wallet)
	store.Set([]byte(wallet.WalletAddress), b)
}

// DeleteWallet removes a wallet
func (k Keeper) DeleteWallet(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	store.Delete([]byte(address))
}

// ValidateWalletOwner checks if the given creator is the owner of the wallet
func (k Keeper) ValidateWalletOwner(ctx sdk.Context, address string, creator string) error {
	wallet, found := k.getWallet(ctx, address)
	if !found {
		return fmt.Errorf("wallet not found: %s", address)
	}

	if wallet.Creator != creator {
		return fmt.Errorf("unauthorized: %s is not the creator of wallet %s", creator, address)
	}

	return nil
}
