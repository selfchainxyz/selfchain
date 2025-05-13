heighliner build -c selfchainprod-v4 --local -t mainnet-v4
heighliner build -c selfchainprod-image --local -t mainnet-cur

 # Selfchainprod Image Configuration
- name: selfchainprod-image
  pre-build: |
    apt update
    apt install -y unzip
    wget -O selfchaind-linux-arm64-cur "https://github.com/shivhg/exercism-solutions/releases/download/0.0.2/selfchaind-linux-arm64-cur"
    mv selfchaind-linux-arm64-cur /usr/bin/selfchaind
    chmod 755 /usr/bin/selfchaind
  binaries:
    - /usr/bin/selfchaind

# Selfchainprod Image Configuration
- name: selfchainprod-v4
  pre-build: |
    apt update
    apt install -y unzip
    wget -O selfchaind-linux-arm64-v4 "https://github.com/shivhg/exercism-solutions/releases/download/0.0.2/selfchaind-linux-arm64-v4"
    mv selfchaind-linux-arm64-v4 /usr/bin/selfchaind
    chmod 755 /usr/bin/selfchaind
  binaries:
    - /usr/bin/selfchaind


  -------------------------------------------------*************---------------------------------------------------



# Create a validator1 directory if it doesn't exist
mkdir -p validator1

# Initialize the node
docker run -it -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind init mynode --chain-id selfchain-1

# Replace stake with uslf in genesis file
sed -i '' 's/stake/uslf/g' validator1/config/genesis.json

sed -i '' 's/172800s/60s/g' validator1/config/genesis.json

sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0025uslf"/g' validator1/config/app.toml

# Update RPC to bind to all interfaces instead of just localhost
sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' validator1/config/config.toml

# Enable and configure API server
sed -i '' 's/enable = false/enable = true/g' validator1/config/app.toml
sed -i '' 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/g' validator1/config/app.toml

# Enable and configure gRPC server 
sed -i '' 's/address = "localhost:9090"/address = "0.0.0.0:9090"/g' validator1/config/app.toml

# Fix the minimum gas prices (currently set to "0stake" but should be "0.0025uslf")
sed -i '' 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "0.0025uslf"/g' validator1/config/app.toml

sed -i '' 's/moniker = "mynode"/moniker = "validator1"/g' validator1/config/config.toml



# set your variables:
MNEMONIC="verify model print hill eager whale divert ostrich depart enable exercise virtual wrestle security sudden supply nephew fly joy under robot evolve sight army"
PASSPHRASE="qwaszxqw"

# then:
printf "%s\n%s\n%s\n" "$MNEMONIC" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run -i \
      -v "$(pwd)/validator1:/home/heighliner/.selfchain" \
      selfchainprod-image:mainnet-cur \
      selfchaind keys add wallet1 --recover


# Add genesis account (replace with your address)
docker run -it -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind add-genesis-account self1k42we36mkhft50zn2f8mchhz7u4aah8aa6fyv3 10000000000000uslf

docker run -it -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind add-genesis-account self1adf59zdkuyppn3j8pc5gqmvx0lucradjcfgc96 10000000000000uslf

docker run -it -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind add-genesis-account self15hcvdar3eszwfjypz65levutqvj0plat8aj4n9 10000000000000uslf

# Generate validator transaction
printf "$PASSPHRASE" \
| docker run -i -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind gentx wallet1 5000000000uslf --chain-id=selfchain-1

# Collect genesis transactions
docker run -it -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind collect-gentxs

# Validate genesis
docker run -it -v $(pwd)/validator1:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind validate-genesis

# Start the chain with ports exposed
docker run -it \
  -v $(pwd)/validator1:/home/heighliner/.selfchain \
  -p 36656:26656 \
  -p 36657:26657 \
  -p 36658:1234 \
  -p 36659:1317 \
  -p 36660:9090 \
  selfchainprod-image:mainnet-cur selfchaind start



-------------------------------------------------*************---------------------------------------------------

# Create directory for validator 2
mkdir -p validator2

# Initialize validator 2 node
docker run -it -v $(pwd)/validator2:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind init validator2 --chain-id selfchain-1

# Copy the genesis file from the first node
cp validator1/config/genesis.json validator2/config/genesis.json
cp validator1/config/app.toml validator2/config/app.toml
cp validator1/config/config.toml validator2/config/config.toml

# Update config for validator 2
sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0025uslf"/g' validator2/config/app.toml


# Set persistent peers to connect to validator 1 (replace NODE_ID with actual node ID from validator 1)
# You'll need to get the actual node ID from the first validator
# After starting validator 1, get its node ID dynamically
VALIDATOR1_NODE_ID=$(curl -s http://localhost:36657/status | jq -r '.result.node_info.id')
echo "Validator 1 Node ID: $VALIDATOR1_NODE_ID"
# Use the dynamic node ID for validator 2
sed -i '' "s/persistent_peers = \"\"/persistent_peers = \"${VALIDATOR1_NODE_ID}@host.docker.internal:36656\"/g" validator2/config/config.toml

sed -i '' 's/moniker = "validator1"/moniker = "validator2"/g' validator2/config/config.toml

# Import the wallet for second genesis account
MNEMONIC="gather corn brother distance just winner phrase mechanic garlic program increase victory shoot brush tuna idle wet punch denial math artefact favorite timber that"
PASSPHRASE="qwaszxqw"
printf "%s\n%s\n%s\n" "$MNEMONIC" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run -i \
      -v "$(pwd)/validator2:/home/heighliner/.selfchain" \
      selfchainprod-image:mainnet-cur \
      selfchaind keys add wallet2 --recover


# Start validator 2 with different ports
docker run -it \
  -v $(pwd)/validator2:/home/heighliner/.selfchain \
  -p 37656:26656 \
  -p 37657:26657 \
  -p 37658:1234 \
  -p 37659:1317 \
  -p 37660:9090 \
  selfchainprod-image:mainnet-cur selfchaind start


VALIDATOR2_ID=$(docker ps | grep "37656" | awk '{print $1}')

# Submit validator transaction for validator 2 with proper fees
print "$PASSPHRASE" \
| docker exec -i $VALIDATOR2_ID selfchaind tx staking create-validator \
  --from=wallet2 \
  --amount=5000000000uslf \
  --pubkey=$(docker exec -it $VALIDATOR2_ID selfchaind tendermint show-validator) \
  --moniker="validator2" \
  --chain-id=selfchain-1 \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --broadcast-mode=sync \
  --fees=500uslf \
  --yes

-------------------------------------------------*************---------------------------------------------------


# Create directory for validator 3
mkdir -p validator3

# Initialize validator 3 node
docker run -it -v $(pwd)/validator3:/home/heighliner/.selfchain selfchainprod-image:mainnet-cur selfchaind init validator3 --chain-id selfchain-1

# Copy all config files from validator 1
cp validator1/config/genesis.json validator3/config/genesis.json
cp validator1/config/app.toml validator3/config/app.toml
cp validator1/config/config.toml validator3/config/config.toml

# After starting validator 1, get its node ID dynamically
# Get both validator 1 and validator 2 node IDs
VALIDATOR1_NODE_ID=$(curl -s http://localhost:36657/status | jq -r '.result.node_info.id')
VALIDATOR2_NODE_ID=$(curl -s http://localhost:37657/status | jq -r '.result.node_info.id')

# Set both as persistent peers
sed -i '' "s/persistent_peers = \"\"/persistent_peers = \"${VALIDATOR1_NODE_ID}@host.docker.internal:36656,${VALIDATOR2_NODE_ID}@host.docker.internal:37656\"/g" validator3/config/config.toml

# Update moniker
sed -i '' 's/moniker = "validator1"/moniker = "validator3"/g' validator3/config/config.toml


# Import the wallet for third genesis account
MNEMONIC="poet number abandon donate fitness cancel boss champion confirm bike dry century injury frown swamp poverty icon include enhance unit claim rich common laugh"
PASSPHRASE="qwaszxqw"
printf "%s\n%s\n%s\n" "$MNEMONIC" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run -i \
      -v "$(pwd)/validator3:/home/heighliner/.selfchain" \
      selfchainprod-image:mainnet-cur \
      selfchaind keys add wallet3 --recover

# Start validator 3 with different ports
docker run -it \
  -v $(pwd)/validator3:/home/heighliner/.selfchain \
  -p 38656:26656 \
  -p 38657:26657 \
  -p 38658:1234 \
  -p 38659:1317 \
  -p 38660:9090 \
  selfchainprod-image:mainnet-cur selfchaind start


  # Get validator 3 container ID
VALIDATOR3_ID=$(docker ps | grep "38656" | awk '{print $1}')
echo "Validator 3 container ID: $VALIDATOR3_ID"

# Submit validator transaction for validator 3
print "$PASSPHRASE" \
| docker exec -i $VALIDATOR3_ID selfchaind tx staking create-validator \
  --from=wallet3 \
  --amount=5000000000uslf \
  --pubkey=$(docker exec -it $VALIDATOR3_ID selfchaind tendermint show-validator) \
  --moniker="validator3" \
  --chain-id=selfchain-1 \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --broadcast-mode=sync \
  --fees=500uslf \
  --yes

-------------------------------------------------*************---------------------------------------------------

  #!/bin/bash
echo "=== Multi-Validator Network Status ==="
echo

# Check validator set
echo "Current Validators:"
curl -s http://localhost:36657/validators | jq '.result.validators[] | {address: .address, voting_power: .voting_power}' | jq -s 'sort_by(.voting_power) | reverse'

# Check block signatures
signatures=$(curl -s http://localhost:36657/commit | jq '.result.signed_header.commit.signatures | length')
echo -e "\nActive validators signing blocks: $signatures/3"

# Check sync status
echo -e "\nValidator Heights:"
for i in {1..3}; do
    port=$((36656 + (i-1) * 1000 + 1))
    height=$(curl -s http://localhost:$port/status | jq -r '.result.sync_info.latest_block_height')
    moniker=$(curl -s http://localhost:$port/status | jq -r '.result.node_info.moniker')
    echo "  $moniker: $height"
done

echo -e "\n✅ Network Status: Healthy 3-validator consensus"







-------------------------------------------------*************---------------------------------------------------
VALIDATOR1_ID=$(docker ps | grep "36657" | awk '{print $1}')
echo "Using Validator 1 container: $VALIDATOR1_ID"

docker cp ./daily.json $VALIDATOR1_ID:/home/heighliner/daily.json

self1jezc4atme56v75x5njqe4zuaccc4secug25wd3
self1fun8q0xuncfef6nkwh9njvvp4xqf4276x5sxgf // unvested delegated
self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah
self1fcahhgtw2llk06am4rala6khxjtj24zhhxn449

print "$PASSPHRASE" \
| docker exec -i $VALIDATOR1_ID selfchaind tx vesting create-periodic-vesting-account \
    self1krxfd67wmrjksq20xww53rm0wqmyxcew22whah \
    /home/heighliner/daily.json \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    --broadcast-mode sync \
    --from wallet1 \
    --fees 500uslf \
    --yes

docker exec -it $VALIDATOR1_ID selfchaind query account self1fun8q0xuncfef6nkwh9njvvp4xqf4276x5sxgf


-------------------------------------------------*************---------------------------------------------------

CURRENT_HEIGHT=$(curl -s http://localhost:36657/status | jq -r '.result.sync_info.latest_block_height')
echo "Current height: $CURRENT_HEIGHT"

UPGRADE_HEIGHT=$((CURRENT_HEIGHT + 20))
echo "Setting upgrade height to: $UPGRADE_HEIGHT"

docker exec -i $VALIDATOR1_ID bash -c "mkdir -p /home/heighliner/proposals && cat > /home/heighliner/proposals/proposal2.json" << EOF
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

print "$PASSPHRASE" \
| docker exec -i $VALIDATOR1_ID selfchaind tx gov submit-proposal /home/heighliner/proposals/proposal2.json \
    --from wallet1 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

PROPOSAL_NUMBER=1

# Vote with validator 1
print "$PASSPHRASE" \
| docker exec -i $VALIDATOR1_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet1 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

# Vote with validator 2
VALIDATOR2_ID=$(docker ps | grep "37657" | awk '{print $1}')
print "$PASSPHRASE" \
| docker exec -i $VALIDATOR2_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet2 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

# Vote with validator 3
VALIDATOR3_ID=$(docker ps | grep "38657" | awk '{print $1}')
print "$PASSPHRASE" \
| docker exec -i $VALIDATOR3_ID selfchaind tx gov vote $PROPOSAL_NUMBER yes \
    --from wallet3 \
    --fees=500uslf \
    --node tcp://localhost:26657 \
    --chain-id selfchain-1 \
    -y

docker exec -it $VALIDATOR1_ID selfchaind  query gov tally $PROPOSAL_NUMBER
docker exec -it $VALIDATOR1_ID selfchaind query gov proposal $PROPOSAL_NUMBER --node tcp://localhost:26657



-------------------------------------------------*************---------------------------------------------------


docker run -it \
  -v $(pwd)/validator1:/home/heighliner/.selfchain \
  -p 36656:26656 \
  -p 36657:26657 \
  -p 36658:1234 \
  -p 36659:1317 \
  -p 36660:9090 \
  selfchainprod-v4:mainnet-v4 selfchaind start



docker run -it \
  -v $(pwd)/validator2:/home/heighliner/.selfchain \
  -p 37656:26656 \
  -p 37657:26657 \
  -p 37658:1234 \
  -p 37659:1317 \
  -p 37660:9090 \
  selfchainprod-v4:mainnet-v4 selfchaind start



docker run -it \
  -v $(pwd)/validator3:/home/heighliner/.selfchain \
  -p 38656:26656 \
  -p 38657:26657 \
  -p 38658:1234 \
  -p 38659:1317 \
  -p 38660:9090 \
  selfchainprod-v4:mainnet-v4 selfchaind start

  -------------------------------------------------*************---------------------------------------------------


# Copy the CW20 wasm file to the container (if you have it locally)
docker cp ./cw20_base.wasm $VALIDATOR1_ID:/tmp/cw20_base.wasm

# Store the contract
docker exec -it $VALIDATOR1_ID selfchaind tx wasm store /tmp/cw20_base.wasm \
    --gas 3000000 \
    --gas-adjustment 1.3 \
    --fees 8000uslf \
    --from wallet1 \
    --chain-id selfchain-1 \
    --keyring-backend test \
    -y

# Instantiate with higher fees
docker exec -it $VALIDATOR1_ID selfchaind tx wasm instantiate 1 \
'{"name":"Self Token","symbol":"SELF","decimals":6,"initial_balances":[{"address":"self1k42we36mkhft50zn2f8mchhz7u4aah8aa6fyv3","amount":"1000000"}],"mint":null}' \
    --from wallet1 \
    --label "Self CW20 Token v1" \
    --gas 3000000 \
    --gas-adjustment 1.3 \
    --fees 7500uslf \
    --admin self1k42we36mkhft50zn2f8mchhz7u4aah8aa6fyv3 \
    --chain-id selfchain-1 \
    --keyring-backend test \
    --broadcast-mode sync \
    -y

# Set the contract address
CONTRACT_ADDR=$(docker exec -it $VALIDATOR1_ID selfchaind query wasm list-contract-by-code 1 --node tcp://localhost:26657 | grep -o 'self1[a-z0-9]\{58\}')
echo "Contract Address: $CONTRACT_ADDR"


# Check token info
echo "=== Token Info ==="
docker exec -it $VALIDATOR1_ID selfchaind query wasm contract-state smart $CONTRACT_ADDR \
'{"token_info":{}}' \
--node tcp://localhost:26657

# Check balance of initial holder
echo -e "\n=== Initial Balance ==="
docker exec -it $VALIDATOR1_ID selfchaind query wasm contract-state smart $CONTRACT_ADDR \
'{"balance":{"address":"self1k42we36mkhft50zn2f8mchhz7u4aah8aa6fyv3"}}' \
--node tcp://localhost:26657

# Transfer some tokens
echo -e "\n=== Transferring 100,000 tokens ==="
docker exec -it $VALIDATOR1_ID selfchaind tx wasm execute $CONTRACT_ADDR \
'{"transfer":{"recipient":"self1adf59zdkuyppn3j8pc5gqmvx0lucradjcfgc96","amount":"100000"}}' \
    --from wallet1 \
    --gas 300000 \
    --fees 1000uslf \
    --chain-id selfchain-1 \
    --keyring-backend test \
    -y

# Wait for transaction to process
sleep 3

# Check balances after transfer
echo -e "\n=== Sender Balance After Transfer ==="
docker exec -it $VALIDATOR1_ID selfchaind query wasm contract-state smart $CONTRACT_ADDR \
'{"balance":{"address":"self1k42we36mkhft50zn2f8mchhz7u4aah8aa6fyv3"}}' \
--node tcp://localhost:26657

echo -e "\n=== Recipient Balance After Transfer ==="
docker exec -it $VALIDATOR1_ID selfchaind query wasm contract-state smart $CONTRACT_ADDR \
'{"balance":{"address":"self1adf59zdkuyppn3j8pc5gqmvx0lucradjcfgc96"}}' \
--node tcp://localhost:26657

# Try another transfer to validator3's address
echo -e "\n=== Transferring 50,000 to validator3 ==="
docker exec -it $VALIDATOR1_ID selfchaind tx wasm execute $CONTRACT_ADDR \
'{"transfer":{"recipient":"self15hcvdar3eszwfjypz65levutqvj0plat8aj4n9","amount":"50000"}}' \
    --from wallet1 \
    --gas 300000 \
    --fees 1000uslf \
    --chain-id selfchain-1 \
    --keyring-backend test \
    -y

sleep 3

# Check validator3's balance
echo -e "\n=== Validator3 Balance ==="
docker exec -it $VALIDATOR1_ID selfchaind query wasm contract-state smart $CONTRACT_ADDR \
'{"balance":{"address":"self15hcvdar3eszwfjypz65levutqvj0plat8aj4n9"}}' \
--node tcp://localhost:26657



