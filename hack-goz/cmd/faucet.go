package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/iqlusioninc/relayer/relayer"

	"github.com/irisnet/hack-goz/nets"
)

var (
	flagTimeout = "timeout"
	// precisionReuse = sdk.NewDecFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil))
)

func FaucetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "faucet",
		Short: "Hack Faucet",
		Long:  `Faucet allows you to bulk request tokens from other zones.`,
	}

	requestCmd := &cobra.Command{
		Use:   "request [prefix]",
		Short: "Bulk request tokens from other zones",
		Args:  cobra.ExactArgs(1),
		RunE:  runRequestCmd,
	}

	requestCmd.Flags().IntP(flagTimeout, "t", 10, "Timeout in seconds for requesting the faucet")

	collectCmd := &cobra.Command{
		Use:   "collect [prefix] [to]",
		Short: "Collect tokens from the addresses",
		Args:  cobra.ExactArgs(2),
		RunE:  runCollectCmd,
	}

	cmd.AddCommand(
		requestCmd,
		collectCmd,
	)

	return cmd
}

func runRequestCmd(cmd *cobra.Command, args []string) error {
	prefix := args[0]
	timeout, err := cmd.Flags().GetInt(flagTimeout)
	if err != nil {
		return err
	}
	kb, err := keyring.New(prefix, keyring.BackendTest, homePath, nil)
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err != nil {
		return err
	}

	for {
		for _, i := range infos {
			addrStatus := requestFaucet(i.GetAddress().String(), time.Duration(timeout)*time.Second)
			fmt.Printf("\n- Tokens for: %s\n", addrStatus.address)
			for _, chain := range addrStatus.chains {
				fmt.Printf("\n	%s: %s\n", chain.chainID, chain.result)
			}
		}
	}
}

func requestFaucet(addr string, timeout time.Duration) addressStatus {
	cc := make(chan chainStatus)
	for _, chain := range config.Chains {
		go func(chain *relayer.Chain, cc chan chainStatus) {
			chainStatus := &chainStatus{
				chainID: chain.ChainID,
				result:  "",
			}
			u, err := url.Parse(chain.RPCAddr)
			if err != nil {
				chainStatus.result = err.Error()
				cc <- *chainStatus
				return
			}

			host, _, err := net.SplitHostPort(u.Host)
			if err != nil {
				chainStatus.result = err.Error()
				cc <- *chainStatus
				return
			}

			urlString := fmt.Sprintf("%s://%s:%d", u.Scheme, host, 8000)

			body, err := json.Marshal(relayer.FaucetRequest{Address: addr, ChainID: chain.ChainID})
			if err != nil {
				chainStatus.result = err.Error()
				cc <- *chainStatus
				return
			}
			resp, err := nets.GetHTTPClient(timeout).Post(urlString, "application/json", bytes.NewBuffer(body))
			if err != nil {
				chainStatus.result = err.Error()
				cc <- *chainStatus
				return
			}
			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				chainStatus.result = err.Error()
				cc <- *chainStatus
				return
			}

			chainStatus.result = string(respBody)
			cc <- *chainStatus
		}(chain, cc)
	}

	chainResults := make([]chainStatus, len(config.Chains))
	for i := range chainResults {
		chainResults[i] = <-cc
	}

	return addressStatus{
		address: addr,
		chains:  chainResults,
	}
}

type addressStatus struct {
	address string
	chains  []chainStatus
}

type chainStatus struct {
	chainID string
	result  string
}

func runCollectCmd(cmd *cobra.Command, args []string) error {

	prefix := args[0]
	to, err := sdk.AccAddressFromBech32(args[1])
	if err != nil {
		return err
	}
	kb, err := keyring.New(prefix, keyring.BackendTest, homePath, nil)
	if err != nil {
		return err
	}

	infos, err := kb.List()
	if err != nil {
		return err
	}

	for _, info := range infos {
		for _, chain := range config.Chains {
			_ = chain.Init(homePath, chain.Cdc.Codec, chain.Amino.Codec, 2*time.Second, false)
			chain.Keybase = kb
			chain.Key = info.GetName()
			bal, err := chain.QueryBalance(chain.Key)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			gp := getGasPrices(chain).AmountOf(chain.DefaultDenom)
			fee := gp.MulInt64(int64(chain.Gas)).TruncateInt()
			amt := bal.AmountOf(chain.DefaultDenom).Sub(fee)

			if amt.GT(sdk.ZeroInt()) {
				msg := bank.NewMsgSend(info.GetAddress(), to, sdk.NewCoins(sdk.NewCoin(chain.DefaultDenom, amt)))
				sb, err := chain.BuildAndSignTx([]sdk.Msg{msg})
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				res, err := chain.BroadcastTxCommit(sb)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fmt.Println(res.String())
			}
		}
	}

	return nil
}
