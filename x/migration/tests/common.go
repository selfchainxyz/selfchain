package test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Alice = "front1wnr2vnx8tjxv8mkp20dagdvn8pylr7x7l0mt9s"
	Bob   = "front1kphjpt7crfmlhjyzvf44ng3ye9rvfy7z3pyajj"
	Carol = "front1skzx07z9nhnv5js3x5ucchey7yyc6enk2l3sqx"
	Migrator_1 = "front16077t285ajgp5a9retp4fkhs7enusr2n3lp4jg"
	Migrator_2 = "front1wnr2vnx8tjxv8mkp20dagdvn8pylr7x7l0mt9s"
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
