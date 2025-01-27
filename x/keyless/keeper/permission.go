package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"selfchain/x/keyless/types"
)

// GrantPermission grants permissions to a grantee for a specific wallet
func (k Keeper) GrantPermission(ctx sdk.Context, grant *types.MsgGrantPermission) (*types.Permission, error) {
	// Check if permission already exists
	existingPerm, err := k.GetPermission(ctx, grant.WalletId, grant.Grantee)
	if err == nil && existingPerm != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidPermission, "permission already exists")
	}

	// Create new permission
	now := time.Now().UTC()
	permission := &types.Permission{
		WalletId:    grant.WalletId,
		Grantee:     grant.Grantee,
		Permissions: make([]string, len(grant.Permissions)),
		GrantedAt:   &now,
		ExpiresAt:   grant.ExpiresAt,
		Revoked:     false,
	}

	// Convert WalletPermission enum values to strings
	for i, p := range grant.Permissions {
		permission.Permissions[i] = types.WalletPermission_name[int32(p)]
	}

	// Store the permission
	k.SetPermission(ctx, permission)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"permission_granted",
			sdk.NewAttribute("wallet_id", grant.WalletId),
			sdk.NewAttribute("grantee", grant.Grantee),
			sdk.NewAttribute("permissions", fmt.Sprintf("%v", permission.Permissions)),
		),
	)

	return permission, nil
}

// RevokePermission revokes permissions from a grantee for a specific wallet
func (k Keeper) RevokePermission(ctx sdk.Context, revoke *types.MsgRevokePermission) error {
	// Check if permission exists
	perm, err := k.GetPermission(ctx, revoke.WalletId, revoke.Grantee)
	if err != nil {
		return err
	}

	// Mark as revoked
	now := time.Now().UTC()
	perm.Revoked = true
	perm.RevokedAt = &now

	// Update the permission
	k.SetPermission(ctx, perm)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"permission_revoked",
			sdk.NewAttribute("wallet_id", revoke.WalletId),
			sdk.NewAttribute("grantee", revoke.Grantee),
		),
	)

	return nil
}

// HasPermission checks if a grantee has a specific permission for a wallet
func (k Keeper) HasPermission(ctx sdk.Context, walletId string, grantee string, permission types.WalletPermission) bool {
	perm, err := k.GetPermission(ctx, walletId, grantee)
	if err != nil {
		return false
	}

	// Check if revoked
	if perm.Revoked {
		return false
	}

	// Check expiry
	if perm.ExpiresAt != nil && perm.ExpiresAt.Before(time.Now().UTC()) {
		return false
	}

	// Convert WalletPermission to string
	permStr := types.WalletPermission_name[int32(permission)]

	// Check if the permission is in the list
	for _, p := range perm.Permissions {
		if p == permStr {
			return true
		}
	}

	return false
}

// GetPermission gets a permission by wallet ID and grantee
func (k Keeper) GetPermission(ctx sdk.Context, walletId string, grantee string) (*types.Permission, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPermissionKey(walletId, grantee)
	
	bz := store.Get(key)
	if bz == nil {
		return nil, types.ErrPermissionNotFound
	}

	var permission types.Permission
	k.cdc.MustUnmarshal(bz, &permission)
	return &permission, nil
}

// SetPermission sets a permission in the store
func (k Keeper) SetPermission(ctx sdk.Context, permission *types.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPermissionKey(permission.WalletId, permission.Grantee)
	bz := k.cdc.MustMarshal(permission)
	store.Set(key, bz)
}

// DeletePermission deletes a permission from the store
func (k Keeper) DeletePermission(ctx sdk.Context, permission *types.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPermissionKey(permission.WalletId, permission.Grantee)
	store.Delete(key)
}

// GetAllPermissions gets all permissions from the store
func (k Keeper) GetAllPermissions(ctx sdk.Context) ([]*types.Permission, error) {
	var permissions []*types.Permission
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.KeyPrefix+"permission/"))

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var permission types.Permission
		k.cdc.MustUnmarshal(iterator.Value(), &permission)
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

// GetPermissionsForWallet gets all permissions for a specific wallet
func (k Keeper) GetPermissionsForWallet(ctx sdk.Context, walletId string) ([]*types.Permission, error) {
	var permissions []*types.Permission
	prefixBytes := types.GetPermissionPrefix(walletId)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixBytes)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var permission types.Permission
		k.cdc.MustUnmarshal(iterator.Value(), &permission)
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

// ValidatePermissions validates a list of permissions
func (k Keeper) ValidatePermissions(permissions []types.WalletPermission) error {
	for _, perm := range permissions {
		if perm == types.WalletPermission_WALLET_PERMISSION_UNSPECIFIED {
			return types.ErrInvalidPermission
		}
	}
	return nil
}
