package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

//Keeper は、storageを管理するstructで、getter/setterのメソッドを実装する
type Keeper struct {
	coinKeeper bank.Keeper //coinの送受信等に必要なモジュール

	namesStoreKey  sdk.StoreKey //sdk.Contextからname store(sdk.KVStore)にアクセスするためのkey
	ownersStoreKey sdk.StoreKey //sdk.Contextからowners store(sdk.KVStore)にアクセスするためのkey
	pricesStoreKey sdk.StoreKey //sdk.Contextからprices store(sdk.KVStore)にアクセスするためのkey

	cdc *codec.Codec //codecはgo-aminoを用いて、byte codeをdecode/encodeしてtendermint側と通信するときに必要となる
}

//NewKeeper 新しいKeeperを生成する
func NewKeeper(coinKeeper bank.Keeper, namesStoreKey sdk.StoreKey, ownersStoreKey sdk.StoreKey, priceStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:     coinKeeper,
		namesStoreKey:  namesStoreKey,
		ownersStoreKey: ownersStoreKey,
		pricesStoreKey: priceStoreKey,
		cdc:            cdc,
	}
}

//ResolveName 名前解決する時に表示する値を取得する
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	store := ctx.KVStore(k.namesStoreKey)
	bz := store.Get([]byte(name))
	return string(bz)
}

//SetName 名前解決する時に表示される値を設定する
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	store := ctx.KVStore(k.namesStoreKey)
	store.Set([]byte(name), []byte(value))
}

//HasOwner nameのownerが設定されているか取得する
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.ownersStoreKey)
	bz := store.Get([]byte(name))
	return bz != nil
}

//GetOwner nameのownerを取得する
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	store := ctx.KVStore(k.ownersStoreKey)
	bz := store.Get([]byte(name))
	return bz
}

//SetOwner nameのownerを設定する
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.ownersStoreKey)
	store.Set([]byte(name), owner)
}

//GetPrice nameの現在価格を取得する
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	//ownerがいない場合は、1mycoinを価格として返す
	if !k.HasOwner(ctx, name) {
		return sdk.Coins{sdk.NewInt64Coin("mycoin", 1)}
	}
	store := ctx.KVStore(k.pricesStoreKey)
	bz := store.Get([]byte(name))
	var price sdk.Coins
	k.cdc.MustUnmarshalBinaryBare(bz, &price)
	return price
}

//SetPrice nameの価格を設定する
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	store := ctx.KVStore(k.pricesStoreKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(price))
}
