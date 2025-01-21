# Selfchain

Self Chain is the first Modular Intent-Centric Access Layer1 blockchain and keyless wallet infrastructure service using MPC-TSS/AA for multi-chain Web3 access. The innovative system simplifies the user experience with its intent-focused approach, using LLM to interpret user intent and discover the most efficient paths.

## Key Features
- **Keyless Wallets**: Effortless onboarding and recovery with complete self-custody over assets
- **Intent-Centric Access**: LLM-powered system for interpreting user intent and finding optimal paths
- **Account Abstraction**: Integrated with MPC-TSS for secure signing and reduced transaction fees
- **Delegated Proof of Stake**: Secured by Tendermint consensus mechanism
- **Cosmos SDK Based**: Built with Cosmos SDK v0.47.10 for robust blockchain development
- **CosmWasm Support**: Smart contract functionality with WASM
- **IBC Enabled**: Inter-Blockchain Communication for cross-chain interactions

## Custom Modules

Selfchain implements custom modules to extend the base Cosmos SDK functionality:

### Migration Module (x/migration)

The Migration module facilitates chain upgrades and state migrations with the following features:

- **Chain Migration Process**:
  - Manages authorized migrators
  - Validates migration transactions
  - Tracks migration records
  - Handles token migration between chains

- **Key Components**:
  - Migrator Management: Add/remove authorized migrators
  - Token Migration: Process token transfers with validation
  - Transaction Verification: Verify cross-chain transactions
  - State Management: Track migration status and records

- **Integration Points**:
  - Interfaces with selfvesting module
  - Bank keeper integration for token management
  - Parameter management for configuration

### Selfvesting Module (x/selfvesting)

The Selfvesting module implements advanced token vesting functionality:

- **Vesting Features**:
  - Custom vesting schedules
  - Cliff period implementation
  - Multiple token support
  - Position tracking
  - Automated distribution

- **Key Components**:
  - Schedule Management: Create and track vesting schedules
  - Token Distribution: Handle automated token releases
  - Position Tracking: Monitor vesting positions
  - Claim Processing: Process token claims

- **Security Features**:
  - Cliff period enforcement
  - Position validation
  - Balance checks
  - Access control

### Smart Contract Support (CosmWasm)

Integrated WASM smart contract functionality with:

- **Features**:
  - Contract deployment and execution
  - Gas metering and limits
  - Query support
  - State management
  - Cross-contract calls

- **Security**:
  - Contract size limits
  - Memory restrictions
  - Access control
  - Gas management

### Interchain Communication (IBC)

Full IBC protocol implementation for cross-chain operations:

- **Features**:
  - Token transfers
  - Cross-chain communication
  - Interchain accounts
  - Packet relay
  - Channel management

- **Security**:
  - Packet verification
  - Timeout handling
  - State verification
  - Channel authentication

## Prerequisites

- [Go](https://go.dev/) 1.22 or later
- [Ignite CLI](https://docs.ignite.com/welcome/install)
- [jq](https://stedolan.github.io/jq/) (optional) for JSON parsing

## Installation

```bash
# Clone the repository
git clone https://github.com/selfchainxyz/selfchain.git
cd selfchain

# Build the chain
make install

# Check installation
selfchaind version
```

## Ignite CLI Commands

Ignite CLI provides a suite of commands to scaffold and manage your blockchain:

```bash
# Create a new chain
ignite scaffold chain selfchain

# Add a new module with keeper and messages
ignite scaffold module migration --dep bank,staking

# Add a new message with fields
ignite scaffold message create-vesting name:string amount:coin start_time:uint end_time:uint --module selfvesting

# Add a new query with response fields
ignite scaffold query get-vesting address:string --response amount:coin,start_time:uint,end_time:uint --module selfvesting

# Generate only proto files
ignite generate proto-go

# Start chain with automatic reloading
ignite chain serve --verbose

# Build for production
ignite chain build --release

# Reset chain data
ignite chain serve --reset-once

# Start with custom config
ignite chain serve --config custom.yml
```

## Selfchaind Commands

The `selfchaind` binary provides various commands for interacting with the blockchain:

### Basic Commands
```bash
# Show version information
selfchaind version

# Show all available commands
selfchaind --help

# Start the node (run the blockchain)
selfchaind start

# Start node with custom config
selfchaind start --home ./custom_node_home

# Start node with Prometheus metrics
selfchaind start --metrics

# Show node status
selfchaind status

# Query node info
selfchaind status --node tcp://localhost:26657
```

### Keys and Accounts
```bash
# Add a new key
selfchaind keys add validator_key

# Add a key with custom keyring backend
selfchaind keys add validator_key --keyring-backend test

# Recover an existing key
selfchaind keys add validator_key --recover

# List all keys
selfchaind keys list

# Delete a key
selfchaind keys delete validator_key

# Export a key (save to file)
selfchaind keys export validator_key

# Import a key
selfchaind keys import validator_key validator_key.backup

# Show account information
selfchaind query account self1... --node tcp://localhost:26657

# Add a genesis account
selfchaind add-genesis-account self1... 1000000000uself,1000000000stake

# Add a genesis account with vesting
selfchaind add-genesis-account self1... 1000000000uself --vesting-amount 500000000uself --vesting-start-time 1624582800 --vesting-end-time 1624669200
```

### Transaction Commands
```bash
# Send tokens
selfchaind tx bank send \
    $(selfchaind keys show my_account -a) \
    self1... \
    1000000uself \
    --chain-id self-1 \
    --gas auto \
    --gas-prices 0.025uself

# Delegate tokens to validator
selfchaind tx staking delegate \
    selfvaloper1... \
    1000000uself \
    --from my_account \
    --chain-id self-1 \
    --gas auto \
    --gas-prices 0.025uself

# Redelegate tokens to another validator
selfchaind tx staking redelegate \
    selfvaloper1... \
    selfvaloper1... \
    1000000uself \
    --from my_account \
    --chain-id self-1

# Unbond tokens
selfchaind tx staking unbond \
    selfvaloper1... \
    1000000uself \
    --from my_account \
    --chain-id self-1

# Submit governance proposal
selfchaind tx gov submit-proposal \
    --title "Test Proposal" \
    --description "This is a test proposal" \
    --type text \
    --deposit 1000000uself \
    --from my_account \
    --chain-id self-1

# Vote on proposal
selfchaind tx gov vote 1 yes \
    --from my_account \
    --chain-id self-1
```

### Query Commands
```bash
# Query bank balances
selfchaind query bank balances self1...

# Query bank total supply
selfchaind query bank total

# Query staking validators
selfchaind query staking validators

# Query specific validator
selfchaind query staking validator selfvaloper1...

# Query delegations to validator
selfchaind query staking delegations-to selfvaloper1...

# Query delegator delegations
selfchaind query staking delegations self1...

# Query governance proposals
selfchaind query gov proposals

# Query specific proposal
selfchaind query gov proposal 1

# Query proposal votes
selfchaind query gov votes 1

# Query proposal tally
selfchaind query gov tally 1
```

### Vesting Commands
```bash
# Create a vesting account
selfchaind tx vesting create-vesting-account \
    self1... \
    1000000uself \
    1624582800 \
    --from my_account \
    --chain-id self-1

# Query vesting account
selfchaind query auth account self1...

# Query vesting period
selfchaind query vesting periods self1...
```

### IBC Commands
```bash
# Query IBC client state
selfchaind query ibc client state <client_id>

# Query IBC client states
selfchaind query ibc client states

# Query IBC connection
selfchaind query ibc connection end <connection_id>

# Query IBC connections
selfchaind query ibc connection connections

# Query IBC channel
selfchaind query ibc channel end <channel_id> <port_id>

# Query IBC channels
selfchaind query ibc channel channels

# Transfer tokens via IBC
selfchaind tx ibc-transfer transfer \
    transfer \
    <dst_channel_id> \
    <recipient_address> \
    1000000uself \
    --from my_account \
    --chain-id self-1 \
    --gas auto \
    --gas-prices 0.025uself
```

### CosmWasm Commands
```bash
# Store WASM code
selfchaind tx wasm store contract.wasm \
    --from my_account \
    --chain-id self-1 \
    --gas auto \
    --gas-prices 0.025uself \
    --instantiate-everybody true

# Instantiate contract
selfchaind tx wasm instantiate 1 \
    '{"name": "My Contract", "symbol": "MC"}' \
    --label "my_contract" \
    --from my_account \
    --chain-id self-1 \
    --gas auto \
    --gas-prices 0.025uself \
    --admin $(selfchaind keys show my_account -a)

# Execute contract
selfchaind tx wasm execute \
    <contract_address> \
    '{"transfer": {"recipient": "self1...", "amount": "1000"}}' \
    --from my_account \
    --chain-id self-1 \
    --gas auto \
    --gas-prices 0.025uself

# Query contract
selfchaind query wasm contract-state smart \
    <contract_address> \
    '{"balance": {"address": "self1..."}}'

# List all contracts
selfchaind query wasm list-contract-by-code <code_id>

# Get contract info
selfchaind query wasm contract <contract_address>

# Get contract history
selfchaind query wasm contract-history <contract_address>
```

## Running the Chain

### Development Mode
To start the chain in development mode, use:
```bash
CHAIN_ENV=development selfchaind start
```

### Production Mode
The chain has undergone a governance proposal that included state changes. When running in production, you'll need to:

1. Use the older binary to sync the chain up to the proposal block
2. Upgrade to the current binary after reaching the proposal block

This ensures proper state synchronization through the upgrade process.

Note: The specific version requirements and upgrade block height will be provided in the network upgrade documentation.

## Command Reference

### Migration Commands
```bash
# Query migration params
selfchaind query migration params

# Query migration status
selfchaind query migration status

# Submit migration proposal
selfchaind tx gov submit-proposal \
    --title "Chain Migration" \
    --description "Migrate chain to new version" \
    --type migration \
    --from my_account \
    --chain-id self-1

# Add authorized migrator
selfchaind tx migration add-migrator \
    self1... \
    --from validator_key \
    --chain-id self-1

# Remove authorized migrator
selfchaind tx migration remove-migrator \
    self1... \
    --from validator_key \
    --chain-id self-1

# Migrate tokens
selfchaind tx migration migrate \
    1000000uself \
    0x123... \
    self1... \
    --tx-hash 0xabc... \
    --from my_account \
    --chain-id self-1

# Query migrator list
selfchaind query migration migrators

# Query migration records
selfchaind query migration records \
    --address self1... \
    --page 1 \
    --limit 100

# Query total migrated amount
selfchaind query migration total-migrated
```

### Vesting Commands
```bash
# Create a vesting schedule
selfchaind tx selfvesting create-vesting-schedule \
    self1... \
    1000000uself \
    1624582800 \
    1624669200 \
    --cliff-time 1624583800 \
    --from my_account \
    --chain-id self-1

# Create multiple vesting schedules
selfchaind tx selfvesting create-vesting-schedules \
    self1... \
    1000000uself,2000000stake \
    1624582800,1624582800 \
    1624669200,1624669200 \
    --from my_account \
    --chain-id self-1

# Query all vesting schedules
selfchaind query selfvesting schedules self1...

# Query specific vesting schedule
selfchaind query selfvesting schedule self1... 1

# Query vesting account balance
selfchaind query selfvesting balance self1...

# Query claimable amount
selfchaind query selfvesting claimable self1... 1

# Claim vested tokens
selfchaind tx selfvesting claim \
    --from my_account \
    --chain-id self-1

# Query vesting parameters
selfchaind query selfvesting params

# Query total vesting amount
selfchaind query selfvesting total-vesting

# Query vesting status
selfchaind query selfvesting status self1...
```

### Validator Management
```bash
# Create validator
selfchaind tx staking create-validator \
    --amount 1000000uself \
    --commission-max-change-rate 0.01 \
    --commission-max-rate 0.2 \
    --commission-rate 0.1 \
    --from validator_key \
    --min-self-delegation 1 \
    --moniker "my_validator" \
    --pubkey $(selfchaind tendermint show-validator) \
    --chain-id self-1

# Edit validator
selfchaind tx staking edit-validator \
    --commission-rate 0.15 \
    --moniker "new_validator_name" \
    --from validator_key \
    --chain-id self-1

# Unjail validator
selfchaind tx slashing unjail \
    --from validator_key \
    --chain-id self-1
```

## Manual Chain Initialization

To start a new chain manually, follow these steps:

```bash
# Build the binary
go build -o ./selfchaind cmd/selfchaind/main.go

# Initialize the chain
# 'speeder' is the moniker (name) for your node
./selfchaind init speeder --chain-id 1

# Create a new key (account)
./selfchaind keys add bartsimp

# Add genesis account with initial tokens
# uslf is the native token denomination (micro SLF, 1 SLF = 1,000,000 uslf)
./selfchaind add-genesis-account bartsimp 10000000000uslf --keyring-backend os

# Create genesis transaction
# This stakes tokens and creates the first validator
./selfchaind gentx bartsimp 100000000uslf --chain-id 1

# Collect genesis transactions
./selfchaind collect-gentxs

# Start the chain
./selfchaind start
```

You can customize these values:
- `speeder`: Change to your preferred node moniker
- `chain-id`: Use any unique identifier for your chain
- `bartsimp`: Replace with your preferred key name
- `10000000000uslf`: Adjust the initial token allocation
- `100000000uslf`: Modify the amount of tokens to stake

After initialization, you can find the configuration files in:
- `~/.selfchain/config/`: Chain configuration
- `~/.selfchain/data/`: Chain data
- `~/.selfchain/keyring-os/`: Keyring (when using os backend)

## Chain Information

- Chain ID: `self-1`
- Binary: `selfchaind`
- Bech32 Prefix: `self`
- Dependencies:
  - Cosmos SDK: v0.47.10
  - CometBFT: v0.37.4
  - CosmWasm: Enabled
  - IBC: Enabled

## Consensus

Self Chain uses Delegated Proof of Stake (DPoS) consensus secured by Tendermint. The consensus process works as follows:

1. A validator (proposer) is chosen to submit a new block of transactions
2. Validators vote in two rounds on whether to accept or reject the proposed block
3. If accepted, the block is signed and added to the chain
4. Transaction fees are distributed as staking rewards to validators and delegators
5. Proposers receive additional rewards for block creation

The top 100 validators by stake participate in consensus. Validators can bond $SLF tokens themselves and also receive delegations from token holders.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.