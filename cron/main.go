package main

import (
	"encoding/json"
	"os"

	"github.com/xerodotc/rangsitcity-waterlevel/waterlevel"
)

const (
	waterLevelJSONFile = "./data/waterlevel.json"
)

func main() {
	waterLevelDataList, err := waterlevel.GetWaterLevelData()
	if err != nil {
		panic(err)
	}

	fJSON, err := os.Create(waterLevelJSONFile)
	if err != nil {
		panic(err)
	}

	defer fJSON.Close()
	if err := json.NewEncoder(fJSON).Encode(waterLevelDataList); err != nil {
		panic(err)
	}
}
