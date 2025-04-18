#!/bin/bash

# Exit on error
set -e

# Parse command line arguments
NODE_ALREADY_RUNNING=false

# Process command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --node-running)
      NODE_ALREADY_RUNNING=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--node-running]"
      exit 1
      ;;
  esac
done

# =====================================================
# FRESH SETUP INSTRUCTIONS (Run these for a new chain)
# =====================================================
# selfchainold init mynode --chain-id selfchain-1
#
# # Replace stake with uslf in genesis file
# sed -i '' 's/stake/uslf/g' ~/.selfchain/config/genesis.json
#
# # Add wallet (skip existing wallet)
# selfchainold keys add wallet
#
# # Add genesis account
# selfchainold add-genesis-account <address_for> 10000000000uslf
#
# # Generate validator transaction
# selfchainold gentx wallet 5000000000uslf --chain-id=selfchain-1
#
# # Collect genesis transactions
# selfchainold collect-gentxs
#
# # Validate genesis
# selfchainold validate-genesis
#
# # Optional: Set minimum gas prices
# # sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0025uslf"/g' ~/.selfchain/config/app.toml
#
# # Start the chain
# # selfchainold start
# =====================================================

echo "=========== RESETTING BLOCKCHAIN ==========="
if [ "$NODE_ALREADY_RUNNING" = false ]; then
  selfchaind tendermint unsafe-reset-all

  echo "=========== STARTING NODE ==========="
  # Start the node in the background
  selfchainold start &
  NODE_PID=$!

  # Wait for the node to start properly (adjust sleep time as needed)
  log_step "Waiting for node to start..."
  sleep 20  # Give the node time to start up

  # Function to check if node is running
  check_node() {
    curl -s http://0.0.0.0:26657/status > /dev/null
    return $?
  }

  # Keep checking until node is accessible
  while ! check_node; do
    log_step "Node not ready yet, waiting..."
    sleep 5
  done

  log_step "Node started successfully!"
else
  echo "Skipping node startup as --node-running flag was provided"
  echo "Using existing running node for operations"
  
  # Quick check if node is actually running
  if ! curl -s http://0.0.0.0:26657/status > /dev/null; then
    echo "ERROR: Node does not appear to be running at http://0.0.0.0:26657"
    echo "Please start the node or remove the --node-running flag"
    exit 1
  fi
fi

# Create daily.json if it doesn't exist
if [ ! -f ./daily.json ]; then
  echo "Creating daily.json..."
  cat > ./daily.json << EOF
{
    "start_time": 1744162927,
    "periods": [
        {
            "coins": "5000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        },
        {
            "coins": "1000000000uslf",
            "length_seconds": 86400
        }
    ]
}
EOF
fi

echo "=========== CREATING VESTING ACCOUNTS ==========="

# Add logging for better visibility
log_step() {
  echo "$(date '+%H:%M:%S') - $1"
}

# Common parameters for transactions
NODE_PARAM="--node http://0.0.0.0:26657"
CHAIN_PARAM="--chain-id selfchain-1"
BROADCAST_PARAM="--broadcast-mode sync"
FROM_ALICE="--from alice"
YES_FLAG="--yes" # Add yes flag to auto-confirm transactions

# Create periodic vesting accounts
log_step "Creating periodic vesting accounts..."
PERIODIC_VESTING_ADDRESSES=(
  "self1hnlvp7z7s24m86utfenwzkwce8nv56hk434fyl"
  "self15qpk886wcrmvzwxxkpz0zw2avmm2yc76uay2jz"
  "self1pf5ffkmvda7d93jyfvxcvep32t2863gcstsu3w"
  "self1yqtry709yamnqsaj0heav7pxz72958a6ll0qc9"
)

for addr in "${PERIODIC_VESTING_ADDRESSES[@]}"; do
  log_step "Creating periodic vesting account for $addr"
  selfchainold tx vesting create-periodic-vesting-account \
    $addr \
    ./daily.json \
    $NODE_PARAM \
    $CHAIN_PARAM \
    $BROADCAST_PARAM \
    $FROM_ALICE \
    $YES_FLAG
  
  # Wait between transactions to avoid overwhelming the node
  sleep 5
done

# Create permanent locked accounts
log_step "Creating permanent locked accounts..."
PERMANENT_LOCKED_ADDRESSES=(
  "self15y0sxd8zrzynuhsv76v6888u2ekef7a7xsj7gs"
  "self18fam8qwdxk70lz2jxvz5wfj6c3pa6dy9pwlj40"
  "self1maxcxghzl09glkqych3hhfcuvza56es7c4t4uy"
)

for addr in "${PERMANENT_LOCKED_ADDRESSES[@]}"; do
  log_step "Creating permanent locked account for $addr"
  selfchainold tx vesting create-permanent-locked-account \
    $addr \
    16000000000uslf \
    $NODE_PARAM \
    $CHAIN_PARAM \
    $BROADCAST_PARAM \
    $FROM_ALICE \
    $YES_FLAG
  
  # Wait between transactions
  sleep 5
done

echo "=========== DELEGATING TOKENS ==========="
# Delegate tokens to validator
VALIDATOR="selfvaloper1aqkakdf0vwlsf24ugh0g4e0jcpwqvsautyltwz"
GAS_PARAMS="--gas auto --gas-adjustment 1.4"

# Define delegations: address, amount, from_account
DELEGATIONS=(
  "self1hnlvp7z7s24m86utfenwzkwce8nv56hk434fyl 2000000000uslf test2"
  "self15y0sxd8zrzynuhsv76v6888u2ekef7a7xsj7gs 16000000000uslf test5"
  "self1yqtry709yamnqsaj0heav7pxz72958a6ll0qc9 17000000000uslf bell"
  "self1maxcxghzl09glkqych3hhfcuvza56es7c4t4uy 6000000000uslf test7"
  "self1pf5ffkmvda7d93jyfvxcvep32t2863gcstsu3w 12000000000uslf test1"
)

for delegation in "${DELEGATIONS[@]}"; do
  read -r addr amount from_account <<< "$delegation"
  log_step "Delegating $amount from $addr (account: $from_account) to validator $VALIDATOR"
  
  selfchainold tx staking delegate $VALIDATOR $amount \
    --from $addr \
    $CHAIN_PARAM \
    $GAS_PARAMS \
    --from $from_account \
    $YES_FLAG
    
  # Wait between delegations
  sleep 5
done

echo "=========== CREATING GOVERNANCE PROPOSAL ==========="

# Set fixed target height instead of calculating dynamically
TARGET_HEIGHT=60
log_step "Setting upgrade height to fixed value: $TARGET_HEIGHT"

# Create proposal2.json with fixed height
cat > ./proposal2.json << EOF
{
 "messages": [
  {
   "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
   "authority": "self10d07y265gmmuvt4z0w9aw880jnsr700jlfwec6",
   "plan": {
    "name": "v2",
    "height": "$TARGET_HEIGHT",
    "info": "",
    "upgraded_client_state": null
   }
  }
 ],
 "metadata": "ipfs://CID",
 "deposit": "10000000uslf",
 "title": "Add CosmWasm Support",
 "summary": "Adding CosmWasm module support",
 "voting_period": "200s"
}
EOF

log_step "Submitting governance proposal..."
selfchainold tx gov submit-proposal proposal2.json \
  $NODE_PARAM \
  $CHAIN_PARAM \
  $BROADCAST_PARAM \
  $FROM_ALICE \
  $YES_FLAG

log_step "Waiting for proposal to be processed..."
sleep 10

log_step "Voting 'yes' on proposal..."
selfchainold tx gov vote 1 yes $FROM_ALICE $CHAIN_PARAM $YES_FLAG

echo "=========== SETUP COMPLETED ==========="
if [ "$NODE_ALREADY_RUNNING" = false ]; then
  echo "Node is running in the background with PID: $NODE_PID"
  echo "To stop the node, run: kill $NODE_PID"
else
  echo "Script completed using existing node"
fi
echo "Upgrade scheduled at block height: $TARGET_HEIGHT"

# Note: The script doesn't kill the node process if it started one
# You can manually stop it when needed with: kill $NODE_PID

#POST this manually need to verify account and balances for the list of address. below is list for reference

#	"self1yqtry709yamnqsaj0heav7pxz72958a6ll0qc9": "self1pw76gr9ag9gv8a6jfg0d5c3ag5mm8pghq5sen2", //periodic vestion - full stake
#	"self1pf5ffkmvda7d93jyfvxcvep32t2863gcstsu3w": "self13yl6nrfs34hu2ayyt2k4l7z36wf5eh2zw4d537", //periodic vestion -   stake entire unvested and more
#	"self1hnlvp7z7s24m86utfenwzkwce8nv56hk434fyl": "self144lz8lxa44w4h6x4g7w2qjz5d6n6qqzw8puzgl", //periodic vestion -  small stake
#	"self15qpk886wcrmvzwxxkpz0zw2avmm2yc76uay2jz": "self184tq49l234842jpjjqp4wuugck7u7ruv28jmzr", //periodic vestion -  no stake
#	"self18fam8qwdxk70lz2jxvz5wfj6c3pa6dy9pwlj40": "self12aqj8s0jkaetndt844z27r2wf2jxtwx7vl4fae", //permanent locked - test 3 no stake
#	"self15y0sxd8zrzynuhsv76v6888u2ekef7a7xsj7gs": "self10gkrql5vj9897k67hyp3nlxqkyzpaxvg3f722k", //permanent locked - full 16xx delegated
#	"self1maxcxghzl09glkqych3hhfcuvza56es7c4t4uy": "self1dmxphe3syl57e5qpwvr735mu5gppatvrwqz5np", //permanent locked - 6xx out of 16xx delegated .