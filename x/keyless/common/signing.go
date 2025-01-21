package common

// SigningParams contains parameters for signing operations
type SigningParams struct {
	NetworkID string
	Message   []byte
	Metadata  map[string]interface{}
}