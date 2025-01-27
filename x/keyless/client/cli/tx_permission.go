package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"selfchain/x/keyless/types"
)

func CmdGrantPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant-permission [wallet-id] [grantee] [permissions] [expires-at]",
		Short: "Grant permissions to a wallet",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			walletID := args[0]
			grantee := args[1]

			// Parse permissions
			var permissions []types.WalletPermission
			if err := json.Unmarshal([]byte(args[2]), &permissions); err != nil {
				return fmt.Errorf("invalid permissions format: %v", err)
			}

			// Parse expiry time
			expiresAt, err := time.Parse(time.RFC3339, args[3])
			if err != nil {
				return fmt.Errorf("invalid expiry time format: %v", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgGrantPermission(
				clientCtx.GetFromAddress().String(),
				walletID,
				grantee,
				permissions,
				&expiresAt,
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

func CmdRevokePermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-permission [wallet-id] [grantee] [permissions]",
		Short: "Revoke permissions from a grantee",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			walletID := args[0]
			grantee := args[1]

			// Parse permissions
			var permissions []types.WalletPermission
			if err := json.Unmarshal([]byte(args[2]), &permissions); err != nil {
				return fmt.Errorf("invalid permissions format: %v", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRevokePermission(
				clientCtx.GetFromAddress().String(),
				walletID,
				grantee,
				permissions,
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
