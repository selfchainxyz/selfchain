package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Alice = "front1pvqswgpwjl8273gzqk98ntr8jgvdmvw5cy2asa"
	Bob   = "front1dv35uwyu5h99x80etl6t0nd3q0425jk8ru5fsy"
	Carol = "front1ruy8kz8tqn9teeg6zqj9fg3e2w6xckxrt6spx8"
	Migrator_1 = "front18alpt2t8vw8sxuf5pe0lz3ktmtc4c2wn47r7nf"
	Migrator_2 = "front1wty56897zvnva2lcz30sp368eg8zj4z506g50r"
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
