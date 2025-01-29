package types

// NewWallet creates a new Wallet instance
func NewWallet(
	creator string,
	publicKey string,
	walletAddress string,
	chainId string,
	status WalletStatus,
	keyVersion uint32,
) *Wallet {
	return &Wallet{
		Id:            walletAddress, // Set id to be same as wallet_address for consistency
		Creator:       creator,
		PublicKey:     publicKey,
		WalletAddress: walletAddress,
		ChainId:      chainId,
		Status:       status,
		KeyVersion:   keyVersion,
	}
}

// ValidateBasic performs basic validation of a Wallet
func (w *Wallet) ValidateBasic() error {
	if w.Creator == "" {
		return ErrInvalidWalletCreator
	}

	if w.PublicKey == "" {
		return ErrInvalidPublicKey
	}

	if w.ChainId == "" {
		return ErrInvalidChainID
	}

	return nil
}
