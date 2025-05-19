#!/bin/bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Define timezones for each validator
declare -A VALIDATOR_TIMEZONES=(
    ["validator1"]="America/New_York"      # UTC-5 (EST)
    ["validator2"]="Europe/London"         # UTC+0 (GMT)
    ["validator3"]="Asia/Tokyo"            # UTC+9 (JST)
    ["validator4"]="Australia/Sydney"      # UTC+11 (AEDT)
    ["validator5"]="America/Los_Angeles"   # UTC-8 (PST)
)

# Define mnemonics for each validator
declare -A VALIDATOR_MNEMONICS=(
    ["validator1"]="verify model print hill eager whale divert ostrich depart enable exercise virtual wrestle security sudden supply nephew fly joy under robot evolve sight army"
    ["validator2"]="gather corn brother distance just winner phrase mechanic garlic program increase victory shoot brush tuna idle wet punch denial math artefact favorite timber that"
    ["validator3"]="poet number abandon donate fitness cancel boss champion confirm bike dry century injury frown swamp poverty icon include enhance unit claim rich common laugh"
    ["validator4"]="pond little question injury green puzzle penalty trial hill diesel more impulse major oyster offer improve ensure lion sound broccoli dune evolve moon bounce"
    ["validator5"]="cannon shock power obey deny actress elephant craft glue direct power siren route cinnamon stay change pistol amazing door patient knife advance layer prosper"
)

# Define port ranges for each validator
declare -A VALIDATOR_PORTS=(
    ["validator1"]="36656:36657:36658:36659:36660"
    ["validator2"]="37656:37657:37658:37659:37660"
    ["validator3"]="38656:38657:38658:38659:38660"
    ["validator4"]="39656:39657:39658:39659:39660"
    ["validator5"]="40656:40657:40658:40659:40660"
)

# Define addresses for genesis accounts
GENESIS_ACCOUNTS=(
    "self1k42we36mkhft50zn2f8mchhz7u4aah8aa6fyv3"
    "self1adf59zdkuyppn3j8pc5gqmvx0lucradjcfgc96"
    "self15hcvdar3eszwfjypz65levutqvj0plat8aj4n9"
    "self16l2nfdknmjtzvvja2u64p64hg8qescss2njeyc"
    "self19x80kkwwwypj4e6tlpfucj6unkzd23qa374mvu"
)

PASSPHRASE="qwaszxqw"

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

# Function to print timezone info
print_timezone_info() {
    local validator_name=$1
    local timezone=${VALIDATOR_TIMEZONES[$validator_name]}
    local current_time=$(TZ=$timezone date "+%Y-%m-%d %H:%M:%S %Z")
    print_status "$validator_name timezone: $timezone ($current_time)"
}

# Function to cleanup existing containers and directories
cleanup_existing() {
    print_warning "Cleaning up existing containers and data..."

    # Stop and remove existing containers if they exist
    for validator in validator1 validator2 validator3 validator4 validator5; do
        if docker ps -a --format 'table {{.Names}}' | grep -q "^${validator}$"; then
            print_status "Stopping and removing existing ${validator} container..."
            docker stop ${validator} 2>/dev/null || true
            docker rm ${validator} 2>/dev/null || true
        fi
    done

    # Remove existing directories
    for dir in validator1 validator2 validator3 validator4 validator5; do
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

# Function to start a validator with specific timezone
start_validator_with_timezone() {
    local validator_name=$1
    local timezone=${VALIDATOR_TIMEZONES[$validator_name]}
    local ports=$2

    print_timezone_info $validator_name

    # Parse ports (format: "p1:p2:p3:p4:p5")
    IFS=':' read -ra ADDR <<< "$ports"
    local p2p_port=${ADDR[0]}
    local rpc_port=${ADDR[1]}
    local pprof_port=${ADDR[2]}
    local api_port=${ADDR[3]}
    local grpc_port=${ADDR[4]}

    print_status "Starting $validator_name with timezone $timezone..."

    if ! docker run -d --name $validator_name \
        -e TZ=$timezone \
        -v $(pwd)/$validator_name:/home/heighliner/.selfchain \
        -p $p2p_port:26656 \
        -p $rpc_port:26657 \
        -p $pprof_port:1234 \
        -p $api_port:1317 \
        -p $grpc_port:9090 \
        selfchainprod-image:mainnet-cur selfchaind start; then
        print_error "Failed to start $validator_name container"
        exit 1
    fi

    # Wait for validator to be ready
    wait_for_service "http://localhost:$rpc_port/status" 60
}

# Function to get RPC port from validator ports string
get_rpc_port() {
    local ports=$1
    IFS=':' read -ra ADDR <<< "$ports"
    echo ${ADDR[1]}
}

# Function to get P2P port from validator ports string
get_p2p_port() {
    local ports=$1
    IFS=':' read -ra ADDR <<< "$ports"
    echo ${ADDR[0]}
}

# Function to build persistent peers string
build_persistent_peers() {
    local current_validator=$1
    local peers=""

    # Get all validator numbers that come before current
    for i in {1..5}; do
        local validator="validator$i"
        if [ "$validator" = "$current_validator" ]; then
            break
        fi

        # Get node ID and P2P port for this validator
        local rpc_port=$(get_rpc_port "${VALIDATOR_PORTS[$validator]}")
        local p2p_port=$(get_p2p_port "${VALIDATOR_PORTS[$validator]}")
        local node_id=$(curl -s "http://localhost:$rpc_port/status" | jq -r '.result.node_info.id')

        if [ -n "$peers" ]; then
            peers="${peers},"
        fi
        peers="${peers}${node_id}@host.docker.internal:${p2p_port}"
    done

    echo "$peers"
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
    echo ""
    echo "Timezone Configuration:"
    for validator in validator1 validator2 validator3 validator4 validator5; do
        echo "  $validator: ${VALIDATOR_TIMEZONES[$validator]}"
    done
    exit 0
else
    # Check for existing containers and prompt user
    existing_containers=()
    for validator in validator1 validator2 validator3 validator4 validator5; do
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

# Display timezone information
echo "=============================================="
echo "Validator Timezone Configuration:"
echo "=============================================="
for validator in validator1 validator2 validator3 validator4 validator5; do
    print_timezone_info $validator
done
echo "=============================================="
echo

# =============================================================================
# VALIDATOR 1 SETUP (Special case - genesis validator)
# =============================================================================
print_status "Setting up Validator 1 in ${VALIDATOR_TIMEZONES[validator1]}..."

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
printf "%s\n%s\n%s\n" "${VALIDATOR_MNEMONICS[validator1]}" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run --rm -i \
      -v "$(pwd)/validator1:/home/heighliner/.selfchain" \
      selfchainprod-image:mainnet-cur \
      selfchaind keys add wallet1 --recover

# Add genesis accounts
for account in "${GENESIS_ACCOUNTS[@]}"; do
    docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind add-genesis-account "$account" 10000000000000uslf
done

# Generate validator transaction
printf "$PASSPHRASE" \
| docker run --rm -i -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind gentx wallet1 5000000000uslf --chain-id=selfchain-1

# Collect genesis transactions
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind collect-gentxs

# Validate genesis
docker run --rm -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind validate-genesis

# Start validator 1 with timezone
start_validator_with_timezone "validator1" "${VALIDATOR_PORTS[validator1]}"

# Wait for blocks to be produced
wait_for_blocks "http://localhost:$(get_rpc_port "${VALIDATOR_PORTS[validator1]}")" 60

print_status "Validator 1 is ready!"

# =============================================================================
# VALIDATORS 2-5 SETUP (Loop through remaining validators)
# =============================================================================
for i in {2..5}; do
    validator="validator$i"
    wallet="wallet$i"

    print_status "Setting up $validator in ${VALIDATOR_TIMEZONES[$validator]}..."

    # Create directory for validator
    mkdir -p "$validator"

    # Initialize validator node
    docker run --rm -v "$(pwd)/$validator:/home/heighliner/.selfchain" selfchainprod-image:mainnet-cur selfchaind init "$validator" --chain-id selfchain-1

    # Copy configuration files from validator 1
    cp validator1/config/genesis.json "$validator/config/genesis.json"
    cp validator1/config/app.toml "$validator/config/app.toml"
    cp validator1/config/config.toml "$validator/config/config.toml"

    # Build persistent peers for this validator
    persistent_peers=$(build_persistent_peers "$validator")
    print_status "$validator persistent peers: $persistent_peers"

    # Set persistent peers
    sed -i.bak "s/persistent_peers = \"\"/persistent_peers = \"$persistent_peers\"/g" "$validator/config/config.toml"

    # Update moniker
    sed -i.bak "s/moniker = \"validator1\"/moniker = \"$validator\"/g" "$validator/config/config.toml"

    # Import wallet
    printf "%s\n%s\n%s\n" "${VALIDATOR_MNEMONICS[$validator]}" "$PASSPHRASE" "$PASSPHRASE" \
      | docker run --rm -i \
          -v "$(pwd)/$validator:/home/heighliner/.selfchain" \
          selfchainprod-image:mainnet-cur \
          selfchaind keys add "$wallet" --recover

    # Start validator with timezone
    start_validator_with_timezone "$validator" "${VALIDATOR_PORTS[$validator]}"

    # Wait for validator to sync with the chain
    print_status "Waiting for $validator to sync..."
    sleep 10

    # Wait for blocks to ensure the chain is producing
    wait_for_blocks "http://localhost:$(get_rpc_port "${VALIDATOR_PORTS[$validator]}")" 60

    # Submit validator transaction
    print_status "Creating $validator..."
    printf "$PASSPHRASE" \
    | docker exec -i "$validator" selfchaind tx staking create-validator \
      --from="$wallet" \
      --amount=5000000000uslf \
      --pubkey="$(docker exec "$validator" selfchaind tendermint show-validator)" \
      --moniker="$validator" \
      --chain-id=selfchain-1 \
      --commission-rate="0.10" \
      --commission-max-rate="0.20" \
      --commission-max-change-rate="0.01" \
      --min-self-delegation="1" \
      --broadcast-mode=sync \
      --fees=500uslf \
      --yes

    print_status "$validator is ready!"

    # Wait before setting up next validator
    sleep 10
done

# =============================================================================
# SUMMARY WITH TIMEZONE INFO
# =============================================================================
echo
echo "=============================================="
echo "All validators are now running in different timezones!"
echo "=============================================="

for validator in validator1 validator2 validator3 validator4 validator5; do
    ports=${VALIDATOR_PORTS[$validator]}
    IFS=':' read -ra PORT_ARRAY <<< "$ports"
    rpc_port=${PORT_ARRAY[1]}
    api_port=${PORT_ARRAY[3]}
    grpc_port=${PORT_ARRAY[4]}
    timezone=${VALIDATOR_TIMEZONES[$validator]}

    echo "$validator ($timezone):"
    echo "  RPC: http://localhost:$rpc_port"
    echo "  API: http://localhost:$api_port"
    echo "  gRPC: http://localhost:$grpc_port"
    echo "  Current time: $(TZ=$timezone date)"
    echo
done

echo "=============================================="
echo
echo "To check validator time and status:"
for validator in validator1 validator2 validator3 validator4 validator5; do
    echo "docker exec $validator date && docker exec $validator selfchaind status"
done
echo
echo "To check timezone settings inside containers:"
for validator in validator1 validator2 validator3 validator4 validator5; do
    echo "docker exec $validator cat /etc/timezone"
done
echo
echo "To stop all validators:"
echo "docker stop validator1 validator2 validator3 validator4 validator5"
echo
echo "To remove all validators:"
echo "docker rm validator1 validator2 validator3 validator4 validator5"
echo
echo "To clean up everything (containers + data):"
echo "./setup-validators.sh --clean"
echo "=============================================="