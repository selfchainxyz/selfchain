package cli

import (
	"strconv"

	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMigrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [tx_hash] [eth-address] [dest-address] [amount] [token] [log_index]",
		Short: "Broadcast message migrate",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argTxHash := args[0]
			argEthAddress := args[1]
			argDestAddress := args[2]
			argAmount := args[3]

			argToken, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			argLogIndex, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMigrate(
				clientCtx.GetFromAddress().String(),
				argTxHash,
				argEthAddress,
				argDestAddress,
				argAmount,
				argToken,
				argLogIndex,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
