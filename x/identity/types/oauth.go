package types

import (
	"fmt"
	"strings"
)

// ValidateProvider checks if the provider is valid
func ValidateProvider(provider string) error {
	provider = strings.ToLower(provider)
	switch provider {
	case "google", "github", "apple":
		return nil
	default:
		return fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
}

// ValidateOAuth validates the OAuth message
func ValidateOAuth(provider string, token string) error {
	if err := ValidateProvider(provider); err != nil {
		return err
	}
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}
	return nil
}

// ValidateSocialIdentity validates the social identity
func ValidateSocialIdentity(did string, provider string, socialId string) error {
	if err := ValidateProvider(provider); err != nil {
		return err
	}
	if socialId == "" {
		return fmt.Errorf("social ID cannot be empty")
	}
	if did == "" {
		return fmt.Errorf("DID cannot be empty")
	}
	return nil
}
