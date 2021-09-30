package waterlevel

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/icholy/digest"
)

const (
	dataURL  = "http://rangsit.org/waterlevel/"
	photoURL = "http://rangsitcity.ddns.net:1029/stw-cgi/video.cgi?msubmenu=stream&action=view&Profile=1"
)

func GetWaterLevelData() ([]WaterLevelDataPoint, error) {
	return GetWaterLevelDataWithClient(http.DefaultClient)
}

func GetWaterLevelDataWithClient(client *http.Client) ([]WaterLevelDataPoint, error) {
	resp, err := client.Get(dataURL)
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

func GetWaterLevelPhoto() ([]byte, error) {
	return GetWaterLevelPhotoWithClient(http.DefaultClient)
}

func GetWaterLevelPhotoWithClient(client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, photoURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return parseJPEGStreaming(resp)
	} else if resp.StatusCode != http.StatusUnauthorized {
		return nil, errors.New("status code " + resp.Status)
	}

	wwwAuthenticate := resp.Header.Get("WWW-Authenticate")
	chal, err := digest.ParseChallenge(wwwAuthenticate)
	if err != nil {
		return nil, err
	}

	cred, err := digest.Digest(chal, digest.Options{
		Username: "guest",
		Password: "guest",
		Method:   req.Method,
		URI:      req.URL.RequestURI(),
		Count:    1,
	})
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", cred.String())

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code " + resp.Status)
	}

	return parseJPEGStreaming(resp)
}

func parseJPEGStreaming(resp *http.Response) ([]byte, error) {
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return io.ReadAll(resp.Body)
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, errors.New("no boundary params")
	}

	multipartReader := multipart.NewReader(resp.Body, boundary)
	defer resp.Body.Close()

	part, err := multipartReader.NextPart()
	if err != nil {
		return nil, err
	}

	defer part.Close()

	partContentType := part.Header.Get("Content-Type")
	if partContentType != "image/jpeg" {
		return nil, errors.New("wrong content type: got " + partContentType)
	}

	var buffer bytes.Buffer
	var smallBuffer [1024]byte

	_, partErr := part.Read(smallBuffer[:])
	for partErr == nil {
		_, err := buffer.Write(smallBuffer[:])
		if err != nil {
			return nil, err
		}
		_, partErr = part.Read(smallBuffer[:])
	}
	if partErr != io.EOF {
		return nil, partErr
	}

	return buffer.Bytes(), nil
}
