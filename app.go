package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	dbm "github.com/tendermint/tendermint/libs/db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/EG-easy/sample-cosmos-app/x/namespace"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const (
	appName = "nameservice"
)

type nameServiceApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey
	keyNSnames *sdk.KVStoreKey
	keyNSowners *sdk.KVStoreKey
	keyNSprices *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams *sdk.KVStoreKey
	tkeyParams *sdk.TransientStoreKey

	accountKeeper auth.AccountKeeper
	bankKeeper bank.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	paramsKeeper params.Keeper
	nsKeeper nameservice.Keeper
}

func NewNameServiceApp(logger log.Logger, db dbm.DB) *nameServiceApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &nameServiceApp{
		BaseApp: bApp,
		cdc: cdc,

		keyMain:sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyNSnames:       sdk.NewKVStoreKey("ns_names"),
		keyNSowners:      sdk.NewKVStoreKey("ns_owners"),
		keyNSprices:      sdk.NewKVStoreKey("ns_prices"),
		keyFeeCollection: sdk.NewKVStoreKey("fee_collection"),
		keyParams:        sdk.NewKVStoreKey("params"),
		tkeyParams:       sdk.NewTransientStoreKey("transient_params"),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)

	app.nsKeeper = nameservice.NewKeeper(
		app.bankKeeper,
		app.keyNSnames,
		app.keyNSowners,
		app.keyNSprices,
		app.cdc,
	)

	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("nameservice", nameservice.NewHandler(app.nsKeeper))

	app.SetInitChainer(app.initChainer)

	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyNSnames,
		app.keyNSowners,
		app.keyNSprices,
		app.keyFeeCollection,
		app.keyParams,
		app.tkeyParams,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}


type GenesisState struct {
	AuthData auth.GenesisState `json:"auth"`
	BankData bank.GenesisState `json:"bank"`
	Accounts []*auth.BaseAccount `json:"accounts"`
}

func (app *nameServiceApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)

	return abci.ResponseInitChain{}
}

func (app *nameServiceApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error){
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*auth.BaseAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address:acc.GetAddress(),
			Coins: acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := GenesisState{
		Accounts:accounts,
		AuthData: auth.DefaultGenesisState(),
		BankData: bank.DefaultGenesisState(),
	}

	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	nameservice.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}