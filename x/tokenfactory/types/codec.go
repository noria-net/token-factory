package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgTokenFactoryCreateDenom{}, "osmosis/tokenfactory/create-denom", nil)
	cdc.RegisterConcrete(&MsgTokenFactoryMint{}, "osmosis/tokenfactory/mint", nil)
	cdc.RegisterConcrete(&MsgTokenFactoryBurn{}, "osmosis/tokenfactory/burn", nil)
	cdc.RegisterConcrete(&MsgTokenFactoryForceTransfer{}, "osmosis/tokenfactory/force-transfer", nil)
	cdc.RegisterConcrete(&MsgTokenFactoryChangeAdmin{}, "osmosis/tokenfactory/change-admin", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgTokenFactoryCreateDenom{},
		&MsgTokenFactoryMint{},
		&MsgTokenFactoryBurn{},
		&MsgTokenFactoryForceTransfer{},
		&MsgTokenFactoryChangeAdmin{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
