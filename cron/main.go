package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	fPhoto, err = os.Open(waterLevelPhotoFile)
	if err != nil {
		panic(err)
	}

	defer fPhoto.Close()

	fLatestPhoto, err := os.Create(waterLevelLatestPhotoFile)
	if err != nil {
		panic(err)
	}

	defer fLatestPhoto.Close()

	if _, err := io.Copy(fLatestPhoto, fPhoto); err != nil {
		panic(err)
	}
}
