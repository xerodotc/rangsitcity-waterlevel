package waterlevel

import "time"

type WaterLevelDataPoint struct {
	WaterLevelCM int              `json:"water_level_cm"`
	Status       WaterLevelStatus `json:"status"`
	RecordTime   time.Time        `json:"record_time"`
}

type WaterLevelStatus int

const (
	WaterLevelStatusGreen  WaterLevelStatus = 1
	WaterLevelStatusYellow WaterLevelStatus = 2
	WaterLevelStatusOrange WaterLevelStatus = 3
	WaterLevelStatusRed    WaterLevelStatus = 4
)

const waterLevelDataTimeFormat = "2006-01-02 15:04:05"
