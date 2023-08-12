package test

import (
	"context"
	"testing"

	keepertest "selfchain/testutil/keeper"
	"selfchain/x/selfvesting"
	"selfchain/x/selfvesting/keeper"
	mocktest "selfchain/x/selfvesting/tests/mock"
	"selfchain/x/selfvesting/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func setup(t testing.TB) (types.MsgServer, context.Context, keeper.Keeper, *gomock.Controller, *mocktest.MockBankKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := mocktest.NewMockBankKeeper(ctrl)
	k, ctx := keepertest.SelfvestingKeeperWithMocks(t, bankMock)

	// setup genesis params for this module
	genesis := *types.DefaultGenesis()
	genesis.VestingPositionsList = []types.VestingPositions{}

	selfvesting.InitGenesis(ctx, *k, genesis)

	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)

	return server, context, *k, ctrl, bankMock
}

func TestShouldFailIfInvalidBeneficiaryAddr(t *testing.T) {
	_, ctx, keeper, ctrl, _ := setup(t)
	defer ctrl.Finish()

	_, err := keeper.AddBeneficiary(sdk.UnwrapSDKContext(ctx), types.AddBeneficiaryRequest{
		Beneficiary: "Invalid Address",
		Cliff:       0,
		Duration:    0,
		Amount:      "",
	})

	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
}
