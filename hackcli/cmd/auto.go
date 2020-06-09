package cmd

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/irisnet/hack-goz/cmd/monitor"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iqlusioninc/relayer/relayer"
)

var (
	dcon = "defaultconnectionid"
	dcha = "defaultchannelid"
	dpor = "defaultportid"
	dord = "ordered"
)

func AutoTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "auto-tx",
		Aliases: []string{"at"},
		Short:   "auto IBC transaction commands",
	}

	cmd.AddCommand(
		autoTransferCmd(),
		autoMultiTransferCmd(),
		autoUpdateClientCmd(),
	)

	return cmd
}

func autoTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [src-chain-id] [dst-chain-id] [path-name] [amount] [source] [dst-addr] [time]",
		Short: "auto raw send",
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

			period, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return err
			}
			if period == 0 {
				return fmt.Errorf("peroid must greater than 0")
			}

			if err := transferMsg(c[src], c[dst], amount, dstAddr, source); err != nil {
				c[src].Log(err.Error())
			}

			ticker := time.NewTicker(time.Second * time.Duration(period))
			for range ticker.C {
				if err := transferMsg(c[src], c[dst], amount, dstAddr, source); err != nil {
					c[src].Log(err.Error())
				}
			}

			return nil
		},
	}

	return cmd
}

func autoMultiTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-transfer [src-chain-id] [dst-chain-id] [path-name] [amount] [source] [dst-addr] [number] [time]",
		Short: "auto multi raw send",
		Args:  cobra.ExactArgs(8),
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

			period, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return err
			}
			if period == 0 {
				return fmt.Errorf("peroid must greater than 0")
			}

			if err := multiTransferMsg(c[src], c[dst], amount, dstAddr, source, number); err != nil {
				c[src].Log(err.Error())
			}

			ticker := time.NewTicker(time.Second * time.Duration(period))
			for range ticker.C {
				if err := multiTransferMsg(c[src], c[dst], amount, dstAddr, source, number); err != nil {
					c[src].Log(err.Error())
				}
			}

			return nil
		},
	}

	return cmd
}

func autoUpdateClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-client [src-chain-id] [dst-chain-id] [src-client-id] [peroid] [timeout]",
		Aliases: []string{"uc"},
		Short:   "auto update client for dst-chain on src-chain",
		Args:    cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			monitor.StartMonitor()

			src, dst := args[0], args[1]

			chains, err := config.Chains.Gets(src, dst)
			if err != nil {
				return err
			}

			if err = chains[src].AddPath(args[2], dcon, dcha, dpor, dord); err != nil {
				return err
			}

			period, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			if period == 0 {
				return fmt.Errorf("peroid must greater than 0")
			}

			timeout, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}
			if timeout == 0 {
				return fmt.Errorf("timeout must greater than 0")
			}

			// ================================================================

			metrics := monitor.PrometheusMetrics()

			success := updateClient(cmd, chains[src], chains[dst])
			metrics.Success.WithLabelValues().Set(success)
			if success != 1 {
				return fmt.Errorf("Update Client Failed")
			}

			ticker := time.NewTicker(time.Second * time.Duration(period))
			for range ticker.C {
				sDone := make(chan struct{})
				wg := new(sync.WaitGroup)
				wg.Add(1)

				go func() {
					success := updateClient(cmd, chains[src], chains[dst])
					metrics.Success.WithLabelValues().Set(success)
					if success != 1 {
						fmt.Printf("\nUpdate Client Failed\n")
					}
					wg.Done()
				}()

				go func() {
					wg.Wait()
					close(sDone)
				}()

				select {
				case <-sDone:
				case <-time.After(time.Duration(timeout) * time.Second):
					metrics.Success.WithLabelValues().Set(0)
					fmt.Printf("\nupdate-client timeout\n")
				}
			}

			return nil
		},
	}

	return cmd
}

func updateClient(
	cmd *cobra.Command,
	src *relayer.Chain,
	dst *relayer.Chain,
) float64 {
	success := float64(1)

	dstHeader, err := dst.UpdateLiteWithHeader()
	if err != nil {
		dst.Log(fmt.Sprintf("Update Header Error: %s", err.Error()))
		return 0
	}

	if res, err := src.SendMsgs([]sdk.Msg{src.PathEnd.UpdateClient(dstHeader, src.MustGetAddress())}); err != nil {
		src.Log(fmt.Sprintf("Update Client Error: %s", err.Error()))
		return 0
	} else if res.Code != 0 {
		src.Log(fmt.Sprintf("Update Client Error: %s", res.String()))
		return 0
	} else {
		src.Log(fmt.Sprintf("Update Client Success: %s", res.String()))
	}

	return success
}

func setPathsFromArgs(src, dst *relayer.Chain, name string) (*relayer.Path, error) {
	// Find any configured paths between the chains
	paths, err := config.Paths.PathsFromChains(src.ChainID, dst.ChainID)
	if err != nil {
		return nil, err
	}

	var path *relayer.Path
	if path, err = paths.Get(name); err != nil {
		return path, err
	}

	if err = src.SetPath(path.End(src.ChainID)); err != nil {
		return nil, err
	}

	if err = dst.SetPath(path.End(dst.ChainID)); err != nil {
		return nil, err
	}

	return path, nil
}
