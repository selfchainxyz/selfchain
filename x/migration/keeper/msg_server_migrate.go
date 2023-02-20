package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"selfchain/x/migration/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) Migrate(goCtx context.Context, msg *types.MsgMigrate) (*types.MsgMigrateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Make sure signer is in the list of migrators
	_, migratorExist := k.GetMigrator(ctx, msg.Creator)
	if !migratorExist {
		return nil, types.ErrUnknownMigrator
	}

	// 2. Create a hash of the message
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

	// 3. Check if message i.e. migration request has been processed already
	_, migrationExists := k.GetTokenMigration(ctx, msgHash)
	if migrationExists {
		return nil, types.ErrMigrationProcessed
	}

	// 4. Calculate the correct amount to mint
	amount := sdkmath.NewUintFromString(msg.Amount)

	var ratio uint64
	switch msg.Token {
	case uint64(types.Front):
		ratio = types.FRONT_RATIO
	case uint64(types.Hotcross):
		ratio = types.HOTCROSS_RATIO
	}

	// WEI has 18 decimals whereas our denomiation is uself thus it has 10^6.
	normalizedAmount := amount.QuoUint64(uint64(math.Pow(10, 12)))
	mintedAmount := normalizedAmount.MulUint64(ratio).Quo(sdkmath.NewUint(100))
	mintedCoins := sdk.NewCoins(sdk.NewCoin(
		types.DENOM,
		sdkmath.NewIntFromBigInt(mintedAmount.BigInt()),
	))

	// 5. Mint new coins
	mintError := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintedCoins)
	if mintError != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "could not mint new coins (%s)", mintError)
	}

	// 6. transfer coins to destAddress

	// We don't need to check the validatity of the address since it's been done in the Msg::ValidateBasic method
	destAddr, _ := sdk.AccAddressFromBech32(msg.DestAddress)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, destAddr, mintedCoins)

	// 7. Store the token migration so it can't be processed again
	k.SetTokenMigration(ctx, types.TokenMigration{
		MsgHash:   msgHash,
		Processed: true,
	})

	return &types.MsgMigrateResponse{}, nil
}
