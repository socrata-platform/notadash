package main

import (
	"fmt"
	lib "github.com/socrata-platform/notadash/lib"
	"github.com/codegangsta/cli"
	"github.com/marpaia/graphite-golang"
	"log"
	"net"
	"os"
	"strings"
)

func connectToGraphite(host string, port int) *graphite.Graphite {
	g, err := graphite.NewGraphite(host, port)
	if err != nil {
		log.Printf("An error occurred while trying to connect to the graphite server: (%s:%d)", host, port)
		log.Fatal(err)
	}
	return g
}

func sendToGraphite(g *graphite.Graphite, name string, path string, metric string) {
	mpath := strings.Join([]string{path, name}, ".")
	err := g.SimpleSend(mpath, metric)
	if err != nil {
		log.Printf("Failed to send %v from %v to Graphite server\n", name, path)
	}
	return
}

func loadSlave(host string) *lib.MesosSlave {
	mesos := &lib.Mesos{
		Host: host,
	}
	mesosClient := mesos.Client()
	slave, err := mesos.LoadSlaveStats(host, mesosClient)
	if err != nil {
		log.Printf("An error occured while gathering Slave stats for %v", host)
		log.Fatal(err)
	}
	return slave
}

func getHostIp() string {
	host, err := os.Hostname()
	out, err := net.LookupIP(host)
	if err != nil {
		log.Fatal("Unable to determin fully qualified name of host...")
	}
	return out[0].String()
}

func runReportSlaveAllocation(ctx *cli.Context) int {
	host_ip := ctx.GlobalString("ip-addr")
	if host_ip == "" {
		host_ip = getHostIp()
		log.Printf("setting host ip to %v", host_ip)
	}

	slave := loadSlave(host_ip)

	g := connectToGraphite(
		ctx.GlobalString("graphite-host"),
		ctx.GlobalInt("graphite-port"),
	)
	defer g.Disconnect()

	ss := slave.Slave.Stats
	m_path := strings.Join(
		[]string{
			"mesos-stats",
			strings.Replace(host_ip, ".", "-", -1),
		},
		".")

	sendToGraphite(g, "CpusPercent", m_path, fmt.Sprintf("%.2f", ss.CpusPercent))
	sendToGraphite(g, "CpusUsed", m_path, fmt.Sprintf("%.2f", ss.CpusUsed))
	sendToGraphite(g, "CpusTotal", m_path, fmt.Sprintf("%.2f", ss.CpusTotal))
	sendToGraphite(g, "MemPercent", m_path, fmt.Sprintf("%.2f", ss.MemPercent))
	sendToGraphite(g, "MemUsed", m_path, fmt.Sprintf("%.2f", ss.MemUsed))
	sendToGraphite(g, "MemTotal", m_path, fmt.Sprintf("%.2f", ss.MemTotal))
	sendToGraphite(g, "DiskPercent", m_path, fmt.Sprintf("%.2f", ss.DiskPercent))
	sendToGraphite(g, "DiskUsed", m_path, fmt.Sprintf("%.2f", ss.DiskUsed))
	sendToGraphite(g, "DiskTotal", m_path, fmt.Sprintf("%.2f", ss.DiskTotal))

	log.Printf("Metrics Sent!")
	return 0
}
