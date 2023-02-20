package cli

import (
	"context"

	"selfchain/x/migration/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListTokenMigration() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-token-migration",
		Short: "list all token-migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllTokenMigrationRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.TokenMigrationAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowTokenMigration() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-token-migration [msg-hash]",
		Short: "shows a token-migration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argMsgHash := args[0]

			params := &types.QueryGetTokenMigrationRequest{
				MsgHash: argMsgHash,
			}

			res, err := queryClient.TokenMigration(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
