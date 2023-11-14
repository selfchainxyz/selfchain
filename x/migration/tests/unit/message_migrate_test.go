package test

import (
	"selfchain/testutil/sample"
	test "selfchain/x/migration/tests"
	"selfchain/x/migration/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgMigrate_ValidateBasic(t *testing.T) {
	test.InitSDKConfig()

	setup(t)

	tests := []struct {
		name string
		msg  types.MsgMigrate
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgMigrate{
				Creator:     "invalid_address",
				Amount:      "2000000000000000000",
				EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
				DestAddress: sample.AccAddress(),
				Token:       uint64(types.Front),
				TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgMigrate{
				Creator:     sample.AccAddress(),
				Amount:      "2000000000000000000",
				EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
				DestAddress: sample.AccAddress(),
				Token:       uint64(types.Front),
				TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
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

func TestValidateBasicDestAddress(t *testing.T) {
	// wrong prefix
	msg := types.MsgMigrate{
		Creator:     sample.AccAddress(),
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "cosmos1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount:      "2000000000000000000",
		Token:       uint64(types.Front),
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}

	err := msg.ValidateBasic()
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg2 := types.MsgMigrate{
		Creator:     sample.AccAddress(),
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "self184l076mmmeumx2cw6eqdct7tchufux0v43cxca",
		Amount:      "2000000000000000000",
		Token:       uint64(types.Front),
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}

	err2 := msg2.ValidateBasic()
	require.ErrorIs(t, err2, sdkerrors.ErrInvalidAddress)

	// correct prefix but invalid address
	msg3 := types.MsgMigrate{
		Creator:     sample.AccAddress(),
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: sample.AccAddress(),
		Amount:      "2000000000000000000",
		Token:       uint64(types.Front),
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}
	err3 := msg3.ValidateBasic()

	require.NoError(t, err3)
}

func TestValidateBasicWrongToken(t *testing.T) {
	msg := types.MsgMigrate{
		Creator:     sample.AccAddress(),
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "self184l076mmmeumx2cw6eqdct7tchufux0v43cxcg",
		Amount:      "2000000000000000000",
		Token:       2,
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}

	err := msg.ValidateBasic()
	require.ErrorIs(t, err, types.ErrTokenNotSupported)
}

func TestValidateBasicEmptyStringValues(t *testing.T) {
	msg := types.MsgMigrate{
		Creator:     sample.AccAddress(),
		EthAddress:  "",
		DestAddress: "self184l076mmmeumx2cw6eqdct7tchufux0v43cxcg",
		Amount:      "2000000000000000000",
		Token:       uint64(types.Hotcross),
		TxHash:      "2683f98e2bc2fb5a36c4064d561121fb5087451e70df03b8593dc427ef228c86",
	}

	err := msg.ValidateBasic()
	require.ErrorIs(t, err, types.ErrEmptyStringValue)

	msg2 := types.MsgMigrate{
		Creator:     sample.AccAddress(),
		EthAddress:  "baf6dc2e647aeb6f510f9e318856a1bcd66c5e19",
		DestAddress: "self184l076mmmeumx2cw6eqdct7tchufux0v43cxcg",
		Amount:      "2000000000000000000",
		Token:       uint64(types.Hotcross),
		TxHash:      "",
	}

	err2 := msg2.ValidateBasic()
	require.ErrorIs(t, err2, types.ErrEmptyStringValue)
}
