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
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WalletKey))
	bz := store.Get([]byte(req.Id))
	if bz == nil {
		return nil, status.Error(codes.NotFound, "wallet not found")
	}

	var wallet types.Wallet
	k.cdc.MustUnmarshal(bz, &wallet)
	return &types.QueryWalletResponse{
		Wallet: wallet,
	}, nil
}

// ListWallets returns all wallets with pagination
func (k Keeper) ListWallets(c context.Context, req *types.QueryListWalletsRequest) (*types.QueryListWalletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	walletStore := prefix.NewStore(store, types.KeyPrefix(types.WalletKey))

	var wallets []types.Wallet
	pageRes, err := k.Paginate(walletStore, req.Pagination, func(key []byte, value []byte) error {
		var wallet types.Wallet
		if err := k.cdc.Unmarshal(value, &wallet); err != nil {
			return err
		}

		wallets = append(wallets, wallet)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListWalletsResponse{
		Wallets:    wallets,
		Pagination: pageRes,
	}, nil
}

// KeyRotation returns information about a specific key rotation
func (k Keeper) KeyRotation(c context.Context, req *types.QueryKeyRotationRequest) (*types.QueryKeyRotationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.KeyRotationKey))
	key := append([]byte(req.WalletId), sdk.Uint64ToBigEndian(req.Version)...)
	bz := store.Get(key)
	if bz == nil {
		return nil, status.Error(codes.NotFound, "key rotation not found")
	}

	var rotation types.KeyRotation
	k.cdc.MustUnmarshal(bz, &rotation)
	return &types.QueryKeyRotationResponse{
		Rotation: rotation,
	}, nil
}

// KeyRotations returns all key rotations for a wallet with pagination
func (k Keeper) KeyRotations(c context.Context, req *types.QueryKeyRotationsRequest) (*types.QueryKeyRotationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	rotationStore := prefix.NewStore(store, types.KeyPrefix(types.KeyRotationKey))

	var rotations []types.KeyRotation
	pageRes, err := k.Paginate(rotationStore, req.Pagination, func(key []byte, value []byte) error {
		var rotation types.KeyRotation
		if err := k.cdc.Unmarshal(value, &rotation); err != nil {
			return err
		}

		rotations = append(rotations, rotation)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryKeyRotationsResponse{
		Rotations:  rotations,
		Pagination: pageRes,
	}, nil
}

// ListAuditEvents returns all audit events with pagination
func (k Keeper) ListAuditEvents(c context.Context, req *types.QueryListAuditEventsRequest) (*types.QueryListAuditEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.WalletId == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	eventStore := prefix.NewStore(store, types.KeyPrefix(types.AuditEventKey))

	var events []types.AuditEvent
	pageRes, err := k.Paginate(eventStore, req.Pagination, func(key []byte, value []byte) error {
		var event types.AuditEvent
		if err := k.cdc.Unmarshal(value, &event); err != nil {
			return err
		}

		// Filter by wallet ID, event type, and success
		if event.WalletId != req.WalletId {
			return nil
		}
		if req.EventType != "" && event.EventType != req.EventType {
			return nil
		}
		if event.Success != req.Success {
			return nil
		}

		events = append(events, event)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListAuditEventsResponse{
		Events:     events,
		Pagination: pageRes,
	}, nil
}
