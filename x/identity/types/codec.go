package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// OAuth messages
	cdc.RegisterConcrete(&MsgLinkSocialIdentity{}, "identity/LinkSocialIdentity", nil)
	cdc.RegisterConcrete(&MsgUnlinkSocialIdentity{}, "identity/UnlinkSocialIdentity", nil)
	cdc.RegisterConcrete(&MsgVerifyOAuthToken{}, "identity/VerifyOAuthToken", nil)
	
	// MFA messages
	cdc.RegisterConcrete(&MsgAddMFA{}, "identity/AddMFA", nil)
	cdc.RegisterConcrete(&MsgRemoveMFA{}, "identity/RemoveMFA", nil)
	cdc.RegisterConcrete(&MsgVerifyMFA{}, "identity/VerifyMFA", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		// OAuth messages
		&MsgLinkSocialIdentity{},
		&MsgUnlinkSocialIdentity{},
		&MsgVerifyOAuthToken{},
		
		// MFA messages
		&MsgAddMFA{},
		&MsgRemoveMFA{},
		&MsgVerifyMFA{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
