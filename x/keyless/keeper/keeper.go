package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	identitytypes "selfchain/x/identity/types"
	"selfchain/x/keyless/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		memKey         storetypes.StoreKey
		paramstore     paramtypes.Subspace
		identityKeeper types.IdentityKeeper
		tssProtocol    types.TSSProtocol
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	identityKeeper types.IdentityKeeper,
	tssProtocol types.TSSProtocol,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	k := &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		paramstore:     ps,
		identityKeeper: identityKeeper,
		tssProtocol:    tssProtocol,
	}

	return k
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetTSSProtocol sets the TSS protocol implementation
func (k *Keeper) SetTSSProtocol(protocol types.TSSProtocol) {
	k.tssProtocol = protocol
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
	// Extract shares from DID document
	var shares [][]byte
	for _, service := range didDoc.Service {
		if service.Type == "KeyRecoveryShare" {
			shares = append(shares, []byte(service.ServiceEndpoint))
		}
	}

	if len(shares) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInsufficientShares, "no recovery shares found in DID document")
	}

	// Get params
	params := k.GetParams(ctx)

	// Validate number of shares
	if len(shares) < int(params.MinRecoveryThreshold) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientShares, "got %d shares, need %d", len(shares), params.MinRecoveryThreshold)
	}

	// Reconstruct key
	reconstructedKey, err := k.tssProtocol.ReconstructKey(ctx, shares)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrKeyReconstruction, err.Error())
	}

	return reconstructedKey, nil
}

// ReconstructWalletFromDID reconstructs a wallet from a DID document
func (k Keeper) ReconstructWalletFromDID(ctx sdk.Context, didDoc identitytypes.DIDDocument) ([]byte, error) {
	// Extract shares from DID document
	var shares [][]byte
	for _, service := range didDoc.Service {
		if service.Type == "KeyRecoveryShare" {
			shares = append(shares, []byte(service.ServiceEndpoint))
		}
	}

	if len(shares) == 0 {
		return nil, sdkerrors.Wrap(types.ErrInsufficientShares, "no recovery shares found in DID document")
	}

	// Get params
	params := k.GetParams(ctx)

	// Validate number of shares
	if len(shares) < int(params.MinRecoveryThreshold) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientShares, "got %d shares, need %d", len(shares), params.MinRecoveryThreshold)
	}

	// Reconstruct key
	reconstructedKey, err := k.tssProtocol.ReconstructKey(ctx, shares)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrKeyReconstruction, err.Error())
	}

	return reconstructedKey, nil
}

// InitiateRecovery initiates the wallet recovery process
func (k Keeper) InitiateRecovery(ctx sdk.Context, walletAddress string, recoveryToken string, recoveryAddress string) error {
	// Validate inputs
	if len(walletAddress) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidWalletID, "wallet address cannot be empty")
	}

	if len(recoveryToken) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidMetadata, "recovery token cannot be empty")
	}

	if len(recoveryAddress) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidMetadata, "recovery address cannot be empty")
	}

	// Get wallet
	wallet, found := k.GetWallet(ctx, walletAddress)
	if !found {
		return sdkerrors.Wrapf(types.ErrWalletNotFound, "wallet not found: %s", walletAddress)
	}

	// Verify wallet is not already in recovery
	if wallet.Status == types.WalletStatus_WALLET_STATUS_RECOVERING {
		return sdkerrors.Wrap(types.ErrWalletBusy, "wallet is already in recovery")
	}

	// Create recovery info
	now := ctx.BlockTime()
	recoveryInfo := &types.RecoveryInfo{
		Did:             walletAddress,
		RecoveryToken:   recoveryToken,
		RecoveryAddress: recoveryAddress,
		Status:          types.RecoveryStatus_RECOVERY_STATUS_PENDING,
		CreatedAt:       now,
	}

	// Save recovery info
	k.SetRecoveryInfo(ctx, recoveryInfo)

	// Update wallet status
	wallet.Status = types.WalletStatus_WALLET_STATUS_RECOVERING
	wallet.UpdatedAt = &now
	k.SetWallet(ctx, wallet)

	return nil
}

// SubmitRecoveryShare submits a recovery share for a wallet
func (k Keeper) SubmitRecoveryShare(ctx sdk.Context, walletAddress string, shareData []byte, signature []byte) error {
	// Validate inputs
	if len(walletAddress) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidWalletID, "wallet address cannot be empty")
	}

	if len(shareData) == 0 {
		return sdkerrors.Wrap(types.ErrInvalidShare, "share data cannot be empty")
	}

	// Get wallet
	wallet, found := k.GetWallet(ctx, walletAddress)
	if !found {
		return sdkerrors.Wrapf(types.ErrWalletNotFound, "wallet not found: %s", walletAddress)
	}

	// Verify wallet is in recovery
	if wallet.Status != types.WalletStatus_WALLET_STATUS_RECOVERING {
		return sdkerrors.Wrap(types.ErrInvalidStatus, "wallet is not in recovery")
	}

	// Verify signature
	pubKey, err := k.GetPublicKey(ctx, walletAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to get public key")
	}

	err = k.tssProtocol.VerifySignature(ctx, pubKey, shareData, signature)
	if err != nil {
		return sdkerrors.Wrap(types.ErrInvalidSignature, err.Error())
	}

	// Store share data
	now := ctx.BlockTime()
	share := &types.ShareData{
		WalletAddress: walletAddress,
		Creator:       sdk.AccAddress(ctx.BlockHeader().ProposerAddress).String(),
		ShareData:     shareData,
		Signature:     signature,
		CreatedAt:     now,
	}

	k.SetShareData(ctx, share)

	return nil
}

// GetPublicKey returns the public key for a given wallet
func (k Keeper) GetPublicKey(ctx sdk.Context, walletAddress string) ([]byte, error) {
	wallet, found := k.GetWallet(ctx, walletAddress)
	if !found {
		return nil, fmt.Errorf("wallet not found: %s", walletAddress)
	}
	return []byte(wallet.PublicKey), nil
}

// SetWallet saves a wallet to the store
func (k Keeper) SetWallet(ctx sdk.Context, wallet *types.Wallet) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.WalletKey))
	key := []byte(wallet.WalletAddress)
	bz := k.cdc.MustMarshal(wallet)
	store.Set(key, bz)
}

// SetRecoveryInfo saves recovery info to the store
func (k Keeper) SetRecoveryInfo(ctx sdk.Context, info *types.RecoveryInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyPrefixRecoveryInfo))
	key := []byte(info.Did)
	bz := k.cdc.MustMarshal(info)
	store.Set(key, bz)
}

// SetShareData saves share data to the store
func (k Keeper) SetShareData(ctx sdk.Context, share *types.ShareData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyShareKey))
	key := []byte(share.WalletAddress)
	bz := k.cdc.MustMarshal(share)
	store.Set(key, bz)
}

// GetIdentityKeeper returns the identity keeper
func (k Keeper) GetIdentityKeeper() types.IdentityKeeper {
	return k.identityKeeper
}
