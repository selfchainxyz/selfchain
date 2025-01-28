#!/usr/bin/env bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

function update_imports() {
  local file="$1"
  echo -e "${BLUE}Processing file: $file${NC}"

  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS version
    sed -i '' "s|github.com/cosmos/cosmos-sdk/store|cosmossdk.io/store|g" "$file"
    sed -i '' "s|github.com/cosmos/cosmos-sdk/x/evidence|cosmossdk.io/x/evidence|g" "$file"
    sed -i '' "s|github.com/cosmos/cosmos-sdk/x/feegrant|cosmossdk.io/x/feegrant|g" "$file"

    # Comment out the line rewriting to x/upgrade/client, because that subpackage is gone:
    # sed -i '' "s|github.com/cosmos/cosmos-sdk/x/upgrade[^/]*/client|github.com/cosmos/cosmos-sdk/x/upgrade/client|g" "$file"

    # This line is correct for rewriting x/upgrade (but no /client):
    sed -i '' "s|github.com/cosmos/cosmos-sdk/x/upgrade|cosmossdk.io/x/upgrade|g" "$file"

    # IBC-Go references from v7 to v8:
    sed -i '' "s|github.com/cosmos/ibc-go/v7|github.com/cosmos/ibc-go/v8|g" "$file"

    # If you once had references to the old modules/capability path:
    # sed -i '' "s|github.com/cosmos/ibc-go/v8/modules/capability|github.com/cosmos/ibc-go/modules/capability|g" "$file"

    # Comet rename:
    sed -i '' "s|client/grpc/tmservice|client/grpc/cmtservice|g" "$file"

    # Snapshots
    sed -i '' "s|github.com/cosmos/cosmos-sdk/snapshots|cosmossdk.io/store/snapshots|g" "$file"
    sed -i '' "s|github.com/cosmos/cosmos-sdk/snapshots/types|cosmossdk.io/store/snapshots/types|g" "$file"

  else
    # Linux version
    sed -i "s|github.com/cosmos/cosmos-sdk/store|cosmossdk.io/store|g" "$file"
    sed -i "s|github.com/cosmos/cosmos-sdk/x/evidence|cosmossdk.io/x/evidence|g" "$file"
    sed -i "s|github.com/cosmos/cosmos-sdk/x/feegrant|cosmossdk.io/x/feegrant|g" "$file"

    # Comment out the line rewriting to x/upgrade/client, because that subpackage is gone:
    # sed -i "s|github.com/cosmos/cosmos-sdk/x/upgrade[^/]*/client|github.com/cosmos/cosmos-sdk/x/upgrade/client|g" "$file"

    sed -i "s|github.com/cosmos/cosmos-sdk/x/upgrade|cosmossdk.io/x/upgrade|g" "$file"

    sed -i "s|github.com/cosmos/ibc-go/v7|github.com/cosmos/ibc-go/v8|g" "$file"

    # If you had references to modules/capability in v7 or v8:
    # sed -i "s|github.com/cosmos/ibc-go/v8/modules/capability|github.com/cosmos/ibc-go/modules/capability|g" "$file"

    # Comet rename:
    sed -i "s|client/grpc/tmservice|client/grpc/cmtservice|g" "$file"

    sed -i "s|github.com/cosmos/cosmos-sdk/snapshots|cosmossdk.io/store/snapshots|g" "$file"
    sed -i "s|github.com/cosmos/cosmos-sdk/snapshots/types|cosmossdk.io/store/snapshots/types|g" "$file"
  fi
}

function main() {
  echo -e "${BLUE}Starting import updates...${NC}"

  # Clean up previous state
  rm -f go.sum
  go clean -modcache

  # Process all Go files
  find . -name "*.go" -type f | while read -r file; do
    update_imports "$file"
  done

  # Retry go mod tidy up to 3 times
  for i in {1..3}; do
    if go mod tidy; then
      echo -e "${GREEN}Successfully updated module dependencies!${NC}"
      exit 0
    fi
    echo -e "${BLUE}Retrying go mod tidy (attempt $i)...${NC}"
    sleep 2
  done

  echo -e "${RED}Failed to update module dependencies after 3 attempts.${NC}"
  exit 1
}

main
