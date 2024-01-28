# LightLink Hummingbird

![Test, Build Image & Push to ECR](https://github.com/pellartech/lightlink-hummingbird/actions/workflows/build_and_publish.yml/badge.svg?branch=main)

![LightLink Hummingbird preview screenshot](<preview.png>)

Hummingbird is a light client for interacting with the [LightLink Protocol](https://lightlink.io).

It is designed to work in unison with the [lightlink-hummingbird-contracts](https://github.com/pellartech/lightlink-hummingbird-contracts) repository.

## Installation

### Prerequisites

- [Golang](https://go.dev/dl/) (v1.21.5 or higher) installed. Go version can be checked with `$ go version`

### Option 1: Install from source (Linux / MacOS)

```bash
git clone git@github.com:pellartech/lightlink-hummingbird.git
cd lightlink-hummingbird
git checkout -b v0.0.1
make
```

### Option 2:  Install from binary (Windows / Linux / MacOS)

Download the latest release from [here](https://github.com/pellartech/lightlink-hummingbird/releases)

### Configuration

Hummingbird requires a configuration file to run. A sample configuration file is provided in the repository [here](config.example.json). Copy this file and rename it to `config.json`. Then fill in the required fields.

```bash
cp config.example.json config.json
```

**Note**: configuration file `config.json` path can be specified with the `--config-path` flag. If not specified, the default path is `./config.json`

```
{
  "storePath": "./storage", // Path to the local storage
  "celestia": {
    "token": "abcd", // Celestia token
    "endpoint": "", // Celestia rpc endpoint
    "namespace": "", // Celestia namespace to identify the blobs
    "tendermint_rpc": "", // Tendermint rpc endpoint
    "grpc": "", // Celestia grpc endpoint
    "gasPrice": 0.01 // Gas price to use when submitting new headers to Celestia
  },
  "ethereum": {
    "httpEndpoint": "", // Ethereum http rpc endpoint
    "wsEndpoint": "", // Ethereum websocket rpc endpoint
    "canonicalStateChain": "", // Canonical state chain contract address
    "daOracle": "", // Data availability oracle contract address
    "challenge": "", // Challenge contract address
    "gasPriceIncreasePercent": 100 // Increase gas price manually when submitting new headers to L1
  },
  "lightlink": {
    "endpoint": "", // LightLink rpc endpoint
    "delay": 100 // Delay in ms between each block query
  },
  "rollup": {
    "bundleSize": 200, // Number of headers to include in each bundle
    "l1pollDelay": 90000, // (90sec) Delay in ms between each L1 block query
    "l2pollDelay": 30000, // (30sec) Delay in ms between each L2 block query
    "storeCelestiaPointers": true, // Store celestia pointers in the local storage
    "storeHeaders": true // Store headers in the local storage
  },
  "defender": {
    "workerDelay": 60000 // (60sec) Delay in ms between scanning for historical challenges to defend 
  }
}
```

## Usage

```bash
hb rollup info  # Get the current rollup state
hb rollup next  # [Publisher Only] Generate the next rollup block
hb rollup start # [Publisher Only] Start the rollup loop to generate and submit bundles
hb challenge challenge-da <block_number> # Challenge data availability
hb defender defend-da <block_hash> # Defend data availability
hd defender info-da <block_hash> # Provides info on an existing challenge
hb defender prove-da <block_hash> # Prove data availability
hb defender start # Start the defender loop to watch and defend challenges
```

The following root flags are available for all commands:

```bash
--config-path <path> # Path to the config file
--log-level <level> # Log level (debug, info, warn, error)
--log-format <format> # Log format (json, text)
--log-source <bool> # Log source file and line
```

see `hb --help` for more information

<p align="center">
  <img src="humming.png" style="size:50%" alt="HummingBird">
</p>
