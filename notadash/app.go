package main

import (
    "github.com/codegangsta/cli"
)


func buildApp() *cli.App {
    app := cli.NewApp()
    app.Name = "notadash-mon"
    app.Usage = "Monitoring utility for the Mesos/Marathon/Docker stack --> decidedly not-a-dash"
    app.EnableBashCompletion = true
    app.Version = "0.1.0"

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "c, config",
            Usage: "Specify a config file (default: ~/.notadash.gcfg)",
            Value: filepath.Join(os.Getenv("HOME"), ".notadash.gcfg"),
            EnvVar: "NOTADASH_CONFIG",
        },
    }
