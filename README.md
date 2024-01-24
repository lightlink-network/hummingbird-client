# LightLink Hummingbird

![Test, Build Image & Push to ECR](https://github.com/pellartech/lightlink-hummingbird/actions/workflows/build_and_publish.yml/badge.svg?branch=main)

![LightLink Hummingbird preview screenshot](<preview.png>)

Hummingbird is a command line tool for interacting with the LightLink protocol.

It is designed to work in unison with the [lightlink-hummingbird-contracts](https://github.com/pellartech/lightlink-hummingbird-contracts) repository.

## Commands

```bash
hb rollup info  # Get the current rollup state
hb rollup next  # Generate the next rollup block
hb rollup start # Start the rollup loop to generate and submit bundles
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