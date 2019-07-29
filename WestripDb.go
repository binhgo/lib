package lib

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"lib/Err"
	"lib/dao"
	"lib/errors"
)

var (
	westripDbOnce sync.Once
	westripDb     *WestripDb
)

type WestripDb struct {
	db             *sqlx.DB
	DataSourceName string
}

func NewWestripDb() (*WestripDb) {
	dataSourceName := "westrip_db_connect_info"
	westripDbOnce.Do(func() {
		db, err := sqlx.Open("mysql", viper.GetString(dataSourceName))
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		db.SetMaxIdleConns(50)
		db.SetMaxOpenConns(200)
		db.SetConnMaxLifetime(time.Duration(1) * time.Hour)
		westripDb = &WestripDb{db, dataSourceName}
	})

	return westripDb
}

func (this *WestripDb) Destroy() {
	westripDbOnce = sync.Once{}
	this.db.Close()
}

func (this *WestripDb) QueryCarTypes(args ...interface{}) ([]dao.DbCarTypesRow, error) {

	sql := "SELECT ct.luggage_num,ct.guest_num,ct.id as car_type_id,ct.car_type_name,ct.introduction,ct.pictures,ct.mapping_car_type,ct.mapping_seat_type,ct.seat_type "
	sql = sql + " FROM basedata_car_type ct "
	sql = sql + " WHERE ct.guest_num >= ? AND ct.status = 1"
	var dbCarInfoRow []dao.DbCarTypesRow
	err := this.db.Select(&dbCarInfoRow, sql, args[0])

	if err != nil {
		logrus.Info(err)
	}

	return dbCarInfoRow, err
}

func (this *WestripDb) QueryCarInfos(carTypeId int) ([]dao.DBCarRow, error) {

	re := regexp.MustCompile("([1-9])?0*([0-9])?")
	match := re.FindStringSubmatch(strconv.Itoa(carTypeId))

	sql := "SELECT id,car_brand_id,car_brand_name,car_model "
	sql += " FROM basedata_car "
	sql += " WHERE `car_type` = ? AND `seat_type` = ? AND  `status` = 1 ORDER BY price_top desc LIMIT 1"

	var dbCarRows []dao.DBCarRow
	err := this.db.Select(&dbCarRows, sql, match[1], match[2])

	if err != nil {
		Log.Info(err)
	}

	return dbCarRows, err
}

func (this *WestripDb) QueryCurrencyByCityId(cityId int) (dao.DbCurrencyRow, error) {

	var dbCurrencyRow []dao.DbCurrencyRow
	err := this.db.Select(&dbCurrencyRow, "SELECT bcs.id as cityId,bcy.currency,bcy.rate as currencyRate,bcs.`timezone` FROM basedata_currency bcy left join basedata_cities bcs  on bcy.country_id = bcs.country_id where bcs.id = ?", cityId)

	if err != nil {
		logrus.Info(err)
	}

	return dbCurrencyRow[0], err
}

func (this *WestripDb) GetGuideData(carTypeId int, serviceCityId int) ([]dao.GuideRowStruct, error) {

	re := regexp.MustCompile("([1-9])?0*([0-9]{1,2})?")
	match := re.FindStringSubmatch(strconv.Itoa(carTypeId))

	guideRows := []dao.GuideRowStruct{}
	// err := this.db.Select(&guideRows, "SELECT distinct ggsc.guide_id FROM gcenter_guide_service_city ggsc LEFT JOIN gcenter_guide_car gg ON gg.guide_id = ggsc.guide_id AND gg.audit_status = 3 AND ggsc.status=1  WHERE service_city_id = ? AND ggsc.status=1 AND (ggsc.guide_id IN (SELECT guide_id FROM gcenter_guide_car WHERE seating_capacity>=? AND car_type_id >= ? AND audit_status = 3) or ggsc.guide_id IN (SELECT id AS guide_id FROM gcenter_guide WHERE role_id = 21 AND status=1)) AND ggsc.guide_id IN (SELECT id AS guide_id FROM gcenter_guide WHERE status= 1)", serviceCityId, match[2], carTypeId)
	// err := this.db.Select(&guideRows, "SELECT distinct ggsc.guide_id FROM gcenter_guide_service_city ggsc LEFT JOIN gcenter_guide_car gg ON gg.guide_id = ggsc.guide_id AND gg.audit_status = 3 AND ggsc.status=1  WHERE service_city_id = ? AND ggsc.status=1 AND (ggsc.guide_id IN (SELECT guide_id FROM gcenter_guide_car WHERE seating_capacity>=? AND car_type_id >= ? AND audit_status = 3) or ggsc.guide_id IN (SELECT id AS guide_id FROM gcenter_guide WHERE role_id = 21 AND status=1))", serviceCityId, match[2], carTypeId)
	err := this.db.Select(&guideRows, "SELECT distinct ggsc.guide_id FROM gcenter_guide_service_city ggsc LEFT JOIN gcenter_guide_car gg ON gg.guide_id = ggsc.guide_id AND ggsc.status=1  WHERE service_city_id = ? AND ggsc.status=1 AND (ggsc.guide_id IN (SELECT guide_id FROM gcenter_guide_car WHERE seating_capacity>=? AND car_type_id >= ? AND audit_status = 3)) AND (ggsc.guide_id IN (SELECT id as guide_id FROM gcenter_guide where status = 1 AND audit_status = 3))", serviceCityId, match[2], carTypeId)
	if err != nil {
		err = errors.NewUserError(Err.ERR_DISPATCH_NO_GUIDES, "guide empty"+err.Error())
	}

	return guideRows, err

}

// 根据车型获取司导可用车辆
func (this *WestripDb) GetGuideCarListData(guideId int, carTypeId int) ([]dao.GuideCarRowStruct, error) {

	re := regexp.MustCompile("([1-9])?0*([0-9]{1,2})?")
	match := re.FindStringSubmatch(strconv.Itoa(carTypeId))

	guideCarRows := []dao.GuideCarRowStruct{}
	err := this.db.Select(&guideCarRows, "SELECT ggc.id as guide_car_id,ggc.guide_id,ggc.brand_id,ggc.brand_name,ggc.car_model_id,ggc.car_model_name,ggc.car_type_id,ggc.car_type_name,ggc.plate_num,gg.mobile as guide_mobile,gg.name as guide_name FROM gcenter_guide_car ggc INNER JOIN gcenter_guide gg ON ggc.guide_id = gg.id AND gg.audit_status = 3 AND gg.status = 1  WHERE ggc.guide_id = ? AND ggc.audit_status = 3 AND  ggc.seating_capacity>=? AND ggc.car_type_id >= ?", guideId, match[2], carTypeId)

	if err != nil {
		err = errors.NewUserError(Err.ERR_DISPATCH_NO_GUIDES, "guide empty"+err.Error())
	}

	return guideCarRows, err

}

// 根据车型获取司导可用车辆
func (this *WestripDb) GetGuideCarData(guideId int, guideCarId int) (dao.GuideCarRowStruct, error) {

	guideCarRows := []dao.GuideCarRowStruct{}
	err := this.db.Select(&guideCarRows, "SELECT ggc.id as guide_car_id,ggc.guide_id,ggc.brand_id,ggc.brand_name,ggc.car_model_id,ggc.car_model_name,ggc.car_type_id,ggc.car_type_name,ggc.plate_num,gg.mobile as guide_mobile,gg.name as guide_name FROM gcenter_guide_car ggc INNER JOIN gcenter_guide gg ON ggc.guide_id = gg.id AND gg.audit_status = 3 AND gg.status = 1  WHERE ggc.guide_id = ? AND ggc.audit_status = 3 AND ggc.id = ?", guideId, guideCarId)

	if err != nil || len(guideCarRows) == 0 {
		return dao.GuideCarRowStruct{}, errors.NewUserError(Err.ERR_DISPATCH_NO_GUIDES, "guide empty")
	} else {
		return guideCarRows[0], err
	}

}

func (this *WestripDb) GetGuideInfo(guideId int) (dao.GuideInfoRowStruct, error) {

	guideInfoRows := []dao.GuideInfoRowStruct{}
	err := this.db.Select(&guideInfoRows, "SELECT id,mobile,`name` FROM gcenter_guide  WHERE id=? ", guideId)

	if err != nil {
		logrus.Info(err)
	}

	return guideInfoRows[0], err
}

func (this *WestripDb) GetOrderInfo(orderId int64) (dao.OrderRowData, error) {

	orderRowData := []dao.OrderRowData{}

	err := this.db.Select(&orderRowData, "SELECT order_id,car_type_id,car_type_name,service_city_id FROM  transaction_order WHERE order_id=?", orderId)

	if err != nil {
		logrus.Info(err)
	}

	return orderRowData[0], err
}

// 根据机场获取机场信息
func (this *WestripDb) GetAirportInfos(airportCode string) ([]dao.AirportRowStruct, error) {

	airportRows := []dao.AirportRowStruct{}
	err := this.db.Select(&airportRows, "SELECT id,country_id,country_name,city_id,city_name,code,`name`,location FROM basedata_airport  WHERE code=?;", strings.ToUpper(airportCode))

	if err != nil {
		logrus.Info(err)
	}

	return airportRows, err
}

// 根据机场获取机场信息
func (this *WestripDb) GetCartypeInfo(cartypeId int) (dao.CartypeRowStruct, error) {

	cartypeRows := []dao.CartypeRowStruct{}
	err := this.db.Select(&cartypeRows, "SELECT `id`,car_type_name FROM basedata_car_type  WHERE id=? AND status = 1;", cartypeId)

	if err != nil {
		logrus.Info(err)
	}

	if len(cartypeRows) != 1 {
		panic(errors.NewUserError(Err.ERROR, "cartypeinfo error"))
	}

	return cartypeRows[0], err
}

// 查找城市信息
func (this *WestripDb) QueryCityInfo(cityIds []int) ([]dao.DbCitiesRow, error) {

	sql := "SELECT id,city_name,country_id,country_name,`location` FROM  basedata_cities "
	sql += "WHERE id in(?) AND status = 1;"

	var query string
	var args []interface{}
	var err error

	query, args, err = sqlx.In(sql, cityIds)

	if err != nil {
		Log.Info(err)
	}

	query = this.db.Rebind(query)

	var dbCitiesRow []dao.DbCitiesRow
	rows, err := this.db.Queryx(query, args...)

	for rows.Next() {
		row := dao.DbCitiesRow{}
		err = rows.StructScan(&row)
		dbCitiesRow = append(dbCitiesRow, row)
	}

	if err != nil {
		logrus.Info(err)
	}

	if len(dbCitiesRow) == 0 {
		return []dao.DbCitiesRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	if err != nil {
		logrus.Info(err)
	}

	return dbCitiesRow, err
}

// 根据携程城市ID获取城市信息
func (this *WestripDb) QueryCityInfoByCtripCityid(cityIds int) ([]dao.DbCitiesRow, error) {

	sql := "SELECT id,city_name,country_id,country_name,`location` FROM  basedata_cities "
	sql += "WHERE ctrip_city_id =? AND status = 1;"

	var query string
	var args []interface{}
	var err error

	query, args, err = sqlx.In(sql, cityIds)

	if err != nil {
		Log.Info(err)
	}

	query = this.db.Rebind(query)

	var dbCitiesRow []dao.DbCitiesRow
	rows, err := this.db.Queryx(query, args...)

	for rows.Next() {
		row := dao.DbCitiesRow{}
		err = rows.StructScan(&row)
		dbCitiesRow = append(dbCitiesRow, row)
	}

	if err != nil {
		logrus.Info(err)
	}

	if len(dbCitiesRow) == 0 {
		return []dao.DbCitiesRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	if err != nil {
		logrus.Info(err)
	}

	return dbCitiesRow, err
}

// 根据城市名称获取城市信息
func (this *WestripDb) QueryCityInfoByLikeName(likeNames []string) ([]dao.DbCitiesRow, error) {

	sql := "SELECT id,city_name,country_id,country_name,continent_id,continent_name,`location` FROM  basedata_cities "
	sql += "WHERE status = 1 AND city_name in (?)"

	var query string
	var args []interface{}
	var err error

	query, args, err = sqlx.In(sql, likeNames)

	if err != nil {
		Log.Info(err)
	}

	query = this.db.Rebind(query)

	var dbCitiesRow []dao.DbCitiesRow
	rows, err := this.db.Queryx(query, args...)

	for rows.Next() {
		row := dao.DbCitiesRow{}
		err = rows.StructScan(&row)
		dbCitiesRow = append(dbCitiesRow, row)
	}

	if err != nil {
		logrus.Info(err)
	}

	if len(dbCitiesRow) == 0 {
		return []dao.DbCitiesRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	if err != nil {
		logrus.Info(err)
	}

	return dbCitiesRow, err
}

// 根据国家名国家信息
func (this *WestripDb) QueryCountryInfoByName(countryName string) (countriesRow dao.CountriesRowStruct, err error) {

	sql := "SELECT id,country_name,continent_id,continent_name,code FROM  basedata_countries "
	sql += "WHERE  country_name = ?"

	countriesRows := []dao.CountriesRowStruct{}
	err = this.db.Select(&countriesRows, sql, countryName)

	if err != nil {
		logrus.Info(err)
	}

	if len(countriesRows) != 1 {
		return countriesRow, errors.NewUserError(Err.ERROR, "countriesRow error")
	}

	return countriesRows[0], err
}
