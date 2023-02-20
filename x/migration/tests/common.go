package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Alice      = "self1pvqswgpwjl8273gzqk98ntr8jgvdmvw5cy2asa"
	Bob        = "self1dv35uwyu5h99x80etl6t0nd3q0425jk8ru5fsy"
	Carol      = "self1ruy8kz8tqn9teeg6zqj9fg3e2w6xckxrt6spx8"
	AclAdmin   = "self1dexv46w9kzqjjr73cpannt74rxq63gtxsrszd8"
	Migrator_1 = "self18alpt2t8vw8sxuf5pe0lz3ktmtc4c2wn47r7nf"
	Migrator_2 = "self1wty56897zvnva2lcz30sp368eg8zj4z506g50r"
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
