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
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// Wallet returns a specific wallet
func (k Keeper) Wallet(goCtx context.Context, req *types.QueryWalletRequest) (*types.QueryWalletResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	wallet, err := k.GetWallet(ctx, req.Address)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wallet not found: %s", err.Error())
	}

	return &types.QueryWalletResponse{Wallet: wallet}, nil
}

// Wallets returns all wallets
func (k Keeper) Wallets(goCtx context.Context, req *types.QueryWalletsRequest) (*types.QueryWalletsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	wallets, err := k.GetAllWalletsFromStore(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wallets: %s", err.Error())
	}

	// Convert []types.Wallet to []*types.Wallet
	walletPtrs := make([]*types.Wallet, len(wallets))
	for i := range wallets {
		walletPtrs[i] = &wallets[i]
	}

	return &types.QueryWalletsResponse{Wallets: walletPtrs}, nil
}

// PartyData returns TSS party data for a wallet
func (k Keeper) PartyData(goCtx context.Context, req *types.QueryPartyDataRequest) (*types.QueryPartyDataResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	partyData, err := k.GetPartyData(ctx, req.WalletAddress)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "party data not found: %s", err.Error())
	}

	return &types.QueryPartyDataResponse{
		WalletAddress: req.WalletAddress,
		PartyData:     partyData,
	}, nil
}

// KeyRotation returns a specific key rotation by wallet ID and version
func (k Keeper) KeyRotation(goCtx context.Context, req *types.QueryKeyRotationRequest) (*types.QueryKeyRotationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	rotation, err := k.GetKeyRotation(ctx, req.WalletId, req.Version)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "key rotation not found: %s", err.Error())
	}

	return &types.QueryKeyRotationResponse{Rotation: rotation}, nil
}

// KeyRotations returns all key rotations for a wallet
func (k Keeper) KeyRotations(goCtx context.Context, req *types.QueryKeyRotationsRequest) (*types.QueryKeyRotationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyRotationKey))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var rotations []*types.KeyRotation
	for ; iterator.Valid(); iterator.Next() {
		var rotation types.KeyRotation
		err := k.cdc.Unmarshal(iterator.Value(), &rotation)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal key rotation: %s", err.Error())
		}
		if rotation.WalletId == req.WalletId {
			rotations = append(rotations, &rotation)
		}
	}

	return &types.QueryKeyRotationsResponse{Rotations: rotations}, nil
}

// KeyRotationStatus returns the status of a key rotation operation
func (k Keeper) KeyRotationStatus(goCtx context.Context, req *types.QueryKeyRotationStatusRequest) (*types.QueryKeyRotationStatusResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	rotationStatus, err := k.GetKeyRotationStatus(ctx, req.WalletId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "key rotation status not found: %s", err.Error())
	}

	return &types.QueryKeyRotationStatusResponse{
		WalletId:   req.WalletId,
		Status:     rotationStatus.Status.String(),
		Version:    rotationStatus.Version,
		NewPubKey:  rotationStatus.NewPublicKey,
	}, nil
}

// BatchSignStatus returns the status of a batch signing operation
func (k Keeper) BatchSignStatus(goCtx context.Context, req *types.QueryBatchSignStatusRequest) (*types.QueryBatchSignStatusResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	signStatus, err := k.GetBatchSignStatus(ctx, req.WalletId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "batch sign status not found: %s", err.Error())
	}

	return &types.QueryBatchSignStatusResponse{
		WalletId:   req.WalletId,
		BatchId:    req.BatchId,
		Status:     signStatus.Status.String(),
		Signatures: signStatus.Signatures,
	}, nil
}

// ListAuditEvents returns all audit events
func (k Keeper) ListAuditEvents(goCtx context.Context, req *types.QueryListAuditEventsRequest) (*types.QueryListAuditEventsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte("audit/"))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var events []*types.AuditEvent
	for ; iterator.Valid(); iterator.Next() {
		var event types.AuditEvent
		err := k.cdc.Unmarshal(iterator.Value(), &event)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal audit event: %s", err.Error())
		}
		if event.WalletId == req.WalletId && (req.EventType == "" || event.EventType == req.EventType) {
			events = append(events, &event)
		}
	}

	return &types.QueryListAuditEventsResponse{Events: events}, nil
}
