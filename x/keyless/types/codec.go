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
	cdc.RegisterConcrete(&MsgGrantPermission{}, "keyless/GrantPermission", nil)
	cdc.RegisterConcrete(&MsgRevokePermission{}, "keyless/RevokePermission", nil)
	cdc.RegisterConcrete(&MsgBatchSignRequest{}, "keyless/BatchSign", nil)
	cdc.RegisterConcrete(&MsgInitiateKeyRotation{}, "keyless/InitiateKeyRotation", nil)
	cdc.RegisterConcrete(&MsgCompleteKeyRotation{}, "keyless/CompleteKeyRotation", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateWallet{},
		&MsgRecoverWallet{},
		&MsgGrantPermission{},
		&MsgRevokePermission{},
		&MsgBatchSignRequest{},
		&MsgInitiateKeyRotation{},
		&MsgCompleteKeyRotation{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
}
