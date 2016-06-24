package main

import (
	//"github.com/DataFoundry/servicebroker_dcos/api"

	"net/http"
)

func main() {

	//m.Get("/v2/catalog", brokerCatalog)
	//m.Put("/v2/service_instances/:service_id", createServiceInstance)
	//m.Delete("/v2/service_instances/:service_id", deleteServiceInstance)
	//m.Put("/v2/service_instances/:service_id/service_bindings/:binding_id", createServiceBinding)
	//m.Delete("/v2/service_instances/:service_id/service_bindings/:binding_id", deleteServiceBinding)
	//


	mux := http.NewServeMux()
	mux.Handle("/v2/catalog", catalogInfo)
	http.ListenAndServe(":5000", mux)
}
