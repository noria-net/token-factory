package bindings

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	bindingstypes "github.com/noria-net/token-factory/x/tokenfactory/bindings/types"
	tokenfactorykeeper "github.com/noria-net/token-factory/x/tokenfactory/keeper"
)

type QueryPlugin struct {
	bankKeeper         *bankkeeper.BaseKeeper
	tokenFactoryKeeper *tokenfactorykeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(b *bankkeeper.BaseKeeper, tfk *tokenfactorykeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		bankKeeper:         b,
		tokenFactoryKeeper: tfk,
	}
}

// GetDenomAdmin is a query to get denom admin.
func (qp CustomQueryHandler) GetDenomAdmin(ctx sdk.Context, denom string) (*bindingstypes.AdminResponse, error) {
	metadata, err := qp.tokenfactory.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin for denom: %s", denom)
	}
	return &bindingstypes.AdminResponse{Admin: metadata.Admin}, nil
}

func (qp CustomQueryHandler) GetDenomsByCreator(ctx sdk.Context, creator string) (*bindingstypes.DenomsByCreatorResponse, error) {
	// TODO: validate creator address
	denoms := qp.tokenfactory.GetDenomsFromCreator(ctx, creator)
	return &bindingstypes.DenomsByCreatorResponse{Denoms: denoms}, nil
}

func (qp CustomQueryHandler) GetMetadata(ctx sdk.Context, denom string) (*bindingstypes.MetadataResponse, error) {
	metadata, found := qp.bankkeeper.GetDenomMetaData(ctx, denom)
	var parsed *bindingstypes.Metadata
	if found {
		parsed = SdkMetadataToWasm(metadata)
	}
	return &bindingstypes.MetadataResponse{Metadata: parsed}, nil
}

func (qp CustomQueryHandler) GetParams(ctx sdk.Context) (*bindingstypes.ParamsResponse, error) {
	params := qp.tokenfactory.GetParams(ctx)
	return &bindingstypes.ParamsResponse{
		Params: bindingstypes.Params{
			DenomCreationFee: ConvertSdkCoinsToWasmCoins(params.DenomCreationFee),
		},
	}, nil
}
