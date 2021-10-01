# Machine-readable Rangsit City Water Level Report

JSON formatted data and camera footages for Rangsit City Water Level report data obtained from http://rangsit.org/waterlevel/.

## Disclaimers

This project is **NOT** related **NOR** affiliated with Rangsit City Municipality.

โปรเจคนี้ **ไม่มี** ความเกี่ยวข้องใด ๆ กับเทศบาลนครรังสิต

## Usage

* JSON file URL:
    * `https://raw.githubusercontent.com/xerodotc/rangsitcity-waterlevel/main/data/waterlevel.json`
* Camera footages:
    * Latest photo: https://raw.githubusercontent.com/xerodotc/rangsitcity-waterlevel/main/data/photos/latest.jpg
    * All photos: https://github.com/xerodotc/rangsitcity-waterlevel/tree/main/data/photos
* Data will be updated every 3 hours.

### Data format

JSON file contains an array of data points similar to these:
```json
[
    {
        "water_level_cm": 110,
        "status": "green",
        "record_time": "2021-09-27T08:00:00+07:00"
    },
    {
        "water_level_cm": 117,
        "status": "green",
        "record_time":"2021-09-27T16:30:00+07:00"
    }
]
```

For each data point, there are fields in the object, which are:
* `water_level_cm`: Water level in centimeters.
* `status`: Flag status, can be `green`, `yellow`, `orange`, `red`.
* `record_time`: Data point recorded time in RFC3339 format.

Data points are ordered from oldest to newest data point.
