# Deceptive And Fix

This article described the method and fix of `deceptive relay` in ICS-20 [Abnormal Scenes](./scenes.md) we write before.

## deceptive relay

ICS20 packets sent through unordered channel can be repeated relayed, causes the number of cross-chain tokens minted on counterparty chain is greater than the number of escrowed native tokens on source chain.

For unordered channel, the chain will save packet acknowledgement after receiving the packet. And the chain should check whether the packet has been received before packet execution.

https://github.com/cosmos/cosmos-sdk/blob/master/x/ibc/04-channel/keeper/packet.go#L243

```go
func (k Keeper) PacketExecuted(
    ctx sdk.Context,
    chanCap *capability.Capability,
    packet exported.PacketI,
    acknowledgement []byte,
) error {
    ...
    if acknowledgement != nil || channel.Ordering == types.UNORDERED {
        k.SetPacketAcknowledgement(
            ctx, packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence(),
            types.CommitAcknowledgement(acknowledgement),
        )
    }
    ...
}
```

In `PacketExecuted()`, it only stores PacketAcknowledgement without check whether the PacketAcknowledgement exists. So all packets can be received repeatedly. We have fixed this bug here: https://github.com/cosmos/cosmos-sdk/pull/6337

### relayer

For unordered channel, cmd `rly tx rly [path]` and `rly start [path]` will relay all packets every time they are executed, whether they have been relayed or not.

In ordered channel, receive-sequence and send-sequence must be sequential, so we can get the unrelayed sequences by `[src-next-send-seq] - [dst-next-recv-seq]`.

But for unordered channel, we can only get `src-next-send-seq`. So we can't get all unrelayed sequences though the method in ordered channel.

We resolved this problem here: https://github.com/iqlusioninc/relayer/pull/271

We filter the sequences by whether the acknowledgement of packet exists, so that relayer can work normally.

```go
func newRlySeq(chain *Chain, start, end uint64, ordered bool, height uint64) []uint64 {
    if end < start {
        return []uint64{}
    }
    s := make([]uint64, 0, 1+(end-start))
    for start < end {
        if ordered {
            s = append(s, start)
            start++
        } else {
            ack, err := chain.QueryPacketAck(int64(height), int64(start))
            if err != nil {
                chain.Log(err.Error())
            } else if ack.Data == nil {
                s = append(s, start)
            }
            start++
        }
    }
    return s
}
```
