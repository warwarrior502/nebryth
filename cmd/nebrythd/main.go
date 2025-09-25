package main

import (
	"io"
	"os"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/spf13/cobra"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	server "github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/client/keys"

	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/nebryth/nebrythd/app"
)

func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	return app.NewNebrythApp(logger, db)
}

func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	return servertypes.ExportedApp{}, nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "nebrythd",
		Short: "Nebryth Daemon (Cosmos SDK v0.47.5) — with staking support",
	}

	// Keys
	rootCmd.AddCommand(keys.Commands(app.DefaultNodeHome))

	// Genesis utilities
	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.GenTxCmd(app.ModuleBasics, app.TxEncodingConfig(), nil, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(nil, app.DefaultNodeHome, nil),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		genutilcli.AddGenesisAccountCmd(app.DefaultNodeHome),
	)

	// Query commands
	rootCmd.AddCommand(
		authcli.GetQueryCmd(),
		bankcli.GetQueryCmd(),
	)

	// Server commands
	rootCmd.AddCommand(
		server.StartCmd(newApp, app.DefaultNodeHome),
		server.ExportCmd(appExport, app.DefaultNodeHome),
	)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome, app.AppName); err != nil {
		os.Exit(1)
	}
}