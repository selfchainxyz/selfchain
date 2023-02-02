package types

import (
	commontest "frontier/testutil"
	"frontier/testutil/sample"
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgMigrate_ValidateBasic(t *testing.T) {
	commontest.InitSDKConfig()

	tests := []struct {
		name string
		msg  MsgMigrate
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgMigrate{
				Creator: "invalid_address",
				Amount: strconv.FormatUint(MIN_MIGRATION_AMOUNT, 10),
				DestAddress: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgMigrate{
				Creator: sample.AccAddress(),
				Amount: strconv.FormatUint(MIN_MIGRATION_AMOUNT, 10),
				DestAddress: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func  TestMsgMigrate_ValidateBasic_min_amount(t *testing.T) {
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: strconv.FormatUint(MIN_MIGRATION_AMOUNT - 1, 10),
	}	

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, ErrInvalidMigrationAmount)
}

func  TestMsgMigrate_ValidateBasic_destAddress(t *testing.T) {
	// wrong prefix
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: "cosmos1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: strconv.FormatUint(MIN_MIGRATION_AMOUNT, 10),
	}	

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg2 := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: "front116r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: strconv.FormatUint(MIN_MIGRATION_AMOUNT, 10),
	}	

	err2 := msg2.ValidateBasic()
	require.ErrorIs(t, err2, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg3 := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: sample.AccAddress(),
		Amount: strconv.FormatUint(MIN_MIGRATION_AMOUNT, 10),
	}
	err3 := msg3.ValidateBasic()
	
	require.NoError(t, err3)
}
