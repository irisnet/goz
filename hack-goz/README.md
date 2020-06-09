# hack-goz

A command line tool for GoZ

## Install

```bash
make install
```

## Init

Same as the [relayer](https://github.com/iqlusioninc/relayer)

## Tx

**multi transfer**

```bash
hackcli tx multi-transfer [src-chain-id] [dst-chain-id] [path-name] [amount] [source] [dst-addr] [number]
```

## Auto Tx

### auto update client

```bash
hackcli auto-tx update-client [src-chain-id] [dst-chain-id] [src-client-id] [peroid] [timeout]
```

### auto transfer

```bash
hackcli auto-tx transfer [src-chain-id] [dst-chain-id] [path-name] [amount] [source] [dst-addr] [time]
```

### auto multi transfer

```bash
hackcli auto-tx multi-transfer [src-chain-id] [dst-chain-id] [path-name] [amount] [source] [dst-addr] [number] [time]
```

## Monitor & Alert

Prometheus Metrics will be exposed at port 8080 when Auto Tx is started
