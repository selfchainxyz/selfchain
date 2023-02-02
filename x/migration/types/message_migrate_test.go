package types

import (
	"testing"

	"frontier/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgMigrate_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgMigrate
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgMigrate{
				Creator: "invalid_address",
				Amount: MIN_MIGRATION_AMOUNT,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgMigrate{
				Creator: sample.AccAddress(),
				Amount: MIN_MIGRATION_AMOUNT,
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
		TxHash: "eacffe9c44c1f4f77537766b772afd9da6d84ad215e43c52057909fb4d9c2488",
		EthAddress: "0x37f1f67955ac36763409377bd2ce64da414c3972",
		DestAddress: "front1k6r2mzwhkn3tr8hz947kqkl7ym9gnrgf0a0g6v",
		Amount: MIN_MIGRATION_AMOUNT - 1,
		Token: 0,
	}	

	err := msg.ValidateBasic()
  require.ErrorIs(t, err, ErrInvalidMigrationAmount)
}
