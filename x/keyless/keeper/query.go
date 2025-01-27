package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"selfchain/x/keyless/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// Wallet returns a wallet by address
func (k Keeper) Wallet(goCtx context.Context, req *types.QueryWalletRequest) (*types.QueryWalletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	wallet, err := k.GetWallet(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.QueryWalletResponse{Wallet: wallet}, nil
}

// Wallets returns all wallets
func (k Keeper) Wallets(goCtx context.Context, req *types.QueryWalletsRequest) (*types.QueryWalletsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	wallets, err := k.ListWallets(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryWalletsResponse{
		Wallets:    wallets,
		Pagination: nil,
	}, nil
}

// PartyData returns party data for a wallet
func (k Keeper) PartyData(goCtx context.Context, req *types.QueryPartyDataRequest) (*types.QueryPartyDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	partyData, err := k.GetPartyData(ctx, req.WalletAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryPartyDataResponse{
		WalletAddress: req.WalletAddress,
		PartyData:     partyData,
	}, nil
}

// KeyRotation returns a specific key rotation
func (k Keeper) KeyRotation(goCtx context.Context, req *types.QueryKeyRotationRequest) (*types.QueryKeyRotationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	rotation, err := k.GetKeyRotation(ctx, req.WalletId, req.Version)
	if err != nil {
		return nil, err
	}

	return &types.QueryKeyRotationResponse{Rotation: rotation}, nil
}

// KeyRotations returns all key rotations for a wallet
func (k Keeper) KeyRotations(goCtx context.Context, req *types.QueryKeyRotationsRequest) (*types.QueryKeyRotationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	rotations := k.GetKeyRotationsForWallet(ctx, req.WalletId)

	return &types.QueryKeyRotationsResponse{
		Rotations:  rotations,
		Pagination: nil,
	}, nil
}

// KeyRotationStatus returns the status of a key rotation
func (k Keeper) KeyRotationStatus(goCtx context.Context, req *types.QueryKeyRotationStatusRequest) (*types.QueryKeyRotationStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	status, err := k.GetKeyRotationStatus(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}

	return &types.QueryKeyRotationStatusResponse{
		WalletId:   req.WalletId,
		Status:     status.GetStatus().String(),
		Version:    status.GetVersion(),
		NewPubKey:  status.GetNewPublicKey(),
		Error:      "",
	}, nil
}

// BatchSignStatus returns the status of a batch signing operation
func (k Keeper) BatchSignStatus(goCtx context.Context, req *types.QueryBatchSignStatusRequest) (*types.QueryBatchSignStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	status, err := k.GetBatchSignStatus(ctx, req.BatchId)
	if err != nil {
		return nil, err
	}

	return &types.QueryBatchSignStatusResponse{
		WalletId:   req.WalletId,
		BatchId:    req.BatchId,
		Status:     status.GetStatus().String(),
		Signatures: status.GetSignatures(),
		Error:      "",
	}, nil
}

// ListAuditEvents returns audit events for a wallet
func (k Keeper) ListAuditEvents(goCtx context.Context, req *types.QueryListAuditEventsRequest) (*types.QueryListAuditEventsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	events := k.GetAuditEventsForWallet(ctx, req.WalletId)

	return &types.QueryListAuditEventsResponse{
		Events:     events,
		Pagination: nil,
	}, nil
}

// Permissions returns all permissions for a wallet
func (k Keeper) Permissions(goCtx context.Context, req *types.QueryPermissionsRequest) (*types.QueryPermissionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	permissions, err := k.GetPermissionsForWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}

	return &types.QueryPermissionsResponse{
		Permissions: permissions,
		Pagination: nil,
	}, nil
}

// Permission returns a specific permission by wallet ID and grantee
func (k Keeper) Permission(goCtx context.Context, req *types.QueryPermissionRequest) (*types.QueryPermissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	permission, err := k.GetPermission(ctx, req.WalletId, req.Grantee)
	if err != nil {
		return nil, err
	}

	return &types.QueryPermissionResponse{Permission: permission}, nil
}
