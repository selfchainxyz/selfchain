package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateWallet{}, "keyless/CreateWallet", nil)
	cdc.RegisterConcrete(&MsgRecoverWallet{}, "keyless/RecoverWallet", nil)
	cdc.RegisterConcrete(&MsgSignTransaction{}, "keyless/SignTransaction", nil)
	cdc.RegisterConcrete(&MsgBatchSign{}, "keyless/BatchSign", nil)
	cdc.RegisterConcrete(&MsgInitiateKeyRotation{}, "keyless/InitiateKeyRotation", nil)
	cdc.RegisterConcrete(&MsgCompleteKeyRotation{}, "keyless/CompleteKeyRotation", nil)
	cdc.RegisterConcrete(&MsgCancelKeyRotation{}, "keyless/CancelKeyRotation", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateWallet{},
		&MsgRecoverWallet{},
		&MsgSignTransaction{},
		&MsgBatchSign{},
		&MsgInitiateKeyRotation{},
		&MsgCompleteKeyRotation{},
		&MsgCancelKeyRotation{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
