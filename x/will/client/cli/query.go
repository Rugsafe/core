package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/CosmWasm/wasmd/x/will/types"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the will module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	queryCmd.AddCommand(
		GetWillCmd(),
	)
	return queryCmd
}

// GetWillCmd will return a will
func GetWillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [will-id]",
		Short: "Query a Will by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			willID := args[0]
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetWill(
				context.Background(),
				&types.QueryGetWillRequest{
					WillId: willID,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
