package api

import (
	"encoding/json"
	"log"
)

const (
	dcos_Marathon_List_Deployment_Api = "/service/marathon/v2/deployments"
)

type DeploymentsInterface interface {
	Deployment() DeploymentInterface
}

type DeploymentInterface interface {
	List() ([]Deployment, error)
}

type deploymentOption struct {
	*dcosOption
}

func (o *deploymentOption) List() ([]Deployment, error) {
	cre_k, cre_v := getCredentials(o.Acs_token)

	data, err := httpGet(o.host+dcos_Marathon_List_Deployment_Api, cre_k, cre_v)
	if err != nil {
		return nil, err
	}

	ds := []Deployment{}
	if err := json.Unmarshal(data, &ds); err != nil {
		return nil, err
	}

	log.Println("[Info] list deployment success.")
	return ds, nil
}
