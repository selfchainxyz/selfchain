# Keyless Module API Documentation

## Message Types

### MsgCreateWallet

Creates a new keyless wallet with distributed key shares.

```protobuf
message MsgCreateWallet {
    string creator = 1;
    string did = 2;
    uint32 threshold = 3;
    SecurityLevel security_level = 4;
}
```

**Parameters:**
- `creator`: The address creating the wallet
- `did`: Decentralized identifier for the wallet
- `threshold`: Number of shares required for reconstruction
- `security_level`: Security level enum (LOW, MEDIUM, HIGH)

**Response:**
```protobuf
message MsgCreateWalletResponse {
    string wallet_address = 1;
    bytes public_key = 2;
    KeyMetadata metadata = 3;
}
```

### MsgSignTransaction

Signs a transaction using the distributed key shares.

```protobuf
message MsgSignTransaction {
    string creator = 1;
    string wallet_address = 2;
    bytes message = 3;
}
```

**Parameters:**
- `creator`: The address initiating the signing
- `wallet_address`: The wallet to sign with
- `message`: The message to sign

**Response:**
```protobuf
message MsgSignTransactionResponse {
    bytes signature = 1;
    SignatureMetadata metadata = 2;
}
```

### MsgRecoverWallet

Initiates wallet recovery using DID-based verification.

```protobuf
message MsgRecoverWallet {
    string creator = 1;
    string wallet_address = 2;
    string recovery_token = 3;
    string recovery_address = 4;
}
```

**Parameters:**
- `creator`: The address initiating recovery
- `wallet_address`: The wallet to recover
- `recovery_token`: OAuth2 or DID-based recovery token
- `recovery_address`: Address to receive recovered wallet

**Response:**
```protobuf
message MsgRecoverWalletResponse {
    string wallet_address = 1;
    bytes public_key = 2;
    RecoveryMetadata metadata = 3;
}
```

## Query Types

### QueryWalletRequest

Queries wallet information.

```protobuf
message QueryWalletRequest {
    string wallet_address = 1;
}
```

**Response:**
```protobuf
message QueryWalletResponse {
    Wallet wallet = 1;
}
```

### QueryRecoveryStatusRequest

Queries recovery status.

```protobuf
message QueryRecoveryStatusRequest {
    string wallet_address = 1;
}
```

**Response:**
```protobuf
message QueryRecoveryStatusResponse {
    RecoveryStatus status = 1;
    RecoveryMetadata metadata = 2;
}
```

## Events

### EventCreateWallet
```protobuf
message EventCreateWallet {
    string wallet_address = 1;
    string did = 2;
    uint32 threshold = 3;
    SecurityLevel security_level = 4;
    int64 created_at = 5;
}
```

### EventSignTransaction
```protobuf
message EventSignTransaction {
    string wallet_address = 1;
    bytes message_hash = 2;
    int64 signed_at = 3;
}
```

### EventRecoverWallet
```protobuf
message EventRecoverWallet {
    string wallet_address = 1;
    string recovery_address = 2;
    RecoveryStatus status = 3;
    int64 recovered_at = 4;
}
```

## Error Codes

```go
const (
    ErrInvalidWallet       = sdkerrors.Register(ModuleName, 1, "invalid wallet")
    ErrInsufficientShares  = sdkerrors.Register(ModuleName, 2, "insufficient shares")
    ErrUnauthorized        = sdkerrors.Register(ModuleName, 3, "unauthorized")
    ErrRateLimitExceeded   = sdkerrors.Register(ModuleName, 4, "rate limit exceeded")
    ErrNetworkPartition    = sdkerrors.Register(ModuleName, 5, "network partition")
    ErrInvalidRecoveryToken = sdkerrors.Register(ModuleName, 6, "invalid recovery token")
    ErrRecoveryInProgress  = sdkerrors.Register(ModuleName, 7, "recovery in progress")
)
```

## CLI Commands

### Create Wallet
```bash
selfchaind tx keyless create-wallet [did] [threshold] [security-level] \
    --from [creator] \
    --chain-id [chain-id]
```

### Sign Transaction
```bash
selfchaind tx keyless sign-tx [wallet-address] [message] \
    --from [creator] \
    --chain-id [chain-id]
```

### Recover Wallet
```bash
selfchaind tx keyless recover-wallet [wallet-address] [recovery-token] [recovery-address] \
    --from [creator] \
    --chain-id [chain-id]
```

### Query Wallet
```bash
selfchaind query keyless wallet [wallet-address]
```

### Query Recovery Status
```bash
selfchaind query keyless recovery-status [wallet-address]
```
