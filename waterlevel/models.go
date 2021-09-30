package waterlevel

import "time"

type WaterLevelDataPoint struct {
	WaterLevelCM int              `json:"water_level_cm"`
	Status       WaterLevelStatus `json:"status"`
	RecordTime   time.Time        `json:"record_time"`
}

type WaterLevelStatus int

func (s WaterLevelStatus) String() string {
	switch s {
	case WaterLevelStatusGreen:
		return "green"
	case WaterLevelStatusYellow:
		return "yellow"
	case WaterLevelStatusOrange:
		return "orange"
	case WaterLevelStatusRed:
		return "red"
	}
	return ""
}

func (s WaterLevelStatus) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *WaterLevelStatus) UnmarshalText(text []byte) error {
	t := string(text)
	switch t {
	case "green":
		*s = WaterLevelStatusGreen
	case "yellow":
		*s = WaterLevelStatusYellow
	case "orange":
		*s = WaterLevelStatusOrange
	case "red":
		*s = WaterLevelStatusRed
	default:
		*s = 0
	}
	return nil
}

const (
	WaterLevelStatusGreen  WaterLevelStatus = 1
	WaterLevelStatusYellow WaterLevelStatus = 2
	WaterLevelStatusOrange WaterLevelStatus = 3
	WaterLevelStatusRed    WaterLevelStatus = 4
)

const waterLevelDataTimeFormat = "2006-01-02 15:04:05"
