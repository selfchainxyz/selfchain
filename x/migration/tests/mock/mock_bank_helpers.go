package test

import (
	"context"
	"frontier/x/migration/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

func coinsOf(amount uint64) sdk.Coins {
	return sdk.Coins{
		sdk.Coin{
			Denom:  sdk.DefaultBondDenom,
			Amount: sdk.NewInt(int64(amount)),
		},
	}
}

func (escrow *MockBankKeeper) ExpectMint(context context.Context, who string, amount uint64) *gomock.Call {
	whoAddr, err := sdk.AccAddressFromBech32(who)
	if err != nil {
		panic(err)
	}

	return escrow.EXPECT().SendCoinsFromModuleToAccount(sdk.UnwrapSDKContext(context), types.ModuleName, whoAddr, coinsOf(amount))
}
