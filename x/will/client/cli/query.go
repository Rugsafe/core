package cli

import (
	"context"
	"fmt"

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
		ListWillsCmd(),
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
				fmt.Println("INSIDE WILL QUERY.GO, ERROR 1")
				return err
			}

			willID := args[0]
			fmt.Printf("GetWillCMD: %s", willID)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetWill(
				context.Background(),
				&types.QueryGetWillRequest{
					WillId: willID,
				},
			)
			fmt.Println("ENTIRE WILL STRUCT QUERY.GO")
			fmt.Println(res)
			fmt.Println("GetWillByID ID: " + res.Will.ID)
			fmt.Println("GetWillByID Beneficiary Name: " + res.Will.Beneficiary)
			fmt.Println("GetWillByID Name: " + res.Will.Name)
			if err != nil {
				fmt.Println("INSIDE WILL QUERY.GO, ERROR 2")
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// func ListWillsCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "list",
// 		Short: "Fetch a list of will by your address",
// 		Args:  cobra.ExactArgs(0),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx, err := client.GetClientQueryContext(cmd)
// 			if err != nil {
// 				fmt.Println("INSIDE LIST WILL QUERY.GO, ERROR 1")
// 				return err
// 			}

// 			address := clientCtx.GetFromAddress()

// 			fmt.Printf("ListWillsCmd address...: %s --- %s", address, clientCtx.FromAddress)
// 			queryClient := types.NewQueryClient(clientCtx)

// 			res, err := queryClient.ListWills(
// 				context.Background(),
// 				&types.QueryListWillsRequest{
// 					Address: address.String(),
// 				},
// 			)
// 			if err != nil {
// 				fmt.Println("INSIDE LIST WILL QUERY.GO, ERROR 2")
// 				return err
// 			}

// 			return clientCtx.PrintProto(res)
// 		},
// 	}
// 	flags.AddQueryFlagsToCmd(cmd)
// 	return cmd
// }

func ListWillsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [address]",
		Short: "Fetch a list of will by the specified address",
		Args:  cobra.ExactArgs(1), // Now requires exactly one argument: the address
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				fmt.Println("Error obtaining client context:", err)
				return err
			}

			address := args[0] // Use the first argument as the address

			fmt.Printf("Fetching wills for address: %s\n", address)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ListWills(
				context.Background(),
				&types.QueryListWillsRequest{
					Address: address,
				},
			)
			if err != nil {
				fmt.Println("Error querying wills:", err)
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
