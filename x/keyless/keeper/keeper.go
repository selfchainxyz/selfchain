package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"selfchain/x/keyless/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
	}
)

const (
	// Store prefixes
	walletPrefix    = "wallet"
	partyDataPrefix = "party_data"
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetPartyDataStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(partyDataPrefix))
}

// getWallet returns a wallet from store
func (k Keeper) getWallet(ctx sdk.Context, address string) (val types.Wallet, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(walletPrefix))

	b := store.Get([]byte(address))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// setWallet sets a wallet in store
func (k Keeper) setWallet(ctx sdk.Context, wallet types.Wallet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(walletPrefix))
	b := k.cdc.MustMarshal(&wallet)
	store.Set([]byte(wallet.WalletAddress), b)
}

// ValidateWalletOwner validates if the given creator is the owner of the wallet
func (k Keeper) ValidateWalletOwner(ctx sdk.Context, address string, creator string) error {
	wallet, found := k.getWallet(ctx, address)
	if !found {
		return fmt.Errorf("wallet not found: %s", address)
	}

	if wallet.Creator != creator {
		return fmt.Errorf("unauthorized: %s is not the owner of wallet %s", creator, address)
	}

	return nil
}
