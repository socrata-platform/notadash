package lib

import (
	"fmt"
	"regexp"

	"github.com/fsouza/go-dockerclient"
)

var (
	containerRe, _ = regexp.Compile(`[^/]+`)
)

func NewDockerClient() *docker.Client {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	return client
}

type DockerClient interface {
	ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error)
	StopContainer(id string, ttl uint) error
}

// TODO (boldfield) :: This is kind of stupid... we shouldn't be loading all running containers
// TODO (boldfield) :: every time we check for the existance of a single container
func ContainerRunning(container string, client DockerClient) (bool, error) {
	if containers, err := client.ListContainers(docker.ListContainersOptions{All: false}); err != nil {
		return false, err
	} else {
		for _, c := range containers {
			for _, n := range c.Names {
				clean := containerRe.FindString(n)
				if clean == container {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// TODO (boldfield) :: Using this method as a data source for the above method would be nice...
func ListRunningContainers(client DockerClient) ([]string, error) {
	runningContainers := make([]string, 0)
	if containers, err := client.ListContainers(docker.ListContainersOptions{All: false}); err != nil {
		return runningContainers, err
	} else {
		for _, c := range containers {
			for _, n := range c.Names {
				clean := containerRe.FindString(n)
				runningContainers = append(runningContainers, clean)
			}
		}
	}
	return runningContainers, nil
}

func StopContainer(name string, timeout uint, client DockerClient) error {
	fmt.Printf("%s: %s\n", PrintRed("Stopping running container"), name)
	err := client.StopContainer(name, timeout) // Give the container 5min to shut down
	return err
}
