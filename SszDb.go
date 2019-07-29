package lib

import (
	"fmt"
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
	sszDbOnce sync.Once
	sszDb     *SszDb
)

type SszDb struct {
	db             *sqlx.DB
	DataSourceName string
}

func NewSszDb() (*SszDb) {
	dataSourceName := "shensuanzi_db_connect_info"
	sszDbOnce.Do(func() {
		db, err := sqlx.Open("mysql", viper.GetString(dataSourceName))
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		db.SetMaxIdleConns(50)
		db.SetMaxOpenConns(200)
		db.SetConnMaxLifetime(time.Duration(1) * time.Hour)
		sszDb = &SszDb{db, dataSourceName}
	})

	return sszDb
}

func (this *SszDb) Destroy() {
	sszDbOnce = sync.Once{}
	this.db.Close()
}

func (this *SszDb) QueryCarTypes(args ...interface{}) ([]dao.DbCarTypesRow, error) {

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

func (this *SszDb) QueryCarInfos(carTypeId int) ([]dao.DBCarRow, error) {

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

func (this *SszDb) QueryCurrencyByCityId(cityId int) (dao.DbCurrencyRow, error) {
	var dbCurrencyRow []dao.DbCurrencyRow
	err := this.db.Select(&dbCurrencyRow, "SELECT bcs.id as cityId,bcy.currency,bcy.rate as currencyRate,bcs.`timezone` FROM basedata_currency bcy left join basedata_cities bcs  on bcy.country_id = bcs.country_id where bcs.id = ?", cityId)

	if err != nil {
		logrus.Info(err)
	}

	return dbCurrencyRow[0], err
}

func (this *SszDb) QueryPriceBasicRuleByCityId(cityId int) (dao.DbPriceBasicRuleRow, error) {
	var dbPriceBasicRuleRow []dao.DbPriceBasicRuleRow
	err := this.db.Select(&dbPriceBasicRuleRow, "SELECT id, city_id, city_name, country_id,country_name,checkin_price, welcome_board_price,childseat_price, night_shift, night_shift_start_time, night_shift_end_time,season_shift_1, season_shift_periods_1, season_shift_2, season_shift_periods_2, season_shift_3, season_shift_periods_3, season_shift_4, season_shift_periods_4,season_daliy_shift_1, season_daliy_shift_periods_1, season_daliy_shift_2, season_daliy_shift_periods_2, season_daliy_shift_3, season_daliy_shift_periods_3, season_daliy_shift_4, season_daliy_shift_periods_4, user_profit, user_shift, guide_profit, guide_shift, urgent_shift,accommodation_fee, urgent_pupdoff_limit_hour, urgent_daliy_limit_hour, booking_pupdoff_limit_hour, booking_daliy_limit_hour, car_unit_parkingfee, service_timeoutfee, city_area_1, city_area_2, city_area_3, city_area_4, city_area_5, city_area_6, city_area_7, city_area_8, city_area_9, city_area_10 FROM price_basic_rule where city_id = ? AND status = 0", cityId)

	if err != nil {
		logrus.Info(err)
	}

	return dbPriceBasicRuleRow[0], err
}

func (this *SszDb) QueryPriceBasicRuleByCountryId(countryId int) ([]dao.DbPriceBasicRuleRow, error) {
	var dbPriceBasicRuleRow []dao.DbPriceBasicRuleRow
	err := this.db.Select(&dbPriceBasicRuleRow, "SELECT id, city_id, city_name, country_id,country_name,checkin_price, welcome_board_price,childseat_price, night_shift, night_shift_start_time, night_shift_end_time,season_shift_1, season_shift_periods_1, season_shift_2, season_shift_periods_2, season_shift_3, season_shift_periods_3, season_shift_4, season_shift_periods_4, season_daliy_shift_1, season_daliy_shift_periods_1, season_daliy_shift_2, season_daliy_shift_periods_2, season_daliy_shift_3, season_daliy_shift_periods_3, season_daliy_shift_4, season_daliy_shift_periods_4,user_profit, user_shift, guide_profit, guide_shift, urgent_shift,accommodation_fee, urgent_pupdoff_limit_hour, urgent_daliy_limit_hour, booking_pupdoff_limit_hour, booking_daliy_limit_hour, car_unit_parkingfee, service_timeoutfee, city_area_1, city_area_2, city_area_3, city_area_4, city_area_5, city_area_6, city_area_7, city_area_8, city_area_9, city_area_10 FROM price_basic_rule where (country_id = ? or 1=1) AND status = 0", countryId)

	if err != nil {
		logrus.Info(err)
	}

	return dbPriceBasicRuleRow, err
}

func (this *SszDb) QueryCityAreasByCityId(cityId int) ([]dao.DbCityAreaRow, error) {
	var dbCityAreaRow []dao.DbCityAreaRow
	err := this.db.Select(&dbCityAreaRow, "SELECT city_area_1,city_area_2,city_area_3,city_area_4,city_area_5,city_area_6,city_area_7,city_area_8,city_area_9,city_area_10 FROM price_basic_rule WHERE city_id = ?", cityId)

	if err != nil {
		logrus.Info(err)
	}

	return dbCityAreaRow, err
}

func (this *SszDb) QueryPupdoff(cityId int, airportCode string, cityAreaId int) ([]dao.DbPupdoffRow, error) {
	sql := "SELECT car_type_id,car_type_name,city_area_id,price "
	sql += "FROM price_pupdoff_rule "
	sql += "WHERE city_id = ? AND airport_code = ? AND city_area_id = ?  AND status = 0;"
	var dbPupdoffRow []dao.DbPupdoffRow
	err := this.db.Select(&dbPupdoffRow, sql, cityId, airportCode, cityAreaId)

	if err != nil {
		logrus.Info(err)
	}

	return dbPupdoffRow, err
}

func (this *SszDb) QueryPupdoffOne(cityId int, airportCode string, cityAreaId int, carTypeId int) (dao.DbPupdoffRow, error) {
	sql := "SELECT car_type_id,car_type_name,city_area_id,price "
	sql += "FROM price_pupdoff_rule "
	sql += "WHERE city_id = ? AND airport_code = ? AND city_area_id = ? AND car_type_id = ?  AND status = 0;"
	var dbPupdoffRow []dao.DbPupdoffRow
	err := this.db.Select(&dbPupdoffRow, sql, cityId, airportCode, cityAreaId, carTypeId)

	if err != nil {
		logrus.Info(err)
	}

	if len(dbPupdoffRow) != 1 {
		return dao.DbPupdoffRow{}, errors.NewUserError(Err.ERROR, "QueryPupdoffOne empty")
	}

	return dbPupdoffRow[0], err
}

func (this *SszDb) QueryPriceChannelRule(cityId int, serviceType int, channelType int) (dao.DbPriceChannelRuleRow, error) {
	sql := "SELECT id, city_id, city_name, service_type, channel_type, channel_profit, activity_shift "
	sql += "FROM price_channel_rule  "
	sql += "WHERE city_id = ? AND channel_type = ? AND service_type = ?  AND status = 0;"

	var dbPupdoffRow []dao.DbPriceChannelRuleRow
	err := this.db.Select(&dbPupdoffRow, sql, cityId, channelType, serviceType)

	if err != nil {
		logrus.Info(err)
	}

	if len(dbPupdoffRow) == 0 {
		return dao.DbPriceChannelRuleRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	return dbPupdoffRow[0], err
}

func (this *SszDb) QueryDailyPriceRule(cityId int) ([]dao.DbPriceCarRuleRow, error) {
	sql := "SELECT id, city_id, city_name, car_type_id, car_type_name, car_shift, daily_km_fee, daily_zone_inner_5, daily_zone_inner_6, daily_zone_inner_7, daily_zone_inner_8, daily_zone_inner_9, daily_zone_inner_10, daily_zone_surronds "
	sql += "FROM price_car_rule  "
	sql += "WHERE city_id = ? AND status = 0;"

	var dbPriceCarRuleRow []dao.DbPriceCarRuleRow
	err := this.db.Select(&dbPriceCarRuleRow, sql, cityId)

	if err != nil {
		logrus.Info(err)
	}

	if len(dbPriceCarRuleRow) == 0 {
		return []dao.DbPriceCarRuleRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	return dbPriceCarRuleRow, err
}

/*
查询该车型在不同城市的包车报价
*/
func (this *SszDb) QueryCitysDailyPriceRuleByCarTypeIds(cityIds []int, carTypeIds []int) ([]dao.DbPriceCarRuleRow, error) {
	sql := "SELECT id, city_id, city_name, car_type_id, car_type_name, car_shift, daily_km_fee, daily_zone_inner_5, daily_zone_inner_6, daily_zone_inner_7, daily_zone_inner_8, daily_zone_inner_9, daily_zone_inner_10, daily_zone_surronds "
	sql += "FROM price_car_rule  "

	var query string
	var args []interface{}
	var err error

	sql += "WHERE city_id in(?) AND car_type_id in(?) AND status = 0;"
	query, args, err = sqlx.In(sql, cityIds, carTypeIds)

	if err != nil {
		Log.Info(err)
	}

	// fmt.Println(args)

	query = this.db.Rebind(query)

	var dbPriceCarRuleRow []dao.DbPriceCarRuleRow
	rows, err := this.db.Queryx(query, args...)

	for rows.Next() {
		row := dao.DbPriceCarRuleRow{}
		err = rows.StructScan(&row)
		dbPriceCarRuleRow = append(dbPriceCarRuleRow, row)
	}

	fmt.Println(rows)
	if err != nil {
		logrus.Info(err)
	}

	if len(dbPriceCarRuleRow) == 0 {
		return []dao.DbPriceCarRuleRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	return dbPriceCarRuleRow, err
}

/*
查询该车型在不同城市的包车报价
*/
func (this *SszDb) QueryCitysDailyPriceRule(cityIds []int, carTypeId int) ([]dao.DbPriceCarRuleRow, error) {
	sql := "SELECT id, city_id, city_name, car_type_id, car_type_name, car_shift, daily_km_fee, daily_zone_inner_5, daily_zone_inner_6, daily_zone_inner_7, daily_zone_inner_8, daily_zone_inner_9, daily_zone_inner_10, daily_zone_surronds "
	sql += "FROM price_car_rule  "

	var query string
	var args []interface{}
	var err error

	if carTypeId == 0 {
		sql += "WHERE city_id in(?) AND status = 0;"
		query, args, err = sqlx.In(sql, cityIds)
	} else {
		sql += "WHERE city_id in(?) AND car_type_id = ? AND status = 0;"
		query, args, err = sqlx.In(sql, cityIds, carTypeId)
	}

	if err != nil {
		Log.Info(err)
	}

	fmt.Println(args)

	query = this.db.Rebind(query)

	var dbPriceCarRuleRow []dao.DbPriceCarRuleRow
	rows, err := this.db.Queryx(query, args...)

	for rows.Next() {
		row := dao.DbPriceCarRuleRow{}
		err = rows.StructScan(&row)
		dbPriceCarRuleRow = append(dbPriceCarRuleRow, row)
	}

	fmt.Println(rows)
	if err != nil {
		logrus.Info(err)
	}

	if len(dbPriceCarRuleRow) == 0 {
		return []dao.DbPriceCarRuleRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	return dbPriceCarRuleRow, err
}

func (this *SszDb) QueryPeriphery(cityId int, periphery_city_id int, car_type_id int) (dao.DbPricePeripheryRow, error) {
	sql := "SELECT periphery_id,`periphery_shift`,city_id,city_name, periphery_city_id,car_type_id,car_type_name,car_price "
	sql += "FROM price_periphery  "
	sql += "WHERE city_id = ? AND periphery_city_id = ? AND car_type_id = ? AND status = 0;"

	var dbPricePeripheryRow []dao.DbPricePeripheryRow
	err := this.db.Select(&dbPricePeripheryRow, sql, cityId, periphery_city_id, car_type_id)

	if err != nil {
		logrus.Info(err)
	}

	if len(dbPricePeripheryRow) == 0 {
		return dao.DbPricePeripheryRow{}, errors.NewUserError(Err.ERROR, "no data")
	}

	return dbPricePeripheryRow[0], err
}

func (this *SszDb) CreatePriceSnapshot(priceSnaphot dao.PriceSnaphot) error {

	tx := this.db.MustBegin()
	sql := "INSERT INTO price_snapshot(`id`, `order_id`, `channel_id`, `price_mark`, `city_id`, `city_name`, `timezone`, `service_start_time`, `city_currtime`, `order_type`, `car_type_id`, `car_type_name`, `car_class`, `channel_car_type_id`, `channel_car_type_name`, `guide_base_price`, `guide_max_price`, `sys_price`, `customer_price`, `guest_num`, `adult_number`, `child_number`, `unit_child_seat_price`, `luggage_num`, `max_luggage_num`, `urgent`, `urgent_shift`, `channel_profit`, `activity_shift`, `season_shift`, `user_shift`, `user_profit`, `city_area_index`, `unit_pupoff_price`, `unit_daily_fee`, `ctrip_pupoff_shift`, `ctrip_day_shift`, `unit_price_line`, `welcome_board_price`, `checkin_price`, `currency`, `currency_rate`, `price_channel_rule_id`, `price_basic_rule_id`, `offer_sign`, `request_param`, `expired_at`, `create_time`, `update_time`, `status`) "
	sql += "VALUES (:id,:order_id,:channel_id,:price_mark,:city_id,:city_name,:timezone,:service_start_time,:city_currtime,:order_type,:car_type_id,:car_type_name,:car_class,:channel_car_type_id,:channel_car_type_name,:guide_base_price,:guide_max_price,:sys_price,:customer_price,:guest_num,:adult_number,:child_number,:unit_child_seat_price,:luggage_num,:max_luggage_num,:urgent,:urgent_shift,:channel_profit,:activity_shift,:season_shift,:user_shift,:user_profit,:city_area_index,:unit_pupoff_price,:unit_daily_fee,:ctrip_pupoff_shift,:ctrip_day_shift,:unit_price_line,:welcome_board_price,:checkin_price,:currency,:currency_rate,:price_channel_rule_id,:price_basic_rule_id,:offer_sign,:request_param,:expired_at,:create_time,:update_time,:status)"

	_, err1 := tx.NamedExec(sql, &priceSnaphot)
	if err1 != nil {
		fmt.Println(err1.Error())
	}

	err := tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err1 != nil {
		return err1
	}

	return nil
}

func (this *SszDb) UpdatePriceSnapshot(id string, orderId int64, customerPrice float64, guideMaxManualPrice float64, updateTime time.Time) error {

	tx := this.db.MustBegin()

	sql := "UPDATE price_snapshot SET order_id = ?,customer_price=?,guide_max_manual_price=?,update_time=NOW() WHERE id = ?"
	res := tx.MustExec(sql, orderId, customerPrice, guideMaxManualPrice, id)
	Log.Info(res)
	err := tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (this *SszDb) UpdatePriceSnapshotInfo(id string, guideBasePrice float64, guideMaxPrice float64, sysPrice float64, adultNumber int, childNumber int, maxLuggageNum int, isCheckIn bool, IsWelcomeboard bool) error {

	tx := this.db.MustBegin()

	sql := "UPDATE price_snapshot SET guide_base_price = ?,guide_max_price=?,sys_price = ?,adult_number=?,child_number=?,max_luggage_num=?,update_time=NOW(),is_check_in=?,is_welcome_board=? WHERE id = ?"
	res := tx.MustExec(sql, guideBasePrice, guideMaxPrice, sysPrice, adultNumber, childNumber, maxLuggageNum, isCheckIn, IsWelcomeboard, id)

	err := tx.Commit()
	if err != nil {
		fmt.Println(err.Error(), res)
	}

	return nil
}

func (this *SszDb) QueryPriceSnapshot(priceSnapshotId string) (dao.PriceSnaphot, error) {
	sql := "SELECT `id`, `order_id`, `channel_id`, `price_mark`, `city_id`, `city_name`, `timezone`, `service_start_time`, `city_currtime`, `order_type`, `car_type_id`, `car_type_name`, `car_class`, `channel_car_type_id`, `channel_car_type_name`, `guide_base_price`, `guide_max_price`,`guide_max_manual_price`, `sys_price`, `customer_price`, `guest_num`, `adult_number`, `child_number`, `unit_child_seat_price`, `luggage_num`, `max_luggage_num`, `urgent`, `urgent_shift`, `channel_profit`, `activity_shift`, `season_shift`, `user_shift`, `user_profit`, `city_area_index`, `unit_pupoff_price`, `unit_daily_fee`, `ctrip_pupoff_shift`, `ctrip_day_shift`, `unit_price_line`, `is_check_in`, `is_welcome_board`, `welcome_board_price`, `checkin_price`, `currency`, `currency_rate`, `price_channel_rule_id`, `price_basic_rule_id`, `offer_sign`, `request_param`, `expired_at`, `create_time`, `update_time`, `status` "
	sql += "FROM price_snapshot  "
	sql += "WHERE id = ? AND status = 0;"

	var dbPriceSnaphotRow []dao.PriceSnaphot
	err := this.db.Select(&dbPriceSnaphotRow, sql, priceSnapshotId)

	if err != nil {
		logrus.Info(err)
	}

	if len(dbPriceSnaphotRow) == 0 {
		return dao.PriceSnaphot{}, errors.NewUserError(Err.ERROR, "no data")
	}

	return dbPriceSnaphotRow[0], err
}

func (this *SszDb) QueryPriceSnapshotByPriceMark(priceMark string) (dao.PriceSnaphot, error) {
	sql := "SELECT id, order_id,channel_id, city_id, city_name, timezone,service_start_time, order_type, car_type_id, car_type_name,guide_max_price, guide_base_price,`guide_max_manual_price`,sys_price, customer_price, currency,currency_rate,adult_number,child_number,luggage_num,max_luggage_num,guest_num,request_param, urgent, price_channel_rule_id, channel_profit, activity_shift, offer_sign, price_basic_rule_id, expired_at, create_time, update_time, status "
	sql += "FROM price_snapshot  "
	sql += "WHERE price_mark = ? AND status = 0;"

	var dbPriceSnaphotRow []dao.PriceSnaphot
	err := this.db.Select(&dbPriceSnaphotRow, sql, priceMark)

	if err != nil {
		logrus.Info(err)
	}

	if len(dbPriceSnaphotRow) == 0 {
		return dao.PriceSnaphot{}, nil
	}

	return dbPriceSnaphotRow[0], err
}

func (this *SszDb) QueryPriceSnapshotByChannelInfo(price_mark string, channel_car_type_id int) (dao.PriceSnaphot, error) {
	sql := "SELECT id, order_id,channel_id, city_id, city_name, timezone,service_start_time, order_type, car_type_id, car_type_name, guide_base_price,guide_max_price,`guide_max_manual_price`, sys_price, customer_price, currency,currency_rate,request_param, urgent, price_channel_rule_id, channel_profit, activity_shift, offer_sign, price_basic_rule_id, expired_at, create_time, update_time, status "
	sql += "FROM price_snapshot  "
	sql += "WHERE price_mark = ? AND channel_car_type_id = ? AND status = 0;"

	var dbPriceSnaphotRow []dao.PriceSnaphot
	err := this.db.Select(&dbPriceSnaphotRow, sql, price_mark, channel_car_type_id)

	if err != nil {
		logrus.Info(err)
	}

	if len(dbPriceSnaphotRow) == 0 {
		return dao.PriceSnaphot{}, errors.NewUserError(Err.ERROR, "no data")
	}

	return dbPriceSnaphotRow[0], err
}

/*
获取携程车型报价系数
*/
func (this *SszDb) QueryCtripChannelCarTypes(cityId int, channelCarTypeId []int, includeCarTypeId []int) ([]dao.DbChannelCarTypesRow, error) {

	var dbChannelCarInfoRow []dao.DbChannelCarTypesRow
	var err error

	sql := "SELECT channel_car_type_id,channel_car_type_name,city_id,city_name,channel_car_class,car_type_id,price_shift,is_chinese_service"
	sql = sql + " FROM price_ctrip"
	sql = sql + " WHERE city_id=? AND status = 1 "

	if len(channelCarTypeId) > 0 {
		sql = sql + " AND channel_car_type_id IN(" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(channelCarTypeId)), ","), "[]") + ")"
	}
	if len(includeCarTypeId) > 0 {
		sql = sql + " AND car_type_id IN(" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(includeCarTypeId)), ","), "[]") + ")"
	}

	err = this.db.Select(&dbChannelCarInfoRow, sql, cityId)

	if err != nil {
		logrus.Info(sql, err)
	}

	return dbChannelCarInfoRow, err
}

func (this *SszDb) QueryPupdoffAirports(airportCode string) ([]dao.DbPupdoffAirportRow, error) {
	sql := "SELECT DISTINCT city_id,city_name,airport_id,airport_name "
	sql += "FROM price_pupdoff_rule "
	sql += "WHERE airport_code = ? AND status = 0;"
	var dbPupdoffAirportRow []dao.DbPupdoffAirportRow
	err := this.db.Select(&dbPupdoffAirportRow, sql, airportCode)

	if err != nil {
		logrus.Info(err)
	}

	return dbPupdoffAirportRow, err
}

func (this *SszDb) QueryExpensiveAirports(cityId int) (dao.DbExpensiveAirportsRow, error) {
	sql := "SELECT  airport_id,airport_name,airport_code,price "
	sql += "FROM price_pupdoff_rule "
	sql += "WHERE city_id = ? AND status = 0 ORDER BY price DESC LIMIT 1;"
	var dbCheapestAirportsRow []dao.DbExpensiveAirportsRow
	err := this.db.Select(&dbCheapestAirportsRow, sql, cityId)

	if err != nil {
		logrus.Info(err)
		return dao.DbExpensiveAirportsRow{}, err
	}

	if len(dbCheapestAirportsRow) == 0 {
		return dao.DbExpensiveAirportsRow{}, errors.NewUserError(Err.ERR_PRICE_CARTYPE_EMPTY, "city pupoff price empty", "城市无接送机报价")
	}

	return dbCheapestAirportsRow[0], err
}

func (this *SszDb) CreatePriceLine(priceLine dao.PriceLine) error {

	tx := this.db.MustBegin()
	sql := "INSERT INTO price_line(price_line_id,line_id,line_name,start_city_id,start_city_name,car_type_id,car_type_name,price,guide_price,bind_guide,status) VALUES (:price_line_id,:line_id,:line_name,:start_city_id,:start_city_name,:car_type_id,:car_type_name,:price,:guide_price,:bind_guide,:status)"
	_, err1 := tx.NamedExec(sql, &priceLine)
	if err1 != nil {
		fmt.Println(err1.Error())
	}

	err := tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err1 != nil {
		return err1
	}

	return nil
}

func (this *SszDb) DeletePriceLine(lineId int) error {

	tx := this.db.MustBegin()
	sql := "DELETE FROM price_line WHERE line_id=?"
	res := tx.MustExec(sql, lineId)
	Log.Info(res)
	err := tx.Commit()
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (this *SszDb) QueryPriceLine(cityId int, lineId int) ([]dao.DBPriceLineRow, error) {
	sql := "SELECT  price_line_id,line_id,line_name,start_city_id,start_city_name,car_type_id,car_type_name,price,guide_price,bind_guide,status "
	sql += "FROM price_line "
	sql += "WHERE start_city_id = ? AND status = 1 AND line_id = ?"
	var dbPriceLineRow []dao.DBPriceLineRow
	err := this.db.Select(&dbPriceLineRow, sql, cityId, lineId)

	if err != nil {
		logrus.Info(err)
	}

	return dbPriceLineRow, err
}

/*
获取携程车型报价系数
*/
func (this *SszDb) QueryCtripChannelDayCarTypes(cityId int, packageTypeId int, channelCarTypeId []int, includeChannelCarTypeId []int) ([]dao.DbPriceCtripDayRow, error) {

	var dbChannelCarInfoRow []dao.DbPriceCtripDayRow
	var err error

	sql := "SELECT price_ctrip_day_id,package_type_id,channel_car_type_id,channel_car_type_name,city_id,city_name,channel_car_class,car_type_id,car_type_name,price_shift,is_chinese_service"
	sql = sql + " FROM price_ctrip_day"
	sql = sql + " WHERE city_id=? AND package_type_id = ? AND status = 1 "

	if len(channelCarTypeId) > 0 {
		sql = sql + " AND channel_car_type_id IN(" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(channelCarTypeId)), ","), "[]") + ")"
	}
	if len(includeChannelCarTypeId) > 0 {
		sql = sql + " AND channel_car_type_id IN(" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(includeChannelCarTypeId)), ","), "[]") + ")"
	}

	err = this.db.Select(&dbChannelCarInfoRow, sql, cityId, packageTypeId)

	if err != nil {
		logrus.Info(sql, err)
	}

	return dbChannelCarInfoRow, err
}

/*
获取单次接送区域信息
*/
func (this *SszDb) QueryPtpArea(cityId int) ([]dao.DbPricePtpAreaRow, error) {
	sql := "SELECT  price_area_id,service_type,service_type_name,city_id,city_name,area_price_type,area_name,area_pois,remark,update_time,create_time,status "
	sql += "FROM price_area "
	sql += "WHERE city_id = ? AND status = 1 ORDER BY area_price_type ASC;"
	var dbPtpAreaRow []dao.DbPricePtpAreaRow
	err := this.db.Select(&dbPtpAreaRow, sql, cityId)

	if err != nil {
		logrus.Info(err)
	}

	return dbPtpAreaRow, err
}

/*
获取粤港专线单次接送区域车型信息
*/
func (this *SszDb) QueryPricePtp(startPriceAreaId int, startAreaPriceType int, endPriceAreaId int, endAreaPriceType int) ([]dao.DbPricePtpRow, error) {
	sql := "SELECT  price_ptp_id,start_city_id,start_city_name,end_city_id,end_city_name,start_price_area_id,start_area_price_type,end_price_area_id,end_area_price_type,car_type_id,car_type_name,price,bind_guide,update_time,create_time,status "
	sql += "FROM price_ptp "
	sql += "WHERE ((start_price_area_id = ? AND start_area_price_type = ? AND end_price_area_id = ? AND end_area_price_type = ?) OR (start_price_area_id = ? AND start_area_price_type = ? AND end_price_area_id = ? AND end_area_price_type = ?))  AND status = 1 ORDER BY car_type_id ASC;"
	var dbPricePtpRow []dao.DbPricePtpRow
	err := this.db.Select(&dbPricePtpRow, sql, startPriceAreaId, startAreaPriceType, endPriceAreaId, endAreaPriceType, endPriceAreaId, endAreaPriceType, startPriceAreaId, startAreaPriceType)

	if err != nil {
		logrus.Info(err)
	}

	return dbPricePtpRow, err
}

/*
获取急单信息
*/
func (this *SszDb) QueryPriceUrgentInfo(cityId int, channelType int, serviceType int, carTypeId int) ([]dao.DbPriceUrgentInfo, error) {

	sql := "SELECT pur.*,purc.car_type_id,purc.car_type_name "
	sql += "FROM price_urgent_rule AS pur INNER JOIN price_urgent_rule_cartype purc "
	sql += "ON pur.price_urgent_rule_id = purc.price_urgent_rule_id AND pur.status = 1 AND purc.status = 1 AND pur.city_id = ? AND channel_type in (?,0) AND service_type = ? AND purc.car_type_id = ? "
	sql += "ORDER BY channel_type desc,rule_index Desc;"

	var dbPriceUrgentInfo []dao.DbPriceUrgentInfo
	err := this.db.Select(&dbPriceUrgentInfo, sql, cityId, channelType, serviceType, carTypeId)

	if err != nil {
		logrus.Info(err)
	}

	return dbPriceUrgentInfo, err
}
