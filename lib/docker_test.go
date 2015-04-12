package lib

import (
	"encoding/json"
    "github.com/fsouza/go-dockerclient"

    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)


type MockDockerClient struct {
    mock.Mock
}

func (c *MockDockerClient) ListContainers(opts docker.ListContainersOptions) ([]docker.APIContainers, error) {
	jsonContainers := `[
     {
             "Id": "8dfafdbc3a40",
             "Image": "base:latest",
             "Command": "echo 1",
             "Created": 1367854155,
             "Ports":[{"PrivatePort": 2222, "PublicPort": 3333, "Type": "tcp"}],
             "Status": "Exit 0",
             "Names": ["cont1"]
     },
     {
             "Id": "9cd87474be90",
             "Image": "base:latest",
             "Command": "echo 222222",
             "Created": 1367854155,
             "Ports":[{"PrivatePort": 2222, "PublicPort": 3333, "Type": "tcp"}],
             "Status": "Exit 0",
             "Names": ["cont2"]
     },
     {
             "Id": "3176a2479c92",
             "Image": "base:latest",
             "Command": "echo 3333333333333333",
             "Created": 1367854154,
             "Ports":[{"PrivatePort": 2221, "PublicPort": 3331, "Type": "tcp"}],
             "Status": "Exit 0",
             "Names": ["cont3"]
     },
     {
             "Id": "4cb07b47f9fb",
             "Image": "base:latest",
             "Command": "echo 444444444444444444444444444444444",
             "Ports":[{"PrivatePort": 2223, "PublicPort": 3332, "Type": "tcp"}],
             "Created": 1367854152,
             "Status": "Exit 0",
             "Names": ["cont4"]
     }
]`
	var expected []docker.APIContainers
	json.Unmarshal([]byte(jsonContainers), &expected)
    return expected, nil
}


func (c *MockDockerClient) StopContainer(id string, ttl uint) (error) {
    return nil
}


func TestContainerExists(t *testing.T) {
    mockDockerClient := new(MockDockerClient)
    exists, err := ContainerRunning("cont4", mockDockerClient)
    assert.True(t, exists, "The container must exist")
    assert.Nil(t, err, "No error should be returned")
}


func TestContainerDoesNotExist(t *testing.T) {
    mockDockerClient := new(MockDockerClient)
    exists, err := ContainerRunning("cont5", mockDockerClient)
    assert.False(t, exists, "The container must not exist")
    assert.Nil(t, err, "No error should be returned")
}


func TestListContainers(t *testing.T) {
    expected := []string{ "cont1", "cont2", "cont3", "cont4" }
    mockDockerClient := new(MockDockerClient)
    exists, err := ListRunningContainers(mockDockerClient)
    assert.Nil(t, err, "No error should be returned")
    for i, cont := range expected {
        assert.Equal(t, cont, exists[i], "Should be running")
    }
}
