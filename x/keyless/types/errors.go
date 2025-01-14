package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/keyless module sentinel errors
var (
	ErrUnauthorized = sdkerrors.Register(ModuleName, 1, "unauthorized")
	ErrInvalidMaxWallets = sdkerrors.Register(ModuleName, 2, "invalid max wallets per DID")
	ErrInvalidMaxShares = sdkerrors.Register(ModuleName, 3, "invalid max shares per wallet")
	ErrInvalidRecoveryThreshold = sdkerrors.Register(ModuleName, 4, "invalid recovery threshold")
	ErrInvalidRecoveryWindow = sdkerrors.Register(ModuleName, 5, "invalid recovery window")
	ErrInvalidMaxAttempts = sdkerrors.Register(ModuleName, 6, "invalid max signing attempts")
	ErrWalletExists = sdkerrors.Register(ModuleName, 7, "wallet already exists")
	ErrWalletNotFound = sdkerrors.Register(ModuleName, 8, "wallet not found")
	ErrInvalidWalletAddress = sdkerrors.Register(ModuleName, 9, "invalid wallet address")
	ErrInvalidWalletDID = sdkerrors.Register(ModuleName, 10, "invalid wallet DID")
	ErrInvalidWalletCreator = sdkerrors.Register(ModuleName, 11, "invalid wallet creator")
	ErrInvalidSignature = sdkerrors.Register(ModuleName, 12, "invalid signature")
	ErrSigningAttemptExceeded = sdkerrors.Register(ModuleName, 13, "signing attempt limit exceeded")
	ErrRecoveryInProgress = sdkerrors.Register(ModuleName, 14, "recovery already in progress")
	ErrRecoveryNotAllowed = sdkerrors.Register(ModuleName, 15, "recovery not allowed")
	ErrInvalidRecoveryProof = sdkerrors.Register(ModuleName, 16, "invalid recovery proof")
	ErrInvalidRequest   = sdkerrors.Register(ModuleName, 1100, "invalid request")
	ErrInvalidResponse  = sdkerrors.Register(ModuleName, 1101, "invalid response")
	ErrKeyGenFailed     = sdkerrors.Register(ModuleName, 1102, "key generation failed")
	ErrEncryptFailed    = sdkerrors.Register(ModuleName, 1103, "encryption failed")
	ErrDecryptFailed    = sdkerrors.Register(ModuleName, 1104, "decryption failed")
	ErrShareNotFound    = sdkerrors.Register(ModuleName, 1105, "share not found")
	ErrInvalidShare     = sdkerrors.Register(ModuleName, 1106, "invalid share")
)
