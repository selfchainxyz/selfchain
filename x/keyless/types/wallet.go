package types

// NewWallet creates a new Wallet instance
func NewWallet(creator, address, pubKey, chainId string) *Wallet {
	return &Wallet{
		Creator:  creator,
		Address:  address,
		PubKey:   pubKey,
		ChainId:  chainId,
		Status:   "active",
	}
}
