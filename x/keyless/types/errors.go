package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/keyless module sentinel errors
var (
	// Basic errors
	ErrInvalidRequest        = sdkerrors.Register(ModuleName, 1, "invalid request")
	ErrInvalidResponse       = sdkerrors.Register(ModuleName, 2, "invalid response")
	ErrUnauthorized         = sdkerrors.Register(ModuleName, 3, "unauthorized")
	ErrInvalidParam         = sdkerrors.Register(ModuleName, 4, "invalid parameter")

	// Wallet errors
	ErrWalletNotFound       = sdkerrors.Register(ModuleName, 10, "wallet not found")
	ErrWalletExists         = sdkerrors.Register(ModuleName, 11, "wallet already exists")
	ErrInvalidWalletStatus  = sdkerrors.Register(ModuleName, 12, "invalid wallet status")
	ErrInvalidWalletAddress = sdkerrors.Register(ModuleName, 13, "invalid wallet address")
	ErrInvalidWalletID      = sdkerrors.Register(ModuleName, 14, "invalid wallet ID")
	ErrInvalidWalletDID     = sdkerrors.Register(ModuleName, 15, "invalid wallet DID")
	ErrInvalidWalletCreator = sdkerrors.Register(ModuleName, 16, "invalid wallet creator")

	// Permission errors
	ErrInvalidPermission    = sdkerrors.Register(ModuleName, 20, "invalid permission")
	ErrPermissionNotFound   = sdkerrors.Register(ModuleName, 21, "permission not found")
	ErrPermissionExpired    = sdkerrors.Register(ModuleName, 22, "permission expired")
	ErrPermissionRevoked    = sdkerrors.Register(ModuleName, 23, "permission revoked")

	// Key management errors
	ErrKeyGenFailed         = sdkerrors.Register(ModuleName, 30, "key generation failed")
	ErrInvalidPublicKey     = sdkerrors.Register(ModuleName, 31, "invalid public key")
	ErrInvalidKeyVersion    = sdkerrors.Register(ModuleName, 32, "invalid key version")
	ErrInvalidKeyRotationStatus = sdkerrors.Register(ModuleName, 33, "invalid key rotation status")
	ErrKeyRotationNotFound  = sdkerrors.Register(ModuleName, 34, "key rotation not found")

	// Share management errors
	ErrInvalidShare         = sdkerrors.Register(ModuleName, 40, "invalid share")
	ErrShareNotFound        = sdkerrors.Register(ModuleName, 41, "share not found")
	ErrEncryptFailed        = sdkerrors.Register(ModuleName, 42, "encryption failed")
	ErrDecryptFailed        = sdkerrors.Register(ModuleName, 43, "decryption failed")

	// Recovery errors
	ErrInvalidRecoveryProof = sdkerrors.Register(ModuleName, 50, "invalid recovery proof")
	ErrRecoveryInProgress   = sdkerrors.Register(ModuleName, 51, "recovery already in progress")
	ErrRecoveryNotAllowed   = sdkerrors.Register(ModuleName, 52, "recovery not allowed")
	ErrRecoveryNotFound     = sdkerrors.Register(ModuleName, 53, "recovery session not found")
	ErrRecoveryExpired      = sdkerrors.Register(ModuleName, 54, "recovery session expired")
	ErrInvalidRecoveryToken = sdkerrors.Register(ModuleName, 55, "invalid recovery token")

	// Security errors
	ErrInvalidSignature     = sdkerrors.Register(ModuleName, 60, "invalid signature")
	ErrInvalidChainID       = sdkerrors.Register(ModuleName, 61, "invalid chain ID")
	ErrInvalidSecurityLevel = sdkerrors.Register(ModuleName, 62, "invalid security level")
	ErrRateLimitExceeded    = sdkerrors.Register(ModuleName, 63, "rate limit exceeded")
	ErrMFARequired          = sdkerrors.Register(ModuleName, 64, "MFA verification required")
	ErrMFAFailed            = sdkerrors.Register(ModuleName, 65, "MFA verification failed")

	// Party data errors
	ErrInvalidPartyData     = sdkerrors.Register(ModuleName, 70, "invalid party data")
	ErrInvalidRotationProof = sdkerrors.Register(ModuleName, 71, "invalid rotation proof")
)
