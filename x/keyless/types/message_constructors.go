package types

// NewMsgCreateWallet creates a new MsgCreateWallet instance
func NewMsgCreateWallet(creator, pubKey, walletAddress, chainId string) *MsgCreateWallet {
	return &MsgCreateWallet{
		Creator:       creator,
		PubKey:       pubKey,
		WalletAddress: walletAddress,
		ChainId:      chainId,
	}
}

// NewMsgRecoverWallet creates a new MsgRecoverWallet instance
func NewMsgRecoverWallet(creator, walletAddress, recoveryProof, newPubKey, signature string) *MsgRecoverWallet {
	return &MsgRecoverWallet{
		Creator:       creator,
		WalletAddress: walletAddress,
		RecoveryProof: recoveryProof,
		NewPubKey:    newPubKey,
		Signature:    signature,
	}
}

// NewMsgSignTransaction creates a new MsgSignTransaction instance
func NewMsgSignTransaction(creator, walletAddress, unsignedTx, chainId string) *MsgSignTransaction {
	return &MsgSignTransaction{
		Creator:       creator,
		WalletAddress: walletAddress,
		UnsignedTx:   unsignedTx,
		ChainId:      chainId,
	}
}
