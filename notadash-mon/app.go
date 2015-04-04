package main

import (
    "github.com/codegangsta/cli"
)

var VERSION = "0.1.0-beta"

func buildApp() *cli.App {
    app := cli.NewApp()
    app.Name = "notadash-mon"
    app.Usage = "Monitoring utility for the Mesos/Marathon/Docker stack --> decidedly not-a-dash"
    app.EnableBashCompletion = true
    app.Version = VERSION

    app.Flags = []cli.Flag{
        cli.BoolFlag{
            Name:  "verbose",
            Usage: "Show more output",
        },
        cli.StringFlag{
            Name:  "marathon-host",
            Usage: "URL to use for Marathon cluster discovery.",
            EnvVar: "NOTADASH_MARATHON_URL",
        },
        cli.StringFlag{
            Name:  "mesos-host",
            Usage: "URL to use for Mesos cluster discovery.",
            EnvVar: "NOTADASH_MESOS_URL",
        },
//        cli.StringFlag{
//            Name:  "c, config",
//            Usage: "Specify a config file (default: ~/.notadash.gcfg)",
//            Value: filepath.Join(os.Getenv("HOME"), ".notadash.gcfg"),
//            EnvVar: "NOTADASH_CONFIG",
//        },
    }

    app.Commands = []cli.Command{
        {
            Name: "tasks",
            Usage:  "Cross-check all tasks registered with Mesos and Marathon.",
            Action: checkTasks,
        },
        {
            Name: "slave",
            Flags: []cli.Flag{
                cli.BoolFlag{
                    Name:  "docker-verify",
                    Usage: "Verify running containers match expectations. Assumes command is being run on a slave",
                },
            },
            Usage:  "Verify all tasks registered for mesos slave are running as expected. Must be run on target mesos slave.",
            Action: checkSlave,
        },
    }


    return app
}
