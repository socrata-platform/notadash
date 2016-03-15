package main

import (
	"bytes"
	"fmt"
	chronos "github.com/behance/go-chronos/chronos"
	lib "github.com/socrata-platform/notadash/lib"
	"github.com/codegangsta/cli"
	"github.com/ryanuber/columnize"
	"os/exec"
)

func runCheckSlave(ctx *cli.Context) int {
	fmt.Println("Discovering running applications and associated tasks...")

	marathon, err := loadMarathon(ctx.GlobalString("marathon-host"))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	slave, err := loadMesos(ctx.GlobalString("marathon-host"))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	slaveFrameworks := slave.Framework("marathon")

	dockerClient := lib.NewDockerClient()
	marathonApps, ignoredImages, err := buildMesosMarathonMatrix(slave.Slave.Id, slave.Slave.HostName, slaveFrameworks, marathon, dockerClient, ctx.GlobalBool("ignore-deploys"))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	discrepancy, containerAccount, output, err := verifyApplications(marathonApps)
	if err != nil {
		fmt.Println(lib.PrintRed("An error occoured while verifying applications!"))
		fmt.Println(err)
		return 1
	}

	orphanedContainers := make(boolmap)
	chronosHost := ctx.GlobalString("chronos-host")
	if chronosHost != "" {
		config := chronos.Config{
			URL: chronosHost,
		}
		client, err := chronos.NewClient(config)
		jobs, err := client.Jobs()
		if err == nil {
			for _, job := range *jobs {
				if job.Container != nil {
					ignoredImages = append(ignoredImages, job.Container.Image)
				}
			}
		}
	}
	containers, err := lib.ListRunningContainers(dockerClient, ignoredImages)
	if err != nil {
		fmt.Println(lib.PrintRed("An error occoured while determining if docker container is running!"))
		fmt.Println(err)
		return 1
	}

	for _, container := range containers {
		if !containerAccount[container] {
			if ctx.Bool("kill-stragglers") {
				if err := lib.StopContainer(container, 300, dockerClient); err != nil {
					fmt.Printf("An error occoured while trying to stop container (%s): %s\n", container, err)
					orphanedContainers[container] = true
				}
			} else {
				orphanedContainers[container] = true
			}
		}
	}

	if discrepancy || len(orphanedContainers) > 0 {
		if discrepancy {
			fmt.Println(lib.PrintYellow("Discrepency in task state found!"))
			result := columnize.SimpleFormat(output)
			fmt.Println(result)
		}
		if len(orphanedContainers) > 0 {
			fmt.Println(lib.PrintYellow("Orphaned docker containers found!"))
			tmp_output := []string{
				"Orphaned Docker Containers | ",
			}
			for c := range orphanedContainers {
				tmp_output = append(tmp_output, fmt.Sprintf(" | %s", lib.PrintRed(c)))
			}
			result := columnize.SimpleFormat(tmp_output)
			fmt.Println(result)
		}
		return 2
	}
	fmt.Println(lib.PrintGreen("Mesos and Marathon agree about running tasks!"))
	return 0
}

func loadMesos(mesosHost string) (*lib.MesosSlave, error) {
	mesos := &lib.Mesos{
		Host: mesosHost,
	}
	mesosClient := mesos.Client()

	cmd := exec.Command("/bin/hostname", "-f")
	host, err := cmd.Output()
	if err != nil {
		fmt.Println("Unable to determin fully qualified name of host...")
		return nil, err
	}

	host = bytes.Trim(host, " \n\t")
	slave, err := mesos.LoadSlaveState(string(host), mesosClient)
	if err != nil {
		return nil, err
	}
	return slave, nil
}

func buildMesosMarathonMatrix(slaveId, slaveHostName string, slaveFrameworks lib.FrameworkMap, marathon *lib.Marathon, dockerClient lib.DockerClient, ignoreDeploys bool) (*lib.MarathonApps, []string, error) {
	ignoredImages := make([]string, 0)
	marathonApps := &lib.MarathonApps{}
	if len(slaveFrameworks) > 0 {
		for _, a := range marathon.Apps {
			if len(a.Deployments) > 0 && ignoreDeploys {
				ignoredImages = append(ignoredImages, a.Container.Docker.Image)
				continue
			}

			if tasks, err := marathon.Client().Tasks(a.ID); err != nil {
				return nil, nil, err
			} else {
				for _, t := range tasks.Tasks {
					if slaveHostName == t.Host {
						marathonApps.AddTask(t.ID, t.AppID, slaveId, slaveHostName, false, true, false)
					}
				}
			}
			for _, f := range slaveFrameworks {
				for _, e := range f.Executors {
					for _, t := range e.Tasks {
						containerRunning, err := lib.ContainerRunning(e.RegisteredContainerName(t), dockerClient)
						if err != nil {
							return nil, nil, err
						}
						mTask := marathonApps.AddTask(t.Id, t.AppId(), slaveId, slaveHostName, true, false, containerRunning)
						mTask.Container = e.RegisteredContainerName(t)
					}
				}
			}
		}
	}
	return marathonApps, ignoredImages, nil
}

func verifyApplications(marathonApps *lib.MarathonApps) (bool, boolmap, []string, error) {
	containerAccount := make(boolmap)
	output := make([]string, 1)
	output[0] = "Application | Task ID | Slave Host | Mesos/Marathon/Docker"
	discrepancy := false
	for _, a := range marathonApps.Apps {
		app_discrepancy := false
		app_output := make([]string, 1)
		app_output[0] = fmt.Sprintf("%s| | | ", a.Id)
		for _, t := range a.Tasks {
			containerAccount[t.Container] = true
			if !(t.Mesos && t.Marathon) {
				app_discrepancy = true
				ln := fmt.Sprintf(
					" | %s | %s | %s/%s/%s",
					t.Id,
					t.SlaveHost,
					lib.PrintBool(t.Mesos),
					lib.PrintBool(t.Marathon),
					lib.PrintBool(t.Docker),
				)
				app_output = append(app_output, ln)
			}
		}
		if app_discrepancy {
			discrepancy = true
			output = append(output, app_output...)
		}
	}
	return discrepancy, containerAccount, output, nil
}
