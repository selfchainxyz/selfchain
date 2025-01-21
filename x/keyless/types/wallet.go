package types

import "time"

// NewWallet creates a new Wallet instance
func NewWallet(id, publicKey string, securityLevel, threshold, parties uint32) *Wallet {
	now := time.Now()
	return &Wallet{
		Id:            id,
		PublicKey:     publicKey,
		KeyVersion:    1,
		Permissions:   make([]string, 0),
		CreatedAt:     &now,
		UpdatedAt:     &now,
		Metadata:      make(map[string]string),
		SecurityLevel: securityLevel,
		Threshold:     threshold,
		Parties:       parties,
	}
}
