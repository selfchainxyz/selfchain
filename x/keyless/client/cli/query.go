package cli

import (
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"selfchain/x/keyless/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group keyless queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		CmdQueryWallet(),
		CmdListWallets(),
		CmdQueryKeyRotation(),
		CmdQueryKeyRotations(),
		CmdListAuditEvents(),
	)

	return cmd
}

func CmdQueryWallet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet [id]",
		Short: "Query a wallet by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Wallet(cmd.Context(), &types.QueryWalletRequest{
				Id: args[0],
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
		Short: "Query all wallets",
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
			res, err := queryClient.ListWallets(cmd.Context(), &types.QueryListWalletsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "wallets")
	return cmd
}

func CmdQueryKeyRotation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-rotation [wallet-id] [version]",
		Short: "Query a key rotation by wallet ID and version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			version, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("version must be a positive integer: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.KeyRotation(cmd.Context(), &types.QueryKeyRotationRequest{
				WalletId: args[0],
				Version:  version,
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

func CmdQueryKeyRotations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-rotations [wallet-id]",
		Short: "Query all key rotations for a wallet",
		Args:  cobra.ExactArgs(1),
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
			res, err := queryClient.KeyRotations(cmd.Context(), &types.QueryKeyRotationsRequest{
				WalletId:   args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "key-rotations")
	return cmd
}

func CmdListAuditEvents() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit-events [wallet-id] [event-type] [success]",
		Short: "Query audit events for a wallet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			success, err := strconv.ParseBool(args[2])
			if err != nil {
				return fmt.Errorf("success must be true or false: %w", err)
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ListAuditEvents(cmd.Context(), &types.QueryListAuditEventsRequest{
				WalletId:   args[0],
				EventType:  args[1],
				Success:    success,
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "audit-events")
	return cmd
}
