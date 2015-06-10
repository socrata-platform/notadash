package main

import (
	"fmt"
	lib "github.com/boldfield/notadash/lib"
	"github.com/codegangsta/cli"
	"github.com/ryanuber/columnize"
)

func runFindFailures(ctx *cli.Context) int {
	fmt.Println("Checking for tasks that have 0 running instances...")

	marathon, err := loadMarathon(ctx.GlobalString("marathon-host"))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	output := make([]string, 1)
	output[0] = "Application | Instances | Running Tasks"
	discrepancy := false
	for _, a := range marathon.Apps {
		if a.Instances > 0 && a.TasksRunning == 0 {
			discrepancy = true
			line := fmt.Sprintf(
				"%s | %d | %d",
				a.ID,
				a.Instances,
				a.TasksRunning,
			)
			output = append(output, line)
		}
	}
	if discrepancy {
		fmt.Println(lib.PrintYellow("Failing Marathon task found!"))
		result := columnize.SimpleFormat(output)
		fmt.Println(result)
		return 2
	}
	return 0
}
