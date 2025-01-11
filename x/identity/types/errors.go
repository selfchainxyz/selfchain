package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/identity module sentinel errors
var (
	ErrInvalidVerificationTimeout = sdkerrors.Register(ModuleName, 1100, "invalid verification timeout")
	ErrInvalidMaxCredentials     = sdkerrors.Register(ModuleName, 1101, "invalid max credentials per DID")
	ErrNoOAuthProvidersAllowed   = sdkerrors.Register(ModuleName, 1102, "no OAuth providers allowed")
	ErrNoCredentialTypesAllowed  = sdkerrors.Register(ModuleName, 1103, "no credential types allowed")
	ErrInvalidDID               = sdkerrors.Register(ModuleName, 1104, "invalid DID")
	ErrDIDNotFound             = sdkerrors.Register(ModuleName, 1105, "DID not found")
	ErrDIDAlreadyExists        = sdkerrors.Register(ModuleName, 1106, "DID already exists")
	ErrInvalidVerificationMethod = sdkerrors.Register(ModuleName, 1107, "invalid verification method")
	ErrInvalidService          = sdkerrors.Register(ModuleName, 1108, "invalid service")
	ErrInvalidCredential       = sdkerrors.Register(ModuleName, 1109, "invalid credential")
	ErrCredentialNotFound      = sdkerrors.Register(ModuleName, 1110, "credential not found")
	ErrCredentialRevoked       = sdkerrors.Register(ModuleName, 1111, "credential is revoked")
	ErrCredentialExpired       = sdkerrors.Register(ModuleName, 1112, "credential has expired")
	ErrVerificationNotFound    = sdkerrors.Register(ModuleName, 1113, "verification record not found")
	ErrVerificationExpired     = sdkerrors.Register(ModuleName, 1114, "verification has expired")
	ErrVerificationIncomplete  = sdkerrors.Register(ModuleName, 1115, "verification not completed")
	ErrUnauthorizedProvider    = sdkerrors.Register(ModuleName, 1116, "unauthorized OAuth provider")
	ErrUnauthorizedCredentialType = sdkerrors.Register(ModuleName, 1117, "unauthorized credential type")
	ErrMaxCredentialsReached   = sdkerrors.Register(ModuleName, 1118, "maximum credentials per DID reached")
	ErrInvalidOAuthToken       = sdkerrors.Register(ModuleName, 1119, "invalid OAuth token")
	ErrInvalidRecoveryAddress  = sdkerrors.Register(ModuleName, 1120, "invalid recovery address")
	ErrInvalidCredentialSchema = sdkerrors.Register(ModuleName, 1121, "invalid credential schema")
)
