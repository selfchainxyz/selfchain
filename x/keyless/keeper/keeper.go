package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cometbft/cometbft/libs/log"
	"selfchain/x/keyless/types"
	identitytypes "selfchain/x/identity/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		memKey        storetypes.StoreKey
		paramstore    paramtypes.Subspace
		identityKeeper types.IdentityKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	identityKeeper types.IdentityKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		identityKeeper: identityKeeper,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetPartyDataStore returns the store for TSS party data
func (k Keeper) GetPartyDataStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, []byte("party-data/"))
}

// GetKeyShare retrieves a key share for a DID
func (k Keeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("key-share/"))
	bz := store.Get([]byte(did))
	if bz == nil {
		return nil, false
	}
	return bz, true
}

// StoreKeyShare stores a key share for a DID
func (k Keeper) StoreKeyShare(ctx sdk.Context, did string, keyShare []byte) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("key-share/"))
	store.Set([]byte(did), keyShare)
	return nil
}

// ReconstructWallet reconstructs a wallet from a DID document
func (k Keeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) ([]byte, error) {
	// TODO: Implement wallet reconstruction logic
	return nil, nil
}

// InitiateRecovery initiates the wallet recovery process
func (k Keeper) InitiateRecovery(ctx sdk.Context, did string, recoveryToken string, recoveryAddress string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("recovery/"))
	recoveryInfo := &types.RecoveryInfo{
		Did:             did,
		RecoveryToken:   recoveryToken,
		RecoveryAddress: recoveryAddress,
		Status:          types.RecoveryStatus_PENDING,
		CreatedAt:       ctx.BlockTime().Unix(),
	}

	bz, err := k.cdc.Marshal(recoveryInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal recovery info: %w", err)
	}

	store.Set([]byte(did), bz)
	return nil
}

// GetIdentityKeeper returns the identity keeper
func (k Keeper) GetIdentityKeeper() types.IdentityKeeper {
	return k.identityKeeper
}
