package main

import (
	"fmt"
	"github.com/asiainfoLDP/servicebroker_dcos/servicebroker"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

//curl 127.0.0.1:5000/v2/catalog
func catalogHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(Catalog_Info_Path)
	if err != nil {
		log.Printf("read catalog err %v\n", err)
		io.WriteString(w, fmt.Sprint(servicebroker.NewServiceProviderError(400, err)))
		return
	}

	io.WriteString(w, string(data))
	w.WriteHeader(200)
	return
}

func getCatalog() (*servicebroker.Catalog, error) {
	c := new(servicebroker.Catalog)
	if err := jsonFileUnMarshal(Catalog_Info_Path, c); err != nil {
		return nil, err
	}

	return c, nil
}
