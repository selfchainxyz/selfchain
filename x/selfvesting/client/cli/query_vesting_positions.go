package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"selfchain/x/selfvesting/types"
)

func CmdListVestingPositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-vesting-positions",
		Short: "list all vestingPositions",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllVestingPositionsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.VestingPositionsAll(context.Background(), params)
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

func CmdShowVestingPositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-vesting-positions [beneficiary]",
		Short: "shows a vestingPositions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argBeneficiary := args[0]

			params := &types.QueryGetVestingPositionsRequest{
				Beneficiary: argBeneficiary,
			}

			res, err := queryClient.VestingPositions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
