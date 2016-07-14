package main

import (
	"github.com/1851616111/util/auth"
	"github.com/gorilla/mux"
	"net/http"
)

var router = mux.NewRouter()

func main() {
	initCatalog()

	admin := auth.NewWrapper("michael", "123456")

	http.Handle("/", router)
	router.HandleFunc("/v2/catalog", middlerCatalog(admin, catalogHandler)).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", middler(admin, provisionHandler)).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", middler(admin, bindHandler)).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}", middler(admin, deProvisionHandler))
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", middler(admin, unbindServiceInstanceHandler)).Methods("DELETE")
	http.ListenAndServe(":5000", router)
}
