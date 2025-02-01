# Keyless Module Documentation

## Overview

The keyless module implements a secure, distributed key management system using threshold signature scheme (TSS) for the Self Chain. It eliminates the need for private keys by leveraging distributed key shares and provides robust wallet recovery mechanisms through decentralized identifiers (DIDs).

## Core Components

### 1. TSS Protocol

The TSS protocol implementation provides:
- Distributed key generation
- Threshold signing
- Key reconstruction
- Share verification

```go
type TSSProtocol interface {
    GenerateKeyShares(ctx sdk.Context, walletAddress string, threshold uint32, securityLevel SecurityLevel) (*KeyGenResponse, error)
    SignMessage(ctx sdk.Context, message []byte, shares [][]byte) ([]byte, error)
    ReconstructKey(ctx sdk.Context, shares [][]byte) ([]byte, error)
    VerifyShare(ctx sdk.Context, share []byte, publicKey []byte) error
}
```

### 2. Keeper Layer

The keeper layer manages:
- Wallet state
- Key share storage
- Recovery processes
- Identity integration

```go
type Keeper interface {
    CreateWallet(ctx sdk.Context, did string, threshold uint32, securityLevel SecurityLevel) (*Wallet, error)
    SignTransaction(ctx sdk.Context, walletAddress string, message []byte) ([]byte, error)
    ReconstructWallet(ctx sdk.Context, didDoc DIDDocument) ([]byte, error)
}
```

## Security Features

### 1. Rate Limiting
- Per-wallet operation limits
- Global rate limiting
- Configurable thresholds

### 2. Multi-Factor Authentication
- OAuth2 integration
- DID-based verification
- Custom MFA providers

### 3. Audit Logging
- Operation tracking
- Security event logging
- Recovery attempt monitoring

## Performance Characteristics

### 1. Wallet Creation
- Average latency: < 1 second
- Gas consumption: ~800,000 gas units
- Success rate: > 99.9%

### 2. Transaction Signing
- Average latency: < 500ms
- Gas consumption: ~400,000 gas units
- Concurrent operations: Up to 100/second

### 3. Wallet Recovery
- Average latency: < 2 seconds
- Success rate: > 99%
- Rate limited to prevent abuse

## Integration Guide

### 1. Module Setup
```go
// In app.go
app.KeylessKeeper = keylesskeeper.NewKeeper(
    appCodec,
    keys[keyless.StoreKey],
    keys[keyless.MemStoreKey],
    app.GetSubspace(keyless.ModuleName),
    app.IdentityKeeper,
    tss.NewTSSProtocolImpl(),
)
```

### 2. Creating a Wallet
```go
// Client code
msg := types.NewMsgCreateWallet(
    creator,
    "did:selfchain:123",
    2,
    types.SecurityLevel_HIGH,
)
```

### 3. Signing Transactions
```go
// Client code
msg := types.NewMsgSignTransaction(
    creator,
    walletAddress,
    txBytes,
)
```

### 4. Wallet Recovery
```go
// Client code
msg := types.NewMsgRecoverWallet(
    creator,
    walletAddress,
    recoveryToken,
)
```

## Error Handling

### Common Error Types
1. `ErrInvalidWallet`: Invalid wallet address or not found
2. `ErrInsufficientShares`: Not enough shares for reconstruction
3. `ErrUnauthorized`: Unauthorized operation attempt
4. `ErrRateLimitExceeded`: Too many operations
5. `ErrNetworkPartition`: Network communication issues

### Recovery Procedures
1. Network Partition
   - Retry with exponential backoff
   - Fallback to secondary nodes
   
2. Share Recovery
   - Verify DID document
   - Reconstruct from backup shares
   - Re-issue shares if needed

## Best Practices

1. Security
   - Always verify DID documents
   - Implement proper rate limiting
   - Use secure communication channels

2. Performance
   - Cache frequently used wallets
   - Batch similar operations
   - Monitor gas consumption

3. Recovery
   - Regular share backups
   - Test recovery procedures
   - Document recovery processes

## Gas Optimization

### Gas Costs
Operation | Average Gas Cost
----------|----------------
Wallet Creation | 800,000
Transaction Signing | 400,000
Share Verification | 100,000
Wallet Recovery | 600,000

### Optimization Tips
1. Batch similar operations
2. Cache frequently used data
3. Optimize state access patterns
4. Use appropriate security levels

## Testing

### Test Coverage
```bash
# Run unit tests
go test ./x/keyless/...

# Run benchmarks
go test -bench=. ./x/keyless/...

# Run integration tests
go test -tags=integration ./x/keyless/...
```

### Performance Testing
```bash
# Benchmark wallet creation
go test -bench=BenchmarkWalletCreation ./x/keyless/...

# Benchmark transaction signing
go test -bench=BenchmarkSignTransaction ./x/keyless/...
```
