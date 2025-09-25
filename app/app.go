package app

import (
	"os"
	"path/filepath"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	// Server API/Config for interface compliance
	serverapi "github.com/cosmos/cosmos-sdk/server/api"
	servercfg "github.com/cosmos/cosmos-sdk/server/config"

	// Auth / Bank / Staking / Genutil
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/x/genutil"
)

const AppName = "nebryth"

var DefaultNodeHome = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".nebrythd"
	}
	return filepath.Join(home, ".nebrythd")
}()

// ModuleBasics defines the basic modules for genesis/init
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	genutil.AppModuleBasic{},
)

type NebrythApp struct {
	*baseapp.BaseApp
	appCodec codec.Codec
	txCfg    client.TxConfig

	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.BaseKeeper
	StakingKeeper stakingkeeper.Keeper

	ModuleManager *module.Manager
}

var maccPerms = map[string][]string{
	authtypes.FeeCollectorName: nil,
	banktypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
	stakingtypes.BondedPoolName:   {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
}

func NewNebrythApp(logger log.Logger, db dbm.DB) *NebrythApp {
	// Codec & TX config
	ir := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(ir)
	txCfg := tx.NewTxConfig(appCodec, tx.DefaultSignModes)

	// Base app
	bapp := baseapp.NewBaseApp(AppName, logger, db, txCfg.TxDecoder())
	bapp.SetVersion("0.1.0")

	app := &NebrythApp{
		BaseApp:  bapp,
		appCodec: appCodec,
		txCfg:    txCfg,
	}

	// Store keys
	authKey := sdk.NewKVStoreKey(authtypes.StoreKey)
	bankKey := sdk.NewKVStoreKey(banktypes.StoreKey)
	stakingKey := sdk.NewKVStoreKey(stakingtypes.StoreKey)

	// Keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		authKey,
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		AppName,
	)

	blockedAddrs := make(map[string]bool)
	for name := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(name).String()] = true
	}

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		bankKey,
		app.AccountKeeper,
		blockedAddrs,
		AppName,
	)

	app.StakingKeeper = *stakingkeeper.NewKeeper(
		appCodec,
		stakingKey,
		app.AccountKeeper,
		app.BankKeeper,
		AppName,
	)

	// Module manager
	app.ModuleManager = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper, nil),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.DeliverTx, app.txCfg),
	)

	// AnteHandler
	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: txCfg.SignModeHandler(),
			FeegrantKeeper:  nil,
		},
	)
	if err != nil {
		panic(err)
	}
	app.SetAnteHandler(anteHandler)

	return app
}

// ABCI hooks (no-ops for now)
func (app *NebrythApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	return abci.ResponseInitChain{}
}
func (app *NebrythApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}
func (app *NebrythApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

// Server Application interface methods (stubs to satisfy interface)
func (app *NebrythApp) RegisterAPIRoutes(_ *serverapi.Server, _ servercfg.APIConfig) {}
func (app *NebrythApp) RegisterTxService(_ client.Context)                           {}
func (app *NebrythApp) RegisterTendermintService(_ client.Context)                   {}
func (app *NebrythApp) RegisterNodeService(_ client.Context)                         {}

// TxEncodingConfig returns the app's TxEncodingConfig for CLI commands
func TxEncodingConfig() client.TxEncodingConfig {
	ir := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(ir)
	return tx.NewTxConfig(cdc, tx.DefaultSign