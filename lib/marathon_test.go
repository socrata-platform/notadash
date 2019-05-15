package lib

import (
	"encoding/json"
	marathon "github.com/gambol99/go-marathon"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
        "net/url"
)

type MockMarathonClient struct {
	mock.Mock
}

func (c *MockMarathonClient) ListApplications(v url.Values) ([]string, error) {
	var thing []string
	return thing, nil
}

func (c *MockMarathonClient) Tasks(string) (*marathon.Tasks, error) {
	jsonTasks := `{
    "tasks": [
        {
            "appId":"/infrastructure/docker-registry",
            "id":"infrastructure_docker-registry.872a3060-dfae-11e4-bdf0-3624e0a93f16",
            "host":"ip-10-120-100-25.us-west-2.compute.internal",
            "ports":[5000],
            "startedAt":"2015-04-10T18:22:52.829Z",
            "stagedAt":"2015-04-10T18:22:21.295Z",
            "version":"2015-04-09T23:40:51.805Z"
        }
    ]
}`
	var expected marathon.Tasks
	json.Unmarshal([]byte(jsonTasks), &expected)
	return &expected, nil
}

func (c *MockMarathonClient) Applications(v url.Values) (*marathon.Applications, error) {
	jsonContainers := `{
    "apps": [
        {
            "args": null,
            "backoffFactor": 1.15,
            "backoffSeconds": 1,
            "cmd": "python3 -m http.server 8080",
            "constraints": [],
            "container": {
                "docker": {
                    "image": "python:3",
                    "network": "BRIDGE",
                    "portMappings": [
                        {
                            "containerPort": 8080,
                            "hostPort": 0,
                            "servicePort": 9000,
                            "protocol": "tcp"
                        },
                        {
                            "containerPort": 161,
                            "hostPort": 0,
                            "protocol": "udp"
                        }
                    ]
                },
                "type": "DOCKER",
                "volumes": []
            },
            "cpus": 0.5,
            "dependencies": [],
            "deployments": [],
            "disk": 0.0,
            "env": {},
            "executor": "",
            "healthChecks": [
                {
                    "command": null,
                    "gracePeriodSeconds": 5,
                    "intervalSeconds": 20,
                    "maxConsecutiveFailures": 3,
                    "path": "/",
                    "portIndex": 0,
                    "protocol": "HTTP",
                    "timeoutSeconds": 20
                }
            ],
            "id": "/fake_app",
            "instances": 2,
            "mem": 64.0,
            "ports": [
                10000,
                10001
            ],
            "requirePorts": false,
            "storeUrls": [],
            "tasksRunning": 2,
            "tasksStaged": 0,
            "upgradeStrategy": {
                "minimumHealthCapacity": 1.0
            },
            "uris": [],
            "user": null,
            "version": "2014-09-25T02:26:59.256Z"
        },
        {
            "args": null,
            "backoffFactor": 1.15,
            "backoffSeconds": 1,
            "cmd": "python3 -m http.server 8080",
            "constraints": [],
            "container": {
                "docker": {
                    "image": "python:3",
                    "network": "BRIDGE",
                    "portMappings": [
                        {
                            "containerPort": 8080,
                            "hostPort": 0,
                            "servicePort": 9000,
                            "protocol": "tcp"
                        },
                        {
                            "containerPort": 161,
                            "hostPort": 0,
                            "protocol": "udp"
                        }
                    ]
                },
                "type": "DOCKER",
                "volumes": []
            },
            "cpus": 1.5,
            "dependencies": [],
            "deployments": [],
            "disk": 0.0,
            "env": {},
            "executor": "",
            "healthChecks": [
                {
                    "command": null,
                    "gracePeriodSeconds": 5,
                    "intervalSeconds": 20,
                    "maxConsecutiveFailures": 3,
                    "path": "/",
                    "portIndex": 0,
                    "protocol": "HTTP",
                    "timeoutSeconds": 20
                }
            ],
            "id": "/fake_app_broken",
            "instances": 2,
            "mem": 64.0,
            "ports": [
                10000,
                10001
            ],
            "requirePorts": false,
            "storeUrls": [],
            "tasksRunning": 2,
            "tasksStaged": 0,
            "upgradeStrategy": {
                "minimumHealthCapacity": 1.0
            },
            "uris": [],
            "user": null,
            "version": "2014-09-25T02:26:59.256Z"
        }
    ]
}`
	var expected marathon.Applications
	json.Unmarshal([]byte(jsonContainers), &expected)
	return &expected, nil
}

func expApplications() *marathon.Applications {
        mem := float64(64)
        disk := float64(0)
	return &marathon.Applications{
		Apps: []marathon.Application{
			marathon.Application{
				ID:   "/fake_app",
				CPUs: float64(0.5),
				Mem:  &mem,
				Disk: &disk,
			},
			marathon.Application{
				ID:   "/fake_app_broken",
				CPUs: float64(1.5),
				Mem:  &mem,
				Disk: &disk,
			},
		},
	}
}

func expMarathonApps() *MarathonApps {
	return &MarathonApps{
		Apps: []*MarathonApp{
			expFakeApp(),
		},
	}
}

func expFakeApp() *MarathonApp {
	return &MarathonApp{
		Id: "/fake_app",
		Tasks: []*MarathonTask{
			&MarathonTask{
				Id:        "fake-app-task",
				Container: "fake-app-container",
				SlaveId:   "slave1",
				SlaveHost: "slave1-host",
			},
		},
	}
}

func TestLoadApps(t *testing.T) {
	m := &Marathon{}
	mockClient := new(MockMarathonClient)
	err := m.LoadApps(mockClient)
	exp := expApplications()
	assert.Nil(t, err)
	assert.Equal(t, len(m.Apps), len(exp.Apps), "The expected number of apps should be returned")
	for i, app := range exp.Apps {
		assert.Equal(t, m.Apps[i].ID, app.ID, "App should be returned")
	}
}

func TestGetAppById(t *testing.T) {
	apps := expMarathonApps()
	app := apps.GetAppById("/fake_app")
	assert.Equal(t, app.Id, "/fake_app")
}

func TestGetAppByIdNotExist(t *testing.T) {
	apps := expMarathonApps()
	app := apps.GetAppById("/no_app")
	assert.Nil(t, app)
}

func TestGetTaskById(t *testing.T) {
	app := expFakeApp()
	task := app.GetTaskById("fake-app-task")
	assert.Equal(t, task.Id, "fake-app-task")
}

func TestAddTaskExistingApp(t *testing.T) {
	apps := expMarathonApps()
	apps.AddTask("new-fake-app-task", "/fake_app", "slave2", "slave2-host", true, false, false)
	assert.Equal(t, apps.Apps[0].Tasks[0].Id, "fake-app-task", "Original task should still exist")
	assert.Equal(t, apps.Apps[0].Tasks[1].Id, "new-fake-app-task", "New Task ID should exist")
}

func TestAddTaskNotExistingApp(t *testing.T) {
	apps := expMarathonApps()
	apps.AddTask("faker-app-task", "/faker_app", "slave1", "slave1-host", true, false, false)
	assert.Equal(t, apps.Apps[1].Id, "/faker_app", "New app should exist")
	assert.Equal(t, apps.Apps[1].Tasks[0].Id, "faker-app-task", "New app task should exist")
}

func TestAddAppExisting(t *testing.T) {
	apps := expMarathonApps()
	apps.AddApp("/fake_app")
	assert.Equal(t, len(apps.Apps), 1, "No new app should be added")
	assert.Equal(t, apps.Apps[0].Id, "/fake_app", "Original app should still exist")
}

func TestAddAppNotExisting(t *testing.T) {
	apps := expMarathonApps()
	apps.AddApp("/faker_app")
	assert.Equal(t, len(apps.Apps), 2, "An additional app should exist")
	assert.Equal(t, apps.Apps[0].Id, "/fake_app", "Original app should still exist")
	assert.Equal(t, apps.Apps[1].Id, "/faker_app", "New app should exist")
}
