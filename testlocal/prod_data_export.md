# Local Blockchain Node Setup with Production State - Internal README

## Overview
This document provides step-by-step instructions for setting up a local blockchain node using exported production state. This process creates a local testnet with a single validator for development and testing purposes.

## Prerequisites
- Access to production validator node
- Docker installed locally
- `jq` and `sponge` utilities installed
- Production validator with minimal stake

## Process Overview

### Step 1: Production Validator Setup
Ensure you have an active validator running in production with minimal stake for easier modifications later.

**Example addresses used in this guide:**
- Delegator Address: `self10hyl92j27zuwvu7rrvr6aa5ems2xjfd9q9zgun`
- Validator Address: `selfvaloper10hyl92j27zuwvu7rrvr6aa5ems2xjfd96mt76c`

*Replace these with your actual addresses throughout the process.*

### Step 2: Export Blockchain State
Stop the running production node (required due to database lock) and export the entire state:

```bash
./selfchaind-linux-amd64 export --for-zero-height --height 5745220 > exported_genesis.json
```

**Notes:**
- Replace `5745220` with your desired block height
- This command dumps all state to genesis for zero-height restart
- Node must be stopped to avoid database lock issues

### Step 3: Copy Required Files
Transfer the exported genesis and key files to your local machine:

```bash
# Copy exported genesis from production server
scp -r -i ~/Downloads/slfaws.pem ec2-user@ec2-3-148-216-19.us-east-2.compute.amazonaws.com:/exported_genesis.json .

# Copy node key and validator private key
cp ~/.selfchain/config/node_key.json .
cp ~/.selfchain/config/priv_validator_key.json .
```

**Required files:**
1. `exported_genesis.json`
2. `node_key.json` 
3. `priv_validator_key.json`

### Step 4: Backup and Format Genesis
Create backup and format JSON for editing:

```bash
cp exported_genesis.json exported_genesis_cp.json
jq '.' exported_genesis.json | sponge exported_genesis.json
```

### Step 5: Jail All Validators Except Test Validator
Edit the genesis file to jail all validators except your test validator:

**Find all validators and change:**
```json
"jailed": false
"status": "BOND_STATUS_BONDED"
```

**To:**
```json
"jailed": true
"status": "BOND_STATUS_UNBONDED"
```

**Exception:** Keep your test validator (`selfvaloper10hyl92j27zuwvu7rrvr6aa5ems2xjfd96mt76c`) as:
```json
"jailed": false
"status": "BOND_STATUS_BONDED"
```

### Step 6: Update Test Validator Holdings
Search for your staker/validator addresses and update holdings to 5 billion tokens in **four locations**:

**Two entries for delegator address (`self10hyl92j27zuwvu7rrvr6aa5ems2xjfd9q9zgun`)**
**Two entries for validator address (`selfvaloper10hyl92j27zuwvu7rrvr6aa5ems2xjfd96mt76c`)**

**Important:** This updates staking power and bond status, not account balance.

**Also update delegator shares:**
```json
"operator_address": "selfvaloper10hyl92j27zuwvu7rrvr6aa5ems2xjfd96mt76c",
"delegator_shares": "5000000000000000.000000000000000000"
```

### Step 7: Update Last Validator Powers
In `staking.last_validator_powers`, delete all entries except your test validator:

```json
"last_validator_powers": [
  {
    "address": "selfvaloper10hyl92j27zuwvu7rrvr6aa5ems2xjfd96mt76c",
    "power": "5000000000"
  }
]
```

### Step 8: Clean Up Root Validators Array
**Important:** There are two "validators" arrays in the genesis file. Only modify the **root-level** validators array (usually at the end of the file).

Delete all validators except your test validator **and update the power**:

```json
"validators": [
  {
    "address": "E6B876A13BC5728962191938DBC9D67EFE635195",
    "name": "GeckoGarage",
    "power": "5000000000",
    "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "f4YSPvMQZB6hXv4FI/tJIWfIOzfCslT6HYvFcvAi3L4="
    }
  }
]
```

**Critical:** Update the power field along with deleting other validators.

### Step 9: Update Staking Total Power
Set `staking.last_total_power` to match your validator's micro-tokens:

```json
"staking": {
  "last_total_power": "5000000000000000"
}
```

**Note:** This must equal your test validator's micro-tokens (5000000000000000).

### Step 10: Optional - Update Account Balance
Optionally increase your account balance for testing:

```json
{
  "address": "self10hyl92j27zuwvu7rrvr6aa5ems2xjfd9q9zgun",
  "coins": [
    {
      "amount": "20000000000000000",
      "denom": "uslf"
    }
  ]
}
```

**Purpose:** Provides sufficient balance for transfers and comprehensive testing.

### Step 11: Initialize and Start Local Chain

#### Create Validator Directory and Initialize:
```bash
mkdir -p validator1

docker run -it -v $(pwd)/validator1:/root/.selfchain selfchain:mainnet \
  selfchaind init mynode --chain-id selfchain-1

cp ./exported_genesis.json ./validator1/config/genesis.json
cp ./node_key.json ./validator1/config
cp ./priv_validator_key.json ./validator1/config
```

#### Set Up Wallet:
```bash
# Set your variables
MNEMONIC="selfchain is cosmos keyless wallet chain"
PASSPHRASE="selfchain@keylesswallet"

# Add wallet with recovery
printf "%s\n%s\n%s\n" "$MNEMONIC" "$PASSPHRASE" "$PASSPHRASE" \
  | docker run -i \
      -v "$(pwd)/validator1:/root/.selfchain" \
      selfchain:mainnet \
      selfchaind keys add wallet1 --recover
```

#### Start the Chain:
```bash
docker run -it \
  -v $(pwd)/validator1:/root/.selfchain \
  -e TZ=America/New_York \
  -p 36656:26656 \
  -p 36657:26657 \
  -p 36658:1234 \
  -p 36659:1317 \
  -p 36660:9090 \
  selfchain:mainnet selfchaind start
```

### Step 12: Fix Genesis Supply Errors
**Expected Issue:** You will usually encounter a "genesis supply is incorrect" error.

**Solution:** 
1. Check the error logs for expected supply values
2. Update the `supply` section in genesis with the correct amounts:

```json
"supply": [
  {
    "amount": "2116051",
    "denom": "ibc/0471F1C4E7AFD3F07702BEF6DC365268D64570F7C1FDC98EA6098DD6DE59817B"
  },
  {
    "amount": "527611",
    "denom": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"
  },
  {
    "amount": "300000000000000005",
    "denom": "ibc/AE650BD48F6712E412D2F982E44C9BB9B232F182F6D8C08B38E56613F53DCC3C"
  },
  {
    "amount": "12499593",
    "denom": "ibc/C0E66D1C81D8AAF0E6896E05190FDFBC222367148F86AC3EA679C28327A763CD"
  },
  {
    "amount": "10000",
    "denom": "ibc/EF48E6B1A1A19F47ECAEA62F5670C37C0580E86A9E88498B7E393EB6F49F33C0"
  },
  {
    "amount": "20396055695022155",
    "denom": "uslf"
  }
]
```

**Note:** Use the exact values from your error logs, not these examples.

### Step 13: Reset Database After Genesis Changes
**Critical:** Always reset the database after every genesis modification:

```bash
docker run -it \
  -v $(pwd)/validator1:/root/.selfchain \
  -e TZ=America/New_York \
  -p 36656:26656 \
  -p 36657:26657 \
  -p 36658:1234 \
  -p 36659:1317 \
  -p 36660:9090 \
  selfchain:mainnet selfchaind tendermint unsafe-reset-all
```

**Important:** Steps 12 and 13 may require multiple iterations until supply errors disappear.

### Step 14: Fix Pool Balance Errors
**Expected Issue:** After fixing supply errors, you may encounter bonded/unbonded pool balance errors.

**Solution:**
1. Check error logs for expected balances
2. Update the corresponding addresses in:
   - `not_bonded_tokens_pool` balances
   - `bonded_tokens_pool` balances
3. Always reset database after every genesis update

### Step 15: Final Balance Iterations
After updating pool balances, you may need to repeat steps 12 and 13 to fix total balance based on new log messages.

**Process:** Continue iterating between balance fixes and database resets until all errors are resolved.

## Critical Success Factors

### Error Monitoring
- **Watch logs carefully** - they provide exact error details and expected values
- **Don't ignore any errors** - each must be resolved for successful startup
- **Save working configurations** - backup successful genesis files

### Common Error Types
1. **Genesis supply mismatches** - Update supply section with exact logged values
2. **Pool balance errors** - Update bonded/unbonded pool addresses
3. **Validator power mismatches** - Ensure consistency across all validator power fields
4. **Address format errors** - Verify correct chain-specific address formats

### Iteration Requirements
- **Expect multiple iterations** - This process rarely works on the first attempt
- **Reset database after each change** - Required to clear cached state
- **Keep detailed logs** - Document errors and solutions for future reference

## Troubleshooting Tips

### If Chain Won't Start:
1. Verify all validator power values match (5000000000)
2. Check that only one validator remains in all validator arrays
3. Ensure supply totals match logged expectations
4. Confirm pool balances are correctly updated

### If Consensus Fails:
1. Verify validator public key matches production validator
2. Check that validator is not jailed in final genesis
3. Ensure validator power is consistent across all sections

### If Transactions Fail:
1. Verify account has sufficient balance
2. Check that validator is properly bonded
3. Ensure chain-id matches initialization

## Final Notes

- **This process is chain-specific** - Commands and configurations may vary
- **Keep backups** - Save working configurations at each successful stage
- **Document customizations** - Note any chain-specific modifications required
- **Test thoroughly** - Verify all functionality before using for development

**Success Indicator:** Chain starts successfully with single validator producing blocks and accepting transactions.