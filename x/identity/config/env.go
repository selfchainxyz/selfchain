package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// OAuthProviderConfig holds configuration for an OAuth provider
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	ProfileURL   string
	Scopes       []string
}

// Config holds all configuration for the identity module
type Config struct {
	GoogleOAuth OAuthProviderConfig
	GithubOAuth OAuthProviderConfig
}

// Required environment variables
const (
	EnvGoogleClientID     = "SELFCHAIN_GOOGLE_OAUTH_CLIENT_ID"
	EnvGoogleClientSecret = "SELFCHAIN_GOOGLE_OAUTH_CLIENT_SECRET"
	EnvGithubClientID     = "SELFCHAIN_GITHUB_OAUTH_CLIENT_ID"
	EnvGithubClientSecret = "SELFCHAIN_GITHUB_OAUTH_CLIENT_SECRET"
)

// Default OAuth endpoints
const (
	GoogleAuthURL    = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL   = "https://oauth2.googleapis.com/token"
	GoogleProfileURL = "https://www.googleapis.com/oauth2/v3/userinfo"
	
	GithubAuthURL    = "https://github.com/login/oauth/authorize"
	GithubTokenURL   = "https://github.com/login/oauth/access_token"
	GithubProfileURL = "https://api.github.com/user"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Check required environment variables
	missingVars := []string{}
	requiredVars := []string{
		EnvGoogleClientID,
		EnvGoogleClientSecret,
		EnvGithubClientID,
		EnvGithubClientSecret,
	}

	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	return &Config{
		GoogleOAuth: OAuthProviderConfig{
			ClientID:     os.Getenv(EnvGoogleClientID),
			ClientSecret: os.Getenv(EnvGoogleClientSecret),
			AuthURL:      GoogleAuthURL,
			TokenURL:     GoogleTokenURL,
			ProfileURL:   GoogleProfileURL,
			Scopes:       []string{"openid", "email", "profile"},
		},
		GithubOAuth: OAuthProviderConfig{
			ClientID:     os.Getenv(EnvGithubClientID),
			ClientSecret: os.Getenv(EnvGithubClientSecret),
			AuthURL:      GithubAuthURL,
			TokenURL:     GithubTokenURL,
			ProfileURL:   GithubProfileURL,
			Scopes:       []string{"read:user", "user:email"},
		},
	}, nil
}
