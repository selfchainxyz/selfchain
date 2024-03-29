package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"selfchain/x/migration/types"
	selfvestingTypes "selfchain/x/selfvesting/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Migrate(goCtx context.Context, msg *types.MsgMigrate) (*types.MsgMigrateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	config, configExists := k.GetConfig(ctx); if !configExists {
		panic("Config does not exist")
	}

	// We don't want to get spammed people who migrate small amounts
	amount := sdkmath.NewUintFromString(msg.Amount)
	if amount.LT(sdkmath.NewUint(config.MinMigrationAmount)) {
		return nil, types.ErrInvalidMigrationAmount
	}

	// Make sure signer is in the list of migrators
	_, migratorExist := k.GetMigrator(ctx, msg.Creator)
	if !migratorExist {
		return nil, types.ErrUnknownMigrator
	}

	// Create a hash of the message
	encodedMsg := fmt.Sprintf(
		"%s|%s|%s|%d|%s|%d",
		msg.EthAddress,
		msg.DestAddress,
		msg.Amount,
		msg.Token,
		msg.TxHash,
		msg.LogIndex,
	)
	msgHash := fmt.Sprintf("%x", sha256.Sum256([]byte(encodedMsg)))

	// Check if message i.e. migration request has been processed already
	_, migrationExists := k.GetTokenMigration(ctx, msgHash)
	if migrationExists {
		return nil, types.ErrMigrationProcessed
	}

	// Calculate the correct amount to mint
	var ratio uint64
	switch msg.Token {
	case uint64(types.Front):
		ratio = types.FRONT_RATIO
	case uint64(types.Hotcross):
		ratio = k.HotcrossRatio(ctx)

		if ratio == 0 {
			return nil, types.ErrHotcrossRatioZero
		}
	}

	// WEI has 18 decimals whereas our denomiation is uslf thus it has 10^6 (6 decimals).
	normalizedAmount := amount.QuoUint64(uint64(math.Pow(10, 12)))
	migrationAmount := normalizedAmount.MulUint64(ratio).Quo(sdkmath.NewUint(100))
	instantlyReleased := types.GetInstantlyReleasedAmount()
	migrationCoins := sdk.NewCoins(sdk.NewCoin(
		types.DENOM,
		sdkmath.NewIntFromBigInt(migrationAmount.BigInt()),
	))
	
	// Mint new coins to the selfvesting module
	mintError := k.bankKeeper.MintCoins(ctx, selfvestingTypes.ModuleName, migrationCoins)
	if mintError != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "could not mint new coins (%s)", mintError)
	}

	// We don't need to check the validatity of the address since it's been done in the Msg::ValidateBasic method
	destAddr, _ := sdk.AccAddressFromBech32(msg.DestAddress)

	// If the migration amount is LTE then instantlyReleased which is a constant 1 SLF then we don't need to
	// do any vesting. This can happen in two cases:
	// 1. we migrate 1 FRONT so we get 1 SLF which is instantly released
	// 2. we migrate X HOTCROSS and since hotcross migration ration will be < 100% then for small amount it will
	// create less than 1 SLF (which is the instantlyReleased). For example, if we migrate 1 HOTCROSS and the ratio
	// is 25% then we will get 0.25 SLF. In this case we don't need to create any vesting and we simply mint 0.25 SLF
	// and transfer to the user.
	if migrationAmount.LTE(instantlyReleased) {
		instantlyReleasedCoins := sdk.NewCoins(sdk.NewCoin(
			types.DENOM,
			sdkmath.NewIntFromBigInt(migrationAmount.BigInt()),
		))

		// send the full migration amount to the user
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, selfvestingTypes.ModuleName, destAddr, instantlyReleasedCoins)
	} else {
		// Transfer a fixed amount to the beneficiary so it can pay gas when releasing tokens from the vesting position
		instantlyReleasedCoins := sdk.NewCoins(sdk.NewCoin(
			types.DENOM,
			sdkmath.NewIntFromBigInt(instantlyReleased.BigInt()),
		))
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, selfvestingTypes.ModuleName, destAddr, instantlyReleasedCoins)

		// Add a new beneficiary
		k.selfvestingKeeper.AddBeneficiary(ctx, selfvestingTypes.AddBeneficiaryRequest{
			Beneficiary: msg.DestAddress,
			Cliff:       config.VestingCliff,
			Duration:    config.VestingDuration,
			Amount:      migrationAmount.Sub(instantlyReleased).String(),
		})
	}

	// Store the token migration so it can't be processed again
	k.SetTokenMigration(ctx, types.TokenMigration{
		MsgHash:   msgHash,
		Processed: true,
	})

	return &types.MsgMigrateResponse{}, nil
}
