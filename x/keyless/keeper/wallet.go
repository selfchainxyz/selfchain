package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"selfchain/x/keyless/types"
)

// IsWalletAuthorized checks if a grantee is authorized for a specific permission
func (k Keeper) IsWalletAuthorized(ctx sdk.Context, address string, grantee string, requiredPerm types.WalletPermission) bool {
	// For new wallets, creator is always authorized
	wallet, found := k.GetWallet(ctx, address)
	if !found {
		return false
	}
	if wallet.Creator == grantee {
		return true
	}

	// Check if grantee has the required permission
	perm, found := k.GetPermission(ctx, address, grantee)
	if !found {
		return false
	}

	// Convert WalletPermission enum to string
	permStr := requiredPerm.String()
	return perm.HasPermission(permStr)
}

// ValidateWalletAccess validates access to a wallet for a specific operation
func (k Keeper) ValidateWalletAccess(ctx sdk.Context, walletAddress string, creator string, operation string) error {
	// Get wallet
	wallet, found := k.GetWallet(ctx, walletAddress)
	if !found {
		return fmt.Errorf("wallet not found: %s", walletAddress)
	}

	// Check wallet status
	if wallet.Status != types.WalletStatus_WALLET_STATUS_ACTIVE {
		return fmt.Errorf("wallet is not active: current status %s", wallet.Status)
	}

	// Map operation to required permission
	var requiredPerm types.WalletPermission
	switch operation {
	case "sign":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_SIGN
	case "update":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_UPDATE
	case "delete":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_DELETE
	case "recover":
		requiredPerm = types.WalletPermission_WALLET_PERMISSION_RECOVER
	default:
		return fmt.Errorf("invalid operation: %s", operation)
	}

	// Check if creator has required permission
	if !k.IsWalletAuthorized(ctx, walletAddress, creator, requiredPerm) {
		return fmt.Errorf("creator %s is not authorized for operation %s on wallet %s", creator, operation, walletAddress)
	}

	return nil
}
