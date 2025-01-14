package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"selfchain/x/keyless/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Wallet returns information about a specific wallet
func (k Keeper) Wallet(c context.Context, req *types.QueryWalletRequest) (*types.QueryWalletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	bz := store.Get([]byte(req.Address))
	if bz == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	var wallet types.Wallet
	k.cdc.MustUnmarshal(bz, &wallet)
	return &types.QueryWalletResponse{
		Wallet: &wallet,
	}, nil
}

// Wallets returns all wallets with pagination
func (k Keeper) Wallets(c context.Context, req *types.QueryWalletsRequest) (*types.QueryWalletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	walletStore := prefix.NewStore(store, types.KeyPrefix(types.WalletKey))

	var wallets []*types.Wallet
	pageRes, err := k.Paginate(walletStore, req.Pagination, func(key []byte, value []byte) error {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(value, &wallet); err != nil {
			return err
		}

		wallets = append(wallets, &wallet)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWalletsResponse{
		Wallets:    wallets,
		Pagination: pageRes,
	}, nil
}

// PartyData returns TSS party data for a wallet
func (k Keeper) PartyData(c context.Context, req *types.QueryPartyDataRequest) (*types.QueryPartyDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.WalletAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	partyDataStore := prefix.NewStore(store, types.KeyPrefix(types.PartyDataKey))

	bz := partyDataStore.Get([]byte(req.WalletAddress))
	if bz == nil {
		return nil, status.Error(codes.NotFound, "party data not found")
	}

	var partyData types.PartyData
	if err := k.cdc.Unmarshal(bz, &partyData); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPartyDataResponse{
		WalletAddress: req.WalletAddress,
		PartyData:    &partyData,
	}, nil
}
