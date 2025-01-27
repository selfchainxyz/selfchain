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
	return prefix.NewStore(store, []byte(types.WalletKey))
}

// SaveWallet saves a wallet to the store
func (k Keeper) SaveWallet(ctx sdk.Context, wallet *types.Wallet) error {
	if wallet == nil {
		return fmt.Errorf("wallet cannot be nil")
	}

	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("wallet/%s", wallet.Id))
	bz, err := k.cdc.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet: %v", err)
	}
	store.Set(key, bz)
	return nil
}

// GetWallet gets a wallet from the store
func (k Keeper) GetWallet(ctx sdk.Context, walletId string) (*types.Wallet, error) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("wallet/%s", walletId))
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	var wallet types.Wallet
	if err := k.cdc.Unmarshal(bz, &wallet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet: %v", err)
	}
	return &wallet, nil
}

// DeleteWallet deletes a wallet from the store
func (k Keeper) DeleteWallet(ctx sdk.Context, walletId string) error {
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf("wallet/%s", walletId))
	if !store.Has(key) {
		return fmt.Errorf("wallet not found")
	}

	store.Delete(key)
	return nil
}

// ListWallets returns all wallets
func (k Keeper) ListWallets(ctx sdk.Context) ([]*types.Wallet, error) {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte("wallet/")
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var wallets []*types.Wallet
	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(iterator.Value(), &wallet); err != nil {
			return nil, fmt.Errorf("failed to unmarshal wallet: %v", err)
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

// GetAllWalletsFromStore returns all wallets from the KVStore
func (k Keeper) GetAllWalletsFromStore(ctx sdk.Context) ([]types.Wallet, error) {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte("wallet/")
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	var wallets []types.Wallet
	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(iterator.Value(), &wallet); err != nil {
			return nil, fmt.Errorf("failed to unmarshal wallet: %v", err)
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// GetWalletStatus returns the current status of a wallet
func (k Keeper) GetWalletStatus(ctx sdk.Context, walletAddress string) (types.WalletStatus, error) {
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return types.WalletStatus_WALLET_STATUS_UNSPECIFIED, fmt.Errorf("wallet not found: %v", err)
	}
	return wallet.Status, nil
}

// SetWalletStatus updates the status of a wallet
func (k Keeper) SetWalletStatus(ctx sdk.Context, walletAddress string, status types.WalletStatus) error {
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %v", err)
	}

	wallet.Status = status
	return k.SaveWallet(ctx, wallet)
}

// SavePartyData stores TSS party data
func (k Keeper) SavePartyData(ctx sdk.Context, walletID string, data *types.PartyData) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PartyDataKey))
	bz, err := k.cdc.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal party data: %w", err)
	}
	store.Set([]byte(walletID), bz)
	return nil
}

// GetPartyData retrieves TSS party data
func (k Keeper) GetPartyData(ctx sdk.Context, walletID string) (*types.PartyData, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PartyDataKey))
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PartyDataKey))
	store.Delete([]byte(walletID))
}

// SaveSigningSession stores a signing session
func (k Keeper) SaveSigningSession(ctx sdk.Context, session *types.SigningSession) error {
	if session == nil {
		return fmt.Errorf("signing session cannot be nil")
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SigningSessionKey))
	bz, err := k.cdc.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal signing session: %w", err)
	}
	store.Set([]byte(session.SessionId), bz)
	return nil
}

// GetSigningSession retrieves a signing session by ID
func (k Keeper) GetSigningSession(ctx sdk.Context, sessionID string) (*types.SigningSession, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SigningSessionKey))
	bz := store.Get([]byte(sessionID))
	if bz == nil {
		return nil, fmt.Errorf("signing session not found: %s", sessionID)
	}

	var session types.SigningSession
	if err := k.cdc.Unmarshal(bz, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signing session: %w", err)
	}
	return &session, nil
}

// DeleteSigningSession removes a signing session
func (k Keeper) DeleteSigningSession(ctx sdk.Context, sessionID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.SigningSessionKey))
	store.Delete([]byte(sessionID))
}
