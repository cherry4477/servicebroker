package main

import (
	"fmt"
	broker "github.com/asiainfoLDP/servicebroker_dcos/servicebroker"
	"log"
	"strings"
)

type BrokerKind string

var (
	brokerKinds          = []string{"mysql"}
	serviceIDKindMapping = map[string]BrokerKind{}
	kindToApiMappings    = map[BrokerKind]broker.ServiceBroker{}
)

func init() {
	cs, err := getCatalog()
	if err != nil {
		log.Fatalf("init servicebroker service err %v\n", err)
	}

	identityBrokersKind(cs)
}

func identityBrokersKind(c *broker.Catalog) {
	defer printServcieKinds()

	categorizeServiceNameToKind := func(svc *broker.Service) {
		name := strings.TrimSpace(svc.Name)

		for _, kind := range brokerKinds {
			if strings.Contains(strings.ToLower(name), kind) {
				serviceIDKindMapping[svc.ID] = BrokerKind(kind)
				return
			}
		}

		log.Printf("init catalog services unknow service %v\n", svc)
	}

	c.RangeCatalogFunc(categorizeServiceNameToKind)
}

func printServcieKinds() {
	for _, kind := range serviceIDKindMapping {
		log.Printf("regist %s servicebroker success.\n", kind)
	}
}

func getServiceBroker(serviceID string) (broker.ServiceBroker, error) {
	kind, ok := serviceIDKindMapping[serviceID]
	if !ok {
		return nil, fmt.Errorf("unkown service_id(%s)'s kind", serviceID)
	}

	var i broker.ServiceBroker
	i, ok = kindToApiMappings[kind]
	if !ok {
		return nil, fmt.Errorf("unkown service_id(%s)'s impliment", serviceID)
	}

	return i, nil
}
