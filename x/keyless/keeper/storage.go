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

// SaveWallet stores a wallet
func (k Keeper) SaveWallet(ctx sdk.Context, wallet *types.Wallet) error {
	if wallet == nil {
		return fmt.Errorf("wallet cannot be nil")
	}

	// Validate required fields
	if wallet.WalletAddress == "" {
		return fmt.Errorf("wallet address cannot be empty")
	}
	if wallet.ChainId == "" {
		return fmt.Errorf("chain ID cannot be empty")
	}
	if wallet.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}

	// Check for existing wallet
	existingWallet, err := k.GetWallet(ctx, wallet.WalletAddress)
	if err == nil && existingWallet != nil {
		// Only allow updates to existing wallet if:
		// 1. The creator matches
		// 2. The wallet is in recovery
		if existingWallet.Creator != wallet.Creator {
			return fmt.Errorf("cannot update wallet: creator mismatch")
		}
		if existingWallet.Status != types.WalletStatus_WALLET_STATUS_INACTIVE {
			return fmt.Errorf("wallet with address %s already exists", wallet.WalletAddress)
		}
	}

	store := k.GetWalletStore(ctx)
	bz, err := k.cdc.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet: %w", err)
	}
	store.Set([]byte(wallet.WalletAddress), bz)
	return nil
}

// GetWallet retrieves a wallet
func (k Keeper) GetWallet(ctx sdk.Context, walletAddress string) (*types.Wallet, error) {
	store := k.GetWalletStore(ctx)
	bz := store.Get([]byte(walletAddress))
	if bz == nil {
		return nil, fmt.Errorf("wallet not found: %s", walletAddress)
	}

	var wallet types.Wallet
	if err := k.cdc.Unmarshal(bz, &wallet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet: %w", err)
	}
	return &wallet, nil
}

// DeleteWallet removes a wallet
func (k Keeper) DeleteWallet(ctx sdk.Context, walletAddress string) {
	store := k.GetWalletStore(ctx)
	store.Delete([]byte(walletAddress))
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
