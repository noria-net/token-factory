package bindings

import (
	"encoding/json"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	bindingstypes "github.com/noria-net/token-factory/x/tokenfactory/bindings/types"
	tokenfactorykeeper "github.com/noria-net/token-factory/x/tokenfactory/keeper"
)

type CustomQueryHandler struct {
	wrapped      wasmkeeper.WasmVMQueryHandler
	tokenfactory *tokenfactorykeeper.Keeper
	bankkeeper   *bankkeeper.BaseKeeper
}

func CustomQueryDecorator(tokenFactory *tokenfactorykeeper.Keeper, bankkeeper *bankkeeper.BaseKeeper) func(wasmkeeper.WasmVMQueryHandler) wasmkeeper.WasmVMQueryHandler {
	return func(old wasmkeeper.WasmVMQueryHandler) wasmkeeper.WasmVMQueryHandler {
		return &CustomQueryHandler{
			wrapped:      old,
			tokenfactory: tokenFactory,
			bankkeeper:   bankkeeper,
		}
	}
}

// CustomQuerier dispatches custom CosmWasm bindings queries.
func (m *CustomQueryHandler) HandleQuery(ctx sdk.Context, caller sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error) {
	if request.Custom == nil {
		return m.wrapped.HandleQuery(ctx, caller, request)
	}
	customQuery := request.Custom

	var tokenQuery bindingstypes.TokenFactoryQuery
	if err := json.Unmarshal(customQuery, &tokenQuery); err != nil {
		return nil, sdkerrors.Wrap(err, "requires 'token' field")
	}
	if tokenQuery.Token == nil {
		return m.wrapped.HandleQuery(ctx, caller, request)
	}

	switch {
	case tokenQuery.Token.FullDenom != nil:
		creator := tokenQuery.Token.FullDenom.CreatorAddr
		subdenom := tokenQuery.Token.FullDenom.Subdenom

		fullDenom, err := GetFullDenom(creator, subdenom)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "osmo full denom query")
		}

		res := bindingstypes.FullDenomResponse{
			Denom: fullDenom,
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "failed to marshal FullDenomResponse")
		}

		return bz, nil

	case tokenQuery.Token.Admin != nil:
		res, err := m.GetDenomAdmin(ctx, tokenQuery.Token.Admin.Denom)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal AdminResponse: %w", err)
		}

		return bz, nil

	case tokenQuery.Token.Metadata != nil:
		res, err := m.GetMetadata(ctx, tokenQuery.Token.Metadata.Denom)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal MetadataResponse: %w", err)
		}

		return bz, nil

	case tokenQuery.Token.DenomsByCreator != nil:
		res, err := m.GetDenomsByCreator(ctx, tokenQuery.Token.DenomsByCreator.Creator)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal DenomsByCreatorResponse: %w", err)
		}

		return bz, nil

	case tokenQuery.Token.Params != nil:
		res, err := m.GetParams(ctx)
		if err != nil {
			return nil, err
		}

		bz, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("failed to JSON marshal ParamsResponse: %w", err)
		}

		return bz, nil

	default:
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown token query variant"}
	}
}

// ConvertSdkCoinsToWasmCoins converts sdk type coins to wasm vm type coins
func ConvertSdkCoinsToWasmCoins(coins []sdk.Coin) wasmvmtypes.Coins {
	var toSend wasmvmtypes.Coins
	for _, coin := range coins {
		c := ConvertSdkCoinToWasmCoin(coin)
		toSend = append(toSend, c)
	}
	return toSend
}

// ConvertSdkCoinToWasmCoin converts a sdk type coin to a wasm vm type coin
func ConvertSdkCoinToWasmCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom: coin.Denom,
		// Note: gamm tokens have 18 decimal places, so 10^22 is common, no longer in u64 range
		Amount: coin.Amount.String(),
	}
}
