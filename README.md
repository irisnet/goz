# ICS-20 异常场景

本文通过一系列连续场景展示了 IBC 协议中 ICS-20 存在的异常状况。通过 Alice、Bob、Charlie、Dave、Eve 在 gozhub、irishub、otherhub 中的交互进行说明。

本场景中所有 hub 均使用 goz 第三阶段官方 gaia 版本，未修改代码。

## 初始状态

**path**

`gozhub` - `irishub`

- gozhub
  - port: `transfer`
  - channel: `irisgozchann`

- irishub
  - port: `transfer`
  - channel: `gozirischann`

`gozhub` - `otherhub`

- gozhub
  - port: `transfer`
  - channel: `othergozchann`

- irishub
  - port: `transfer`
  - channel: `gozotherchann`

**balance**:

- Alice
  - irishub: `10uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Bob
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Charlie
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Dave
  - irishub: `10uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Eve
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`

## 场景

### 场景 1

Alice 在 irishub 上有 10 个 `uiris`，Bob 想要购买三个。irishub 和 gozhub 已经完成了 IBC 连接，为了方便资产管理，Bob 决定所有的资产都存放在 gozhub 中，交付方式为 Alice 向 Bob 在 gozhub 中的账户进行跨链转账。

Alice 实际在跨链转账中只发送了 `1uiris`，但是通过欺骗性的relay方式，成功向 Bob 在 gozhub 的账户中转入了 `3transfer/gozirischann/uiris`。但 Bob 并不知情，认为自己的 token 是没问题的。

![scene1](/asset/scene1.png)

场景结束后余额状态：

- Alice
  - irishub: `9uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Bob
  - irishub: `[]`
  - gozhub: `3transfer/gozirischann/uiris`
  - otherhub: `[]`
- Charlie
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Dave
  - irishub: `10uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Eve
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`

### 场景 2

一段时间之后，市场产生变化，Bob 和 Charlie 达成交易将手中的 `3transfer/gozirischann/uiris` 进行转让。Charlie 的资产主要在 otherhub 中，此时 otherhub 也已经和 gozhub 成功建立了连接，交付方式为 Bob 向 Charlie 在 otherhub 中的账户 进行跨链转账。

交易结束后，Charlie 在 otherhub 的中户中成功收到了 `3transfer/othergozchann/transfer/gozirischann/uiris`。Charlie 对 Alice 和 Bob 之间的交易并不知情，并认为收到的 token 和 irishub 中的 `3uiris` 是等值的。

![scene2](/asset/scene2.png)

场景结束后余额状态

- Alice
  - irishub: `9uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Bob
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Charlie
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `3transfer/othergozchann/transfer/gozirischann/uiris`
- Dave
  - irishub: `10uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Eve
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`

### 场景 3

又过一段时间，Bob 决定将转让给 Charlie 的 token 买回，交付方式为通过跨链转账的方式将 Charlie 在 otherhub 中的 `3transfer/othergozchann/transfer/gozirischann/uiris` 转回 Bob 在 gozhub 中的账户。

交易顺利完成，场景1中伪造的 token 在 IBC 网络中能够正常流通。

![scene3](/asset/scene3.png)

场景结束后余额状态

- Alice
  - irishub: `9uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Bob
  - irishub: `[]`
  - gozhub: `3transfer/gozirischann/uiris`
  - otherhub: `[]`
- Charlie
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Dave
  - irishub: `10uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Eve
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`

### 场景 4

Dave 在 irishub 中有 `10uiris`，某天他的好友 Eve 想要借 `2uiris` 转入 gozhub 使用，两人约定通过跨链转账的方式交付。交付完成后，Eve 在 gozhub 中的账户收到 `2transfer/gozirischann/uiris`。

![scene4](/asset/scene4.png)

场景结束后余额状态

- Alice
  - irishub: `9uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Bob
  - irishub: `[]`
  - gozhub: `3transfer/gozirischann/uiris`
  - otherhub: `[]`
- Charlie
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Dave
  - irishub: `8uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Eve
  - irishub: `[]`
  - gozhub: `2transfer/gozirischann/uiris`
  - otherhub: `[]`

### 场景 5

一段时间后，市场变化，Alice 同 Bob 达成交易买回 Bob 手中所有的 uiris，通过跨链转账的方式将 Bob 在 gozhub 中的 `3transfer/gozirischann/uiris` 转到 Alice 在 irishub 中的账户。交易完成后 Alice 在 irishub 中拥有 `12uiris`。

不久之后，Eve 同 Dave 约定还款，Eve 将其在 gozhub 中的 `2transfer/gozirischann/uiris` 通过跨链方式转到 Dave 在 irishub 中的账户。Eve 在 gozhub 发送交易成功，账户余额减少 `2transfer/gozirischann/uiris`，但是 Dave 却无法收到 token，Relay packet 成功，但 execute packet 时托管账户余额不足导致失败。

![scene5](/asset/scene5.png)

场景结束后余额状态

- Alice
  - irishub: `12uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Bob
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Charlie
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`
- Dave
  - irishub: `8uiris`
  - gozhub: `[]`
  - otherhub: `[]`
- Eve
  - irishub: `[]`
  - gozhub: `[]`
  - otherhub: `[]`

## 结论

最终结果表现为 Dave 的资产无法返回，Alice 侵占了 本应属于 Dave 的 `2uiris`。

场景1中 Alice 在 gozhub 中增发的虚假 token 同正常的跨链 token 完全一致，可以在整个 IBC 网络中正常流通。

当跨链 token 返回原链时，能够从委托账户中释放的原生 token 数量小于对方链发行的跨链 token 总数，首先进行兑换的能够跨链返回成功，最后兑换的部分将无法成功。
