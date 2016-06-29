package main

import (
	"fmt"
	dcosapi "github.com/asiainfoLDP/servicebroker_dcos/api"
	broker "github.com/asiainfoLDP/servicebroker_dcos/servicebroker"
	"github.com/asiainfoLDP/servicebroker_dcos/util/rand"
)

func init() {
	kindToApiMappings[BrokerKind("mysql")] = &mysqlBroker{}
}

type mysqlBroker struct{}

func (s *mysqlBroker) Provision(instanceID string, details broker.ProvisionDetails, asyncAllowed bool) (broker.ProvisionedServiceSpec, error) {
	serviceSpec := broker.ProvisionedServiceSpec{IsAsync: true}

	if err := validate(&details); err != nil {
		return serviceSpec, err
	}

	catalog, err := getCatalog()
	if err != nil {
		return serviceSpec, err
	}

	var plan *broker.Plan
	if svc := catalog.GetService(details.ServiceID); svc == nil {
		return serviceSpec, errorMappings[serviceFieldMissingErrorKey]
	} else if plan = svc.GetPlan(details.PlanID); plan == nil {
		return serviceSpec, errorMappings[bindingFieldMissingErrorKey]
	}

	costs := plan.Metadata.Costs[0].Unit
	if len(costs) == 0 {
		return serviceSpec, fmt.Errorf("missing mysql plan cost info")
	}

	planCostUnit := parsePlanUnit(costs)

	app := newMysqlApp(instanceID, planCostUnit)
	mysqlApp, err := dcosClient.Application().Create(app)
	if err != nil {
		if err == dcosapi.ErrConflictInstance {
			return serviceSpec, broker.ErrInstanceAlreadyExists
		}
		return serviceSpec, err
	}

	temCache[instanceId(instanceID)] = mysqlApp

	return serviceSpec, nil
}

func (s *mysqlBroker) Bind(instanceID, bindingID string, details broker.BindDetails) (broker.Binding, error) {
	binding := broker.Binding{}

	app, ok := temCache[instanceId(instanceID)]
	if !ok {
		return binding, broker.ErrInstanceDoesNotExist
	}

	task, err := dcosClient.Task().Get(app.Id)
	if err != nil {
		return binding, fmt.Errorf("missing mysql instance %s ", instanceID)
	}

	binding.Credentials = map[string]string{
		"uri":      fmt.Sprintf("mysql://%s:%s@%s:%d/%s", app.Env["MYSQL_USER"], app.Env["MYSQL_PASSWORD"], task.Host, task.Ports[0], app.Env["MYSQL_DATABASE"]),
		"host":     task.Host,
		"port":     fmt.Sprintf("%d", task.Ports[0]),
		"username": app.Env["MYSQL_USER"],
		"password": app.Env["MYSQL_PASSWORD"],
		"database": app.Env["MYSQL_DATABASE"],
	}

	return binding, nil
}

func (s *mysqlBroker) Deprovision(instanceID string, details broker.DeprovisionDetails, asyncAllowed bool) (broker.IsAsync, error) {
	asynFlag := broker.IsAsync(true)
	//if asyncAllowed == false {
	//	return asynFlag, errors.New("Sync mode is not supported")
	//}

	mysqlApp, ok := temCache[instanceId(instanceID)]
	if !ok {
		return asynFlag, broker.ErrInstanceDoesNotExist
	}

	if err := dcosClient.Application().Delete(mysqlApp.Id); err != nil {
		return asynFlag, err
	}

	return asynFlag, nil
}

func newMysqlApp(id string, unit *planCostUnit) *dcosapi.App {
	return &dcosapi.App{
		Id:        fmt.Sprintf("/mysql/%s", id),
		Cpus:      unit.Cpu,
		Mem:       unit.Mem,
		Disk:      unit.Disk,
		Instances: 1,
		Container: &dcosapi.Container{
			Type: "DOCKER",
			Docker: dcosapi.Docker{
				Image:   "mysql:5.7.12",
				NetWork: "BRIDGE",
				PortMappings: []dcosapi.PortMapping{
					dcosapi.PortMapping{
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
