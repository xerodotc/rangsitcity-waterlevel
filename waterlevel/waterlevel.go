package waterlevel

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const url = "http://rangsit.org/waterlevel/"

func GetWaterLevelData() ([]WaterLevelDataPoint, error) {
	return GetWaterLevelDataWithClient(http.DefaultClient)
}

func GetWaterLevelDataWithClient(client *http.Client) ([]WaterLevelDataPoint, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code " + resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	waterLevelDataList := make([]WaterLevelDataPoint, 0)

	node := doc.Find("#tables > tbody").Find("tr").First()
	for node.Length() > 0 {
		var waterLevelDataPoint WaterLevelDataPoint
		waterLevelText := strings.TrimSpace(node.Find("td:nth-child(2)").Text())
		waterLevelTextPart := strings.Split(waterLevelText, " ")
		waterLevelFloat, err := strconv.ParseFloat(waterLevelTextPart[0], 64)
		if err != nil {
			return nil, err
		}
		waterLevelDataPoint.WaterLevelCM = int(waterLevelFloat * 100)

		recordTimeText := strings.TrimSpace(node.Find("td:nth-child(4)").Text())
		recordTime, err := time.ParseInLocation(waterLevelDataTimeFormat, recordTimeText, time.FixedZone("ICT", 3600*7))
		if err != nil {
			return nil, err
		}
		waterLevelDataPoint.RecordTime = recordTime

		statusImageSrc, ok := node.Find("td:nth-child(5) > img").Attr("src")
		if !ok {
			return nil, errors.New("img tag has no src")
		}
		statusImageFileName := filepath.Base(statusImageSrc)
		statusLevelFileNamePart := strings.Split(statusImageFileName, ".")
		statusLevel, err := strconv.Atoi(statusLevelFileNamePart[0][len("flag"):])
		if err != nil {
			return nil, err
		}
		waterLevelDataPoint.Status = WaterLevelStatus(statusLevel)

		waterLevelDataList = append(waterLevelDataList, waterLevelDataPoint)

		node = node.Next()
	}

	return waterLevelDataList, nil
}
