package main

import (
	"os"
	"encoding/json"
)

func jsonFileUnMarshal(path string, c interface{})  error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(c)

}
