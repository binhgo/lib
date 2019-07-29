package dao

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type PriceSnaphot struct {
	Id                     string
	Order_id               int
	Channel_id             int    // 渠道号
	Price_mark             string // 报价批量ID
	City_id                int
	City_name              string
	Timezone               float64
	Service_start_time     mysql.NullTime
	City_currtime          mysql.NullTime
	Order_type             int
	Car_type_id            int
	Car_type_name          string
	Car_class              int // 车等级
	Channel_car_type_id    int
	Channel_car_type_name  string
	Guide_base_price       float64
	Guide_max_price        float64
	Guide_max_manual_price float64 // 手动指定司导接单价
	Sys_price              float64
	Customer_price         float64
	Guest_num              int // 可座人数，除司机外
	Adult_number           int
	Child_number           int
	Unit_child_seat_price  float64
	Luggage_num            int // 行李数
	Max_luggage_num        int // 最大行李数
	Urgent                 int
	Urgent_shift           float64 // 急单浮动
	Channel_profit         float64 // 渠道浮动
	Activity_shift         float64
	Season_shift           float64
	User_shift             float64
	User_profit            float64
	City_area_index        int
	Unit_pupoff_price      float64
	Unit_daily_fee         float64
	Ctrip_pupoff_shift     float64
	Ctrip_day_shift        float64
	Unit_price_line        float64
	Is_Welcome_board       bool
	Welcome_board_price    float64
	Is_Check_in            bool
	Checkin_price          float64
	Currency               string
	Currency_rate          float64
	Price_channel_rule_id  int
	Price_basic_rule_id    int
	Offer_sign             string
	Request_param          string
	Expired_at             mysql.NullTime
	Create_time            mysql.NullTime
	Update_time            mysql.NullTime
	Status                 int
}

type DbCarTypesRow struct {
	Luggage_num   int
	Guest_num     int
	Car_type_id   int
	Car_type_name sql.NullString
	// Id                int
	Introduction      sql.NullString
	Pictures          sql.NullString
	Mapping_car_type  int
	Mapping_seat_type int
	Brand_name        sql.NullString
	Car_model_name    sql.NullString
	Seat_type         int
}

type DbChannelCarTypesRow struct {
	// Seat_type             int
	City_id               int
	City_name             sql.NullString
	Car_type_id           int
	Channel_car_type_id   int
	Channel_car_type_name sql.NullString
	Channel_car_class     int
	Price_shift           float64
	Is_chinese_service    bool // true 中文司导 false 当地司机
}

type DBCarRow struct {
	Id             int            `db:"id"`
	Car_brand_id   int            `db:"car_brand_id"`
	Car_brand_name sql.NullString `db:"car_brand_name"`
	Car_model      string         `db:"car_model"`
}

type DbCurrencyRow struct {
	CityId       int     `db:"cityId"`
	Currency     string  `db:"currency"`
	CurrencyRate float64 `db:"currencyRate"`
	Timezone     float64 `db:"timezone"`
}

type DbPriceBasicRuleRow struct {
	Id                           int            `db:"id"`
	City_id                      int            `db:"city_id"`
	City_name                    string         `db:"city_name"`
	Country_id                   int            `db:"country_id"`
	Country_name                 sql.NullString `db:"country_name"`
	Checkin_price                float64        `db:"checkin_price"`
	Welcome_board_price          float64        `db:"welcome_board_price"`
	Childseat_price              float64        `db:"childseat_price"`
	Night_shift                  float64        `db:"night_shift"`
	Night_shift_start_time       string         `db:"night_shift_start_time"`
	Night_shift_end_time         string         `db:"night_shift_end_time"`
	Season_shift_1               float64        `db:"season_shift_1"`
	Season_shift_periods_1       sql.NullString `db:"season_shift_periods_1"`
	Season_shift_2               float64        `db:"season_shift_2"`
	Season_shift_periods_2       sql.NullString `db:"season_shift_periods_2"`
	Season_shift_3               float64        `db:"season_shift_3"`
	Season_shift_periods_3       sql.NullString `db:"season_shift_periods_3"`
	Season_shift_4               float64        `db:"season_shift_4"`
	Season_shift_periods_4       sql.NullString `db:"season_shift_periods_4"`
	Season_daliy_shift_1         float64        `db:"season_daliy_shift_1"`
	Season_daliy_shift_periods_1 sql.NullString `db:"season_daliy_shift_periods_1"`
	Season_daliy_shift_2         float64        `db:"season_daliy_shift_2"`
	Season_daliy_shift_periods_2 sql.NullString `db:"season_daliy_shift_periods_2"`
	Season_daliy_shift_3         float64        `db:"season_daliy_shift_3"`
	Season_daliy_shift_periods_3 sql.NullString `db:"season_daliy_shift_periods_3"`
	Season_daliy_shift_4         float64        `db:"season_daliy_shift_4"`
	Season_daliy_shift_periods_4 sql.NullString `db:"season_daliy_shift_periods_4"`
	User_profit                  float64        `db:"user_profit"`
	User_shift                   float64        `db:"user_shift"`
	Guide_profit                 float64        `db:"guide_profit"`
	Guide_shift                  float64        `db:"guide_shift"`
	Urgent_shift                 float64        `db:"urgent_shift"`
	Accommodation_fee            float64        `db:"accommodation_fee"`
	Urgent_pupdoff_limit_hour    int            `db:"urgent_pupdoff_limit_hour"`
	Urgent_daliy_limit_hour      int            `db:"urgent_daliy_limit_hour"`
	Booking_pupdoff_limit_hour   int            `db:"booking_pupdoff_limit_hour"` // 接送机最低限制时间
	Booking_daliy_limit_hour     int            `db:"booking_daliy_limit_hour"`   // 包车最低限制时间
	Car_unit_parkingfee          float64        `db:"car_unit_parkingfee"`
	Service_timeoutfee           float64        `db:"service_timeoutfee"`
	City_area_1                  sql.NullString `db:"city_area_1"`
	City_area_2                  sql.NullString `db:"city_area_2"`
	City_area_3                  sql.NullString `db:"city_area_3"`
	City_area_4                  sql.NullString `db:"city_area_4"`
	City_area_5                  sql.NullString `db:"city_area_5"`
	City_area_6                  sql.NullString `db:"city_area_6"`
	City_area_7                  sql.NullString `db:"city_area_7"`
	City_area_8                  sql.NullString `db:"city_area_8"`
	City_area_9                  sql.NullString `db:"city_area_9"`
	City_area_10                 sql.NullString `db:"city_area_10"`
	Create_time                  mysql.NullTime `db:"create_time"`
	Update_time                  mysql.NullTime `db:"update_time"`
}

type DbCityAreaRow struct {
	City_area_1  sql.NullString `db:"city_area_1"`
	City_area_2  sql.NullString `db:"city_area_2"`
	City_area_3  sql.NullString `db:"city_area_3"`
	City_area_4  sql.NullString `db:"city_area_4"`
	City_area_5  sql.NullString `db:"city_area_5"`
	City_area_6  sql.NullString `db:"city_area_6"`
	City_area_7  sql.NullString `db:"city_area_7"`
	City_area_8  sql.NullString `db:"city_area_8"`
	City_area_9  sql.NullString `db:"city_area_9"`
	City_area_10 sql.NullString `db:"city_area_10"`
}

// car_type_id,car_type_name,city_area_id,price
type DbPupdoffRow struct {
	CarTypeId   int     `db:"car_type_id"`
	CarTypeName string  `db:"car_type_name"`
	CityAreaId  int     `db:"city_area_id"`
	PriceShift  float64 `db:"price_shift"`
	Price       float64 `db:"price"`
}

type DbPriceChannelRuleRow struct {
	ChannelId     int     `db:"id"`
	CityId        int     `db:"city_id"`
	CityName      string  `db:"city_name"`
	ServiceType   int     `db:"service_type"`
	ChannelType   int     `db:"channel_type"`
	ChannelProfit float64 `db:"channel_profit"`
	ActivityShift float64 `db:"activity_shift"`
}

// id, city_id, city_name, car_type_id, car_type_name, car_shift, daily_zone_inner_5, daily_zone_inner_6, daily_zone_inner_7, daily_zone_inner_8, daily_zone_inner_9, daily_zone_inner_10, daily_zone_surronds, price_basic_periphery
type DbPriceCarRuleRow struct {
	PriceCarRuleId    int     `db:"id"`
	CityId            int     `db:"city_id"`
	CityName          string  `db:"city_name"`
	CarTypeId         int     `db:"car_type_id"`
	CarTypeName       string  `db:"car_type_name"`
	CarShift          float64 `db:"car_shift"`
	DaliyKmFee        float64 `db:"daily_km_fee"`
	DailyZoneInner5   float64 `db:"daily_zone_inner_5"`
	DailyZoneInner6   float64 `db:"daily_zone_inner_6"`
	DailyZoneInner7   float64 `db:"daily_zone_inner_7"`
	DailyZoneInner8   float64 `db:"daily_zone_inner_8"`
	DailyZoneInner9   float64 `db:"daily_zone_inner_9"`
	DailyZoneInner10  float64 `db:"daily_zone_inner_10"`
	DailyZoneSurronds float64 `db:"daily_zone_surronds"`
}

// periphery_id, city_id, periphery_city_id, car_type_id,price
type DbPricePeripheryRow struct {
	PeripheryId     int     `db:"periphery_id"`
	CityId          int     `db:"city_id"`
	CityName        string  `db:"city_name"`
	PeripheryCityId string  `db:"periphery_city_id"`
	CarTypeId       int     `db:"car_type_id"`
	CarTypeName     string  `db:"car_type_name"`
	CarPrice        float64 `db:"car_price"`
	PeripheryShift  float64 `db:"periphery_shift"`
}

type GuideRowStruct struct {
	Guide_id int
}

type GuideCarRowStruct struct {
	GuideId      int    `db:"guide_id"`       // 司导ID
	GuideName    string `db:"guide_name"`     // 司导姓名
	GuideMobile  string `db:"guide_mobile"`   // 司导手机号
	GuideCarId   int    `db:"guide_car_id"`   // 司导车辆ID
	PlateNum     string `db:"plate_num"`      // 车牌号
	CarTypeId    int    `db:"car_type_id"`    // 车型ID
	CarTypeName  string `db:"car_type_name"`  // 车型ID
	BrandId      int    `db:"brand_id"`       // 名牌ID
	BrandName    string `db:"brand_name"`     // 品牌名
	CarModelName string `db:"car_model_name"` // 车辆型号
	CarModelId   int    `db:"car_model_id"`   // 车辆型号Id
}

type GuideInfoRowStruct struct {
	Guide_id int    `db:"id"`
	Mobile   string `db:"mobile"`
	Name     string `db:"name"`
}

type OrderRowData struct {
	OrderId       string `db:"order_id"`
	CarTypeId     int    `db:"car_type_id"`
	CarTypeName   string `db:"car_type_name"`
	ServiceCityId string `db:"service_city_id"`
}

type AirportRowStruct struct {
	Airport_id   int    `db:"id"`
	Country_id   int    `db:"country_id"`
	Country_name string `db:"country_name"`
	City_id      int    `db:"city_id"`
	City_name    string `db:"city_name"`
	Code         string `db:"code"`
	Name         string `db:"name"`
	Location     string `db:"location"`
}

type DbPupdoffAirportRow struct {
	CityId      int    `db:"city_id"`
	CityName    string `db:"city_name"`
	AirportId   int    `db:"airport_id"`
	AirportName string `db:"airport_name"`
}

// airport_id,airport_name,airport_code,price
type DbExpensiveAirportsRow struct {
	AirportId   int     `db:"airport_id"`
	AirportName string  `db:"airport_name"`
	AirportCode string  `db:"airport_code"`
	Price       float64 `db:"price"`
}

type CartypeRowStruct struct {
	Car_type_id   int    `db:"id"`
	Car_type_name string `db:"car_type_name"`
}

type CountriesRowStruct struct {
	CountryId     int    `db:"id"`
	CountryName   string `db:"country_name"`
	ContinentId   int    `db:"continent_id"`
	ContinentName string `db:"continent_name"`
	CountryCode   string `db:"code"`
}

type DbCitiesRow struct {
	CityId        int    `db:"id"`
	CityName      string `db:"city_name"`
	CountryId     int    `db:"country_id"`
	CountryName   string `db:"country_name"`
	ContinentId   int    `db:"continent_id"`
	ContinentName string `db:"continent_name"`
	Location      string `db:"location"`
}

type PriceLine struct {
	PriceLineId   int     `db:"price_line_id"`
	LineId        int     `db:"line_id"`
	LineName      string  `db:"line_name"`
	StartCityId   int     `db:"start_city_id"`
	StartCityName string  `db:"start_city_name"`
	CarTypeId     int     `db:"car_type_id"`
	CarTypeName   string  `db:"car_type_name"`
	Price         float64 `db:"price"`
	GuidePrice    float64 `db:"guide_price"`
	BindGuide     string  `db:"bind_guide"`
	Status        int     `db:"status"`
}

type DBPriceLineRow struct {
	PriceLineId   int     `db:"price_line_id"`
	LineId        int     `db:"line_id"`
	LineName      string  `db:"line_name"`
	StartCityId   int     `db:"start_city_id"`
	StartCityName string  `db:"start_city_name"`
	CarTypeId     int     `db:"car_type_id"`
	CarTypeName   string  `db:"car_type_name"`
	Price         float64 `db:"price"`
	GuidePrice    float64 `db:"guide_price"`
	BindGuide     string  `db:"bind_guide"`
	Status        int     `db:"status"`
}

type DbPriceCtripDayRow struct {
	// Seat_type             int
	Price_ctrip_day_id    int
	Package_type_id       int
	Channel_car_type_id   int
	Channel_car_type_name sql.NullString
	City_id               int
	City_name             sql.NullString
	Channel_car_class     int
	Car_type_id           int
	Car_type_name         sql.NullString
	Price_shift           float64
	Is_chinese_service    bool // true 中文司导 false 当地司机
}

type DbPricePtpAreaRow struct {
	Price_area_id     int
	Service_type      int
	Service_type_name sql.NullString
	City_id           int
	City_name         sql.NullString
	Area_price_type   int
	Area_name         sql.NullString
	Area_pois         sql.NullString
	Remark            sql.NullString
	Update_time       mysql.NullTime
	Create_time       mysql.NullTime
	Status            int
}

type DbPricePtpRow struct {
	Price_ptp_id          int
	Start_city_id         int
	Start_city_name       sql.NullString
	End_city_id           int
	End_city_name         sql.NullString
	Start_price_area_id   int
	Start_area_price_type int
	End_price_area_id     int
	End_area_price_type   int
	Car_type_id           int
	Car_type_name         sql.NullString
	Price                 float64
	Bind_guide            sql.NullString
	Update_time           mysql.NullTime
	Create_time           mysql.NullTime
	Status                int
}

type DbPriceUrgentInfo struct {
	Price_urgent_rule_id int
	Rule_index           int
	City_id              int
	City_name            sql.NullString
	Channel_type         int
	Channel_type_name    sql.NullString
	Service_type         int
	Service_type_name    sql.NullString
	Start_hour           int
	End_hour             int
	Urgent_shift         float64
	Create_time          mysql.NullTime
	Update_time          mysql.NullTime
	Remark               sql.NullString
	Status               int
	Car_type_id          int
	Car_type_name        sql.NullString
}
