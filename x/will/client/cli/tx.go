package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

func CreateWillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name] [beneficiary] [height]",
		Short: "Create a Will",
		Args:  cobra.MinimumNArgs(3), // Expecting at least 3 arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parsing height from args
			height, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse height '%s' into int64: %w", args[2], err)
			}

			// Parsing components
			componentsArgs, _ := cmd.Flags().GetStringArray("component")
			var components []*types.ExecutionComponent
			for _, compArg := range componentsArgs {
				component, err := parseComponentFromString(compArg)
				fmt.Println(component)
				if err != nil {
					return fmt.Errorf("failed to parse component: %w", err)
				}
				components = append(components, component)
			}

			msg := types.MsgCreateWillRequest{
				Creator:     clientCtx.GetFromAddress().String(),
				Name:        args[0],
				Beneficiary: args[1],
				Height:      height,
				Components:  components,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringArray("component", []string{}, "Add components to the will. Format: --component <type> <params>. Can be used multiple times for different components.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func parseComponentFromString(compArg string) (*types.ExecutionComponent, error) {
	// Split the input to separate the type from the parameters
	parts := strings.SplitN(compArg, ":", 2)
	fmt.Println(parts)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid component format, expected '<type>:<params...>' but got: %s", compArg)
	}

	componentType, paramsStr := parts[0], parts[1]
	fmt.Println(componentType)
	fmt.Println(paramsStr)
	component := &types.ExecutionComponent{}

	switch componentType {
	case "transfer":
		// 'transfer' expects 'to,amount', but let's keep it flexible
		params := strings.Split(paramsStr, ",")
		if len(params) < 2 {
			return nil, fmt.Errorf("transfer component expects at least 'to,amount', but got: %s", paramsStr)
		}

		to, amountStr := params[0], params[1]
		amount, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount format for transfer component: %s", amountStr)
		}
		amountCoin := sdk.NewInt64Coin("w3ll", amount) // Ensure "w3ll" matches your denomination

		component.ComponentType = &types.ExecutionComponent_Transfer{
			Transfer: &types.TransferComponent{
				To:     to,
				Amount: &amountCoin,
			},
		}

	// Example for a component type that might have a variable number of parameters
	case "customType":
		// 'customType' expects a variable number of parameters
		params := strings.Split(paramsStr, ",")
		// Process params based on your custom logic here
		// For example, you might only need to check if at least one parameter is provided
		if len(params) < 1 {
			return nil, fmt.Errorf("customType component expects at least one parameter, but got: %s", paramsStr)
		}

		// Assuming a hypothetical structure for CustomTypeComponent that takes a slice of strings as params
		// component.ComponentType = &types.ExecutionComponent_CustomType{
		//     CustomType: &types.CustomTypeComponent{
		//         Params: params,
		//     },
		// }

	default:
		return nil, fmt.Errorf("unsupported component type: %s", componentType)
	}

	return component, nil
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
			// willId, err := strconv.ParseUint(args[0], 10, 64)
			willId := args[0]
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
