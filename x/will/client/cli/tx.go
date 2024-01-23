package cli

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/will/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/spf13/cobra"
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
	)
	return txCmd
}

// GetWillCmd will return a will
func CreateWillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-will [will-id]",
		Short:   "Create a Will",
		Aliases: []string{"get"},
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
			msg := types.MsgCreateWill{
				Creator:     clientCtx.GetFromAddress().String(),
				Id:          args[0],
				Name:        "test will",
				Beneficiary: "benefiary 1",
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
		SilenceUsage: true,
	}

	// addInstantiatePermissionFlags(cmd)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
