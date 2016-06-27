package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

const Catalog_Info_Path = "catalog.json"

var router = mux.NewRouter()

func main() {

	http.Handle("/", router)
	router.HandleFunc("/v2/catalog", catalogHandler).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", httpHandlerMaker(createServiceInstanceHandler)).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", httpHandlerMaker(bindServiceInstanceHandler)).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}", httpHandlerMaker(deleteServiceInstanceHandler))
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", httpHandlerMaker(unbindServiceInstanceHandler)).Methods("DELETE")
	http.ListenAndServe(":5000", router)
}
