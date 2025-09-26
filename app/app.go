package app

import (
    upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
    upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

    "github.com/cosmos/cosmos-sdk/baseapp"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/codec"
    codectypes "github.com/cosmos/cosmos-sdk/codec/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// App defines the Nebryth application structure.
type App struct {
    *baseapp.BaseApp

    UpgradeKeeper upgradekeeper.Keeper

    InterfaceRegistry codectypes.InterfaceRegistry
    LegacyAmino       *codec.LegacyAmino
    AppCodec          codec.Codec
}

// NewApp creates and configures a new Nebryth App instance.
func NewApp(
    logger sdk.Logger,
    db sdk.DB,
    traceStore sdk.TraceStore,
    loadLatest bool,
    skipUpgradeHeights map[int64]bool,
    homePath string,
    invCheckPeriod uint,
    encodingConfig client.EncodingConfig,
    appOpts servertypes.AppOptions,
) *App {
    app := &App{
        BaseApp: baseapp.NewBaseApp("nebryth", logger, db, traceStore),
    }

    app.InterfaceRegistry = encodingConfig.InterfaceRegistry
    app.LegacyAmino = encodingConfig.Amino
    app.AppCodec = encodingConfig.Marshaler

    app.UpgradeKeeper.SetUpgradeHandler("v0.2.0", func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
        return vm, nil
    })

    return app
}


