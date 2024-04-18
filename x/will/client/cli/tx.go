package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
		ClaimCmd(),
	)
	return txCmd
}

func CreateWillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name] [beneficiary] [height] [components]",
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

			componentNames, err := cmd.Flags().GetStringArray("component-name")
			if err != nil {
				return fmt.Errorf("failed to parse component names: %w", err)
			}

			componentArgs, err := cmd.Flags().GetStringArray("component-args")
			if err != nil {
				return fmt.Errorf("failed to parse component arguments: %w", err)
			}

			outputs, err := cmd.Flags().GetStringArray("component-output-type")
			if err != nil {
				return fmt.Errorf("failed to parse component arguments: %w", err)
			}

			outputsArgs, err := cmd.Flags().GetStringArray("component-output-args")
			if err != nil {
				return fmt.Errorf("failed to parse component arguments: %w", err)
			}

			if len(componentNames) != len(componentArgs) {
				return fmt.Errorf("mismatch between component names and arguments count")
			}

			if len(componentNames) != len(outputs) {
				return fmt.Errorf("mismatch between component names and outputs count")
			}

			if len(outputs) != len(outputsArgs) {
				return fmt.Errorf("mismatch between component outputs and output args count")
			}

			var sender string = clientCtx.GetFromAddress().String()
			var components []*types.ExecutionComponent
			for i, componentName := range componentNames {
				componentArg := componentArgs[i]
				output := outputs[i]
				outputArgs := outputsArgs[i]
				// component, err := parseComponent(componentName, componentArg)
				component, err := parseComponentFromString(componentName, componentArg, output, outputArgs, sender)
				if err != nil {
					return fmt.Errorf("failed to parse component: %w", err)
				}
				components = append(components, component)
			}

			msg := types.MsgCreateWillRequest{
				Creator:     sender,
				Name:        args[0],
				Beneficiary: args[1],
				Height:      height,
				Components:  components,
			}

			fmt.Println(msg)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringArray("component-name", []string{}, "Names of the components. Use multiple --component-name flags for multiple components.")
	cmd.Flags().StringArray("component-args", []string{}, "Arguments for the components. Use multiple --component-args flags for multiple components. Must match the order of --component-name flags.")
	cmd.Flags().StringArray("component-output-type", []string{}, "Arguments for the outputs of each component. Use multiple --component-output-type flags for multiple component output. Must match the order of --component-output-type flags.")
	cmd.Flags().StringArray("component-output-args", []string{}, "Arguments for the arguments of each component output. Use multiple --component-output-args flags for multiple components. Must match the order of --component-output-args flags.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func generateUniqueComponentID() string {
	return uuid.New().String()
}

func parseComponentFromString(componentName string, componentData string, outputType string, outputArgs string, sender string) (*types.ExecutionComponent, error) {
	// The componentName is already separated, now just need to parse componentData
	typeParts := strings.SplitN(componentData, ":", 2)
	// if len(typeParts) != 2 {
	// 	return nil, fmt.Errorf("invalid component data format, expected 'componentType:componentParams'")
	// }

	fmt.Println("TYPE PARTS: ", typeParts)

	outputParts := strings.SplitN(outputArgs, ":", 2)
	fmt.Println("OUTPUT PARTS: ", outputParts)

	componentType, params := typeParts[0], typeParts[1]
	componentID := generateUniqueComponentID() // Function to generate a unique ID for each component

	var component types.ExecutionComponent
	component.Name = componentName
	component.Id = componentID
	component.Status = "inactive"

	outputParams := outputParts[0]
	fmt.Println("outputType: ", outputType)
	fmt.Println("outputParams: ", outputParams)
	// panic(99)
	switch componentType {
	case "transfer":
		dataParts := strings.Split(params, ",")
		if len(dataParts) != 3 {
			return nil, fmt.Errorf("invalid transfer component params, expected 'to,amount'")
		}
		to, amountStr, denom := dataParts[0], dataParts[1], dataParts[2]
		amount, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount format for transfer component: %s", amountStr)
		}
		amountCoin := sdk.NewInt64Coin(denom, amount) // Ensure "will" matches your denomination

		component.ComponentType = &types.ExecutionComponent_Transfer{
			Transfer: &types.TransferComponent{
				// From:   sender,
				To:     to,
				Amount: &amountCoin,
			},
		}

	case "schnorr":
		dataParts := strings.Split(params, ",")
		if len(dataParts) != 3 {
			return nil, fmt.Errorf("invalid schnorr component params, expected 'public key, signature, message'")
		}
		signature, public_key, message := dataParts[0], dataParts[1], dataParts[2]
		fmt.Println("INSIDE parseComponentFromString for schnorr: ", public_key, signature, message)
		component.ComponentType = &types.ExecutionComponent_Claim{
			Claim: &types.ClaimComponent{
				SchemeType: &types.ClaimComponent_Schnorr{
					Schnorr: &types.SchnorrSignature{
						PublicKey: []byte(public_key),
						Signature: []byte(signature),
						Message:   message,
					},
				},
			},
		}
	case "pedersen":
		dataParts := strings.Split(params, ",")
		if len(dataParts) != 4 {
			return nil, fmt.Errorf("invalid pedersen component params, expected 'commitment,random factor, value, blinding factor'")
		}
		commitment, random_factor, value, blinding_factor := dataParts[0], dataParts[1], dataParts[2], dataParts[3]
		fmt.Println("INSIDE parseComponentFromString for pedersen: ", commitment, random_factor, value, blinding_factor)
		component.ComponentType = &types.ExecutionComponent_Claim{
			Claim: &types.ClaimComponent{
				SchemeType: &types.ClaimComponent_Pedersen{
					Pedersen: &types.PedersenCommitment{
						// note: removed these from the cli
						// BlindingFactor: []byte(blinding_factor),
						// Value:          []byte(value),
						Commitment: []byte(commitment),
					},
				},
			},
		}
	case "gnark":
		dataParts := strings.Split(params, ",")
		if len(dataParts) != 3 {
			return nil, fmt.Errorf("invalid gnark component params, expected 'verification key, public inputs, proof'")
		}
		verification_key, public_inputs, proof := dataParts[0], dataParts[1], dataParts[2]
		fmt.Println("INSIDE parseComponentFromString for pedersen: ", verification_key, public_inputs, proof)
		component.ComponentType = &types.ExecutionComponent_Claim{
			Claim: &types.ClaimComponent{
				SchemeType: &types.ClaimComponent_Gnark{
					Gnark: &types.GnarkZkSnark{
						VerificationKey: []byte(verification_key),
						PublicInputs:    []byte(public_inputs),
					},
				},
			},
		}
	default:
		return nil, fmt.Errorf("unsupported component type: %s", componentType)
	}

	return &component, nil
}

func ClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [will-id] [claim-type] [claim-data]",
		Short: "Submit a claim for a will",
		Long: `Submit a claim for a will with specific data based on the claim type.
Example:
./build/wasmd tx will claim "will-id" "component-id" "schnorr" "signature:data" --from alice --chain-id willchain-mainnet -y
./build/wasmd tx will claim "will-id" "component-id" "pedersen" "commitment:blinding_factor:value" --from alice --chain-id willchain-mainnet -y
./build/wasmd tx will claim "will-id" "component-id" "gnark" "proof:public_inputs" --from alice --chain-id willchain-mainnet -y`,
		Args: cobra.ExactArgs(4), // Ensuring exactly 3 arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			willID := args[0]
			componentID := args[1]
			claimType := args[2]
			claimData := args[3]

			// Construct the claim message based on the claim type
			var msg *types.MsgClaimRequest
			switch claimType {
			case "schnorr":
				// Parse the claim data for SchnorrClaim
				parts := strings.Split(claimData, ":")
				fmt.Println(parts)
				if len(parts) != 3 {
					return fmt.Errorf("invalid data format for Schnorr claim, expected 'signature:data'")
				}

				// pubKeyEncoded, _ := hex.DecodeString(parts[0])
				// signatureEncoded, _ := hex.DecodeString(parts[1])
				// messageEncoded, _ := hex.DecodeString(parts[2])

				msg = &types.MsgClaimRequest{
					WillId:      willID,
					Claimer:     clientCtx.GetFromAddress().String(),
					ComponentId: componentID,
					ClaimType: &types.MsgClaimRequest_SchnorrClaim{
						SchnorrClaim: &types.SchnorrClaim{
							Signature: []byte(parts[0]),
							PublicKey: []byte(parts[1]),
							Message:   parts[2],

							// PublicKey: pubKeyEncoded,
							// Signature: signatureEncoded,
							// Message:   messageEncoded,
						},
					},
				}

			case "pedersen":
				// Parse the claim data for PedersenClaim
				parts := strings.Split(claimData, ":")
				if len(parts) != 3 {
					return fmt.Errorf("invalid data format for Pedersen claim, expected 'commitment:blinding_factor:value'")
				}
				// Additional parsing and validation of parts[0], parts[1], and parts[2] needed here
				msg = &types.MsgClaimRequest{
					WillId:  willID,
					Claimer: clientCtx.GetFromAddress().String(),
					ClaimType: &types.MsgClaimRequest_PedersenClaim{
						PedersenClaim: &types.PedersenClaim{
							Commitment:     []byte(parts[0]),
							BlindingFactor: []byte(parts[1]),
							Value:          []byte(parts[2]),
						},
					},
				}

			case "gnark":
				// Parse the claim data for GnarkClaim
				parts := strings.Split(claimData, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid data format for Gnark claim, expected 'proof:public_inputs'")
				}
				// Additional parsing and validation of parts[0] and parts[1] needed here
				msg = &types.MsgClaimRequest{
					WillId:  willID,
					Claimer: clientCtx.GetFromAddress().String(),
					ClaimType: &types.MsgClaimRequest_GnarkClaim{
						GnarkClaim: &types.GnarkClaim{
							Proof:        []byte(parts[0]),
							PublicInputs: []byte(parts[1]),
						},
					},
				}

			default:
				return fmt.Errorf("unsupported claim type: %s", claimType)
			}

			// Submit the transaction
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

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
