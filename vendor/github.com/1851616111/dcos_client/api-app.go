package api

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	dcos_Create_Marathon_App_Api = "/service/marathon/v2/apps"
	dcos_Get_Marathon_App_Api    = dcos_Create_Marathon_App_Api + "/%s"
	dcos_Delete_Marathon_App_Api = dcos_Get_Marathon_App_Api
)

type ApplicationsInterface interface {
	Application() ApplicationInterface
}

type ApplicationInterface interface {
	Create(app *App) (*App, error)
	//Get(id string) (*App, error)
	Delete(id string) error
}

type applicationOption struct {
	*dcosOption
}

func (o *applicationOption) Create(app *App) (*App, error) {
	var body, data []byte
	var err error
	if body, err = json.Marshal(app); err != nil {
		return nil, err
	}

	cre_k, cre_v := getCredentials(o.Acs_token)
	data, err = httpPost(o.host+dcos_Create_Marathon_App_Api, ContentType_Json, body, cre_k, cre_v)
	if err != nil {
		return nil, err
	}

	ret := new(App)
	err = json.Unmarshal(data, ret)
	if err != nil {
		return nil, err
	}

	log.Printf("[Info] create app %s success\n", app.Id)
	return ret, nil
}

//
//func (o *applicationOption) Get(id string) (*App, error) {
//	url := o.host + fmt.Sprintf(dcos_Get_Marathon_App_Api, id)
//	cre_k, cre_v := getCredentials(o.Acs_token)
//	data, err := httpGet(url, cre_k, cre_v)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Printf("%s", data)
//	app := new(App)
//	if err := json.Unmarshal(data, app); err != nil {
//		return nil, err
//	}
//	fmt.Printf("------> %v", app)
//	log.Printf("[Info] get app %s success.\n", id)
//	return app, nil
//}

func (o *applicationOption) Delete(id string) error {
	url := o.host + fmt.Sprintf(dcos_Delete_Marathon_App_Api, id)
	cre_k, cre_v := getCredentials(o.Acs_token)
	b, err := httpDelete(url, cre_k, cre_v)
	if err != nil {
		return nil
	}

	log.Printf("[Info] delete app %s success. message %v\n", id, string(b))
	return nil
}
