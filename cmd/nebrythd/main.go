package main

import (
    "os"

    "github.com/cosmos/cosmos-sdk/server"
    svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
    "github.com/cosmos/cosmos-sdk/types/module"

    "github.com/warwarrior502/nebryth/app"
)

func main() {
    // Create the root command for nebrythd
    rootCmd, _ := svrcmd.NewRootCmd(
        app.AppName,
        app.MakeEncodingConfig,
        app.NewNebrythApp,
    )

    // Add module init flags (optional, but useful for CLI)
    module.AddModuleInitFlags(rootCmd)

    // Execute the command
    if err := svrcmd.Execute(rootCmd, app.AppName, os.ExpandEnv("$HOME/.nebrythd")); err != nil {
        os.Exit(1)
    }
}

