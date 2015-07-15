package main

import (
	"fmt"
	lib "github.com/socrata-platform/notadash/lib"
	"github.com/codegangsta/cli"
	"github.com/ryanuber/columnize"
	"time"
)

func runFailedTasks(ctx *cli.Context) int {
	fmt.Println("Discovering completed mesos tasks")

	mesos := &lib.Mesos{
		Host: ctx.GlobalString("mesos-host"),
	}
	mesosClient := mesos.Client()
	if err := mesos.LoadClusterInfo(mesosClient); err != nil {
		fmt.Println(err)
		return 1
	}

	if ctx.GlobalBool("only-leader") {
		if err := mesos.ErrIfNotLeader(); err != nil {
			fmt.Println(err)
			return 0
		}
	}

	frameworks := mesos.Framework("marathon")

	timeWindow := ctx.Int("time-window")
	failureLimit := ctx.Int("failure-limit")

	// map of all containers to number of failures in the last 30 min
	failingMap := make(map[string]int)
	now := float64(time.Now().Unix())
	if len(frameworks) > 0 {
		for _, f := range frameworks {
			for _, t := range f.CompletedTasks {
				failure := false
				for _, s := range t.Statuses {
					if s.State == "TASK_FAILED" {
						age := now - s.Timestamp
						failure = age < float64(60*timeWindow) // minutes
					}
				}
				if failure {
					failingMap[t.Name]++
				}
			}
		}
	}

	if len(failingMap) == 0 {
		fmt.Println(lib.PrintGreen(fmt.Sprintf("No failures found in the last %d minutes", timeWindow)))
		return 0
	}

	returnError := false
	output := make([]string, 1)
	output[0] = "Application | Num Failures"
	for app, numFailures := range failingMap {
		line := fmt.Sprintf("%s | %d", app, numFailures)
		output = append(output, line)
		if numFailures > failureLimit {
			returnError = true
		}
	}

	result := columnize.SimpleFormat(output)
	fmt.Println(result)

	if returnError {
		fmt.Println(lib.PrintRed(fmt.Sprintf("Found errors above the failure limit of %d", failureLimit)))
		return 2
	}

	fmt.Println(lib.PrintYellow(fmt.Sprintf("Errors found, but not above the failure limit of %d", failureLimit)))
	return 0
}
