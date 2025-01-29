package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"selfchain/x/keyless/types"
)

// GetWalletStore returns the store for wallet data
func (k Keeper) GetWalletStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, []byte(types.WalletKey))
}

// GetWalletKey returns the key for a wallet
func (k Keeper) GetWalletKey(walletAddress string) []byte {
	return []byte("wallet/" + walletAddress)
}

// SaveWallet saves a wallet to the store
func (k Keeper) SaveWallet(ctx sdk.Context, wallet *types.Wallet) error {
	// Basic validation
	if wallet.WalletAddress == "" {
		return status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}
	if wallet.Creator == "" {
		return status.Error(codes.InvalidArgument, "creator cannot be empty")
	}
	if wallet.ChainId == "" {
		return status.Error(codes.InvalidArgument, "chain ID cannot be empty")
	}

	// Check for duplicate wallet
	existingWallet, err := k.GetWallet(ctx, wallet.WalletAddress)
	if err == nil && existingWallet != nil {
		return status.Error(codes.AlreadyExists, "wallet already exists")
	}

	store := k.GetWalletStore(ctx)
	key := k.GetWalletKey(wallet.WalletAddress)
	bz := k.cdc.MustMarshal(wallet)
	store.Set(key, bz)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWalletCreated,
			sdk.NewAttribute("wallet_address", wallet.WalletAddress),
			sdk.NewAttribute("creator", wallet.Creator),
			sdk.NewAttribute("chain_id", wallet.ChainId),
		),
	)

	return nil
}

// GetWallet retrieves a wallet from the store
func (k Keeper) GetWallet(ctx sdk.Context, walletAddress string) (*types.Wallet, error) {
	if walletAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}

	store := k.GetWalletStore(ctx)
	key := k.GetWalletKey(walletAddress)
	bz := store.Get(key)
	if bz == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("wallet not found: %s", walletAddress))
	}

	var wallet types.Wallet
	if err := k.cdc.Unmarshal(bz, &wallet); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal wallet: %v", err))
	}

	return &wallet, nil
}

// DeleteWallet deletes a wallet from the store
func (k Keeper) DeleteWallet(ctx sdk.Context, walletAddress string) error {
	if walletAddress == "" {
		return status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}

	store := k.GetWalletStore(ctx)
	key := k.GetWalletKey(walletAddress)

	if !store.Has(key) {
		return status.Error(codes.NotFound, fmt.Sprintf("wallet not found: %s", walletAddress))
	}

	store.Delete(key)
	return nil
}

// GetWalletByCreator retrieves a wallet by its creator address
func (k Keeper) GetWalletByCreator(ctx sdk.Context, creator string) ([]*types.Wallet, error) {
	if creator == "" {
		return nil, status.Error(codes.InvalidArgument, "creator address cannot be empty")
	}

	var wallets []*types.Wallet
	store := k.GetWalletStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(iterator.Value(), &wallet); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal wallet: %v", err))
		}

		if wallet.Creator == creator {
			wallets = append(wallets, &wallet)
		}
	}

	return wallets, nil
}

// GetWalletsByChainId retrieves all wallets for a specific chain ID
func (k Keeper) GetWalletsByChainId(ctx sdk.Context, chainId string) ([]*types.Wallet, error) {
	if chainId == "" {
		return nil, status.Error(codes.InvalidArgument, "chain ID cannot be empty")
	}

	var wallets []*types.Wallet
	store := k.GetWalletStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(iterator.Value(), &wallet); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal wallet: %v", err))
		}

		if wallet.ChainId == chainId {
			wallets = append(wallets, &wallet)
		}
	}

	return wallets, nil
}

// ListWallets returns all wallets
func (k Keeper) ListWallets(ctx sdk.Context) ([]*types.Wallet, error) {
	store := k.GetWalletStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var wallets []*types.Wallet
	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(iterator.Value(), &wallet); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal wallet: %v", err))
		}
		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

// GetAllWalletsFromStore returns all wallets from the KVStore
func (k Keeper) GetAllWalletsFromStore(ctx sdk.Context) ([]types.Wallet, error) {
	store := k.GetWalletStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var wallets []types.Wallet
	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(iterator.Value(), &wallet); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal wallet: %v", err))
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
