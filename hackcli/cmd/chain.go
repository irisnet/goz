package cmd

import (
	"time"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	libclient "github.com/tendermint/tendermint/rpc/lib/client"

	sdkCtx "github.com/cosmos/cosmos-sdk/client/context"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/iqlusioninc/relayer/relayer"
)

func getGasPrices(src *relayer.Chain) sdk.DecCoins {
	gp, _ := sdk.ParseDecCoins(src.GasPrices)
	return gp
}

// GetTrustingPeriod returns the trusting period for the chain
func GetTrustingPeriod(src *relayer.Chain) time.Duration {
	tp, _ := time.ParseDuration(src.TrustingPeriod)
	return tp
}

func NewRPCClient(addr string, timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	// TODO: Replace with the global timeout value?
	httpClient.Timeout = timeout
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}

// SendMsg wraps the msg in a stdtx, signs and sends it
func SendMsg(src *relayer.Chain, datagram sdk.Msg) (sdk.TxResponse, error) {
	return src.SendMsgs([]sdk.Msg{datagram})
}

// SendMsgs wraps the msgs in a stdtx, signs and sends it
func SendMsgs(src *relayer.Chain, datagrams []sdk.Msg) (res sdk.TxResponse, err error) {
	var out []byte
	if out, err = src.BuildAndSignTx(datagrams); err != nil {
		return res, err
	}
	return src.BroadcastTxCommit(out)
}

// BuildAndSignTx takes messages and builds, signs and marshals a sdk.Tx to prepare it for broadcast
func BuildAndSignTx(src *relayer.Chain, datagram []sdk.Msg) ([]byte, error) {
	// Fetch account and sequence numbers for the account
	acc, err := auth.NewAccountRetriever(src.Cdc, src).GetAccount(src.MustGetAddress())
	if err != nil {
		return nil, err
	}

	defer src.UseSDKContext()()
	return auth.NewTxBuilder(
		auth.DefaultTxEncoder(src.Amino.Codec), acc.GetAccountNumber(),
		acc.GetSequence(), src.Gas, src.GasAdjustment, false, src.ChainID,
		src.Memo, sdk.NewCoins(), getGasPrices(src)).WithKeybase(src.Keybase).
		BuildAndSign(src.Key, ckeys.DefaultKeyPass, datagram)
}

// BroadcastTxCommit takes the marshaled transaction bytes and broadcasts them
func BroadcastTxCommit(src *relayer.Chain, txBytes []byte) (sdk.TxResponse, error) {
	res, err := sdkCtx.CLIContext{Client: src.Client}.BroadcastTxCommit(txBytes)
	return res, err
}
