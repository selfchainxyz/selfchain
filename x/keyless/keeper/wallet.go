package keeper

import (
    "crypto/rand"
    "encoding/hex"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "selfchain/x/keyless/types"
)

// CreateWallet creates a new keyless wallet and stores it in state
func (k Keeper) CreateWallet(ctx sdk.Context, creator string, did string) (types.Wallet, error) {
    // Generate a new address for the wallet
    address, err := k.generateWalletAddress()
    if err != nil {
        return types.Wallet{}, err
    }

    // Check if wallet with this address already exists
    if k.HasWallet(ctx, address) {
        return types.Wallet{}, types.ErrWalletExists
    }

    // Create a new wallet
    wallet := types.NewWallet(address, did, creator)

    // Store the wallet in state
    k.SetWallet(ctx, wallet)

    return wallet, nil
}

// HasWallet checks if a wallet exists
func (k Keeper) HasWallet(ctx sdk.Context, address string) bool {
    store := ctx.KVStore(k.storeKey)
    return store.Has(types.WalletKey(address))
}

// SetWallet stores a wallet in state
func (k Keeper) SetWallet(ctx sdk.Context, wallet types.Wallet) {
    store := ctx.KVStore(k.storeKey)
    b := k.cdc.MustMarshal(&wallet)
    store.Set(types.WalletKey(wallet.Address), b)
}

// GetWalletState retrieves a wallet from state
func (k Keeper) GetWalletState(ctx sdk.Context, address string) (types.Wallet, error) {
    store := ctx.KVStore(k.storeKey)
    b := store.Get(types.WalletKey(address))
    if b == nil {
        return types.Wallet{}, types.ErrWalletNotFound
    }

    var wallet types.Wallet
    k.cdc.MustUnmarshal(b, &wallet)
    return wallet, nil
}

// GetWalletStateByDID retrieves a wallet by its DID
func (k Keeper) GetWalletStateByDID(ctx sdk.Context, did string) (types.Wallet, error) {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, types.WalletKeyPrefix)
    defer iterator.Close()

    for ; iterator.Valid(); iterator.Next() {
        var wallet types.Wallet
        k.cdc.MustUnmarshal(iterator.Value(), &wallet)
        if wallet.Did == did {
            return wallet, nil
        }
    }

    return types.Wallet{}, types.ErrWalletNotFound
}

// generateWalletAddress generates a new unique wallet address
func (k Keeper) generateWalletAddress() (string, error) {
    // Generate 20 random bytes for the address
    b := make([]byte, 20)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }

    // Convert to hex string with "kw" prefix (keyless wallet)
    return "kw" + hex.EncodeToString(b), nil
}
