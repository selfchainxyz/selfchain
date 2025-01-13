package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"selfchain/x/keyless/types"
)

var _ types.QueryServer = Keeper{}

// GetWallet returns a wallet by its address
func (k Keeper) GetWallet(goCtx context.Context, req *types.QueryGetWalletRequest) (*types.QueryGetWalletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	wallet, found := k.getWallet(ctx, req.WalletAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	return &types.QueryGetWalletResponse{
		Wallet: wallet,
	}, nil
}

// getWallet returns a wallet by its address
func (k Keeper) getWallet(ctx sdk.Context, address string) (types.Wallet, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	
	bz := store.Get([]byte(address))
	if bz == nil {
		return types.Wallet{}, false
	}

	var wallet types.Wallet
	k.cdc.MustUnmarshal(bz, &wallet)
	return wallet, true
}

// ListWallets returns all wallets
func (k Keeper) ListWallets(goCtx context.Context, req *types.QueryListWalletsRequest) (*types.QueryListWalletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	wallets := k.GetAllWallets(ctx)

	return &types.QueryListWalletsResponse{
		Wallets: wallets,
	}, nil
}

// GetAllWallets returns all wallets
func (k Keeper) GetAllWallets(ctx sdk.Context) []types.Wallet {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	wallets := []types.Wallet{}
	for ; iterator.Valid(); iterator.Next() {
		var wallet types.Wallet
		k.cdc.MustUnmarshal(iterator.Value(), &wallet)
		wallets = append(wallets, wallet)
	}

	return wallets
}

// Params returns the module parameters
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
