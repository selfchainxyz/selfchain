package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgMigrate{}, "migration/Migrate", nil)
	cdc.RegisterConcrete(&MsgAddMigrator{}, "migration/AddMigrator", nil)
	cdc.RegisterConcrete(&MsgRemoveMigrator{}, "migration/RemoveMigrator", nil)
	cdc.RegisterConcrete(&MsgUpdateConfig{}, "migration/UpdateConfig", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMigrate{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddMigrator{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveMigrator{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateConfig{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
