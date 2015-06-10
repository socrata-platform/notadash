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

	output := make([]string, 1)
	output[0] = "Hostname | Cpu % | Cpu Ratio | Mem % | Mem Ratio | Disk % | Disk Ratio"

	for _, s := range mesos.Cluster.Slaves {
		ss := s.Stats
		ln := fmt.Sprintf(
			"%s | %.2f | %.1f/%d | %.2f | %d/%d | %.2f| %d/%d",
			s.HostName,
			ss.CpusPercent,
			ss.CpusUsed,
			ss.CpusTotal,
			ss.MemPercent,
			ss.MemUsed,
			ss.MemTotal,
			ss.DiskPercent,
			ss.DiskUsed,
			ss.DiskTotal,
		)
		output = append(output, ln)
	}

	result := columnize.SimpleFormat(output)
	fmt.Println(result)
	return 0
}
