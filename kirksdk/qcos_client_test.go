package kirksdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"
)

func TestStacks(t *testing.T) {
	expectedUrl := "/v3/stacks"
	expectedMethod := "GET"
	expectedRet := StackInfo{
		IsDeployed: true,
		Metadata:   []string{"key=value"},
		Name:       "qiniu-app",
		Services:   []string{"nginx", "mongo"},
		Status:     StatusRunning,
	}
	ret := `[{
"name": "qiniu-app",
"services": [
  "nginx",
  "mongo"
],
"metadata": [
  "key=value"
],
"status": "RUNNING",
"isDeployed": true
}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedUrl, r.URL.Path)
		assert.Equal(t, expectedMethod, r.Method)
		fmt.Fprintln(w, ret)
	}))
	defer ts.Close()
	client := NewQcosClient(QcosConfig{
		Host: ts.URL,
	})
	stacks, err := client.ListStacks(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stacks))
	assert.Equal(t, expectedRet, stacks[0])
}

func TestServicesCreate(t *testing.T) {
	args := CreateServiceArgs{
		Name: "s1",
		Spec: ServiceSpec{
			Image:       "nginx",
			Command:     []string{"echo", "hello"},
			AutoRestart: "always",
			Envs:        []string{"env1=v1", "env2=v2"},
			Hosts:       []string{"earth:1.1.1.1", "mars:2.2.2.2"},
			EntryPoint:  []string{"/bin/sh", "-c"},
			LogCollectors: []LogCollectorSpec{
				{
					Directory: "/var/log/",
					Patterns:  []string{"*.log", "*.txt"},
				},
				{
					Directory: "/run/log/",
					Patterns:  []string{"*.log"},
				},
			},
			StopGraceSec: 5,
			WorkDir:      "/home/",
			UnitType:     "S3_1U2G",
		},
		InstanceNum:       3,
		UpdateParallelism: 2,
		Metadata:          []string{"m1=v1", "m2=v2"},
		Stateful:          true,
		Volumes: []VolumeSpec{
			{
				FsType:    "ext4",
				UnitType:  "SSD1_16G",
				MountPath: "/mnt",
				Name:      "v1",
			},
			{
				FsType:    "ext4",
				UnitType:  "SSD1_16G",
				MountPath: "/mnt2",
				Name:      "v2",
			},
		},
	}
	expecetdUrl := "/v3/stacks/default/services"
	expectedMethod := "POST"
	expectedArgs := `{
  "instanceNum": 3,
	"updateParallelism": 2,
  "metadata": [
      "m1=v1",
      "m2=v2"
  ],
  "name": "s1",
  "spec": {
      "autoRestart": "always",
      "command": [
          "echo",
          "hello"
      ],
      "entryPoint": [
          "/bin/sh",
          "-c"
      ],
      "envs": [
          "env1=v1",
          "env2=v2"
      ],
      "hosts": [
          "earth:1.1.1.1",
          "mars:2.2.2.2"
      ],
      "image": "nginx",
      "logCollectors": [
          {
              "directory": "/var/log/",
              "patterns": [
                  "*.log",
                  "*.txt"
              ]
          },
          {
              "directory": "/run/log/",
              "patterns": [
                  "*.log"
              ]
          }
      ],
      "stopGraceSec": 5,
      "workDir": "/home/",
	  "unitType": "S3_1U2G"
  },
  "stateful": true,
  "volumes": [
      {
          "fsType": "ext4",
          "mountPath": "/mnt",
          "name": "v1",
		  "unitType": "SSD1_16G"
      },
      {
          "fsType": "ext4",
          "mountPath": "/mnt2",
          "name": "v2",
		  "unitType": "SSD1_16G"
      }
  ]
}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expecetdUrl, r.URL.Path)
		assert.Equal(t, expectedMethod, r.Method)
		var (
			actual   CreateServiceArgs
			expected CreateServiceArgs
		)
		b, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)
		err = json.Unmarshal(b, &actual)
		assert.NoError(t, err)
		err = json.Unmarshal([]byte(expectedArgs), &expected)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}))
	defer ts.Close()
	client := NewQcosClient(QcosConfig{
		Host: ts.URL,
	})
	err := client.CreateService(context.TODO(), "default", args)
	assert.NoError(t, err)
}

func TestServicesInspect(t *testing.T) {
	expectedUrl := "/v3/stacks/default/services/s1/inspect"
	expectedMethod := "GET"
	expectedRet := ServiceInfo{
		ContainerIPs:      []string{"1.1.1.1", "1.1.1.2"},
		InstanceNum:       5,
		UpdateParallelism: 2,
		Metadata:          []string{},
		Name:              "spaceship",
		Revision:          1,
		Spec: ServiceSpecExport{
			Image:       "nginx",
			Command:     []string{"echo", "hello"},
			AutoRestart: "always",
			Envs:        []string{"env1=v1", "env2=v2"},
			Hosts:       []string{"earth:1.1.1.1", "mars:2.2.2.2"},
			EntryPoint:  []string{"/bin/sh", "-c"},
			LogCollectors: []LogCollectorSpec{
				{
					Directory: "/var/log/",
					Patterns:  []string{"*.txt", "*.log"},
				},
				{
					Directory: "/run/log/",
					Patterns:  []string{"*.txt"},
				},
			},
			StopGraceSec: 5,
			WorkDir:      "/home/",
		},
		Stack:    "universe",
		State:    StateDeployed,
		Stateful: true,
		Status:   StatusRunning,
		UpdateSpec: ServiceSpecExport{
			Image:       "nginx:v2",
			Command:     []string{"echo", "hello"},
			AutoRestart: "always",
			Envs:        []string{"env1=v1", "env2=v2"},
			Hosts:       []string{"earth:1.1.1.1", "mars:2.2.2.2"},
			EntryPoint:  []string{"/bin/sh", "-c"},
			LogCollectors: []LogCollectorSpec{
				{
					Directory: "/var/log/",
					Patterns:  []string{"*.txt", "*.log"},
				},
				{
					Directory: "/run/log/",
					Patterns:  []string{"*.txt"},
				},
			},
			StopGraceSec: 5,
			WorkDir:      "/home/",
		},
		Volumes: []VolumeSpec{
			{
				FsType:    "ext4",
				UnitType:  "SSD1_16G",
				MountPath: "/mnt",
				Name:      "v1",
			},
			{
				FsType:    "ext4",
				UnitType:  "SSD1_16G",
				MountPath: "/mnt2",
				Name:      "v2",
			},
		},
	}
	ret := `{
    "containerIps": [
        "1.1.1.1",
        "1.1.1.2"
    ],
    "instanceNum": 5,
		"updateParallelism": 2,
    "metadata": [],
    "name": "spaceship",
    "revision": 1,
    "spec": {
        "autoRestart": "always",
        "command": [
            "echo",
            "hello"
        ],
        "entryPoint": [
            "/bin/sh",
            "-c"
        ],
        "envs": [
            "env1=v1",
            "env2=v2"
        ],
        "hosts": [
            "earth:1.1.1.1",
            "mars:2.2.2.2"
        ],
        "image": "nginx",
        "logCollectors": [
            {
                "directory": "/var/log/",
                "patterns": [
                    "*.txt",
                    "*.log"
                ]
            },
            {
                "directory": "/run/log/",
                "patterns": [
                    "*.txt"
                ]
            }
        ],
        "stopGraceSec": 5,
        "workDir": "/home/"
    },
    "stack": "universe",
    "state": "DEPLOYED",
    "stateful": true,
    "status": "RUNNING",
    "updateSpec": {
        "autoRestart": "always",
        "command": [
            "echo",
            "hello"
        ],
        "entryPoint": [
            "/bin/sh",
            "-c"
        ],
        "envs": [
            "env1=v1",
            "env2=v2"
        ],
        "hosts": [
            "earth:1.1.1.1",
            "mars:2.2.2.2"
        ],
        "image": "nginx:v2",
        "logCollectors": [
            {
                "directory": "/var/log/",
                "patterns": [
                    "*.txt",
                    "*.log"
                ]
            },
            {
                "directory": "/run/log/",
                "patterns": [
                    "*.txt"
                ]
            }
        ],
        "stopGraceSec": 5,
        "workDir": "/home/"
    },
    "volumes": [
        {
            "fsType": "ext4",
            "mountPath": "/mnt",
            "name": "v1",
			"unitType": "SSD1_16G"
        },
        {
            "fsType": "ext4",
            "mountPath": "/mnt2",
            "name": "v2",
			"unitType": "SSD1_16G"
        }
    ]
}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedUrl, r.URL.Path)
		assert.Equal(t, expectedMethod, r.Method)
		fmt.Fprintln(w, ret)
	}))
	defer ts.Close()
	client := NewQcosClient(QcosConfig{
		Host: ts.URL,
	})
	info, err := client.GetServiceInspect(context.TODO(), "default", "s1")
	assert.NoError(t, err)
	assert.Equal(t, expectedRet, info)
}

func TestServices(t *testing.T) {
	expectedUrl := "/v3/stacks/default/services"
	expectedMethod := "GET"
	expectedRet := []ServiceInfo{{
		ContainerIPs:      []string{"1.1.1.1", "1.1.1.2"},
		InstanceNum:       5,
		UpdateParallelism: 2,
		Metadata:          []string{},
		Name:              "spaceship",
		Revision:          1,
		Spec: ServiceSpecExport{
			Image:       "nginx",
			Command:     []string{"echo", "hello"},
			AutoRestart: "always",
			Envs:        []string{"env1=v1", "env2=v2"},
			Hosts:       []string{"earth:1.1.1.1", "mars:2.2.2.2"},
			EntryPoint:  []string{"/bin/sh", "-c"},
			LogCollectors: []LogCollectorSpec{
				{
					Directory: "/var/log/",
					Patterns:  []string{"*.txt", "*.log"},
				},
				{
					Directory: "/run/log/",
					Patterns:  []string{"*.txt"},
				},
			},
			StopGraceSec: 5,
			WorkDir:      "/home/",
		},
		Stack:    "universe",
		State:    StateDeployed,
		Stateful: true,
		Status:   StatusRunning,
		UpdateSpec: ServiceSpecExport{
			Image:       "nginx:v2",
			Command:     []string{"echo", "hello"},
			AutoRestart: "always",
			Envs:        []string{"env1=v1", "env2=v2"},
			Hosts:       []string{"earth:1.1.1.1", "mars:2.2.2.2"},
			EntryPoint:  []string{"/bin/sh", "-c"},
			LogCollectors: []LogCollectorSpec{
				{
					Directory: "/var/log/",
					Patterns:  []string{"*.txt", "*.log"},
				},
				{
					Directory: "/run/log/",
					Patterns:  []string{"*.txt"},
				},
			},
			StopGraceSec: 5,
			WorkDir:      "/home/",
		},
		Volumes: []VolumeSpec{
			{
				FsType:    "ext4",
				UnitType:  "SSD1_16G",
				MountPath: "/mnt",
				Name:      "v1",
			},
			{
				FsType:    "ext4",
				UnitType:  "SSD1_16G",
				MountPath: "/mnt2",
				Name:      "v2",
			},
		},
	}}
	ret := `[{
    "containerIps": [
        "1.1.1.1",
        "1.1.1.2"
    ],
    "instanceNum": 5,
		"updateParallelism": 2,
    "metadata": [],
    "name": "spaceship",
    "revision": 1,
    "spec": {
        "autoRestart": "always",
        "command": [
            "echo",
            "hello"
        ],
        "entryPoint": [
            "/bin/sh",
            "-c"
        ],
        "envs": [
            "env1=v1",
            "env2=v2"
        ],
        "hosts": [
            "earth:1.1.1.1",
            "mars:2.2.2.2"
        ],
        "image": "nginx",
        "logCollectors": [
            {
                "directory": "/var/log/",
                "patterns": [
                    "*.txt",
                    "*.log"
                ]
            },
            {
                "directory": "/run/log/",
                "patterns": [
                    "*.txt"
                ]
            }
        ],
        "stopGraceSec": 5,
        "workDir": "/home/"
    },
    "stack": "universe",
    "state": "DEPLOYED",
    "stateful": true,
    "status": "RUNNING",
    "updateSpec": {
        "autoRestart": "always",
        "command": [
            "echo",
            "hello"
        ],
        "entryPoint": [
            "/bin/sh",
            "-c"
        ],
        "envs": [
            "env1=v1",
            "env2=v2"
        ],
        "hosts": [
            "earth:1.1.1.1",
            "mars:2.2.2.2"
        ],
        "image": "nginx:v2",
        "logCollectors": [
            {
                "directory": "/var/log/",
                "patterns": [
                    "*.txt",
                    "*.log"
                ]
            },
            {
                "directory": "/run/log/",
                "patterns": [
                    "*.txt"
                ]
            }
        ],
        "stopGraceSec": 5,
        "workDir": "/home/"
    },
    "volumes": [
        {
            "fsType": "ext4",
            "mountPath": "/mnt",
            "name": "v1",
			"unitType": "SSD1_16G"
        },
        {
            "fsType": "ext4",
            "mountPath": "/mnt2",
            "name": "v2",
			"unitType": "SSD1_16G"
        }
    ]
}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedUrl, r.URL.Path)
		assert.Equal(t, expectedMethod, r.Method)
		fmt.Fprintln(w, ret)
	}))
	defer ts.Close()
	client := NewQcosClient(QcosConfig{
		Host: ts.URL,
	})
	info, err := client.ListServices(context.TODO(), "default")
	assert.NoError(t, err)
	assert.Equal(t, expectedRet, info)
}
