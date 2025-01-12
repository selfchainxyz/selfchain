package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateBasic performs basic validation of the credential
func (c *Credential) ValidateBasic() error {
	if c.Id == "" {
		return sdkerrors.Register(ModuleName, 1100, "invalid credential ID")
	}
	if c.Type == "" {
		return sdkerrors.Register(ModuleName, 1101, "invalid credential type")
	}
	if c.Issuer == "" {
		return sdkerrors.Register(ModuleName, 1102, "invalid credential issuer")
	}
	if c.Subject == "" {
		return sdkerrors.Register(ModuleName, 1103, "invalid credential subject")
	}
	if len(c.Claims) == 0 {
		return sdkerrors.Register(ModuleName, 1104, "invalid credential claims")
	}
	if c.IssuanceDate == nil {
		return sdkerrors.Register(ModuleName, 1105, "invalid credential expiry")
	}
	now := time.Date(2025, 1, 12, 14, 56, 5, 0, time.FixedZone("IST", 5*60*60+30*60))
	if c.ExpirationDate != nil && c.ExpirationDate.Before(now) {
		return sdkerrors.Register(ModuleName, 1105, "invalid credential expiry")
	}
	return nil
}
