#!/bin/bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to cleanup existing containers and directories
cleanup_existing() {
    print_warning "Cleaning up existing containers and data..."
    
    # Stop and remove existing containers if they exist
    for validator in validator1 validator2 validator3; do
        if docker ps -a --format 'table {{.Names}}' | grep -q "^${validator}$"; then
            print_status "Stopping and removing existing ${validator} container..."
            docker stop ${validator} 2>/dev/null || true
            docker rm ${validator} 2>/dev/null || true
        fi
    done
    
    # Remove existing directories
    for dir in validator1 validator2 validator3; do
        if [ -d "$dir" ]; then
            print_status "Removing existing ${dir} directory..."
            rm -rf "$dir"
        fi
    done
    
    print_status "Cleanup completed!"
}

# Function to wait for a service to be ready
wait_for_service() {
    local url=$1
    local timeout=${2:-30}
    local interval=2
    local count=0
    
    print_status "Waiting for service at $url to be ready..."
    
    while ! curl -s "$url" > /dev/null 2>&1; do
        if [ $count -ge $timeout ]; then
            print_error "Timeout waiting for service at $url"
            return 1
        fi
        sleep $interval
        count=$((count + interval))
    done
    
    print_status "Service at $url is ready!"
}

# Function to wait for block production
wait_for_blocks() {
    local rpc_url=$1
    local timeout=${2:-60}
    local interval=2
    local count=0
    
    print_status "Waiting for blocks to be produced..."
    
    while true; do
        if [ $count -ge $timeout ]; then
            print_error "Timeout waiting for blocks to be produced"
            return 1
        fi
        
        # Check if we can get block info and height > 0
        local height=$(curl -s "${rpc_url}/status" | jq -r '.result.sync_info.latest_block_height' 2>/dev/null)
        if [[ "$height" =~ ^[0-9]+$ ]] && [ "$height" -gt 0 ]; then
            print_status "Blocks are being produced! Current height: $height"
            return 0
        fi
        
        sleep $interval
        count=$((count + interval))
    done
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    print_error "jq is required but not installed. Please install jq first."
    exit 1
fi

# Handle command line arguments
if [ "$1" = "--clean" ] || [ "$1" = "-c" ]; then
    cleanup_existing
    print_status "Starting fresh setup..."
elif [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Usage: $0 [OPTION]"
    echo "Options:"
    echo "  --clean, -c    Clean up existing containers and directories before setup"
    echo "  --help, -h     Show this help message"
    exit 0
else
    # Check for existing containers and prompt user
    existing_containers=()
    for validator in validator1 validator2 validator3; do
        if docker ps -a --format 'table {{.Names}}' | grep -q "^${validator}$"; then
            existing_containers+=("$validator")
        fi
    done
    
    if [ ${#existing_containers[@]} -gt 0 ]; then
        print_warning "Found existing containers: ${existing_containers[*]}"
        print_warning "This script will fail if containers already exist."
        echo -n "Do you want to clean up existing containers? (y/N): "
        read -r response
        case "$response" in
            [yY][eE][sS]|[yY])
                cleanup_existing
                ;;
            *)
                print_error "Please clean up existing containers first or run with --clean option"
                exit 1
                ;;
        esac
    fi
fi

# Set variables
MNEMONIC1="verify model print hill eager whale divert ostrich depart enable exercise virtual wrestle security sudden supply nephew fly joy under robot evolve sight army"
MNEMONIC2="gather corn brother distance just winner phrase mechanic garlic program increase victory shoot brush tuna idle wet punch denial math artefact favorite timber that"
MNEMONIC3="poet number abandon donate fitness cancel boss champion confirm bike dry century injury frown swamp poverty icon include enhance unit claim rich common laugh"
PASSPHRASE="qwaszxqw"

# =============================================================================
# VALIDATOR 1 SETUP
# =============================================================================
print_status "Setting up Validator 1..."

# Create validator1 directory
mkdir -p validator1

# Initialize the node
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind init mynode --chain-id selfchain-1

# Replace stake with uslf in genesis file
sed -i.bak 's/stake/uslf/g' validator1/config/genesis.json

# Update unbonding time to 60 seconds for testing
sed -i.bak 's/172800s/60s/g' validator1/config/genesis.json

# Set minimum gas prices
sed -i.bak 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0025uslf"/g' validator1/config/app.toml

# Update RPC to bind to all interfaces
sed -i.bak 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' validator1/config/config.toml

# Enable and configure API server
sed -i.bak 's/enable = false/enable = true/g' validator1/config/app.toml
sed -i.bak 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/g' validator1/config/app.toml

# Enable and configure gRPC server
sed -i.bak 's/address = "localhost:9090"/address = "0.0.0.0:9090"/g' validator1/config/app.toml

# Fix the minimum gas prices if they're still set to stake
sed -i.bak 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "0.0025uslf"/g' validator1/config/app.toml

# Update moniker
sed -i.bak 's/moniker = "mynode"/moniker = "validator1"/g' validator1/config/config.toml

# Import wallet1
printf "%s\n%s\n%s\n" "$MNEMONIC1" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run --rm -i \
      -v "$(pwd)/validator1:/home/heighliner/.selfchain" \
      selfchainprod-image:mainnet-cur \
      selfchaind keys add wallet1 --recover

# Add genesis accounts
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind add-genesis-account self1k42we36mkhft50zn2f8mchhz7u4aah8aa6fyv3 10000000000000uslf
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind add-genesis-account self1adf59zdkuyppn3j8pc5gqmvx0lucradjcfgc96 10000000000000uslf
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind add-genesis-account self15hcvdar3eszwfjypz65levutqvj0plat8aj4n9 10000000000000uslf

# Generate validator transaction
printf "$PASSPHRASE" \
| docker run --rm -i -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind gentx wallet1 5000000000uslf --chain-id=selfchain-1

# Collect genesis transactions
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind collect-gentxs

# Validate genesis
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind validate-genesis

# Start validator 1
print_status "Starting Validator 1..."
if ! docker run -d --name validator1 \
  -v $(pwd)/validator1:/home/heighliner/.selfchain \
  -p 36656:26656 \
  -p 36657:26657 \
  -p 36658:1234 \
  -p 36659:1317 \
  -p 36660:9090 \
  selfchainprod-image:mainnet-cur selfchaind start; then
    print_error "Failed to start validator1 container"
    exit 1
fi

# Wait for validator 1 to be ready
wait_for_service "http://localhost:36657/status" 60

# Wait for blocks to be produced
wait_for_blocks "http://localhost:36657" 60

print_status "Validator 1 is ready!"

# =============================================================================
# VALIDATOR 2 SETUP
# =============================================================================
print_status "Setting up Validator 2..."

# Create directory for validator 2
mkdir -p validator2

# Initialize validator 2 node
docker run --rm -v $(pwd)/validator2:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind init validator2 --chain-id selfchain-1

# Copy configuration files from validator 1
cp validator1/config/genesis.json validator2/config/genesis.json
cp validator1/config/app.toml validator2/config/app.toml
cp validator1/config/config.toml validator2/config/config.toml

# Get validator 1 node ID
VALIDATOR1_NODE_ID=$(curl -s http://localhost:36657/status | jq -r '.result.node_info.id')
print_status "Validator 1 Node ID: $VALIDATOR1_NODE_ID"

# Set persistent peers to connect to validator 1
sed -i.bak "s/persistent_peers = \"\"/persistent_peers = \"${VALIDATOR1_NODE_ID}@host.docker.internal:36656\"/g" validator2/config/config.toml

# Update moniker
sed -i.bak 's/moniker = "validator1"/moniker = "validator2"/g' validator2/config/config.toml

# Import wallet2
printf "%s\n%s\n%s\n" "$MNEMONIC2" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run --rm -i \
      -v "$(pwd)/validator2:/home/heighliner/.selfchain" \
      selfchainprod-image:mainnet-cur \
      selfchaind keys add wallet2 --recover

# Start validator 2
print_status "Starting Validator 2..."
if ! docker run -d --name validator2 \
  -v $(pwd)/validator2:/home/heighliner/.selfchain \
  -p 37656:26656 \
  -p 37657:26657 \
  -p 37658:1234 \
  -p 37659:1317 \
  -p 37660:9090 \
  selfchainprod-image:mainnet-cur selfchaind start; then
    print_error "Failed to start validator2 container"
    exit 1
fi

# Wait for validator 2 to be ready
wait_for_service "http://localhost:37657/status" 60

# Wait for validator 2 to sync with the chain
print_status "Waiting for validator 2 to sync..."
sleep 10

# Wait for blocks to ensure the chain is producing
wait_for_blocks "http://localhost:37657" 60

# Submit validator transaction for validator 2
print_status "Creating validator 2..."
printf "$PASSPHRASE" \
| docker exec -i validator2 selfchaind tx staking create-validator \
  --from=wallet2 \
  --amount=5000000000uslf \
  --pubkey=$(docker exec validator2 selfchaind tendermint show-validator) \
  --moniker="validator2" \
  --chain-id=selfchain-1 \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --broadcast-mode=sync \
  --fees=500uslf \
  --yes

print_status "Validator 2 is ready!"

# Wait for block production
sleep 10

# =============================================================================
# VALIDATOR 3 SETUP
# =============================================================================
print_status "Setting up Validator 3..."

# Create directory for validator 3
mkdir -p validator3

# Initialize validator 3 node
docker run --rm -v $(pwd)/validator3:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind init validator3 --chain-id selfchain-1

# Copy all config files from validator 1
cp validator1/config/genesis.json validator3/config/genesis.json
cp validator1/config/app.toml validator3/config/app.toml
cp validator1/config/config.toml validator3/config/config.toml

# Get validator 2 node ID
VALIDATOR2_NODE_ID=$(curl -s http://localhost:37657/status | jq -r '.result.node_info.id')
print_status "Validator 2 Node ID: $VALIDATOR2_NODE_ID"

# Set both validators as persistent peers
sed -i.bak "s/persistent_peers = \"\"/persistent_peers = \"${VALIDATOR1_NODE_ID}@host.docker.internal:36656,${VALIDATOR2_NODE_ID}@host.docker.internal:37656\"/g" validator3/config/config.toml

# Update moniker
sed -i.bak 's/moniker = "validator1"/moniker = "validator3"/g' validator3/config/config.toml

# Import wallet3
printf "%s\n%s\n%s\n" "$MNEMONIC3" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run --rm -i \
      -v "$(pwd)/validator3:/home/heighliner/.selfchain" \
      selfchainprod-image:mainnet-cur \
      selfchaind keys add wallet3 --recover

# Start validator 3
print_status "Starting Validator 3..."
if ! docker run -d --name validator3 \
  -v $(pwd)/validator3:/home/heighliner/.selfchain \
  -p 38656:26656 \
  -p 38657:26657 \
  -p 38658:1234 \
  -p 38659:1317 \
  -p 38660:9090 \
  selfchainprod-image:mainnet-cur selfchaind start; then
    print_error "Failed to start validator3 container"
    exit 1
fi

# Wait for validator 3 to be ready
wait_for_service "http://localhost:38657/status" 60

# Wait for validator 3 to sync with the chain
print_status "Waiting for validator 3 to sync..."
sleep 10

# Wait for blocks to ensure the chain is producing
wait_for_blocks "http://localhost:38657" 60

# Submit validator transaction for validator 3
print_status "Creating validator 3..."
printf "$PASSPHRASE" \
| docker exec -i validator3 selfchaind tx staking create-validator \
  --from=wallet3 \
  --amount=5000000000uslf \
  --pubkey=$(docker exec validator3 selfchaind tendermint show-validator) \
  --moniker="validator3" \
  --chain-id=selfchain-1 \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --broadcast-mode=sync \
  --fees=500uslf \
  --yes

print_status "Validator 3 is ready!"

# =============================================================================
# SUMMARY
# =============================================================================
echo
echo "=============================================="
echo "All validators are now running!"
echo "=============================================="
echo "Validator 1:"
echo "  RPC: http://localhost:36657"
echo "  API: http://localhost:36659"
echo "  gRPC: http://localhost:36660"
echo
echo "Validator 2:"
echo "  RPC: http://localhost:37657"
echo "  API: http://localhost:37659"
echo "  gRPC: http://localhost:37660"
echo
echo "Validator 3:"
echo "  RPC: http://localhost:38657"
echo "  API: http://localhost:38659"
echo "  gRPC: http://localhost:38660"
echo "=============================================="
echo
echo "To check validator status:"
echo "docker exec validator1 selfchaind status"
echo "docker exec validator2 selfchaind status"
echo "docker exec validator3 selfchaind status"
echo
echo "To stop all validators:"
echo "docker stop validator1 validator2 validator3"
echo
echo "To remove all validators:"
echo "docker rm validator1 validator2 validator3"
echo
echo "To clean up everything (containers + data):"
echo "./setup-validators.sh --clean"
echo "=============================================="