package main

import (
	"encoding/json"
	"errors"
	"fmt"
	dcosapi "github.com/asiainfoLDP/servicebroker_dcos/api"
	"github.com/asiainfoLDP/servicebroker_dcos/servicebroker"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	dcosClient dcosapi.Interface
	temCache   = map[instanceId]*dcosapi.App{}
)

func init() {
	dcosHost := os.Getenv("Dcos_Host_Addr")
	if dcosHost == "" {
		log.Fatal("env Dcos_Host_Addr must not be nil.")
	}

	dcosToken := os.Getenv("Dcos_Token")
	if dcosToken == "" {
		log.Fatal("env Dcos_Token must not be nil.")
	}

	var err error
	dcosClient, err = dcosapi.NewClientInterface(dcosHost, dcosToken)
	if err != nil {
		log.Fatalf("init dcos(%s) client err %v\n", dcosHost, err)
	}

	log.Printf("init dcos(%s) client success.", dcosHost)
}

type instanceId string

//curl -XPUT 127.0.0.1:5000/v2/service_instances/123 -d '{"organization_guid": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
func createServiceInstanceHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	id := vars["instance_id"]
	if len(id) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, errors.New("request param service_instances must not be nil."))))
		return
	}

	newInstanceReq := new(servicebroker.InstanceRequest)
	if err := json.NewDecoder(r.Body).Decode(newInstanceReq); err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(500, err)))
		return
	}

	catalog, err := getCatalog()
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(500, err)))
		return
	}

	if err := newInstanceReq.Validate(catalog); err != nil {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, err)))
		return
	}

	var svc *servicebroker.Service
	if svc = catalog.GetService(newInstanceReq.ServiceId); svc == nil {
		w.WriteHeader(400)
		err := fmt.Errorf("no such service_id %s", newInstanceReq.ServiceId)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, err)))
		return
	} else if plan := svc.GetPlan(newInstanceReq.PlanId); plan == nil {
		w.WriteHeader(400)
		err := fmt.Errorf("no such plan_id %s", newInstanceReq.PlanId)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, err)))
		return
	}

	a := newMysqlApp(id)
	b, _ := json.Marshal(a)
	fmt.Printf("%s\n", b)
	mysqlApp, err := dcosClient.Application().Create(a)
	if err != nil {
		w.WriteHeader(400)
		err := fmt.Errorf("create app(%v) err %v", mysqlApp, err)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, err)))
		return
	}

	temCache[instanceId(id)] = mysqlApp

	io.WriteString(w, "{}")
	return
}

//curl -XPUT 127.0.0.1:5000/v2/service_instances/123/service_bindings/456 -d '{"app_gui": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
func bindServiceInstanceHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	instance_id, binding_id := vars["instance_id"], vars["binding_id"]
	if len(instance_id) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, errors.New("request param service_instances must not be nil."))))
		return
	}

	if len(binding_id) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, errors.New("request param binding_id must not be nil."))))
		return
	}

	bindInstanceReq := new(servicebroker.ServiceBindingRequest)
	if err := json.NewDecoder(r.Body).Decode(bindInstanceReq); err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(500, err)))
		return
	}

	app, ok := temCache[instanceId(instance_id)]
	if !ok {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(500, errors.New("no found service instance %s."))))
		return
	}

	task, err := dcosClient.Task().Get(app.Id)
	if err != nil {
		w.WriteHeader(500)
		err := fmt.Errorf("dcos get task(%s) err %v", app.Id, err)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(500, err)))
		return
	}

	rep := servicebroker.ServiceBindingResponse{
		Credentials: map[string]string{
			"uri":      fmt.Sprintf("mysql://%s:%s@%s:%d/%s", app.Env["MYSQL_USER"], app.Env["MYSQL_PASSWORD"], task.Host, task.Ports[0], app.Env["MYSQL_DATABASE"]),
			"host":     task.Host,
			"port":     fmt.Sprintf("%d", task.Ports[0]),
			"username": app.Env["MYSQL_USER"],
			"password": app.Env["MYSQL_PASSWORD"],
			"database": app.Env["MYSQL_DATABASE"],
		},
	}

	if err := json.NewEncoder(w).Encode(rep); err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(500, err)))
		return
	}

	return
}

//curl -XDELETE 127.0.0.1:5000/v2/service_instances/123/service_bindings/456
func unbindServiceInstanceHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	io.WriteString(w, "{}")
	return
}

//curl -XDELETE 127.0.0.1:5000/v2/service_instances/123
func deleteServiceInstanceHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	instance_id := vars["instance_id"]
	if len(instance_id) == 0 {
		w.WriteHeader(400)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, errors.New("request param service_instances must not be nil."))))
		return
	}

	mysqlApp, ok := temCache[instanceId(instance_id)]
	if !ok {
		w.WriteHeader(410)
		io.WriteString(w, "{}")
		return
	}

	if err := dcosClient.Application().Delete(mysqlApp.Id); err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(500, err)))
		return
	}

	io.WriteString(w, "{}")
	return
}
