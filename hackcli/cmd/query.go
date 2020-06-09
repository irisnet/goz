package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	flagJSON = "json"
)

// queryCmd represents the chain command
func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "IBC Query Commands",
		Long:    "Commands to query IBC primatives, and other useful data on configured chains.",
	}

	cmd.AddCommand(
		queryBalanceCmd(),
	)

	return cmd
}

func queryBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "balance [chain-id] [[key-name]]",
		Aliases: []string{"bal"},
		Short:   "Query the account balances",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			jsn, err := cmd.Flags().GetBool(flagJSON)
			if err != nil {
				return err
			}

			var keyName string
			if len(args) == 2 {
				keyName = args[1]
			}

			coins, err := chain.QueryBalance(keyName)
			if err != nil {
				return err
			}

			var out string
			if jsn {
				byt, err := json.Marshal(coins)
				if err != nil {
					return err
				}
				out = string(byt)
			} else {
				out = coins.String()
			}

			fmt.Println(out)
			return nil
		},
	}

	cmd.Flags().BoolP(flagJSON, "j", false, "returns the response in json format")

	return cmd
}
