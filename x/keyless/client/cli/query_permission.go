package cli

import (
	"selfchain/x/keyless/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-permissions [wallet-id]",
		Short: "List all permissions for a wallet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPermissionsRequest{
				WalletId: args[0],
			}

			res, err := queryClient.Permissions(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetPermission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-permission [wallet-id] [grantee]",
		Short: "Get permission for a specific grantee",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPermissionRequest{
				WalletId: args[0],
				Grantee:  args[1],
			}

			res, err := queryClient.Permission(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
