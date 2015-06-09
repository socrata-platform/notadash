package lib

import (
	mesos "github.com/boldfield/go-mesos"
	"log"
)

type FrameworkMap map[string]*mesos.Framework

type Mesos struct {
	Host        string
	Cluster     *mesos.Cluster
	_client     *mesos.Client
	initialized bool
}

type MesosSlave struct {
	Slave *mesos.Slave
}

func (m *Mesos) Client() *mesos.Client {
	if m._client != nil {
		return m._client
	}

	config := mesos.NewDefaultConfig()
	config.DiscoveryURL = m.Host
	m._client = mesos.NewClient(config)
	return m._client
}

func (m *Mesos) LoadCluster(c *mesos.Client) error {
	if cluster, err := mesos.DiscoverCluster(c); err != nil {
		log.Println(err)
		return err
	} else {
		m.Cluster = cluster
	}

	if err := m.Cluster.LoadSlaveStates(c); err != nil {
		log.Printf("An error was encountered loading slave states: %s", err)
		return err
	}

	if err := m.Cluster.LoadSlaveStats(c); err != nil {
		log.Printf("An error was encountered loading slave states: %s", err)
		return err
	}
	return nil
}

func (m *Mesos) Framework(framework string) FrameworkMap {
	return m.Cluster.GetFramework(framework)
}

func (m *Mesos) LoadSlaveState(host string, c *mesos.Client) (*MesosSlave, error) {
	slave := &mesos.Slave{HostName: host}
	err := slave.LoadState(c)
	return &MesosSlave{Slave: slave}, err
}

func (m *Mesos) LoadSlaveStats(host string, c *mesos.Client) (*MesosSlave, error) {
	slave := &mesos.Slave{HostName: host}
	err := slave.LoadStats(c)
	return &MesosSlave{Slave: slave}, err
}

func (s *MesosSlave) Framework(framework string) FrameworkMap {
	return s.Slave.GetFramework(framework)
}
