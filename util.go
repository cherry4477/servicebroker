package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"

	"errors"
	"github.com/1851616111/util/auth"
	broker "github.com/asiainfoLDP/servicebroker"
	"path/filepath"
	"strings"
)

const Catalog_Info_Path = "../catalog_mysql.json"

const instanceIDLogKey = "instance-id"
const instanceDetailsLogKey = "instance-details"
const bindingIDLogKey = "binding-id"

const invalidServiceDetailsErrorKey = "invalid-service-details"
const invalidBindDetailsErrorKey = "invalid-bind-details"
const invalidUnbindDetailsErrorKey = "invalid-unbind-details"
const invalidDeprovisionDetailsErrorKey = "invalid-deprovision-details"
const instanceLimitReachedErrorKey = "instance-limit-reached"
const instanceAlreadyExistsErrorKey = "instance-already-exists"
const bindingAlreadyExistsErrorKey = "binding-already-exists"
const instanceMissingErrorKey = "instance-missing"
const serviceFieldMissingErrorKey = "service-missing"
const bindingFieldMissingErrorKey = "binding-missing"
const asyncRequiredKey = "async-required"
const planChangeNotSupportedKey = "plan-change-not-supported"
const unknownErrorKey = "unknown-error"

const statusUnprocessableEntity = 422

var errorMappings = map[string]error{
	instanceMissingErrorKey:     errors.New("request missing param instance id"),
	serviceFieldMissingErrorKey: errors.New("request missing param service id"),
	bindingFieldMissingErrorKey: errors.New("request missing param binding id"),
}

func jsonFileUnMarshal(path string, c interface{}) error {
	var err error
	if path, err = filepath.Abs(path); err != nil {
		return err

	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(c)
}

func middler(auth *auth.Wrapper, f func(http.ResponseWriter, *http.Request, map[string]string)) func(w http.ResponseWriter, r *http.Request) {
	return auth.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := mux.Vars(r)["instance_id"]; id == "" {
			respondUnprocessable(w, errorMappings[instanceMissingErrorKey])
			return
		}

		f(w, r, mux.Vars(r))
	})
}

func middlerCatalog(auth *auth.Wrapper, f func(http.ResponseWriter, *http.Request, map[string]string)) func(w http.ResponseWriter, r *http.Request) {
	return auth.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		f(w, r, mux.Vars(r))
	})
}

func respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		fmt.Printf("response being attempted %d %#v\n", status, response)
	}
}

func respondUnprocessable(w http.ResponseWriter, err error) {
	respond(w, statusUnprocessableEntity, ErrorResponse{
		Description: err.Error(),
	})
}

type ErrorResponse struct {
	Error       string `json:"error, omitempty"`
	Description string `json:"description"`
}

func validate(details *broker.ProvisionDetails) error {
	if len(strings.TrimSpace(details.ServiceID)) == 0 {
		return fmt.Errorf("invalid json field service_id")
	}
	if len(strings.TrimSpace(details.PlanID)) == 0 {
		return fmt.Errorf("invalid json field plan_id")
	}
	if len(strings.TrimSpace(details.OrganizationGUID)) == 0 {
		return fmt.Errorf("invalid json field organization_guid")
	}

	return nil
}
