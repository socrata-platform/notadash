package lib

import (
    "os"
    "fmt"
    "github.com/fsouza/go-dockerclient"
)


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
                if n == container {
                    return true
                }
            }
        }
    }
    return false
}
