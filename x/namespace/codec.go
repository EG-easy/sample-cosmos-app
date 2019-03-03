package nameservice

import "github.com/cosmos/cosmos-sdk/codec"

//codecでencode/decodeするためにtypesを登録する
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "nameservice/BuyName", nil)
}