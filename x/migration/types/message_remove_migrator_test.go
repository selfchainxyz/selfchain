package types

import (
	"github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"testing"

	"selfchain/testutil/sample"

	"github.com/stretchr/testify/require"
)

func TestMsgRemoveMigrator_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRemoveMigrator
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRemoveMigrator{
				Creator: "invalid_address",
			},
			err: types.ErrInvalidAccountAddress,
		}, {
			name: "valid address",
			msg: MsgRemoveMigrator{
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
