package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

var VERSION = "0.1.0-beta"

type boolmap map[string]bool

var ctRequired = []string{
	"marathon-host",
	"mesos-host",
}

var csRequired = []string{
	"marathon-host",
	"mesos-host",
}

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
			Name:   "marathon-host",
			Usage:  "URL to use for Marathon cluster discovery.",
			EnvVar: "NOTADASH_MARATHON_URL",
		},
		cli.StringFlag{
			Name:   "mesos-host",
			Usage:  "URL to use for Mesos cluster discovery.",
			EnvVar: "NOTADASH_MESOS_URL",
		},
		cli.BoolFlag{
			Name:  "ignore-deploys",
			Usage: "Ignore active deployments when checking consensus",
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
			Name:   "resources",
			Usage:  "Show resource allocation per node across the cluster.",
			Action: showAllocation,
		},
		{
			Name:   "tasks",
			Usage:  "Cross-check all tasks registered with Mesos and Marathon.",
			Action: checkTasks,
		},
		{
			Name: "slave",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "kill-stragglers",
					Usage: "Kill containers which are still running and registered with Mesos, but Marathon has inconveniently forgotten.",
				},
			},
			Usage:  "Verify all tasks registered for mesos slave are running as expected. Must be run on target mesos slave.",
			Action: checkSlave,
		},
	}

	return app
}

func showAllocation(ctx *cli.Context) {
	if missing, err := validateContext(ctx, ctRequired); err != nil {
		fmt.Println(err)
		fmt.Printf("The following parameters must be defined: %s\n", missing)
		os.Exit(2)
	}

	exitStatus := runShowAllocation(ctx)
	os.Exit(exitStatus)
}

func checkTasks(ctx *cli.Context) {
	if missing, err := validateContext(ctx, ctRequired); err != nil {
		fmt.Println(err)
		fmt.Printf("The following parameters must be defined: %s\n", missing)
		os.Exit(2)
	}

	exitStatus := runCheckTasks(ctx)
	os.Exit(exitStatus)
}

func checkSlave(ctx *cli.Context) {
	if missing, err := validateContext(ctx, csRequired); err != nil {
		fmt.Println(err)
		fmt.Printf("The following parameters must be defined: %s\n", missing)
		os.Exit(1)
	}

	exitStatus := runCheckSlave(ctx)
	os.Exit(exitStatus)
}
