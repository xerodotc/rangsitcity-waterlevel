package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/xerodotc/rangsitcity-waterlevel/waterlevel"
)

const (
	waterLevelJSONFile = "./data/waterlevel.json"
)

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	waterLevelDataList, err := waterlevel.GetWaterLevelDataWithClient(client)
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
