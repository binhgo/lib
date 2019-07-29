package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/spf13/viper"
	"io/ioutil"
	"lib/Err"
	"lib/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TimeConvert struct {
}

/**
"bjTime":"2018-02-24 18:14:50"
 */
func (this *TimeConvert) ConvertToBeijingTime(cityId int, beijingTime string) (delaySeconds int, err error) {
	type Params struct {
		CityId    int    `json:"cityId"`
		BjTime    string `json:"bjTime"`
		TimeRange string `json:"timeRange"`
	}

	pushParamsVal, _ := json.Marshal(&Params{cityId, beijingTime, viper.GetString("work_range")})

	fmt.Printf(string(pushParamsVal))
	//下发APP通知
	req, err := http.NewRequest("POST", viper.GetString("time_get_url"), bytes.NewBuffer(pushParamsVal))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	fmt.Println("response body:", bodyStr)

	type Response struct {
		Code int    `json:"code"`
		Desc string `json:"desc"`
		Data struct {
			Start int `json:"start"`
			End   int `json:"end"`
		}
	}
	response := Response{}
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &response)
	if err != nil {
		return delaySeconds, err
	} else if (response.Code == 200) {
		if err != nil {
			return delaySeconds, err
		} else {
			return response.Data.Start, nil
		}

	} else {
		return delaySeconds, errors.NewUserError(Err.ERR_HTTP_REQEUST, bodyStr)
	}
}

/**
此时JsonTime应为srcTimezone时区对应的时间
返回结果：返回JsonTime对应targetTimezone时区的时间
 */
func (this *TimeConvert) GetRunStartTime(bjTime time.Time, cityId int) (time.Time, error) {
	//获得服务城市当前时间
	delaySecondes, err := this.ConvertToBeijingTime(cityId, bjTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return time.Time{}, err
	}
	runStartTime := bjTime.Add(time.Duration(delaySecondes) * time.Second) //将服务城市时间转为北京时间

	fmt.Println(runStartTime,bjTime)

	//当地休息时间,此范围内
	restHoursArr := viper.GetStringSlice("rest_time")
	restHours := []int{}
	for _, v := range restHoursArr {
		h, _ := strconv.Atoi(v)
		restHours = append(restHours, h)
	}

	isRest := false
	for _, v := range restHours {
		if v == runStartTime.Hour() {
			isRest = true
			break
		}
	}

	if !isRest {
		runStartTime = bjTime
	}

	return runStartTime, err
}

/**
"bjTime":"2018-02-24 18:14:50"
 */
func (this *TimeConvert) GetLocaltime(cityId int, beijingTime time.Time) (localtime time.Time, err error) {
	type Params struct {
		CityId int    `json:"cityId"`
		BjTime string `json:"bjTime"`
	}

	pushParamsVal, _ := json.Marshal(&Params{cityId, beijingTime.Format("2006-01-02 15:04:05")})

	fmt.Printf(string(pushParamsVal), viper.GetString("localtime_get_url"))
	//下发APP通知
	req, err := http.NewRequest("POST", viper.GetString("localtime_get_url"), bytes.NewBuffer(pushParamsVal))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	fmt.Println("response body:", bodyStr)

	type Response struct {
		Code int    `json:"code"`
		Desc string `json:"desc"`
		Data string `json:"data"`
	}

	response := Response{}
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &response)
	if err != nil {
		return localtime, err
	} else if (response.Code == 200) {
		if err != nil {
			return localtime, err
		} else {
			t, err := time.Parse("2006-01-02 15:04:05", response.Data)
			if err != nil {
				return localtime, err
			} else {
				return t, nil
			}
		}
	} else {
		return localtime, errors.NewUserError(Err.ERR_HTTP_REQEUST, bodyStr)
	}
}

/*
判断该时间是否为夜间时间范围内
nightStartTimeStr 夜间费用开时间
nightEndTimeStr 夜间费用结束时间
true 在范围内 false 不在范围内
 */
func (this *TimeConvert) IsNightShift(nightStartTimeStr string, nightEndTimeStr string, serviceTime time.Time) bool {

	year := serviceTime.Year()
	month := serviceTime.Month()
	day := serviceTime.Day()

	nightShiftStartTimeArr := strings.Split(nightStartTimeStr, ":")
	nightShiftEndTimeArr := strings.Split(nightEndTimeStr, ":")

	startHour, err := strconv.Atoi(nightShiftStartTimeArr[0])
	if err != nil {
		return false
	}
	startMinute, err := strconv.Atoi(nightShiftStartTimeArr[1])
	if err != nil {
		return false
	}

	endHour, err := strconv.Atoi(nightShiftEndTimeArr[0])
	if err != nil {
		return false
	}
	endMinute, err := strconv.Atoi(nightShiftEndTimeArr[1])
	if err != nil {
		return false
	}

	if startHour <= endHour {
		panic(errors.NewUserError(Err.ERROR, "night shifth params unvalidate"))
	}

	nightShiftStartTime := time.Date(year, month, day, startHour, startMinute, 0, 0, time.UTC)
	nightShiftStartTimeMax := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)

	nightShiftEndTime := time.Date(year, month, day, endHour, endMinute, 0, 0, time.UTC)
	nightShiftEndTimeMin := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	if serviceTime.Equal(nightShiftStartTime) || serviceTime.Equal(nightShiftStartTimeMax) || serviceTime.Equal(nightShiftEndTime) || serviceTime.Equal(nightShiftEndTimeMin) {
		return true
	}

	if (serviceTime.After(nightShiftStartTime) && serviceTime.Before(nightShiftStartTimeMax)) {
		return true
	}

	if (serviceTime.After(nightShiftEndTimeMin) && serviceTime.Before(nightShiftEndTime)) {
		return true
	}

	return false
}

func (this *TimeConvert) IsSeasonPeriod(seasonPeriods map[string]float64, serviceTime time.Time) (float64,error) {
	if len(seasonPeriods) == 0 {
		return  float64(0),nil
	}

	for seasonPeriod, seasonPeriodShift := range seasonPeriods {
		seasonPeriodArr := strings.Split(seasonPeriod, ":")
		if len(seasonPeriodArr) != 2 {
			continue
		}

		seasonPeriodStart, err := time.Parse("2006-01-02 15:04:00", seasonPeriodArr[0]+" 00:00:00")
		if err != nil {
			continue
		}

		seasonPeriodEnd, err := time.Parse("2006-01-02 15:04:05", seasonPeriodArr[1]+" 23:59:59")
		if err != nil {
			continue
		}


		if serviceTime.Before(seasonPeriodEnd) && serviceTime.After(seasonPeriodStart) {
			return seasonPeriodShift,nil
		}

	}

	return float64(0),nil
}
