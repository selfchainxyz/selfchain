package types

import (
	commontest "frontier/testutil"
	"frontier/testutil/sample"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/holiman/uint256"
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
				Amount: getMinMigrationAmount().Hex(),
				DestAddress: sample.AccAddress(),
				Token: uint64(Front),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgMigrate{
				Creator: sample.AccAddress(),
				Amount: getMinMigrationAmount().Hex(),
				DestAddress: sample.AccAddress(),
				Token: uint64(Front),
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

func  TestMsgMigrate_ValidateBasic_minAmount(t *testing.T) {
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: uint256.NewInt(0).Sub(getMinMigrationAmount(), uint256.NewInt(1)).Hex(),
		Token: uint64(Hotcross),
	}

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, ErrInvalidMigrationAmount)
}

func  TestMsgMigrate_ValidateBasic_destAddress(t *testing.T) {
	// wrong prefix
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: "cosmos1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Front),
	}	

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg2 := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: "front116r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Front),
	}	

	err2 := msg2.ValidateBasic()
	require.ErrorIs(t, err2, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg3 := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: sample.AccAddress(),
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Front),
	}
	err3 := msg3.ValidateBasic()
	
	require.NoError(t, err3)
}

func  TestMsgMigrate_ValidateBasic_WrongToken(t *testing.T) {
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: 2,
	}

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, ErrTokenNotSupported)
}
