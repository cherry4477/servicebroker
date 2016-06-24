package main

import (
	"github.com/asiainfoLDP/servicebroker_dcos/servicebroker"
	"io"
	"net/http"
	"log"
)

var catalogInfo = new(catalog)

func init() {
	if err := jsonFileUnMarshal("catalog.json", catalogInfo); err != nil  {
		log.Fatalf("init catalog config err %v\n", err)
	}

	if len(catalogInfo.Services) == 0 {
		log.Fatal("load catalog config nil.")
	}
}

type catalog servicebroker.Catalog

func (s *catalog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL"+r.URL.String())
}
