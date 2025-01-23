package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/identity module sentinel errors
var (
	// Common errors (1-99)
	ErrInvalidDID         = sdkerrors.Register(ModuleName, 1, "invalid DID")
	ErrDIDNotFound        = sdkerrors.Register(ModuleName, 2, "DID not found")
	ErrInvalidRequest     = sdkerrors.Register(ModuleName, 3, "invalid request")
	ErrDIDAlreadyExists   = sdkerrors.Register(ModuleName, 4, "DID already exists")
	ErrInvalidDIDDocument = sdkerrors.Register(ModuleName, 5, "invalid DID document")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 6, "unauthorized")
	ErrRateLimitExceeded  = sdkerrors.Register(ModuleName, 7, "rate limit exceeded")
	ErrDIDInactive        = sdkerrors.Register(ModuleName, 8, "DID is inactive")

	// Credential errors (100-199)
	ErrInvalidCredential        = sdkerrors.Register(ModuleName, 100, "invalid credential")
	ErrCredentialNotFound       = sdkerrors.Register(ModuleName, 101, "credential not found")
	ErrInvalidPresentation      = sdkerrors.Register(ModuleName, 102, "invalid presentation")
	ErrPresentationNotFound     = sdkerrors.Register(ModuleName, 103, "presentation not found")
	ErrInvalidProof             = sdkerrors.Register(ModuleName, 104, "invalid proof")
	ErrInvalidSignature         = sdkerrors.Register(ModuleName, 105, "invalid signature")
	ErrInvalidStatus            = sdkerrors.Register(ModuleName, 106, "invalid status")
	ErrInvalidType              = sdkerrors.Register(ModuleName, 107, "invalid type")
	ErrInvalidIssuer            = sdkerrors.Register(ModuleName, 108, "invalid issuer")
	ErrInvalidSubject           = sdkerrors.Register(ModuleName, 109, "invalid subject")
	ErrCredentialAlreadyRevoked = sdkerrors.Register(ModuleName, 110, "credential already revoked")
	ErrInvalidExpiryDate        = sdkerrors.Register(ModuleName, 111, "invalid expiry date")
	ErrInvalidCredentialID      = sdkerrors.Register(ModuleName, 112, "invalid credential ID")
	ErrInvalidCredentialType    = sdkerrors.Register(ModuleName, 113, "invalid credential type")
	ErrInvalidCredentialSubject = sdkerrors.Register(ModuleName, 114, "invalid credential subject")
	ErrInvalidCredentialClaims  = sdkerrors.Register(ModuleName, 115, "invalid credential claims")
	ErrInvalidCredentialIssuer  = sdkerrors.Register(ModuleName, 116, "invalid credential issuer")

	// MFA errors (200-299)
	ErrInvalidMFAMethod       = sdkerrors.Register(ModuleName, 200, "invalid MFA method")
	ErrMFAMethodNotFound      = sdkerrors.Register(ModuleName, 201, "MFA method not found")
	ErrMFAMethodInactive      = sdkerrors.Register(ModuleName, 202, "MFA method is not active")
	ErrInvalidMFAChallenge    = sdkerrors.Register(ModuleName, 203, "invalid MFA challenge")
	ErrMFAChallengeExpired    = sdkerrors.Register(ModuleName, 204, "MFA challenge expired")
	ErrInvalidMFAConfig       = sdkerrors.Register(ModuleName, 205, "invalid MFA configuration")
	ErrMFAVerificationFailed  = sdkerrors.Register(ModuleName, 206, "MFA verification failed")
	ErrInvalidMFACode         = sdkerrors.Register(ModuleName, 207, "invalid MFA code")
	ErrMFAMethodAlreadyExists = sdkerrors.Register(ModuleName, 208, "MFA method already exists")

	// OAuth/Social errors (300-399)
	ErrInvalidSocialIdentity = sdkerrors.Register(ModuleName, 300, "invalid social identity")
	ErrInvalidToken          = sdkerrors.Register(ModuleName, 301, "invalid token")

	// Audit log errors (400-499)
	ErrInvalidAuditLogID = sdkerrors.Register(ModuleName, 400, "invalid audit log ID")
	ErrInvalidAction     = sdkerrors.Register(ModuleName, 401, "invalid action")
	ErrInvalidActor      = sdkerrors.Register(ModuleName, 402, "invalid actor")
	ErrInvalidController = sdkerrors.Register(ModuleName, 403, "invalid controller")
	ErrInvalidService    = sdkerrors.Register(ModuleName, 404, "invalid service")

	// Recovery errors (400-499)
	ErrInvalidRecoverySession = sdkerrors.Register(ModuleName, 400, "invalid recovery session")

	// Rate limit errors (500-599)
	ErrInvalidRateLimit = sdkerrors.Register(ModuleName, 500, "invalid rate limit configuration")
)
