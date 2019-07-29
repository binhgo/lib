package util

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"lib/Err"
	errors2 "lib/errors"
)

type JsonTime struct {
	time.Time
}

func (this *JsonTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	t, err := time.Parse(time.RFC3339Nano, s[1:len(s)-1])
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z0700", s[1:len(s)-1])

	}
	this.Time = t
	return
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func Contain(target interface{}, obj interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)

	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

func ParseLatLng(latLng string) (lat float64, lng float64, err error) {

	for {

		if latLng == "" {
			break
		}

		latLngArr := strings.Split(latLng, ",")

		lat, err = strconv.ParseFloat(strings.TrimSpace(latLngArr[0]), 10)
		if err != nil {
			break
		}

		lng, err = strconv.ParseFloat(strings.TrimSpace(latLngArr[1]), 10)
		if err != nil {
			break
		}

		return lat, lng, err
	}

	return lat, lng, errors2.NewUserError(Err.ERR_VALIDATE_FOMATER, "latLng split fail")

}
