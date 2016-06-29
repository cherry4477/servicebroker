package main

import (
	broker "github.com/asiainfoLDP/servicebroker_dcos/servicebroker"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//curl 127.0.0.1:5000/v2/catalog
func catalogHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	l, err := getCatalog()
	if err != nil {
		log.Println("[Get] /v2/catalog err %v\n", err)
	}

	respond(w, http.StatusOK, l)
}

func getCatalog() (*broker.Catalog, error) {
	c := new(broker.Catalog)
	if err := jsonFileUnMarshal(Catalog_Info_Path, c); err != nil {
		return nil, err
	}

	return c, nil
}

//1 Cpu;1000 Mem;50 Disk
type planCostUnit struct {
	Cpu  float32
	Mem  uint32
	Disk uint32
}

func parsePlanUnit(units string) *planCostUnit {
	if len(units) == 0 {
		return nil
	}

	m := new(planCostUnit)
	s := strings.SplitN(units, ";", 3)
	for _, substr := range s {
		if strings.Contains(strings.ToLower(substr), "mem") {
			m.Mem = uint32(extractNumStr(substr))
			continue
		}

		if strings.Contains(strings.ToLower(substr), "cpu") {
			m.Cpu = float32(extractNumStr(substr))
			continue
		}

		if strings.Contains(strings.ToLower(substr), "disk") {
			m.Disk = uint32(extractNumStr(substr))
			continue
		}
	}

	return m
}

func extractNumStr(str string) float64 {
	if len(str) == 0 {
		return 0
	}

	rgp := regexp.MustCompile(`\d+(\.\d+)?`)
	s := rgp.FindString(str)
	num, _ := strconv.ParseFloat(s, 64)

	return num
}
