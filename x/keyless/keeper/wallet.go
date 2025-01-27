package keeper

import (
	"fmt"

	"selfchain/x/keyless/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateWalletAccess validates if the permission exists in the wallet's permissions list
func (k Keeper) ValidateWalletAccess(ctx sdk.Context, walletAddress string, permission string) error {
	// Convert string permission to enum
	var walletPerm types.WalletPermission
	switch permission {
	case "sign":
		walletPerm = types.WalletPermission_WALLET_PERMISSION_SIGN
	case "recover":
		walletPerm = types.WalletPermission_WALLET_PERMISSION_RECOVER
	case "rotate":
		walletPerm = types.WalletPermission_WALLET_PERMISSION_ROTATE
	case "admin":
		walletPerm = types.WalletPermission_WALLET_PERMISSION_ADMIN
	default:
		return fmt.Errorf("invalid permission: %s", permission)
	}

	// Get transaction sender
	sender := sdk.AccAddress(ctx.BlockHeader().ProposerAddress).String()

	// Check if sender has the required permission
	if !k.HasPermission(ctx, walletAddress, sender, walletPerm) {
		return fmt.Errorf("sender %s does not have permission %s for wallet %s", sender, permission, walletAddress)
	}

	return nil
}

// IsWalletAuthorized checks if the creator is authorized to operate on the wallet
func (k Keeper) IsWalletAuthorized(ctx sdk.Context, creator string, walletAddress string) (bool, error) {
	// Get wallet by address
	wallet, err := k.GetWallet(ctx, walletAddress)
	if err != nil {
		return false, fmt.Errorf("wallet not found: %v", err)
	}

	// Check if creator is wallet owner
	return wallet.Creator == creator, nil
}

// ValidateRecoveryProof validates a recovery proof for a wallet
func (k Keeper) ValidateRecoveryProof(ctx sdk.Context, walletId, recoveryProof string) error {
	// TODO: Implement recovery proof validation
	// For now, just return nil for testing
	return nil
}

// SignWithTSS signs a transaction using TSS
func (k Keeper) SignWithTSS(ctx sdk.Context, wallet *types.Wallet, unsignedTx string) ([]byte, error) {
	// TODO: Implement TSS signing
	// For now, just return empty signature for testing
	return []byte{}, nil
}

// updateWalletAfterSigning updates wallet metadata after successful signing
func (k Keeper) updateWalletAfterSigning(ctx sdk.Context, wallet *types.Wallet) error {
	// TODO: Implement wallet metadata update
	// For now, just return nil for testing
	return nil
}

// CreateWallet creates a new wallet
func (k Keeper) CreateWallet(ctx sdk.Context, msg *types.MsgCreateWallet) error {
	// Validate creator
	if msg.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}

	// Create new wallet
	wallet := &types.Wallet{
		Id:            msg.PubKey,
		Creator:       msg.Creator,
		WalletAddress: msg.WalletAddress,
		ChainId:      msg.ChainId,
		Status:       types.WalletStatus_WALLET_STATUS_ACTIVE,
	}

	// Save wallet
	if err := k.SaveWallet(ctx, wallet); err != nil {
		return fmt.Errorf("failed to save wallet: %v", err)
	}

	// Emit wallet created event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateWallet,
			sdk.NewAttribute(types.AttributeKeyWalletID, wallet.Id),
			sdk.NewAttribute(types.AttributeKeyCreator, wallet.Creator),
			sdk.NewAttribute(types.AttributeKeyWalletAddress, wallet.WalletAddress),
			sdk.NewAttribute(types.AttributeKeyChainID, wallet.ChainId),
			sdk.NewAttribute(types.AttributeKeyStatus, wallet.Status.String()),
		),
	)

	return nil
}
