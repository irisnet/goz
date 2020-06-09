package cmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iqlusioninc/relayer/relayer"
)

func transferMsg(src *relayer.Chain, dst *relayer.Chain, amount sdk.Coin, dstAddr sdk.AccAddress, source bool) error {
	if source {
		amount.Denom = fmt.Sprintf("%s/%s/%s", dst.PathEnd.PortID, dst.PathEnd.ChannelID, amount.Denom)
	} else {
		amount.Denom = fmt.Sprintf("%s/%s/%s", src.PathEnd.PortID, src.PathEnd.ChannelID, amount.Denom)
	}

	// Properly render the address string
	done := dst.UseSDKContext()
	dstAddrString := dstAddr.String()
	done()

	// MsgTransfer will call SendPacket on src chain
	txs := relayer.RelayMsgs{
		Src: []sdk.Msg{src.PathEnd.MsgTransfer(
			dst.PathEnd, uint64(9000000000), sdk.NewCoins(amount), dstAddrString, src.MustGetAddress(),
		)},
		Dst: []sdk.Msg{},
	}

	if txs.Send(src, dst); !txs.Success() {
		return fmt.Errorf("failed to send transfer message")
	}
	return nil
}

func multiTransferMsg(src *relayer.Chain, dst *relayer.Chain, amount sdk.Coin, dstAddr sdk.AccAddress, source bool, number uint64) error {
	if source {
		amount.Denom = fmt.Sprintf("%s/%s/%s", dst.PathEnd.PortID, dst.PathEnd.ChannelID, amount.Denom)
	} else {
		amount.Denom = fmt.Sprintf("%s/%s/%s", src.PathEnd.PortID, src.PathEnd.ChannelID, amount.Denom)
	}

	// Properly render the address string
	done := dst.UseSDKContext()
	dstAddrString := dstAddr.String()
	done()

	msgs := []sdk.Msg{}
	for i := uint64(0); i < number; i++ {
		msgs = append(msgs, src.PathEnd.MsgTransfer(
			dst.PathEnd,
			uint64(9000000000),
			sdk.NewCoins(amount),
			dstAddrString,
			src.MustGetAddress(),
		))
	}

	// MsgTransfer will call SendPacket on src chain
	txs := relayer.RelayMsgs{
		Src: msgs,
		Dst: []sdk.Msg{},
	}

	if txs.Send(src, dst); !txs.Success() {
		return fmt.Errorf("failed to send transfer message")
	}
	return nil
}
