package api

import (
	"encoding/json"
	"log"
)

const (
	dcos_Marathon_List_Task_Api = "/service/marathon/v2/tasks"
)

type TasksInterface interface {
	Task() TaskInterface
}

type TaskInterface interface {
	List() (*Tasks, error)
}

type taskOption struct {
	*dcosOption
}

func (o *taskOption) List() (*Tasks, error) {

	cre_k, cre_v := getCredentials(o.Acs_token)
	data, err := httpGet(o.host+dcos_Marathon_List_Task_Api, cre_k, cre_v)
	if err != nil {
		return nil, err
	}

	ret := new(Tasks)
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}

	log.Println("[Info] list task success.\n")
	return ret, nil
}
