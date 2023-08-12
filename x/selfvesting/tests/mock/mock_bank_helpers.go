package test

import (
	"context"
	"selfchain/x/selfvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

func coinsOf(amount uint64) sdk.Coins {
	return sdk.Coins{
		sdk.Coin{
			Denom:  types.DENOM,
			Amount: sdk.NewInt(int64(amount)),
		},
	}
}

func (escrow *MockBankKeeper) ExpectReceiveCoins(context context.Context, who string, amount uint64) *gomock.Call {
	whoAddr, err := sdk.AccAddressFromBech32(who)
	if err != nil {
		panic(err)
	}

	return escrow.EXPECT().SendCoinsFromModuleToAccount(sdk.UnwrapSDKContext(context), types.ModuleName, whoAddr, coinsOf(amount))
}

func (escrow *MockBankKeeper) ExpectMintToModule(context context.Context, amount uint64) *gomock.Call {
	return escrow.EXPECT().MintCoins(sdk.UnwrapSDKContext(context), types.ModuleName, coinsOf(amount))
}
