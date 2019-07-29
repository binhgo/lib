package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/geo/s2"
	"github.com/json-iterator/go"
	"github.com/spf13/viper"
	"io/ioutil"
	"lib/Err"
	"lib/errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type DistanceStruct struct {
	Distance    float64 //距离 单位米
	DistanceKm  float64 //距离 单位千米
	Duration    float64 //时间 单位分
	DurationMin int     //时间 单位秒
}

type AddressComponent struct {
	LongName  string
	ShortName string
	Types     []string
}

type PlaceSearchResultStruct struct {
	Description       string //位置描述
	LatLng            string //坐标经纬度
	MainText          string
	PlaceId           string //google 位置唯一标示
	SecondaryText     string
	CityId            int    //西游计对应城市ID，不存在为0
	CityName          string //西游计对应城市名称,不存在默认为MainText
	NearCityId        int
	NearCityName      string
	ContinentId       int
	ContinentName     string
	CountryId         int
	CountryName       string
	AddressComponents []AddressComponent `json:"placeAddressComponents"`
}

func GetDistance(startLocation string, endLocation string) (distanceStruct DistanceStruct, err error) {

	distanceKm := float64(0)
	url := "http://150.109.62.175/maps/direction?origin=" + startLocation + "&destination=" + endLocation;
	defer func() {
		if r := recover(); r != nil {
			err = errors.NewUserError(Err.ERR_HTTP_REQEUST, "map request error")
			distanceKm = 0
		}

	}()

	//下发APP通知
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return distanceStruct, errors.NewUserError(Err.ERR_PRICE_CARTYPE_EMPTY, fmt.Sprintf("request map err(%s),url(%s)", err.Error(), url))
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	code := jsoniter.Get(body, "code").ToInt()
	//should.Nil()
	if code == 200 {
		distanceStruct.DistanceKm = jsoniter.Get(body, "data", "distanceKm").ToFloat64()
		distanceStruct.Distance = jsoniter.Get(body, "data", "distance").ToFloat64()
		distanceStruct.Duration = jsoniter.Get(body, "data", "duration").ToFloat64()
		distanceStruct.DurationMin = jsoniter.Get(body, "data", "durationMin").ToInt()
	} else {
		distanceKm = 0
		err = errors.NewUserError(Err.ERR_HTTP_REQEUST, "map request error")
	}

	return distanceStruct, err
}

func SearchPlace(input string, location string, radius int) (placeSearchResultStructArr []PlaceSearchResultStruct, err error) {

	type Params struct {
		Input        string `json:"input"`        //搜索内容
		Location     string `json:"location"`     //中心点,纬度,经度
		StrictBounds bool   `json:"strictBounds"` //限制中心点和搜索范围
		Radius       int    `json:"radius"`       //搜索地理范围
		Types        string `json:"types"`        //搜索类型
	}

	strictBounds := true
	if len(location) == 0 {
		strictBounds = false
	}
	pushParamsVal, _ := json.Marshal(&Params{input, location, strictBounds, radius, "GEOCODE"})
	url := viper.GetString("google.map.proxy_api_url")

	fmt.Printf(string(pushParamsVal), url)

	defer func() {
		if r := recover(); r != nil {
			err = errors.NewUserError(Err.ERR_HTTP_REQEUST, "map request error")
		}

	}()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pushParamsVal))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return placeSearchResultStructArr, errors.NewUserError(Err.ERR_PRICE_CARTYPE_EMPTY, fmt.Sprintf("request map err(%s),url(%s)", err.Error(), url))
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	code := jsoniter.Get(body, "code").ToInt()

	if code == 200 {
		jsoniter.Unmarshal([]byte(jsoniter.Get(body, "data").ToString()), &placeSearchResultStructArr)
	} else {
		err = errors.NewUserError(Err.ERR_HTTP_REQEUST, "map request error")
	}

	return placeSearchResultStructArr, err
}

type CityAreas struct {
	areas map[int]*s2.Polygon
}

func (this *CityAreas) AddRect(areaId int, points []s2.Point) {

	if this.areas == nil {
		this.areas = make(map[int]*s2.Polygon)
	}

	loops := []*s2.Loop{}
	loops = append(loops, s2.LoopFromPoints(points))
	this.areas[areaId] = s2.PolygonFromLoops(loops)
}

func (this *CityAreas) PoiInclusion(poi s2.Point) (int, error) {
	var areaId int
	var err error

	//解决mapfor循环无序问题，以及判断区域从小到大
	var keys []int
	for k, _ := range this.areas {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, id := range keys {
		if this.areas[id].ContainsPoint(poi) {
			areaId = id
			break
		}
	}

	if areaId == 0 {
		err = errors.NewUserError(Err.ERROR, "")
	}
	return areaId, err
}

func AreaExists(area string, poi s2.Point) (bool, error) {

	cityAreas := CityAreas{}

	if area == "" || !strings.Contains(area, ";") {
		return false, errors.NewUserError(Err.ERR_PRICE_POI_OUTBOUNDCITY, "address unvalidate", "超出服务范围")
	}

	points := []s2.Point{}
	poiStrs := strings.Split(area, ";")
	for _, poiStr := range poiStrs {
		poiArr := strings.Split(poiStr, ",")
		if len(poiArr) == 2 {
			X, err := strconv.ParseFloat(strings.TrimSpace(poiArr[0]), 10)
			if err != nil {
				panic("poi err")
			}
			Y, err := strconv.ParseFloat(strings.TrimSpace(poiArr[1]), 10)
			if err != nil {
				panic("poi err")
			}
			points = append(points, s2.PointFromLatLng(s2.LatLngFromDegrees(X, Y)))
		} else {
			return false, errors.NewUserError(Err.ERR_PRICE_POI_OUTBOUNDCITY, "address unvalidate", "超出服务范围")
		}
	}
	cityAreas.AddRect(1, points)

	//计算服务时间点所在区域
	cityAreaId, err := cityAreas.PoiInclusion(poi)

	if err != nil || cityAreaId != 1 {
		return false, errors.NewUserError(Err.ERR_PRICE_POI_OUTBOUNDCITY, "address unvalidate", "超出服务范围")
	} else {
		return true, nil
	}

}
