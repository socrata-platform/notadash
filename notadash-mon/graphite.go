package main

import (
	"fmt"
	"strings"
	"github.com/codegangsta/cli"
	"github.com/marpaia/graphite-golang"
	lib "github.com/boldfield/notadash/lib"
)


func connectToGraphite(host string, port int) (g *graphite.Graphite, err error) {
	if g, err = graphite.NewGraphite(host, port); err != nil {
		fmt.Println("An error occurred while trying to connect to the graphite server: (%s:%d)", host, port)
		fmt.Println(err)
	}
	fmt.Printf("Loaded Graphite connection: %#v\n", g)
	return
}

func sendToGraphite(g *graphite.Graphite, name string, path string, metric string) (err error) {
	mpath := strings.Join([]string{path, name}, ".")
	err = g.SimpleSend(mpath, metric)
	if err != nil {
		fmt.Println("Failed to send %v from %v to Graphite server", name, path)
	}
	return
}

func loadSlave(host string) (slave *lib.MesosSlave, err error) {
	mesos := &lib.Mesos{
		Host: host,
	}
	mesosClient := mesos.Client()
	slave, err = mesos.LoadSlaveStats(host, mesosClient)

	return
}

func runReportSlaveAllocation(ctx *cli.Context) int {
	hostname := ctx.GlobalString("hostname")
	slave, err := loadSlave(hostname)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	g, err := connectToGraphite(
		ctx.GlobalString("graphite-host"),
		ctx.GlobalInt("graphite-port"),
	)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer g.Disconnect()

	ss := slave.Slave.Stats
	m_path := strings.Join([]string{"mesos-stats", hostname}, ".")

	sendToGraphite(g, "CpusPercent", m_path, fmt.Sprintf("%.2f", ss.CpusPercent))
	sendToGraphite(g, "CpusUsed", m_path, fmt.Sprintf("%.2f", ss.CpusUsed))
	sendToGraphite(g, "CpusTotal", m_path, fmt.Sprintf("%d", ss.CpusTotal))
	sendToGraphite(g, "MemPercent", m_path, fmt.Sprintf("%.2f", ss.MemPercent))
	sendToGraphite(g, "MemUsed", m_path, fmt.Sprintf("%d", ss.MemUsed))
	sendToGraphite(g, "MemTotal", m_path, fmt.Sprintf("%d", ss.MemTotal))
	sendToGraphite(g, "DiskPercent", m_path, fmt.Sprintf("%.2f", ss.DiskPercent))
	sendToGraphite(g, "DiskUsed", m_path, fmt.Sprintf("%d", ss.DiskUsed))
	sendToGraphite(g, "DiskTotal", m_path, fmt.Sprintf("%d", ss.DiskTotal))

	return 0
}
