package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	"github.com/CosmWasm/wasmd/x/will/types"
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Will transaction subcommand",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	txCmd.AddCommand(
		CreateWillCmd(),
		CheckInCmd(),
	)
	return txCmd
}

// GetWillCmd will return a will
func CreateWillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [will-id]",
		Short:   "Create a Will",
		Aliases: []string{"cw"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// msg, err := parseStoreCodeArgs(args[0], clientCtx.GetFromAddress().String(), cmd.Flags())
			// if err != nil {
			// 	return err
			// }
			fmt.Println("inside tx CreateWill command")
			//
			// logger := log.Logger{}
			// logger := log.NewTestLogger(t)
			logger := network.NewCLILogger(cmd)
			logger.Log("inside tx CreateWill command")
			logger.Log(string(args[0]))
			logger.Log(string(args[1]))
			msg := types.MsgCreateWillRequest{
				Creator:     clientCtx.GetFromAddress().String(),
				Id:          1,
				Name:        args[0],
				Beneficiary: args[1],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
		SilenceUsage: true,
	}

	// addInstantiatePermissionFlags(cmd)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CheckInCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "Checkin [will-id]",
		Short:   "Submit a checkin to the will",
		Aliases: []string{"cw"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			// msg, err := parseStoreCodeArgs(args[0], clientCtx.GetFromAddress().String(), cmd.Flags())
			// if err != nil {
			// 	return err
			// }
			fmt.Println("inside tx Checkin command 1")
			//
			// logger := log.Logger{}
			// logger := log.NewTestLogger(t)
			logger := network.NewCLILogger(cmd)
			logger.Log("inside tx Checkin command 2")
			logger.Log(string(args[0]))
			willId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse will ID: %w", err)
			}

			msg := types.MsgCheckInRequest{
				Creator: clientCtx.GetFromAddress().String(),
				Id:      willId,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
		SilenceUsage: true,
	}

	// addInstantiatePermissionFlags(cmd)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
