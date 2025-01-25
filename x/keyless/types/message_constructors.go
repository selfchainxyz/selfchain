package types

// NewMsgCreateWallet creates a new MsgCreateWallet instance
func NewMsgCreateWallet(creator, pubKey, walletAddress, chainId string) *MsgCreateWallet {
	return &MsgCreateWallet{
		Creator:       creator,
		PubKey:        pubKey,
		WalletAddress: walletAddress,
		ChainId:      chainId,
	}
}

// NewMsgRecoverWallet creates a new MsgRecoverWallet instance
func NewMsgRecoverWallet(creator, walletAddress, newPubKey, recoveryProof string) *MsgRecoverWallet {
	return &MsgRecoverWallet{
		Creator:       creator,
		WalletAddress: walletAddress,
		NewPubKey:     newPubKey,
		RecoveryProof: recoveryProof,
	}
}

// NewMsgSignTransaction creates a new MsgSignTransaction instance
func NewMsgSignTransaction(creator, walletAddress, unsignedTx string) *MsgSignTransaction {
	return &MsgSignTransaction{
		Creator:       creator,
		WalletAddress: walletAddress,
		UnsignedTx:    unsignedTx,
	}
}

// NewMsgBatchSignRequest creates a new MsgBatchSignRequest instance
func NewMsgBatchSignRequest(creator, walletAddress string, messages [][]byte, parties []string) *MsgBatchSignRequest {
	return &MsgBatchSignRequest{
		Creator:       creator,
		WalletAddress: walletAddress,
		Messages:      messages,
		Parties:       parties,
	}
}

// NewMsgInitiateKeyRotation creates a new MsgInitiateKeyRotation instance
func NewMsgInitiateKeyRotation(creator, walletAddress, newPubKey, signature string) *MsgInitiateKeyRotation {
	return &MsgInitiateKeyRotation{
		Creator:       creator,
		WalletAddress: walletAddress,
		NewPubKey:     newPubKey,
		Signature:     signature,
	}
}

// NewMsgCompleteKeyRotation creates a new MsgCompleteKeyRotation instance
func NewMsgCompleteKeyRotation(creator, walletAddress, version, signature, newPubKey string) *MsgCompleteKeyRotation {
	return &MsgCompleteKeyRotation{
		Creator:       creator,
		WalletAddress: walletAddress,
		Version:       version,
		Signature:     signature,
		NewPubKey:     newPubKey,
	}
}
