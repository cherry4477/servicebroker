package api

import (
	"encoding/json"
	"fmt"
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
	Get(string) (*task, error)
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

func (o *taskOption) Get(id string) (*task, error) {
	tasks, err := o.List()
	if err != nil {
		return nil, err
	}

	isIDTask := func(t *task) bool {
		if t.AppId == id {
			return true
		}
		return false
	}

	if t := tasks.filterTasksFunc(isIDTask); t == nil {
		return nil, fmt.Errorf("no found")
	} else {
		log.Printf("[Info] get task %s success.\n", id)
		return t, nil
	}
}

func (t *Tasks) filterTasksFunc(filter func(*task) bool) *task {
	if len(t.Tasks) == 0 {
		return nil
	}

	for _, tsk := range t.Tasks {
		if filter(&tsk) {
			return &tsk
		}
	}

	return nil
}
