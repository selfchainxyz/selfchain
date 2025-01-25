package types

import (
	"fmt"
)

// ErrInvalidWalletId defines the error for an invalid wallet ID
var ErrInvalidWalletId = fmt.Errorf("wallet ID cannot be empty")

// ErrInvalidCreator defines the error for an invalid creator
var ErrInvalidCreator = fmt.Errorf("wallet creator cannot be empty")

// ErrInvalidPublicKey defines the error for an invalid public key
var ErrInvalidPublicKey = fmt.Errorf("wallet public key cannot be empty")

// ErrInvalidChainId defines the error for an invalid chain ID
var ErrInvalidChainId = fmt.Errorf("chain ID cannot be empty")

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
		return ErrInvalidCreator
	}

	if w.PublicKey == "" {
		return ErrInvalidPublicKey
	}

	if w.ChainId == "" {
		return ErrInvalidChainId
	}

	return nil
}
