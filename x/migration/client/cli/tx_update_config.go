package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"selfchain/x/migration/types"
)

var _ = strconv.Itoa(0)

func CmdUpdateConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-config [vesting-duration] [vesting-cliff] [min-migration-amount]",
		Short: "Broadcast message update-config",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argVestingDuration, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argVestingCliff, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			argMinMigrationAmount, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateConfig(
				clientCtx.GetFromAddress().String(),
				argVestingDuration,
				argVestingCliff,
				argMinMigrationAmount,
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
