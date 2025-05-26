## Local Development Setup

### Branch: `develop2(Current Mainnet)`

1. **Modify Dockerfile**:
   Remove the entrypoint line (`ENTRYPOINT ["selfchaind"]`) from the Dockerfile:

   ```
   operations/Dockerfile-ubuntu-prod
   ```

2. **Build Binary**:

   ```bash
   GOOS=linux GOARCH=arm64 go build -tags ledger -o ./build/selfchaind-linux-arm64 ./cmd/selfchaind/main.go
   ```

3. **Build Docker Image**:

   ```bash
   docker build --build-arg BUILDARCH=arm64 -f operations/Dockerfile-ubuntu-prod . -t selfchain:mainnet
   ```

---

### Branch: `upgrade_cosmos_sdk_50`

1. **Build Docker Image** (no manual Dockerfile edits required):

   ```bash
   docker buildx build --platform linux/arm64 --tag selfchain:develop -f operations/Dockerfile .
   ```

2. **Handle Local Testing Configurations**:

    * Ensure `vestingAddresses` and `addressReplacements` are properly configured in `handler.go`.

3. **Run Setup and Upgrade Scripts**:

   ```bash
   bash ./testlocal/setup_validators.sh --clean
   ./testlocal/submit-upgrade-proposal.sh
   ./testlocal/upgrade-validators.sh
   ```
