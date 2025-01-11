package types

// Query types for social identities
const (
	QueryGetSocialIdentity = "social_identity"
	QueryGetLinkedDID      = "linked_did"
)

// NewQuerySocialIdentityRequest creates a new QuerySocialIdentityRequest instance
func NewQuerySocialIdentityRequest(did string, provider string) *QuerySocialIdentityRequest {
	return &QuerySocialIdentityRequest{
		Did:      did,
		Provider: provider,
	}
}

// NewQueryLinkedDIDRequest creates a new QueryLinkedDIDRequest instance
func NewQueryLinkedDIDRequest(provider string, socialId string) *QueryLinkedDIDRequest {
	return &QueryLinkedDIDRequest{
		Provider:  provider,
		SocialId: socialId,
	}
}
