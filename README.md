# Starbase Liquidity


## Env
* go1.22
* mysql

## Install dep tool
```shell
make dep
```

## Build binary
```shell
make build
```
  
## Test code
```shell
make test
```

## Format code
```shell
make fmt
```

## Auto fix code
```shell
make fix
```

## Lint code
```shell
make lint
```

## Module
* mysql sql
* redis(cache)
* ethscan
* watcher
* sender
* pathfinder


## Server

### Orderd
Provide public API in limit order list and details for Executor bots.
Free and limit-credit pull but charge for push.

### Keeperd
Integrate with the existing DEX aggregator Router entry to access liquidity and execute trades.
Utilize the DEX aggregator's API to fetch the best available prices and path for order execution.


## knowledge
* base chain(https://docs.base.org/docs/)
* eth abi(https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings)

## Q&A


## compile problems
kafka needs cgo, so make should cgo is enabled.
