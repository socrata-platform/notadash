package lib

import (
    "log"
    mesos "github.com/boldfield/go-mesos"
)


type Mesos struct {
    Host string
    Cluster *mesos.Cluster
    _client *mesos.Client
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

func (m *Mesos) LoadCluster() error {
    if cluster, err := mesos.DiscoverCluster(m.Client()); err != nil {
        log.Println(err)
        return err
    } else {
        m.Cluster = cluster
    }

    if err := m.Cluster.LoadSlaveStates(m.Client()); err != nil {
        log.Printf("An error was encountered loading slave states: %s", err)
        return err
    }

    if err := m.Cluster.LoadSlaveStats(m.Client()); err != nil {
        log.Printf("An error was encountered loading slave states: %s", err)
        return err
    }
    return nil
}


func (m *Mesos) Framework(framework string) (map[string]*mesos.Framework) {
    return m.Cluster.GetFramework(framework)
}


func (m *Mesos) LoadSlave(host string) (*MesosSlave){
    slave := &mesos.Slave{ HostName: host }
    slave.LoadState(m.Client())
    return &MesosSlave{ Slave: slave }
}


func (s *MesosSlave) Framework(framework string) (map[string]*mesos.Framework) {
    return s.Slave.GetFramework(framework)
}
