package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"encoding/json"
)

/*
Msgsは、stateの変更のtransactionを実行する時に渡すパラメータで、Txsにラップされてnetworkにbroadcastされる
cosmos-sdkを使うことで、Msgsをラップしたり、Txsからのアンラップしたりしてくれるので、実質的にMsgsの定義だけすれば良い
Msgsは、次のinterfaceを満たすように実装すれば良い。

type Msg interface {
	Type() string
	Route() string
	ValidateBasic() Error
	GetSignBytes() []byte
	GetSigners() []AccAddress
}
*/

//SetNameに関するMsg interfaceの実装
type MsgSetName struct {
	Name string
	Value string
	Owner sdk.AccAddress
}

func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name: name,
		Value: value,
		Owner: owner,
	}
}

//module名を定義する
func (msg MsgSetName) Route() string {
	return "nameservice"
}

//action名を決める
func (msg MsgSetName) Type() string {
	return "set_name"
}

//Msgsの中身のチェックをする
func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty(){
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Value cannot be empty")
	}
	return nil
}

//署名するためのMsgデータを取得する
func (msg MsgSetName) GetSignBytes()[]byte{
	b, err := json.Marshal(msg)
	if err != nil{
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//誰の署名が必要か定義する
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//BuyNameに関するMsg interfaceの実装
type MsgBuyName struct {
	Name string
	Bid sdk.Coins
	Buyer sdk.AccAddress
}

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:name,
		Bid:bid,
		Buyer:buyer,
	}
}

func (msg MsgBuyName) Route() string{
	return "nameservice"
}

func (msg MsgBuyName) Type() string {
	return "buy_name"
}

func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty(){
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

func (msg MsgBuyName) GetSignBytes() []byte{
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}