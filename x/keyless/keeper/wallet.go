package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// IsWalletAuthorized checks if a grantee is authorized for a specific permission
func (k Keeper) IsWalletAuthorized(ctx sdk.Context, address string, grantee string, requiredPerm types.WalletPermission) bool {
	// For new wallets, creator is always authorized
	wallet, err := k.GetWallet(ctx, address)
	if err != nil {
		return false
	}
	if wallet == nil {
		return false
	}
	if wallet.Creator == grantee {
		return true
	}

	// Check if grantee has the required permission
	perm, err := k.GetPermission(ctx, address, grantee)
	if err != nil {
		return false
	}
	if perm == nil {
		return false
	}

	// Convert WalletPermission enum to string
	permStr := requiredPerm.String()
	return perm.HasPermission(permStr)
}

// ValidateWalletAccess validates access to a wallet for a specific operation
func (k Keeper) ValidateWalletAccess(ctx sdk.Context, walletAddress string, operation string) error {
	// Get wallet
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return fmt.Errorf("wallet not found: %s", walletAddress)
	}

	// Check wallet status
	if wallet.Status != types.WalletStatus_WALLET_STATUS_ACTIVE {
		return fmt.Errorf("wallet is not active: current status %s", wallet.Status)
	}

	// Map operation to required permission
	var requiredPerm types.WalletPermission
	switch operation {
	case "recovery":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_RECOVER
	case "sign":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_SIGN
	case "rotate":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_ROTATE
	case "admin":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_ADMIN
	default:
		return fmt.Errorf("unknown operation: %s", operation)
	}

	// Check if wallet is authorized for the operation
	if !k.IsWalletAuthorized(ctx, walletAddress, wallet.Creator, requiredPerm) {
		return fmt.Errorf("not authorized for operation: %s", operation)
	}

	return nil
}
