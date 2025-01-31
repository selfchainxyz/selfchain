package keeper

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"selfchain/x/keyless/types"
)

// GrantPermission grants a permission to a grantee for a specific wallet
func (k Keeper) GrantPermission(ctx sdk.Context, grant *types.Permission) error {
	// Validate the permission first
	if err := k.ValidateAndGrantPermission(ctx, grant); err != nil {
		return err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePermissionGranted,
			sdk.NewAttribute("wallet_address", grant.WalletAddress),
			sdk.NewAttribute("grantee", grant.Grantee),
			sdk.NewAttribute("permissions", strings.Join(grant.Permissions, ",")),
		),
	)

	return nil
}

// RevokePermission revokes a permission from a grantee for a specific wallet
func (k Keeper) RevokePermission(ctx sdk.Context, walletAddress string, grantee string) error {
	// Basic validation
	if walletAddress == "" {
		return status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}

	if grantee == "" {
		return status.Error(codes.InvalidArgument, "grantee cannot be empty")
	}

	// Check if wallet exists
	if _, err := k.GetWallet(ctx, walletAddress); err != nil {
		return err
	}

	// Check if permission exists
	existingPerm, err := k.GetPermission(ctx, walletAddress, grantee)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return status.Error(codes.NotFound, "permission does not exist")
		}
		return err
	}

	// Check if already revoked
	if existingPerm.IsRevoked() {
		return status.Error(codes.FailedPrecondition, "permission is already revoked")
	}

	// Delete the permission instead of storing it with revoked status
	k.DeletePermission(ctx, existingPerm)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePermissionRevoked,
			sdk.NewAttribute("wallet_address", walletAddress),
			sdk.NewAttribute("grantee", grantee),
			sdk.NewAttribute("permissions", strings.Join(existingPerm.Permissions, ",")),
		),
	)

	return nil
}

// HasPermission checks if a grantee has a specific permission for a wallet
func (k Keeper) HasPermission(ctx sdk.Context, walletAddress string, grantee string, permission string) (bool, error) {
	// Basic validation
	if walletAddress == "" || grantee == "" {
		return false, status.Error(codes.InvalidArgument, "wallet address and grantee cannot be empty")
	}

	// Check if wallet exists
	if _, err := k.GetWallet(ctx, walletAddress); err != nil {
		return false, nil
	}

	perm, err := k.GetPermission(ctx, walletAddress, grantee)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}

	// Check if revoked
	if perm.IsRevoked() {
		return false, nil
	}

	// Check if expired
	if perm.IsExpired() {
		return false, nil
	}

	// Check if permission exists
	return perm.HasPermission(permission), nil
}

// GetPermission gets a permission by wallet address and grantee
func (k Keeper) GetPermission(ctx sdk.Context, walletAddress string, grantee string) (*types.Permission, error) {
	if walletAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}

	if grantee == "" {
		return nil, status.Error(codes.InvalidArgument, "grantee cannot be empty")
	}

	store := ctx.KVStore(k.storeKey)
	permissionStore := prefix.NewStore(store, []byte(types.PermissionKey))
	key := k.GetPermissionKey(walletAddress, grantee)

	bz := permissionStore.Get(key)
	if bz == nil {
		return nil, status.Error(codes.NotFound, "permission not found")
	}

	var permission types.Permission
	if err := k.cdc.Unmarshal(bz, &permission); err != nil {
		return nil, err
	}

	return &permission, nil
}

// GetPermissionKey returns the key for storing a permission
func (k Keeper) GetPermissionKey(walletAddress string, grantee string) []byte {
	return []byte(fmt.Sprintf("%s/%s", walletAddress, grantee))
}

// StorePermission stores a permission in the store
func (k Keeper) StorePermission(ctx sdk.Context, permission *types.Permission) {
	if permission == nil {
		panic("cannot store nil permission")
	}

	if permission.WalletAddress == "" {
		panic("cannot store permission with empty wallet address")
	}

	if permission.Grantee == "" {
		panic("cannot store permission with empty grantee")
	}

	store := ctx.KVStore(k.storeKey)
	permissionStore := prefix.NewStore(store, []byte(types.PermissionKey))
	key := k.GetPermissionKey(permission.WalletAddress, permission.Grantee)

	bz, err := k.cdc.Marshal(permission)
	if err != nil {
		panic(err)
	}

	permissionStore.Set(key, bz)
}

// DeletePermission deletes a permission from the store
func (k Keeper) DeletePermission(ctx sdk.Context, permission *types.Permission) {
	store := ctx.KVStore(k.storeKey)
	permissionStore := prefix.NewStore(store, []byte(types.PermissionKey))
	key := k.GetPermissionKey(permission.WalletAddress, permission.Grantee)
	permissionStore.Delete(key)
}

// GetAllPermissions gets all permissions from the store
func (k Keeper) GetAllPermissions(ctx sdk.Context) ([]*types.Permission, error) {
	store := ctx.KVStore(k.storeKey)
	permissionStore := prefix.NewStore(store, []byte(types.PermissionKey))

	var permissions []*types.Permission
	iterator := permissionStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var permission types.Permission
		if err := k.cdc.Unmarshal(iterator.Value(), &permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

// GetPermissionsForWallet gets all permissions for a specific wallet
func (k Keeper) GetPermissionsForWallet(ctx sdk.Context, walletAddress string) ([]*types.Permission, error) {
	if walletAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "wallet address cannot be empty")
	}

	store := ctx.KVStore(k.storeKey)
	permissionStore := prefix.NewStore(store, []byte(types.PermissionKey))
	prefix := []byte(fmt.Sprintf("%s/", walletAddress))

	var permissions []*types.Permission
	iterator := sdk.KVStorePrefixIterator(permissionStore, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var permission types.Permission
		if err := k.cdc.Unmarshal(iterator.Value(), &permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

// ValidateAndGrantPermission validates and grants a permission
func (k Keeper) ValidateAndGrantPermission(ctx sdk.Context, permission *types.Permission) error {
	// Basic validation
	if err := permission.Validate(); err != nil {
		return err
	}

	// Check if wallet exists
	if _, err := k.GetWallet(ctx, permission.WalletAddress); err != nil {
		return err
	}

	// Check if permission already exists
	existingPerm, err := k.GetPermission(ctx, permission.WalletAddress, permission.Grantee)
	if err != nil && status.Code(err) != codes.NotFound {
		return err
	}

	if existingPerm != nil {
		// Check if existing permission is revoked or expired
		if !existingPerm.IsRevoked() && !existingPerm.IsExpired() {
			return status.Error(codes.AlreadyExists, "permission already exists and is still valid")
		}
	}

	// Store the permission
	k.StorePermission(ctx, permission)

	return nil
}

// CheckPermission checks if a grantee has a specific permission for a wallet
func (k Keeper) CheckPermission(ctx sdk.Context, walletAddress string, grantee string, requiredPermission string) error {
	// Check if permission exists and is valid
	hasPermission, err := k.HasPermission(ctx, walletAddress, grantee, requiredPermission)
	if err != nil {
		return err
	}

	if !hasPermission {
		return status.Errorf(codes.PermissionDenied, "grantee %s does not have %s permission for wallet %s",
			grantee, requiredPermission, walletAddress)
	}

	return nil
}
