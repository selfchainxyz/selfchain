package test

import (
	"context"

	selfvestingTypes "selfchain/x/selfvesting/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

func (vesting *MockSelfvestingKeeper) ExpectAddBeneficiary(context context.Context, req selfvestingTypes.AddBeneficiaryRequest) *gomock.Call {
	return vesting.EXPECT().AddBeneficiary(sdk.UnwrapSDKContext(context), req)
}
