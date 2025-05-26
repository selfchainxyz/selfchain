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

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    print_error "jq is required but not installed. Please install jq first."
    exit 1
fi

# Set variables
PASSPHRASE="qwaszxqw"
PROPOSAL_NUMBER=1
VALIDATOR1_PORT="36657"
VALIDATOR2_PORT="37657"
VALIDATOR3_PORT="38657"
VALIDATOR4_PORT="39657"
VALIDATOR5_PORT="40657"

# =============================================================================
# GET VALIDATOR CONTAINER IDs
# =============================================================================
print_header "Getting validator container information..."

# Get container IDs dynamically
VALIDATOR1_ID=$(docker ps --format 'table {{.ID}}\t{{.Names}}' | grep validator1 | awk '{print $1}')
VALIDATOR2_ID=$(docker ps --format 'table {{.ID}}\t{{.Names}}' | grep validator2 | awk '{print $1}')
VALIDATOR3_ID=$(docker ps --format 'table {{.ID}}\t{{.Names}}' | grep validator3 | awk '{print $1}')
VALIDATOR4_ID=$(docker ps --format 'table {{.ID}}\t{{.Names}}' | grep validator4 | awk '{print $1}')
VALIDATOR5_ID=$(docker ps --format 'table {{.ID}}\t{{.Names}}' | grep validator5 | awk '{print $1}')

if [ -z "$VALIDATOR1_ID" ] || [ -z "$VALIDATOR2_ID" ] || [ -z "$VALIDATOR3_ID" ] || [ -z "$VALIDATOR4_ID" ] || [ -z "$VALIDATOR5_ID" ]; then
    print_error "Could not find all validator containers. Please ensure all validators are running."
    print_status "Found containers:"
    echo "  Validator 1: ${VALIDATOR1_ID:-NOT FOUND}"
    echo "  Validator 2: ${VALIDATOR2_ID:-NOT FOUND}"
    echo "  Validator 3: ${VALIDATOR3_ID:-NOT FOUND}"
    echo "  Validator 4: ${VALIDATOR4_ID:-NOT FOUND}"
    echo "  Validator 5: ${VALIDATOR5_ID:-NOT FOUND}"
    exit 1
fi

print_status "Found all validator containers:"
echo "  Validator 1: $VALIDATOR1_ID"
echo "  Validator 2: $VALIDATOR2_ID"
echo "  Validator 3: $VALIDATOR3_ID"
echo "  Validator 4: $VALIDATOR4_ID"
echo "  Validator 5: $VALIDATOR5_ID"

# =============================================================================
# GET CURRENT HEIGHT AND SET UPGRADE HEIGHT
# =============================================================================
print_header "Calculating upgrade height..."

CURRENT_HEIGHT=$(curl -s http://localhost:${VALIDATOR1_PORT}/status | jq -r '.result.sync_info.latest_block_height')
print_status "Current height: $CURRENT_HEIGHT"

UPGRADE_HEIGHT=$((CURRENT_HEIGHT + 20))
print_status "Setting upgrade height to: $UPGRADE_HEIGHT"

# =============================================================================
# CREATE AND SUBMIT PROPOSAL
# =============================================================================
print_header "Creating governance proposal..."

# Create the proposal JSON
docker exec -i $VALIDATOR1_ID bash -c "mkdir -p /root/proposals && cat > /root/proposals/upgrade.json" << EOF
{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "authority": "self10d07y265gmmuvt4z0w9aw880jnsr700jlfwec6",
      "plan": {
        "name": "v4",
        "height": "$UPGRADE_HEIGHT",
        "info": "",
        "upgraded_client_state": null
      }
    }
  ],
  "metadata": "ipfs://CID",
  "deposit": "10000000uslf",
  "title": "Add CosmWasm Support",
  "summary": "Adding CosmWasm module support",
  "voting_period": "60s"
}
EOF

print_status "Proposal file created successfully"

# Submit the proposal
print_status "Submitting governance proposal..."
printf "$PASSPHRASE" \
| docker exec -i $VALIDATOR1_ID selfchaind tx gov submit-proposal /root/proposals/upgrade.json \
    --from wallet1 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

print_status "Proposal submitted successfully"

# Wait for proposal to be included in a block
print_status "Waiting for proposal to be included in a block..."
sleep 5

# =============================================================================
# VOTE ON PROPOSAL
# =============================================================================
print_header "Voting on proposal with all validators..."

# Vote with validator 1
print_status "Validator 1 voting YES..."
printf "$PASSPHRASE" \
| docker exec -i $VALIDATOR1_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet1 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

sleep 2

# Vote with validator 2
print_status "Validator 2 voting YES..."
printf "$PASSPHRASE" \
| docker exec -i $VALIDATOR2_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet2 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

sleep 2

# Vote with validator 3
print_status "Validator 3 voting YES..."
printf "$PASSPHRASE" \
| docker exec -i $VALIDATOR3_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet3 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

sleep 2

# Vote with validator 4
print_status "Validator 4 voting YES..."
printf "$PASSPHRASE" \
| docker exec -i $VALIDATOR4_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet4 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

sleep 2

# Vote with validator 5
print_status "Validator 5 voting YES..."
printf "$PASSPHRASE" \
| docker exec -i $VALIDATOR5_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet5 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

print_status "All 5 validators have voted on the proposal"

# Wait for votes to be processed
sleep 3

# =============================================================================
# QUERY PROPOSAL STATUS
# =============================================================================
print_header "Checking proposal status..."

print_status "Checking vote tally..."
docker exec -i $VALIDATOR1_ID selfchaind query gov tally $PROPOSAL_NUMBER --node tcp://localhost:26657

echo

print_status "Checking proposal details..."
docker exec -i $VALIDATOR1_ID selfchaind query gov proposal $PROPOSAL_NUMBER --node tcp://localhost:26657

echo

# =============================================================================
# WAIT FOR PROPOSAL TO PASS
# =============================================================================
print_header "Monitoring proposal status..."

# Function to check proposal status
check_proposal_status() {
    local status=$(docker exec -i $VALIDATOR1_ID selfchaind query gov proposal $PROPOSAL_NUMBER --node tcp://localhost:26657 -o json | jq -r '.status')
    echo "$status"
}

# Wait for proposal to reach final state
print_status "Waiting for proposal to reach final state..."
while true; do
    status=$(check_proposal_status)
    print_status "Current proposal status: $status"

    case $status in
        "PROPOSAL_STATUS_PASSED")
            print_status "✅ Proposal has PASSED!"
            break
            ;;
        "PROPOSAL_STATUS_REJECTED")
            print_error "❌ Proposal has been REJECTED!"
            break
            ;;
        "PROPOSAL_STATUS_FAILED")
            print_error "❌ Proposal has FAILED!"
            break
            ;;
        "PROPOSAL_STATUS_VOTING_PERIOD")
            print_status "⏳ Still in voting period..."
            ;;
        "PROPOSAL_STATUS_DEPOSIT_PERIOD")
            print_status "⏳ Still in deposit period..."
            ;;
        *)
            print_status "⏳ Unknown status: $status"
            ;;
    esac

    sleep 5
done

# =============================================================================
# SUMMARY
# =============================================================================
print_header "Summary"
echo "==========================================="
echo "✅ Governance proposal submitted and voted on"
echo "🔗 Proposal Number: $PROPOSAL_NUMBER"
echo "📊 Upgrade Height: $UPGRADE_HEIGHT"
echo "🏗️ Upgrade Name: v4"
echo "📋 Title: Add CosmWasm Support"
echo "👥 Validators voted: All 5 validators"
echo "==========================================="
echo
echo "To check the proposal status later:"
echo "docker exec -i $VALIDATOR1_ID selfchaind query gov proposal $PROPOSAL_NUMBER --node tcp://localhost:26657"
echo
echo "To check when the upgrade will happen:"
echo "docker exec -i $VALIDATOR1_ID selfchaind query upgrade plan --node tcp://localhost:26657"
echo "==========================================="