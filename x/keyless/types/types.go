package types

// NewWallet creates a new Wallet instance
func NewWallet(address string, did string, creator string) Wallet {
    return Wallet{
        Address:       address,
        Did:          did,
        Status:       "active",
        PersonalShare: "", // Will be set during MPC-TSS setup
        RemoteShare:   "", // Will be set during MPC-TSS setup
        Creator:       creator,
    }
}

// ValidateBasic performs basic validation of wallet fields
func (w Wallet) ValidateBasic() error {
    if w.Address == "" {
        return ErrInvalidWalletAddress
    }
    if w.Did == "" {
        return ErrInvalidWalletDID
    }
    if w.Creator == "" {
        return ErrInvalidWalletCreator
    }
    return nil
}
