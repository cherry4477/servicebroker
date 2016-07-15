package main

import (
	"log"
	"os"
	"strconv"

	esapi "github.com/1851616111/elastic_client"
	"github.com/1851616111/util/rand"
	broker "github.com/asiainfoLDP/servicebroker"
	"time"

	"fmt"
	"net"
	"net/url"
)

var (
	elasticAddr             string
	elasticHost             string
	elasticPort             string
	elasticClient           *esapi.Client
	elasticCatalogFilePath  string
	elasticInstanceTmpCache map[instanceId]elasticInfo
)

type elasticInfo struct {
	username string
	password string
	role     string
	index    string
}

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

	u, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("parse ElasticSearch_Addr(%s) err %v", addr, err)
	}

	if elasticHost, elasticPort, err = net.SplitHostPort(u.Host); err != nil {
		log.Fatalf("parse ElasticSearch_Addr(%s) err %v", addr, err)
	}

	username := os.Getenv("ElasticSearch_Credential_Username")
	if username == "" {
		log.Fatal("env ElasticSearch_Credential_Username is nil")
	}

	password := os.Getenv("ElasticSearch_Credential_Password")
	if password == "" {
		log.Fatal("env ElasticSearch_Credential_Password is nil")
	}

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

	elasticInstanceTmpCache = map[instanceId]elasticInfo{}
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

	indexName, roleName, userName, userPassword := rand.LowerString(20), rand.String(20), rand.String(20), rand.String(24)

	newIndexOptions := *esapi.DefaultIndexOption
	newIndexOptions.IndexName = indexName
	if err := elasticClient.CreateIndex(&newIndexOptions); err != nil {
		return serviceSpec, err
	}

	newRoleOptions := esapi.NewDefaultRoleOption(roleName, indexName)
	if err := elasticClient.CreateRole(newRoleOptions); err != nil {
		return serviceSpec, err
	}

	newUserOptions := &esapi.CreateUserOptions{
		Name: userName,
		User: esapi.User{
			Password: userPassword,
			Roles:    []string{roleName},
			Metadata: map[string]string{
				"time": time.Now().String(),
				"kind": "servicebroker",
			},
		},
	}
	if err := elasticClient.CreateUser(newUserOptions); err != nil {
		return serviceSpec, err
	}

	elasticInstanceTmpCache[instanceId(instanceID)] = elasticInfo{
		username: userName,
		password: userPassword,
		role:     roleName,
		index:    indexName,
	}

	fmt.Printf("[Info] elastic generate servicebroker %v\n", elasticInstanceTmpCache[instanceId(instanceID)])
	return serviceSpec, nil
}

func (s *elasticBroker) Bind(instanceID, bindingID string, details broker.BindDetails) (broker.Binding, error) {
	binding := broker.Binding{}

	brokerInfo, ok := elasticInstanceTmpCache[instanceId(instanceID)]
	if !ok {
		return binding, broker.ErrInstanceDoesNotExist
	}

	binding.Credentials = map[string]string{
		"uri":      fmt.Sprintf("%s:%s@%s:%s/%s", brokerInfo.username, brokerInfo.password, elasticHost, elasticPort, brokerInfo.index),
		"host":     elasticHost,
		"port":     elasticPort,
		"username": brokerInfo.username,
		"password": brokerInfo.password,
		"database": brokerInfo.index,
	}

	return binding, nil
}

func (s *elasticBroker) Deprovision(instanceID string, details broker.DeprovisionDetails, asyncAllowed bool) (broker.IsAsync, error) {
	asynFlag := broker.IsAsync(true)
	//if asyncAllowed == false {
	//	return asynFlag, errors.New("Sync mode is not supported")
	//}

	brokerInfo, ok := elasticInstanceTmpCache[instanceId(instanceID)]
	if !ok {
		return asynFlag, broker.ErrInstanceDoesNotExist
	}

	var err error
	if err = elasticClient.DeleteIndex(brokerInfo.index); err != nil {
		return asynFlag, err
	}
	if err = elasticClient.DeleteUser(brokerInfo.username); err != nil {
		return asynFlag, err
	}
	if err = elasticClient.DeleteRole(brokerInfo.role); err != nil {
		return asynFlag, err
	}

	return asynFlag, nil
}
