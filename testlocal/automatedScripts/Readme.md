# Validator Management Scripts Quick Guide

## Setup Validators
```bash
./setup_validators.sh
# Sets up 3 validators with proper configuration and networking.
# Use --clean flag to remove existing containers before setup.
```

## Submit Upgrade Proposal
```bash
./submit-upgrade-proposal.sh  
# Submits an upgrade proposal (v4) and votes with all validators.
# Automatically calculates upgrade height and monitors proposal status.
```

## Upgrade Validators
```bash
./upgrade-validators.sh
# Upgrades all validators to new Docker image (selfchainprod-v4:mainnet-v4).
# Safely stops old containers and restarts with new image while maintaining data.
```

Each script provides colored output with status updates and handles errors gracefully. Run them in order: setup → propose → upgrade.