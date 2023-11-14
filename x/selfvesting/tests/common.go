package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Alice = "self16uhrfxdj8drxvlal8q6xhy99jh6nt9uqur6pl2"
	Bob   = "self184l076mmmeumx2cw6eqdct7tchufux0v43cxcg"
	Carol = "self19qrygeag248redm708jtt4ks94sp8p34v7l66k"
)

func InitSDKConfig() {
	accountPubKeyPrefix := "selfpub"
	validatorAddressPrefix := "selfvaloper"
	validatorPubKeyPrefix := "selfvaloperpub"
	consNodeAddressPrefix := "selfvalcons"
	consNodePubKeyPrefix := "selfvalconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("self", accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}
