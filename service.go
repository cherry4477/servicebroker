package main

import (
	"encoding/json"
	dcosapi "github.com/asiainfoLDP/servicebroker_dcos/api"
	broker "github.com/asiainfoLDP/servicebroker_dcos/servicebroker"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

//curl -XPUT michael:123456@127.0.0.1:5000/v2/service_instances/123 -d '{"organization_guid": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
func provisionHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	id := vars["instance_id"]

	var details broker.ProvisionDetails
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		respondUnprocessable(w, err)
		return
	}

	acceptsIncompleteFlag, _ := strconv.ParseBool(r.URL.Query().Get("accepts_incomplete"))

	serviceBroker, err := getServiceBroker(details.ServiceID)
	if err != nil {
		respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}

	provisionResponse, err := serviceBroker.Provision(id, details, acceptsIncompleteFlag)
	if err != nil {
		switch err {
		case broker.ErrInstanceAlreadyExists:
			respond(w, http.StatusConflict, struct{}{})
		case broker.ErrInstanceLimitMet:
			respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		case broker.ErrAsyncRequired:
			respond(w, 422, ErrorResponse{
				Error:       "AsyncRequired",
				Description: err.Error(),
			})
		default:
			respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}
		return
	}

	if provisionResponse.IsAsync {
		respond(w, http.StatusAccepted, broker.ProvisioningResponse{
			DashboardURL: provisionResponse.DashboardURL,
		})
	} else {
		respond(w, http.StatusCreated, broker.ProvisioningResponse{
			DashboardURL: provisionResponse.DashboardURL,
		})
	}
}

//curl -XPUT 127.0.0.1:5000/v2/service_instances/123/service_bindings/456 -d '{"app_gui": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
func bindHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	instanceID := vars["instance_id"]
	bindingID := vars["binding_id"]

	var details broker.BindDetails
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		respond(w, statusUnprocessableEntity, ErrorResponse{
			Description: err.Error(),
		})
		return
	}
	serviceBroker, err := getServiceBroker(details.ServiceID)
	if err != nil {
		respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}

	binding, err := serviceBroker.Bind(instanceID, bindingID, details)
	if err != nil {
		switch err {
		case broker.ErrInstanceDoesNotExist:
			respond(w, http.StatusNotFound, ErrorResponse{
				Description: err.Error(),
			})
		case broker.ErrBindingAlreadyExists:
			respond(w, http.StatusConflict, ErrorResponse{
				Description: err.Error(),
			})
		default:
			respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}
		return
	}

	respond(w, http.StatusCreated, binding)
}

//curl -XDELETE 127.0.0.1:5000/v2/service_instances/123/service_bindings/456
func unbindServiceInstanceHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	io.WriteString(w, "{}")
	return
}

//curl -XDELETE 127.0.0.1:5000/v2/service_instances/123?service_id=04EB4D8F-15BF-43F2-B4DA-E7A243E21C83\&plan_id=3C7BAF72-0DF4-420D-B365-B5CF09409C70
func deProvisionHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	instanceID := vars["instance_id"]

	details := broker.DeprovisionDetails{
		PlanID:    r.FormValue("plan_id"),
		ServiceID: r.FormValue("service_id"),
	}

	asyncAllowed := r.FormValue("accepts_incomplete") == "true"

	serviceBroker, err := getServiceBroker(details.ServiceID)
	if err != nil {
		respond(w, http.StatusInternalServerError, ErrorResponse{
			Description: err.Error(),
		})
		return
	}

	isAsync, err := serviceBroker.Deprovision(instanceID, details, asyncAllowed)
	if err != nil {
		switch err {
		case broker.ErrInstanceDoesNotExist:
			respond(w, http.StatusGone, struct{}{})
		case broker.ErrAsyncRequired:
			respond(w, 422, struct{}{})
		default:
			respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}
		return
	}

	if isAsync {
		respond(w, http.StatusAccepted, struct{}{})
	} else {
		respond(w, http.StatusOK, struct{}{})
	}
}
