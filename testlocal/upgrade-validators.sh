#!/bin/bash

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_header() {
    echo -e "${BLUE}[HEADER]${NC} $1"
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

# Set the new Docker image
OLD_IMAGE="selfchainprod-image:mainnet-cur"
NEW_IMAGE="selfchainprod-v4:mainnet-v4"

# =============================================================================
# PREPARATION
# =============================================================================
print_header "Upgrading validators to new Docker image..."
echo "Old image: $OLD_IMAGE"
echo "New image: $NEW_IMAGE"
echo

# Check if validator directories exist
for dir in validator1 validator2 validator3; do
    if [ ! -d "$dir" ]; then
        print_error "Directory $dir not found. Please ensure validators were previously set up."
        exit 1
    fi
done

# =============================================================================
# STOP CURRENT VALIDATORS
# =============================================================================
print_header "Stopping current validators..."

for validator in validator1 validator2 validator3; do
    if docker ps --format 'table {{.Names}}' | grep -q "^${validator}$"; then
        print_status "Stopping ${validator}..."
        docker stop ${validator}
    else
        print_warning "${validator} is not running"
    fi
done

# =============================================================================
# START VALIDATOR 1 WITH NEW IMAGE
# =============================================================================
print_header "Starting Validator 1 with new image..."

# Remove old container if it exists
docker rm validator1 2>/dev/null || true

# Start validator 1
print_status "Starting Validator 1..."
if ! docker run -d --name validator1 \
  -v $(pwd)/validator1:/home/heighliner/.selfchain \
  -p 36656:26656 \
  -p 36657:26657 \
  -p 36658:1234 \
  -p 36659:1317 \
  -p 36660:9090 \
  ${NEW_IMAGE} selfchaind start; then
    print_error "Failed to start validator1 container"
    exit 1
fi

# Wait for validator 1 to be ready
wait_for_service "http://localhost:36657/status" 60
wait_for_blocks "http://localhost:36657" 60

print_status "Validator 1 is ready!"

# =============================================================================
# START VALIDATOR 2 WITH NEW IMAGE  
# =============================================================================
print_header "Starting Validator 2 with new image..."

# Remove old container if it exists
docker rm validator2 2>/dev/null || true

# Start validator 2
print_status "Starting Validator 2..."
if ! docker run -d --name validator2 \
  -v $(pwd)/validator2:/home/heighliner/.selfchain \
  -p 37656:26656 \
  -p 37657:26657 \
  -p 37658:1234 \
  -p 37659:1317 \
  -p 37660:9090 \
  ${NEW_IMAGE} selfchaind start; then
    print_error "Failed to start validator2 container"
    exit 1
fi

# Wait for validator 2 to be ready
wait_for_service "http://localhost:37657/status" 60
print_status "Waiting for validator 2 to sync..."
sleep 10
wait_for_blocks "http://localhost:37657" 60

print_status "Validator 2 is ready!"

# =============================================================================
# START VALIDATOR 3 WITH NEW IMAGE
# =============================================================================
print_header "Starting Validator 3 with new image..."

# Remove old container if it exists
docker rm validator3 2>/dev/null || true

# Start validator 3
print_status "Starting Validator 3..."
if ! docker run -d --name validator3 \
  -v $(pwd)/validator3:/home/heighliner/.selfchain \
  -p 38656:26656 \
  -p 38657:26657 \
  -p 38658:1234 \
  -p 38659:1317 \
  -p 38660:9090 \
  ${NEW_IMAGE} selfchaind start; then
    print_error "Failed to start validator3 container"
    exit 1
fi

# Wait for validator 3 to be ready
wait_for_service "http://localhost:38657/status" 60
print_status "Waiting for validator 3 to sync..."
sleep 10
wait_for_blocks "http://localhost:38657" 60

print_status "Validator 3 is ready!"

# =============================================================================
# VERIFICATION
# =============================================================================
print_header "Verification"

# Check all validators are running
print_status "Checking validator status..."
for port in 36657 37657 38657; do
    height=$(curl -s "http://localhost:${port}/status" | jq -r '.result.sync_info.latest_block_height')
    echo "  Port ${port}: Block height ${height}"
done

# Check validator set
print_status "Checking validator set..."
docker exec validator1 selfchaind query staking validators --node tcp://localhost:26657 -o json | jq -r '.validators[] | .moniker'

# =============================================================================
# SUMMARY
# =============================================================================
print_header "Upgrade Complete!"
echo "=============================================="
echo "✅ All validators upgraded successfully"
echo "🔄 Old image: $OLD_IMAGE"
echo "🆕 New image: $NEW_IMAGE"
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
echo "=============================================="

# =============================================================================
# POST-UPGRADE CHECKS
# =============================================================================
print_header "Post-Upgrade Verification"

# Check if upgrade was successful
print_status "Checking upgrade status..."
upgrade_info=$(docker exec validator1 selfchaind query upgrade applied v4 --node tcp://localhost:26657 2>/dev/null || echo "not found")
if [ "$upgrade_info" != "not found" ]; then
    print_status "✅ Upgrade 'v4' has been applied successfully!"
    echo "$upgrade_info"
else
    print_warning "⚠️  Upgrade 'v4' status not confirmed. This might be expected if the upgrade is still in progress."
fi

# Final status check
print_status "Final validator status check..."
sleep 5

for validator in validator1 validator2 validator3; do
    if docker ps --format 'table {{.Names}}' | grep -q "^${validator}$"; then
        status=$(docker exec $validator selfchaind status --node tcp://localhost:26657 2>/dev/null | jq -r '.sync_info.catching_up' 2>/dev/null || echo "unknown")
        print_status "${validator}: catching_up = $status"
    else
        print_warning "${validator} is not running!"
    fi
done