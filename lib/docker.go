package lib

import (
    "os"
    "fmt"
    "regexp"
    "github.com/fsouza/go-dockerclient"
)

var (
    containerRe, _ = regexp.Compile(`[^/]+`)
)

// TODO (boldfield) :: This is kind of stupid... we shouldn't be loading all running containers 
// TODO (boldfield) :: every time we check for the existance of a single container
func ContainerRunning(container string) (bool) {
    endpoint := "unix:///var/run/docker.sock"
    client, _ := docker.NewClient(endpoint)
    if containers, err := client.ListContainers(docker.ListContainersOptions{ All: false }); err != nil {
        fmt.Println(PrintRed("An error occoured while determining if docker container is running!"))
        fmt.Println(err)
        os.Exit(1)
    } else {
        for _, c := range containers {
            for _, n := range c.Names {
                clean := containerRe.FindString(n)
                if clean == container {
                    return true
                }
            }
        }
    }
    return false
}


// TODO (boldfield) :: Using this method as a data source for the above method would be nice...
func ListRunningContainers() ([]string) {
    runningContainers := make([]string, 0)
    endpoint := "unix:///var/run/docker.sock"
    client, _ := docker.NewClient(endpoint)
    if containers, err := client.ListContainers(docker.ListContainersOptions{ All: false }); err != nil {
        fmt.Println(PrintRed("An error occoured while determining if docker container is running!"))
        fmt.Println(err)
        os.Exit(1)
    } else {
        for _, c := range containers {
            for _, n := range c.Names {
                clean := containerRe.FindString(n)
                runningContainers = append(runningContainers, clean)
            }
        }
    }
    return runningContainers
}
