# Deceptive Relay And Fix

This article describes the problem and fix of `deceptive relay` in ICS-20 [Abnormal Scenes](./scenes.md) we wrote before.

## Deceptive relay

ICS20 packets sent through unordered channel can be repeatedly relayed, causing more tokens minted on the counterparty chain than what's escrowed on the  source chain.

For unordered channel, a blockchain will save a corresponding acknowledgement after receiving a packet. And the blockchain should check whether a received packet has already been processed before executing the packet.

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

In `PacketExecuted()`, it only stores PacketAcknowledgement without checking whether a PacketAcknowledgement already exists. So all packets can be executed repeatedly. We have fixed this bug here: https://github.com/cosmos/cosmos-sdk/pull/6337

### Relayer

For unordered channel, cmd `rly tx rly [path]` and `rly start [path]` will relay all packets no matter they have been relayed or not.

In ordered channel, receive-sequence and send-sequence must be sequential, so we can get the unrelayed sequences by `[src-next-send-seq] - [dst-next-recv-seq]`.

But for unordered channel, we can only get `src-next-send-seq`. So we can't get all unrelayed sequences though the method in ordered channel.

We resolved this problem here: https://github.com/iqlusioninc/relayer/pull/271

We filter the sequences by whether the packet acknowledgements exist, so that relayer can work correctly.

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
