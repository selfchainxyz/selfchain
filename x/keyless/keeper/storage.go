package keeper

import (
	"fmt"

	"selfchain/x/keyless/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetWalletStore returns the store for wallet data
func (k Keeper) GetWalletStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, []byte(walletPrefix))
}

// SaveWallet stores a wallet
func (k Keeper) SaveWallet(ctx sdk.Context, wallet *types.Wallet) error {
	store := k.GetWalletStore(ctx)
	bz, err := k.cdc.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet: %w", err)
	}
	store.Set([]byte(wallet.Id), bz)
	return nil
}

// GetWallet retrieves a wallet
func (k Keeper) GetWallet(ctx sdk.Context, id string) (*types.Wallet, error) {
	store := k.GetWalletStore(ctx)
	bz := store.Get([]byte(id))
	if bz == nil {
		return nil, fmt.Errorf("wallet not found: %s", id)
	}

	var wallet types.Wallet
	if err := k.cdc.Unmarshal(bz, &wallet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet: %w", err)
	}
	return &wallet, nil
}

// DeleteWallet removes a wallet
func (k Keeper) DeleteWallet(ctx sdk.Context, id string) {
	store := k.GetWalletStore(ctx)
	store.Delete([]byte(id))
}

// SavePartyData stores TSS party data
func (k Keeper) SavePartyData(ctx sdk.Context, walletID string, data *types.PartyData) error {
	store := k.GetPartyDataStore(ctx)
	bz, err := k.cdc.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal party data: %w", err)
	}
	store.Set([]byte(walletID), bz)
	return nil
}

// GetPartyData retrieves TSS party data
func (k Keeper) GetPartyData(ctx sdk.Context, walletID string) (*types.PartyData, error) {
	store := k.GetPartyDataStore(ctx)
	bz := store.Get([]byte(walletID))
	if bz == nil {
		return nil, fmt.Errorf("party data not found: %s", walletID)
	}

	var data types.PartyData
	if err := k.cdc.Unmarshal(bz, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal party data: %w", err)
	}
	return &data, nil
}

// DeletePartyData removes TSS party data
func (k Keeper) DeletePartyData(ctx sdk.Context, walletID string) {
	store := k.GetPartyDataStore(ctx)
	store.Delete([]byte(walletID))
}
