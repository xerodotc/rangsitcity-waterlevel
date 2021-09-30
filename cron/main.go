package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/xerodotc/rangsitcity-waterlevel/waterlevel"
)

const (
	waterLevelJSONFile        = "./data/waterlevel.json"
	waterLevelLatestPhotoFile = "./data/photos/latest.jpg"
	waterLevelPhotoFileFormat = "./data/photos/%s.jpg"
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

	fJSON.Close()

	waterLevelPhotoData, err := waterlevel.GetWaterLevelPhotoWithClient(client)
	if err != nil {
		fmt.Printf("cannot get water level photo, skipping: %v\n", err)
		return
	}

	waterLevelPhotoFile := fmt.Sprintf(waterLevelPhotoFileFormat, time.Now().Format("20060102_150405"))
	fPhoto, err := os.Create(waterLevelPhotoFile)
	if err != nil {
		fmt.Printf("cannot save water level photo, skipping: %v\n", err)
		return
	}

	defer fPhoto.Close()
	if _, err := fPhoto.Write(waterLevelPhotoData); err != nil {
		fPhoto.Close()
		os.Remove(waterLevelPhotoFile)
		fmt.Printf("cannot download water level photo, skipping: %v\n", err)
		return
	}

	fPhoto.Close()

	if err := os.Remove(waterLevelLatestPhotoFile); err != nil {
		panic(err)
	}

	if err := os.Symlink(path.Base(waterLevelPhotoFile), waterLevelLatestPhotoFile); err != nil {
		panic(err)
	}
}
