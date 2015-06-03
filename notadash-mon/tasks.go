package main

import (
	"fmt"
	lib "github.com/boldfield/notadash/lib"
	"github.com/codegangsta/cli"
	"github.com/ryanuber/columnize"
)

func runCheckTasks(ctx *cli.Context) int {
	fmt.Println("Discoving running applications and associated tasks...")

	marathon, err := loadMarathon(ctx.GlobalString("marathon-host"))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	mesos := &lib.Mesos{
		Host: ctx.GlobalString("mesos-host"),
	}
	mesosClient := mesos.Client()
	mesos.LoadCluster(mesosClient)

	mesosFrameworks := mesos.Framework("marathon")
	marathonApps := &lib.MarathonApps{}

	marathonClient := marathon.Client()

	if len(mesosFrameworks) > 0 {
		for _, a := range marathon.Apps {
			if len(a.DeploymentID) > 0 && ctx.GlobalBool("ignore-deploys") {
				continue
			}

			if tasks, err := marathonClient.Tasks(a.ID); err != nil {
				fmt.Println(err)
				return 1
			} else {
				for _, t := range tasks.Tasks {
					taskSlave := mesos.Cluster.GetSlaveByHostName(t.Host)
					marathonApps.AddTask(t.ID, t.AppID, taskSlave.Id, taskSlave.HostName, false, true, false)
				}
			}
			for _, f := range mesosFrameworks {
				for _, t := range f.Tasks {
					taskSlave := mesos.Cluster.GetSlaveById(t.SlaveId)
					marathonApps.AddTask(t.Id, t.AppId(), taskSlave.Id, taskSlave.HostName, true, false, false)
				}
			}
		}
	}

	output := make([]string, 1)
	output[0] = "Application | Task ID | Slave Host | Mesos/Marathon"
	discrepancy := false

	for _, a := range marathonApps.Apps {
		app_discrepancy := false
		app_output := make([]string, 1)
		app_output[0] = fmt.Sprintf("%s| | | ", a.Id)

		for _, t := range a.Tasks {
			if !(t.Mesos && t.Marathon) {
				app_discrepancy = true
				ln := fmt.Sprintf(
					" | %s | %s | %s/%s",
					t.Id,
					t.SlaveHost,
					lib.PrintBool(t.Mesos),
					lib.PrintBool(t.Marathon),
				)
				app_output = append(app_output, ln)
			}
		}
		if app_discrepancy {
			discrepancy = true
			output = append(output, app_output...)
		}
	}

	if discrepancy {
		fmt.Println(lib.PrintYellow("Discrepency in task state found!"))
		result := columnize.SimpleFormat(output)
		fmt.Println(result)
		return 2
	}

	fmt.Println(lib.PrintGreen("Mesos and Marathon agree about running tasks!"))
	return 0
}
