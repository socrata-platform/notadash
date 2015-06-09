package main

import (
	"fmt"
	lib "github.com/boldfield/notadash/lib"
	"github.com/codegangsta/cli"
	"github.com/ryanuber/columnize"
)

func runShowAllocation(ctx *cli.Context) int {
	fmt.Println("Collecting worker resource allocation data...")

	mesos := &lib.Mesos{
		Host: ctx.GlobalString("mesos-host"),
	}
	mesosClient := mesos.Client()
	mesos.LoadCluster(mesosClient)

	printAllocations(mesos)
	return 0
}

func printAllocations(mesos *lib.Mesos) int {
	output := make([]string, 1)
	output[0] = "Hostnamme | Cpu % | Cpu Ratio | Mem % | Mem Ratio | Disk % | Disk Ratio"

	for _, s := range mesos.Cluster.Slaves {
		ss := s.Stats
		ln := fmt.Sprintf(
			"%s | %.0f | %.1f/%d | %.0f | %d/%d | %.0f| %d/%d",
			s.HostName,
			ss.CpusPercent*100,
			ss.CpusUsed,
			ss.CpusTotal,
			ss.MemPercent*100,
			ss.MemUsed,
			ss.MemTotal,
			ss.DiskPercent*100,
			ss.DiskUsed,
			ss.DiskTotal,
		)
		output = append(output, ln)
	}

	result := columnize.SimpleFormat(output)
	fmt.Println(result)
	return 0
}
