package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

//querierはapplicationの状態を確認するためのメソッドを定義する


const (
	QueryResolve = "resolve"
	QueryWhois = "whois"
)

//NewQuerierは、querierのroutingを行う
//queryには、Msgのようなinterfaceはないので、自分で場合分けする
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0]{
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, keeper)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper)(res []byte, err sdk.Error){
	name := path[0]

	value := keeper.ResolveName(ctx, name)

	if value == "" {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve name")
	}

	return []byte(value), nil
}

type Whois struct {
	Value string `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins `json:"price"`
}

func queryWhois(ctx sdk.Context, path[]string, req abci.RequestQuery, keeper Keeper)(res []byte, err sdk.Error) {
	name := path[0]

	whois := Whois{}

	whois.Value = keeper.ResolveName(ctx, name)
	whois.Owner = keeper.GetOwner(ctx, name)
	whois.Price = keeper.GetPrice(ctx, name)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, whois)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}