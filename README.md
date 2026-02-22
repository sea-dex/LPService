# Starbase Liquidity Service (lpservice)

Real-time DEX liquidity tracking service that monitors on-chain events and reconstructs pool liquidity state off-chain. Designed for high-performance arbitrage detection and swap quoting across multiple decentralized exchanges.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Supported DEXes & Pool Types](#supported-dexes--pool-types)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Build](#build)
- [Run](#run)
- [Docker](#docker)
- [Development](#development)
- [Project Structure](#project-structure)
- [License](#license)

## Overview

`lpservice` subscribes to blockchain events (pool creation, swaps, mints, burns, fee changes) via RPC/WSS and reconstructs the full liquidity state of DEX pools in real-time. This enables:

- **Off-chain liquidity tracking** — Maintain accurate reserves (AMM) and tick-level liquidity (CAMM) without querying on-chain for every operation
- **Arbitrage detection** — Identify cross-pool arbitrage opportunities with optimal input amount calculation
- **Swap quoting** — Provide accurate swap quotes using locally maintained pool state
- **Pool discovery** — Automatically detect and track new pools as they are created on-chain

The service supports two event ingestion modes:
1. **Direct subscription** — Subscribe to blockchain events via RPC/WSS (LPTracker)
2. **Kafka pipeline** — Produce events to Kafka via `eventd`, consume and process via `liquid`

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Blockchain  │────▶│   eventd     │────▶│    Kafka    │
│  (RPC/WSS)  │     │ (event sync) │     │   (topics)  │
└──────┬───────┘     └──────────────┘     └──────┬──────┘
       │                                         │
       │  Direct mode                  Pipeline mode
       ▼                                         ▼
┌──────────────┐                        ┌──────────────┐
│  LPTracker   │                        │    liquid     │
│  (lpservice) │                        │  (consumer)  │
└──────┬───────┘                        └──────┬───────┘
       │                                       │
       ▼                                       ▼
┌──────────────────────────────────────────────────────┐
│                   Event Handler                       │
│  ┌────────────┐  ┌────────────┐  ┌────────────────┐  │
│  │  Factory    │  │   Pool     │  │  Arbitrage     │  │
│  │  Tracking   │  │  Liquidity │  │  Detection     │  │
│  └────────────┘  └────────────┘  └────────────────┘  │
└──────────────────────┬───────────────────────────────┘
                       │
                       ▼
                 ┌───────────┐
                 │   Redis   │
                 │  (state)  │
                 └───────────┘
```

### Services

| Binary | Entry Point | Description |
|--------|-------------|-------------|
| `lpservice` | [`cmd/LPTracker/main.go`](cmd/LPTracker/main.go) | Main service — subscribes to blockchain events directly, tracks pool liquidity, detects arbitrage opportunities. Exposes pprof on `:6060`. |
| `eventd` | [`cmd/eventd/main.go`](cmd/eventd/main.go) | Event daemon — polls blockchain for logs and publishes them to Kafka topics. |
| `liquid` | [`cmd/liquid/main.go`](cmd/liquid/main.go) | Kafka consumer — reads events from Kafka and processes them through the same event handler pipeline. |

## Supported DEXes & Pool Types

### AMM Pools (UniswapV2-style, constant product)

| Type ID | Name | Description |
|---------|------|-------------|
| 200 | AMM | Standard UniswapV2 / SushiswapV2 / BaseSwap pairs |
| 201 | AeroAMM | Aerodrome V1/V2 stable & volatile pools |
| 202 | InfusionAMM | Infusion Finance AMM pools |

### CAMM Pools (UniswapV3-style, concentrated liquidity)

| Type ID | Name | Description |
|---------|------|-------------|
| 300 | CAMM | UniswapV3 / SushiswapV3 / BaseSwapV3 / AlienBaseV3 pools |
| 301 | AeroCAMM | Aerodrome V3 (Slipstream) concentrated pools |
| 302 | PancakeCAMM | PancakeSwapV3 concentrated pools |

### Curve Pools

| Type ID | Name | Description |
|---------|------|-------------|
| 250 | Curve | Curve-style stable swap pools |

### Tracked Events

- **Pool Creation** — `PairCreated`, `PoolCreated` events from factory contracts
- **Sync** — Reserve updates for AMM pools
- **Swap** — Swap events across all pool types
- **Mint / Burn** — Liquidity add/remove for both AMM and CAMM pools
- **Collect** — Fee collection in CAMM pools
- **Initialize** — Initial price setting for CAMM pools
- **SetCustomFee** — Dynamic fee updates (Aerodrome, Algebra-style)

## Prerequisites

- **Go** ≥ 1.25
- **CGO** enabled (required for Kafka via librdkafka)
- **librdkafka** — C library for Kafka client
- **Redis** — State persistence and caching
- **Kafka** — Event streaming (required for `eventd`/`liquid` pipeline mode)
- **Ethereum RPC/WSS** — Blockchain node access

### macOS

```bash
brew install librdkafka
```

### Linux (Alpine)

```bash
apk add librdkafka-dev pkgconf
```

## Configuration

Copy the example configuration and customize:

```bash
cp config-example.toml config.toml
```

Key configuration sections in [`config-example.toml`](config-example.toml):

| Section | Description |
|---------|-------------|
| `env` | Environment: `local`, `dev`, `staging`, `prod` |
| `log` | Log level and format |
| `chain` | Chain ID, factory contracts, RPC providers, pool query contract |
| `chain.factory` | List of DEX factory contracts with vendor, type, fee, and address |
| `chain.providers` | RPC endpoint URLs with rate limits |
| `chain.pool_query` | Multicall contract for batch pool queries |
| `kafka` | Broker addresses, topic, and consumer group |
| `redis` | Redis connection (address, password, DB) |
| `http` | HTTP server listen address |
| `metrics` | Prometheus metrics endpoint |
| `sentry` | Sentry DSN for error reporting |

### Factory Configuration Example

```toml
[[chain.factory]]
vendor = "aerov2"
type = 201
fee = 0
addr = "0x420DD381b31aEf6683db6B902084cB0FFECe40Da"
```

## Build

```bash
# Install dependencies
make dep

# Build the lpservice binary
make build
```

The [`Makefile`](Makefile) produces the `lpservice` binary from [`cmd/LPTracker`](cmd/LPTracker/main.go) with version info injected via ldflags.

### Build Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the `lpservice` binary |
| `make dep` | Install Go dependencies and tools (golangci-lint, swag) |
| `make test` | Run tests |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code with gofumpt |
| `make fix` | Auto-fix lint issues |
| `make abi` | Generate Go bindings from Solidity ABI files |
| `make clean` | Remove build artifacts |

## Run

### LPTracker (Direct Mode)

```bash
./lpservice --config config.toml
```

CLI flags:

| Flag | Description |
|------|-------------|
| `--config`, `-c` | Path to config file |
| `--mode`, `-m` | Run mode: `reset` (re-sync from scratch) or `resume` (continue from last checkpoint) |

### eventd (Event Producer)

```bash
go run cmd/eventd/main.go --config config.toml
```

### liquid (Kafka Consumer)

```bash
go run cmd/liquid/main.go --config config.toml
```

## Docker

### Build Image

```bash
docker build -t lpservice .
```

The [`Dockerfile`](Dockerfile) uses a multi-stage build:
1. **Build stage** — `golang:1.23-alpine` with `librdkafka-dev` for CGO compilation
2. **Runtime stage** — Minimal Alpine image with the compiled binary

### Docker Compose

```bash
docker-compose up -d
```

See [`docker-compose.yml`](docker-compose.yml) for service definitions.

## Development

### Code Quality

```bash
# Format code
make fmt

# Lint
make lint

# Auto-fix lint issues
make fix

# Run tests
make test
```

### ABI Bindings

Generate Go contract bindings from ABI JSON files in [`contracts/abi/`](contracts/abi/):

```bash
make abi
```

This generates typed Go bindings for contracts including UniswapV2Pair, UniswapV3Pool, UniswapV3Factory, AeroV2Pool, PancakeV3Pool, Multicall3, ERC20/721/1155, and more.

### Key Internal Packages

| Package | Description |
|---------|-------------|
| [`liquid/swapor`](liquid/swapor/) | Core event handler — processes blockchain events, manages factories and pools |
| [`liquid/pool`](liquid/pool/) | Pool data structures, AMM/CAMM math, tick/bitmap management, swap simulation |
| [`liquid/events`](liquid/events/) | Blockchain event subscription (RPC/WSS) and Kafka producer/consumer |
| [`liquid/arb`](liquid/arb/) | Arbitrage pair detection, optimal amount calculation, cross-pool analysis |
| [`liquid/common`](liquid/common/) | Shared types (PoolType, Token, SwapFactory), constants, Redis helpers |
| [`liquid/handlers`](liquid/handlers/) | AMM and CAMM event handlers for liquidity state updates |
| [`liquid/quoter`](liquid/quoter/) | Swap quoting service |
| [`contracts`](contracts/) | Go bindings for on-chain contracts (generated from ABI) |
| [`config`](config/) | Configuration loading and validation |

## Project Structure

```
├── cmd/
│   ├── LPTracker/       # Main binary (lpservice)
│   ├── eventd/          # Event daemon (blockchain → Kafka)
│   ├── liquid/          # Kafka consumer mode
│   └── cmd.go           # Shared CLI runner
├── liquid/
│   ├── swapor/          # Core event processing engine
│   ├── pool/            # Pool types, math, liquidity tracking
│   ├── events/          # Event subscription & Kafka integration
│   ├── arb/             # Arbitrage detection & calculation
│   ├── common/          # Shared types & constants
│   ├── handlers/        # AMM/CAMM event handlers
│   └── quoter/          # Swap quote service
├── contracts/           # Solidity ABI bindings
│   └── abi/             # ABI JSON files
├── config/              # Configuration
├── deploy/              # CI/CD & deployment scripts
├── docs/                # API docs (Swagger) & design docs
├── version/             # Version info
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── config-example.toml
```

## License

This project is licensed under the [GNU General Public License v3.0](LICENSE).
