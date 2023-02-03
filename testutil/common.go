package commontest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitSDKConfig() {
	accountPubKeyPrefix := "frontpub"
	validatorAddressPrefix := "frontvaloper"
	validatorPubKeyPrefix := "frontvaloperpub"
	consNodeAddressPrefix := "frontvalcons"
	consNodePubKeyPrefix := "frontvalconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("front", accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}
