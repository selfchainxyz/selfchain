package cli

import (
	"fmt"
	"encoding/hex"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"selfchain/x/keyless/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Transaction commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdCreateWallet(),
		CmdRecoverWallet(),
		CmdSignTransaction(),
		CmdBatchSign(),
		CmdInitiateKeyRotation(),
		CmdCompleteKeyRotation(),
	)

	return cmd
}

func CmdCreateWallet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-wallet [pub-key] [wallet-address] [chain-id]",
		Short: "Create a new wallet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateWallet(
				clientCtx.GetFromAddress().String(),
				args[0], // pubKey
				args[1], // walletAddress
				args[2], // chainId
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRecoverWallet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recover-wallet [wallet-address] [new-pub-key] [recovery-proof]",
		Short: "Recover a wallet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRecoverWallet(
				clientCtx.GetFromAddress().String(),
				args[0], // walletAddress
				args[1], // newPubKey
				args[2], // recoveryProof
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSignTransaction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign-transaction [wallet-address] [unsigned-tx]",
		Short: "Sign a transaction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSignTransaction(
				clientCtx.GetFromAddress().String(),
				args[0], // walletAddress
				args[1], // unsignedTx
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdBatchSign() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-sign [wallet-address] [messages] [parties]",
		Short: "Batch sign multiple messages",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse hex-encoded messages
			messages := make([][]byte, 0)
			hexMessages := args[1]
			msgBytes, err := hex.DecodeString(hexMessages)
			if err != nil {
				return fmt.Errorf("failed to decode messages: %v", err)
			}
			messages = append(messages, msgBytes)

			// Parse parties
			parties := []string{args[2]}

			msg := types.NewMsgBatchSignRequest(
				clientCtx.GetFromAddress().String(),
				args[0], // walletAddress
				messages,
				parties,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdInitiateKeyRotation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initiate-key-rotation [wallet-address] [new-pub-key] [signature]",
		Short: "Initiate key rotation",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgInitiateKeyRotation(
				clientCtx.GetFromAddress().String(),
				args[0], // walletAddress
				args[1], // newPubKey
				args[2], // signature
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCompleteKeyRotation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete-key-rotation [wallet-address] [new-pub-key] [signature] [recovery-proof]",
		Short: "Complete key rotation",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCompleteKeyRotation(
				clientCtx.GetFromAddress().String(),
				args[0], // walletAddress
				args[1], // newPubKey
				args[2], // signature
				args[3], // recoveryProof
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
