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
				EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
				DestAddress: sample.AccAddress(),
				Token: uint64(Front),
				TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgMigrate{
				Creator: sample.AccAddress(),
				Amount: getMinMigrationAmount().Hex(),
				EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
				DestAddress: sample.AccAddress(),
				Token: uint64(Front),
				TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
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
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: uint256.NewInt(0).Sub(getMinMigrationAmount(), uint256.NewInt(1)).Hex(),
		Token: uint64(Hotcross),
		TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, ErrInvalidMigrationAmount)
}

func  TestMsgMigrate_ValidateBasic_destAddress(t *testing.T) {
	// wrong prefix
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "cosmos1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Front),
		TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}	

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg2 := MsgMigrate{
		Creator: sample.AccAddress(),
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "front116r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Front),
		TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}	

	err2 := msg2.ValidateBasic()
	require.ErrorIs(t, err2, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg3 := MsgMigrate{
		Creator: sample.AccAddress(),
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: sample.AccAddress(),
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Front),
		TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}
	err3 := msg3.ValidateBasic()
	
	require.NoError(t, err3)
}

func  TestMsgMigrate_ValidateBasic_WrongToken(t *testing.T) {
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: 2,
		TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, ErrTokenNotSupported)
}

func  TestMsgMigrate_ValidateBasic_EmptyStringValues(t *testing.T) {
	msg := MsgMigrate{
		Creator: sample.AccAddress(),
		EthAddress: "",
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Hotcross),
		TxHash: "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, ErrEmptyStringValue)

	msg2 := MsgMigrate{
		Creator: sample.AccAddress(),
		EthAddress: "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: getMinMigrationAmount().Hex(),
		Token: uint64(Hotcross),
		TxHash: "",
	}

	err2 := msg2.ValidateBasic()
  require.ErrorIs(t, err2, ErrEmptyStringValue)
}
