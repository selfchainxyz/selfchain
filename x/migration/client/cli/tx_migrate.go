package cli

import (
	"strconv"

	"frontier/x/migration/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMigrate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [eth-address] [dest-address] [amount] [token]",
		Short: "Broadcast message migrate",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argEthAddress := args[0]
			argDestAddress := args[1]
			argAmount, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}
			argToken, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgMigrate(
				clientCtx.GetFromAddress().String(),
				argEthAddress,
				argDestAddress,
				argAmount,
				argToken,
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
