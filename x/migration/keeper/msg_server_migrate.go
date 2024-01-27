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

	// 1. we don't want to get spammed people who migrate small amounts
	amount := sdkmath.NewUintFromString(msg.Amount)
	if amount.LT(sdkmath.NewUint(config.MinMigrationAmount)) {
		return nil, types.ErrInvalidMigrationAmount
	}

	// 2. Make sure signer is in the list of migrators
	_, migratorExist := k.GetMigrator(ctx, msg.Creator)
	if !migratorExist {
		return nil, types.ErrUnknownMigrator
	}

	// 3. Create a hash of the message
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

	// 4. Check if message i.e. migration request has been processed already
	_, migrationExists := k.GetTokenMigration(ctx, msgHash)
	if migrationExists {
		return nil, types.ErrMigrationProcessed
	}

	// 5. Calculate the correct amount to mint
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

	// WEI has 18 decimals whereas our denomiation is uself thus it has 10^6 (6 decimals).
	normalizedAmount := amount.QuoUint64(uint64(math.Pow(10, 12)))
	lockedAmount := normalizedAmount.MulUint64(ratio).Quo(sdkmath.NewUint(100))

	lockedCoins := sdk.NewCoins(sdk.NewCoin(
		types.DENOM,
		sdkmath.NewIntFromBigInt(lockedAmount.BigInt()),
	))

	instantlyReleasedCoins := sdk.NewCoins(sdk.NewCoin(
		types.DENOM,
		sdkmath.NewIntFromBigInt(types.GetInstantlyReleasedAmount().BigInt()),
	))

	// 5. Mint new coins to the selfvesting module
	mintError := k.bankKeeper.MintCoins(ctx, selfvestingTypes.ModuleName, lockedCoins)
	if mintError != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "could not mint new coins (%s)", mintError)
	}

	// 6. Transfer a fixed amount to the beneficiary so it can pay gas when releasing tokens from the vesting position

	// We don't need to check the validatity of the address since it's been done in the Msg::ValidateBasic method
	destAddr, _ := sdk.AccAddressFromBech32(msg.DestAddress)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, selfvestingTypes.ModuleName, destAddr, instantlyReleasedCoins)

	// 7. Add a new beneficiary
	k.selfvestingKeeper.AddBeneficiary(ctx, selfvestingTypes.AddBeneficiaryRequest{
		Beneficiary: msg.DestAddress,
		Cliff:       config.VestingCliff,
		Duration:    config.VestingDuration,
		Amount:      lockedAmount.Sub(types.GetInstantlyReleasedAmount()).String(),
	})

	// 8. Store the token migration so it can't be processed again
	k.SetTokenMigration(ctx, types.TokenMigration{
		MsgHash:   msgHash,
		Processed: true,
	})

	return &types.MsgMigrateResponse{}, nil
}
