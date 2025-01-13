package types

// Wallet represents a keyless wallet in the system
type Wallet struct {
    // Address is the wallet's address on the blockchain
    Address string `json:"address"`
    
    // DID is the decentralized identifier associated with this wallet
    DID string `json:"did"`
    
    // Status represents the current state of the wallet (e.g., active, recovering)
    Status string `json:"status"`
    
    // PersonalShare is the encrypted key share stored on the user's device
    PersonalShare string `json:"personal_share,omitempty"`
    
    // RemoteShare is the key share stored on the chain (encrypted)
    RemoteShare string `json:"remote_share,omitempty"`
    
    // Creator is the address that created this wallet
    Creator string `json:"creator"`
}

// NewWallet creates a new Wallet instance
func NewWallet(address, did, creator string) Wallet {
    return Wallet{
        Address: address,
        DID:     did,
        Status:  "active",
        Creator: creator,
    }
}

// ValidateBasic performs basic validation of wallet fields
func (w Wallet) ValidateBasic() error {
    if w.Address == "" {
        return ErrInvalidWalletAddress
    }
    if w.DID == "" {
        return ErrInvalidDID
    }
    if w.Creator == "" {
        return ErrInvalidCreator
    }
    return nil
}
