package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"selfchain/x/keyless/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryWallet())
	cmd.AddCommand(CmdListWallets())
	cmd.AddCommand(CmdQueryPartyData())
	cmd.AddCommand(CmdQueryKeyRotationStatus())
	cmd.AddCommand(CmdQueryBatchSignStatus())

	return cmd
}

func CmdQueryWallet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet [address]",
		Short: "Query wallet by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Wallet(cmd.Context(), &types.QueryWalletRequest{
				Address: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdListWallets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-wallets",
		Short: "List all wallets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Wallets(cmd.Context(), &types.QueryWalletsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "wallets")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryPartyData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "party-data [wallet-address]",
		Short: "Query party data by wallet address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PartyData(cmd.Context(), &types.QueryPartyDataRequest{
				WalletAddress: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryKeyRotationStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-rotation-status [wallet-id]",
		Short: "Query key rotation status by wallet ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.KeyRotationStatus(cmd.Context(), &types.QueryKeyRotationStatusRequest{
				WalletId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryBatchSignStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-sign-status [wallet-id] [batch-id]",
		Short: "Query batch sign status by wallet ID and batch ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.BatchSignStatus(cmd.Context(), &types.QueryBatchSignStatusRequest{
				WalletId: args[0],
				BatchId:  args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
