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
	ErrInvalidWalletCreator = sdkerrors.Register(ModuleName, 14, "invalid wallet creator")
	ErrWalletLocked         = sdkerrors.Register(ModuleName, 1202, "wallet is locked")
	ErrWalletInactive       = sdkerrors.Register(ModuleName, 1203, "wallet is inactive")
	ErrWalletBusy           = sdkerrors.Register(ModuleName, 1204, "wallet is busy")

	// Permission errors
	ErrInvalidPermission    = sdkerrors.Register(ModuleName, 20, "invalid permission")
	ErrPermissionNotFound   = sdkerrors.Register(ModuleName, 21, "permission not found")
	ErrPermissionExpired    = sdkerrors.Register(ModuleName, 22, "permission expired")
	ErrPermissionRevoked    = sdkerrors.Register(ModuleName, 23, "permission revoked")
	ErrEmptyPermissions     = sdkerrors.Register(ModuleName, 24, "permissions cannot be empty")

	// Key management errors
	ErrKeyGenFailed         = sdkerrors.Register(ModuleName, 30, "key generation failed")
	ErrInvalidPublicKey     = sdkerrors.Register(ModuleName, 1103, "invalid public key")
	ErrInvalidKeyVersion    = sdkerrors.Register(ModuleName, 32, "invalid key version")
	ErrInvalidKeyRotationStatus = sdkerrors.Register(ModuleName, 33, "invalid key rotation status")
	ErrKeyRotationNotFound  = sdkerrors.Register(ModuleName, 34, "key rotation not found")
	ErrKeyGeneration        = sdkerrors.Register(ModuleName, 1400, "key generation failed")
	ErrKeyRotation          = sdkerrors.Register(ModuleName, 1401, "key rotation failed")
	ErrKeyBackup            = sdkerrors.Register(ModuleName, 1402, "key backup failed")
	ErrKeyRecovery          = sdkerrors.Register(ModuleName, 1403, "key recovery failed")
	ErrKeyImport            = sdkerrors.Register(ModuleName, 1404, "key import failed")
	ErrKeyExport            = sdkerrors.Register(ModuleName, 1405, "key export failed")

	// Share management errors
	ErrInvalidShare         = sdkerrors.Register(ModuleName, 1300, "invalid share")
	ErrShareNotFound        = sdkerrors.Register(ModuleName, 41, "share not found")
	ErrEncryptFailed        = sdkerrors.Register(ModuleName, 42, "encryption failed")
	ErrDecryptFailed        = sdkerrors.Register(ModuleName, 43, "decryption failed")
	ErrExpiredShare         = sdkerrors.Register(ModuleName, 1301, "share expired")
	ErrInvalidShareSignature = sdkerrors.Register(ModuleName, 1302, "invalid share signature")
	ErrInsufficientShares   = sdkerrors.Register(ModuleName, 1303, "insufficient shares")
	ErrInsufficientValidShares = sdkerrors.Register(ModuleName, 1304, "insufficient valid shares")

	// Recovery errors
	ErrInvalidRecoveryProof = sdkerrors.Register(ModuleName, 50, "invalid recovery proof")
	ErrRecoveryInProgress   = sdkerrors.Register(ModuleName, 51, "recovery already in progress")
	ErrRecoveryNotAllowed   = sdkerrors.Register(ModuleName, 52, "recovery not allowed")
	ErrRecoveryNotFound     = sdkerrors.Register(ModuleName, 53, "recovery session not found")
	ErrRecoveryExpired      = sdkerrors.Register(ModuleName, 54, "recovery session expired")
	ErrInvalidRecoveryToken = sdkerrors.Register(ModuleName, 55, "invalid recovery token")

	// Security errors
	ErrInvalidSignature     = sdkerrors.Register(ModuleName, 1102, "invalid signature")
	ErrInvalidChainID       = sdkerrors.Register(ModuleName, 1101, "invalid chain ID")
	ErrInvalidSecurityLevel = sdkerrors.Register(ModuleName, 62, "invalid security level")
	ErrRateLimitExceeded    = sdkerrors.Register(ModuleName, 63, "rate limit exceeded")
	ErrMFARequired          = sdkerrors.Register(ModuleName, 64, "MFA verification required")
	ErrMFAFailed            = sdkerrors.Register(ModuleName, 65, "MFA verification failed")

	// Party data errors
	ErrInvalidPartyData     = sdkerrors.Register(ModuleName, 1600, "invalid party data")
	ErrInvalidRotationProof = sdkerrors.Register(ModuleName, 1601, "invalid rotation proof")
	ErrInvalidWalletID      = sdkerrors.Register(ModuleName, 1100, "invalid wallet ID")
	ErrInvalidKeyType       = sdkerrors.Register(ModuleName, 1104, "invalid key type")
	ErrInvalidKeyFormat     = sdkerrors.Register(ModuleName, 1105, "invalid key format")
	ErrInvalidMetadata      = sdkerrors.Register(ModuleName, 1106, "invalid metadata")
	ErrInvalidStatus        = sdkerrors.Register(ModuleName, 1107, "invalid status")
	ErrInvalidThreshold     = sdkerrors.Register(ModuleName, 1108, "invalid threshold")
	ErrInvalidParties       = sdkerrors.Register(ModuleName, 1109, "invalid parties")

	// Protocol errors
	ErrProtocolNotFound     = sdkerrors.Register(ModuleName, 1500, "protocol not found")
	ErrProtocolFailed       = sdkerrors.Register(ModuleName, 1501, "protocol failed")
	ErrInvalidProtocol      = sdkerrors.Register(ModuleName, 1502, "invalid protocol")
	ErrProtocolTimeout      = sdkerrors.Register(ModuleName, 1503, "protocol timeout")

	// Additional errors
	ErrKeyReconstruction      = sdkerrors.Register(ModuleName, 1114, "key reconstruction failed")
	ErrStoringKeyShare        = sdkerrors.Register(ModuleName, 1115, "failed to store key share")
	ErrTSSProtocolNotInit     = sdkerrors.Register(ModuleName, 1116, "TSS protocol not initialized")
)
