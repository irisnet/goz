# IRISnet's work in GoZ

## Preparation stage

Before the start of the GoZ, the IRISnet team had conducted in-depth research on IBC and relayer. We fixed the issue of inaccurate error message in relayer: https://github.com/iqlusioninc/relayer/pull/108.
Then, we submit three important bugs to bug bounty of GoZ.

Bug list:

- `Consensus bug`
This bug will cause lost of some memory data after restart of gaia, resulting in consensus failure. To fix this problem, phase-1a has been delayed. After the end of phase-1a, we found that the bug appeared again, so we submit it to the organizer of GoZ agian. The organizer of GoZ decide to restart phase-1b after this bug was fixed to ensure that the game could run smoothly, and no break change in next updates during all phases.

- `ICS-20 Event bug`
If there are multiple messages in one transaction, the events data will be abnormal. The following msg events will contain all the previous msg events. So the events data will be very large and cause a series of problems. At the same time, relayer can't get the correct packet data by query.

- `Unordered channel bug`
ICS20 packets sent through unordered channel can be repeatedly relayed, causing more tokens minted on the counterparty chain than what's escrowed on the source chain.

In addition, the incentive mechanism is very important in the whole IBC ecosystem. We have proposed a draft of the incentive mechanism, hoping to provide some help to the development of the IBC ecosystem: https://github.com/cosmos/ics/issues/411.

## Phase 1a

After studying the rules of the game, we have developed automated tools and monitoring procedures, which mainly include the following points:

- We started our own private full node instead of using public node to broadcast
- Our automated program will automatically retry if the transaction fails
- We have real-time monitoring and alarm
- We have a plan to handle exceptions manually

We ended up in third place in `phase 1a`.

## Phase 1b

In phase2, we conducted an in-depth test of the ICS-20 event bug discovered previously, found a series of problems more serious. This bug will affect the performance of gozhub nodes seriously, leading to serious problems such as memory overflow, node stop, and timeout of query. And hackers can attack any channels easily at all stages of the game with this bug.

After communicating with the organizer, we attacked the public node after the end of the phase-1b with this bug, cause most of public nodes down. After that, we described the problems caused by this bug in details and fixed it: https://github.com/cosmos/cosmos-sdk/pull/6269

## Phase 2

In phase 2, we made many improvements and optimizations to the relayer, mainly in the following aspects:

- We have implemented automated batch sending of cross-chain transactions
- We have implemented automated batch relaying packets
- We improved query efficiency by using multi-goroutine in relaying packets
- we Added handling of exception such as timeout and transaction failure
- We fixed the issue that can't get the proof of the packets created in the latest block

According to the unofficial P2P [leaderboard](https://dash-goz.p2p.org/public/dashboards/qmf48DlWlQHpnuHg3dLvt7My1MkY7UoE5ru1Iljk?org_slug=default), we won the fourth place, though our statistics of `packets_from_hub` and `packets_from_zones` are missing.

## Phase 3

In phase3, we implemented a deceptive relay between gozhub and irishub, and simulated a series of abnormal scenes to illustrate the consequences and effects:

- https://github.com/irisnet/goz/blob/master/phase-3/scenes.md

We have written up how we did this:

- https://github.com/irisnet/goz/blob/master/phase-3/deceptive.md

and submit two pull requests on cosmos-sdk and relayer to fix it:

- https://github.com/cosmos/cosmos-sdk/pull/6337
- https://github.com/iqlusioninc/relayer/pull/271

## Cumulative Contest Challenge

**Rainbow-GoZ Wallet**

- https://ibc-goz.irisplorer.io/#/download-rainbowgoz

In anticipation of the GoZ Phase 1b launch, the Rainbow-GoZ Wallet is currently ready for users to try out and start experimenting with multi-chain support in wallets, crosschain transfers, and atomic token swaps via IRIS’ proprietary solution Coinswap.Check more details:

- https://medium.com/@irisnet/get-in-line-experience-crosschain-transfers-and-atomic-coinswaps-with-rainbow-goz-61cfc57365f9

**GoZ Network State Visualizer**

- https://ibc-goz.irisplorer.io/#/

With the numerous connections of data in the IBC testnet before the challenge, users are available to view an interchain universe! Check more details:

- https://medium.com/irisnet-blog/irisnet-team-updated-a-new-version-of-goz-network-state-visualizer-62e3d79486f5

## Other

**Participants:**

We are continuing to participate in GoZ, and  participants of  IRISnet Dev Team work on 3 directions: the adversarial competition, Rainbow-GoZ Wallet, and GoZ Network State Visualizer. Everyone has paid a lot for Game of Zones.

**Challenging:**

It's a global game that lasts for 3 weeks. Participants come from different time zones and there may be various temporary situations happening in this game, so it requires us to highly focus our attention ,to keep thinking. It's definitely a challenge.

Cause the data amount, the load carrying by each node of GoZ Hub  is very large, so how to ensure that our program can run stably in this case is also a huge challenge.

At the same time, if we want to compete with various excellent teams, we need to find ways to optimize our tools to obtain a better ranking. This is also a very exciting challenge for us.

**Expectations:**

GoZ  is a very good opportunity to participate and learn how to use the IBC protocol. In the competition, we competed with many excellent teams. We saw many excellent strategies and ideas, and learned a lot.  We do hope GoZ will build a solid foundation for the establishment and improvement of the whole Cosmos ecosystem.
