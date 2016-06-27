package main

import (
	"fmt"
	dcos "github.com/asiainfoLDP/servicebroker_dcos/api"
	"github.com/asiainfoLDP/servicebroker_dcos/util/rand"
)

func newMysqlApp(id string) *dcos.App {
	return &dcos.App{
		Id:        fmt.Sprintf("/mysql/%s", id),
		Cpus:      1,
		Mem:       1000,
		Disk:      5000,
		Instances: 1,
		Container: &dcos.Container{
			Type: "DOCKER",
			Docker: dcos.Docker{
				Image:   "mysql:5.7.12",
				NetWork: "BRIDGE",
				PortMappings: []dcos.PortMapping{
					dcos.PortMapping{
						ContainerPort: 3306,
						HostPort:      uint32(10000 + rand.Intn(50000)),
						ServicePort:   0,
						Protocol:      "tcp",
					},
				},
				ForcePullImage: false,
			},
		},
		Env: map[string]string{
			"MYSQL_USER":          rand.String(10),
			"MYSQL_PASSWORD":      rand.String(16),
			"MYSQL_ROOT_PASSWORD": rand.String(16),
			"MYSQL_DATABASE":      rand.String(12),
		},
	}
}
