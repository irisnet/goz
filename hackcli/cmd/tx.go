package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

////////////////////////////////////////
////  RAW IBC TRANSACTION COMMANDS  ////
////////////////////////////////////////

func TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tx",
		Aliases: []string{"t"},
		Short:   "raw IBC transaction commands",
	}

	cmd.AddCommand(
		multiTransfer(),
	)

	return cmd
}

func multiTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-transfer [src-chain-id] [dst-chain-id] [path-name] [amount] [source] [dst-addr] [number]",
		Short: "multi transfer",
		Long:  "This sends tokens from a relayers configured wallet on chain src to a dst addr on dst",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]
			c, err := config.Chains.Gets(src, dst)
			if err != nil {
				return err
			}

			pth := args[2]
			if _, err = setPathsFromArgs(c[src], c[dst], pth); err != nil {
				return err
			}

			amount, err := sdk.ParseCoin(args[3])
			if err != nil {
				return err
			}

			source, err := strconv.ParseBool(args[4])
			if err != nil {
				return err
			}

			done := c[dst].UseSDKContext()
			dstAddr, err := sdk.AccAddressFromBech32(args[5])
			if err != nil {
				return err
			}
			done()

			number, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return err
			}
			if number == 0 {
				return fmt.Errorf("number must greater than 0")
			}

			return multiTransferMsg(c[src], c[dst], amount, dstAddr, source, number)
		},
	}
	return cmd
}

