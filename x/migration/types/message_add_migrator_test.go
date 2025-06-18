package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	"testing"

	"selfchain/testutil/sample"

	"github.com/stretchr/testify/require"
)

func TestMsgAddMigrator_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddMigrator
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddMigrator{
				Creator: "invalid_address",
			},
			err: errors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgAddMigrator{
				Creator: sample.AccAddress(),
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
