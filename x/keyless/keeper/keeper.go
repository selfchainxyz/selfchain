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

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetPartyDataStore returns the store for TSS party data
func (k Keeper) GetPartyDataStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(partyDataPrefix))
}

// GetKeyShare retrieves a key share for a DID
func (k Keeper) GetKeyShare(ctx sdk.Context, did string) ([]byte, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(walletPrefix))
	keyShare := store.Get([]byte(did))
	if keyShare == nil {
		return nil, false
	}
	return keyShare, true
}

// StoreKeyShare stores a key share for a DID
func (k Keeper) StoreKeyShare(ctx sdk.Context, did string, keyShare []byte) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(walletPrefix))
	store.Set([]byte(did), keyShare)
	return nil
}

// ReconstructWallet reconstructs a wallet from a DID document
func (k Keeper) ReconstructWallet(ctx sdk.Context, didDoc identitytypes.DIDDocument) ([]byte, error) {
	// For now, just return the key share associated with the DID
	keyShare, found := k.GetKeyShare(ctx, didDoc.Id)
	if !found {
		return nil, fmt.Errorf("no key share found for DID %s", didDoc.Id)
	}
	return keyShare, nil
}

// InitiateRecovery initiates the wallet recovery process
func (k Keeper) InitiateRecovery(ctx sdk.Context, did string, recoveryToken string, recoveryAddress string) error {
	// Store recovery information
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("recovery"))
	recoveryInfo := types.RecoveryInfo{
		Did:             did,
		RecoveryToken:   recoveryToken,
		RecoveryAddress: recoveryAddress,
		Status:          types.RecoveryStatus_PENDING,
		CreatedAt:       ctx.BlockTime().Unix(),
	}
	
	bz, err := k.cdc.Marshal(&recoveryInfo)
	if err != nil {
		return err
	}
	
	store.Set([]byte(did), bz)
	return nil
}
