package main

import (
	"log"
	"os"
	"strconv"

	esapi "github.com/1851616111/elastic-client"
)

var (
	elasticClient          *esapi.Client
	elasticCatalogFilePath string
)

func init() {
	start, _ := strconv.ParseBool(os.Getenv("ElasticSearch_Start"))
	if start {
		initElasticSearch()
	}
}

func initElasticSearch() {
	initElasticClient()
	registElasticBroker()

	log.Printf("init elastic search success.")
}

func initElasticClient() {
	addr := os.Getenv("ElasticSearch_Addr")
	if addr == "" {
		log.Fatal("env ElasticSearch_Addr is nil")
	}

	username := os.Getenv("ElasticSearch_Credential_Username")
	if username == "" {
		log.Fatal("env ElasticSearch_Credential_Username is nil")
	}

	password := os.Getenv("ElasticSearch_Credential_Password")
	if password == "" {
		log.Fatal("env ElasticSearch_Credential_Password is nil")
	}

	var err error
	if elasticClient, err = esapi.NewClient(addr); err != nil {
		log.Fatalf("init elastic search client err %v\n", err)
	}

	elasticClient.SetBasicAuth(username, password)
}

func registElasticBroker() {
	elasticCatalogFilePath = os.Getenv("ElasticSearch_Catalog_Path")
	if elasticCatalogFilePath == "" {
		log.Fatal("env ElasticSearch_Catalog_Path is nil")
	}
	allCatalogFilePaths = append(allCatalogFilePaths, elasticCatalogFilePath)

	brokerKinds = append(brokerKinds, "elasticsearch")
	kindToApiMappings[BrokerKind("elasticsearch")] = &elasticBroker{}
}

type elasticBroker struct{}

func (s *elasticBroker) Provision(instanceID string, details broker.ProvisionDetails, asyncAllowed bool) (broker.ProvisionedServiceSpec, error) {
	serviceSpec := broker.ProvisionedServiceSpec{IsAsync: true}

	if err := validate(&details); err != nil {
		return serviceSpec, err
	}

	catalog, err := getCatalog(elasticCatalogFilePath)
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

	instanceTmpCache[instanceId(instanceID)] = mysqlApp

	return serviceSpec, nil
}

func (s *elasticBroker) Bind(instanceID, bindingID string, details broker.BindDetails) (broker.Binding, error) {
	binding := broker.Binding{}

	app, ok := instanceTmpCache[instanceId(instanceID)]
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

func (s *elasticBroker) Deprovision(instanceID string, details broker.DeprovisionDetails, asyncAllowed bool) (broker.IsAsync, error) {
	asynFlag := broker.IsAsync(true)
	//if asyncAllowed == false {
	//	return asynFlag, errors.New("Sync mode is not supported")
	//}

	mysqlApp, ok := instanceTmpCache[instanceId(instanceID)]
	if !ok {
		return asynFlag, broker.ErrInstanceDoesNotExist
	}

	if err := dcosClient.Application().Delete(mysqlApp.Id); err != nil {
		return asynFlag, err
	}

	return asynFlag, nil
}
