package common

const (
	DISPATCH_TYPE_PICKUP  = 1
	DISPATCH_TYPE_DROPOFF = 2
	DISPATCH_TYPE_DAILY   = 3
	DISPATCH_TYPE_PTP     = 4
)

// 派单任务状态
const (
	DISPATCH_TASK_CREATE  = 100001
	DISPATCH_TASK_RUNNING = 100002
	DISPATCH_TASK_CANCEL  = 100003
	DISPATCH_TASK_PAUSE   = 100004
	DISPATCH_TASK_FAIL    = 100005
	DISPATCH_TASK_END     = 100006
)

// 司导派单状态status
const (
	DISPATCH_STATUS_CREATE        = 100105 // 派单创建
	DISPATCH_STATUS_RECIVING      = 100101 // 等待表态
	DISPATCH_STATUS_WAITING       = 100102 // 派单中，等待PK
	DISPATCH_STATUS_SUCCESS       = 100100 // PK成功
	DISPATCH_STATUS_FAIL          = 100103 // PK失败
	DISPATCH_STATUS_CANCEL        = 100104 // 订单取消
	DISPATCH_STATUS_PRECIPITATION = 100105 // 沉淀订单 在沉淀有效期内，用户依然可以抢单
	DISPATCH_STATUS_END           = 100106 // 结束
)

// 报价status
const (
	QUOTATION_STATUS_CREATE        = 200101 // 派单创建
	QUOTATION_STATUS_DOING         = 200102 // 正在报价
	QUOTATION_STATUS_SUCCESS       = 200100 // 报价结束 已找到司导
	QUOTATION_STATUS_NOGUIDE       = 200103 // 报价结束 未找到司导
	QUOTATION_STATUS_FAIL          = 200104 // 报价失败
	QUOTATION_STATUS_CANCEL        = 200106 // 订单取消
	QUOTATION_STATUS_PRECIPITATION = 200107 // 沉淀订单 在沉淀有效期内，用户依然可以抢单
	QUOTATION_STATUS_END           = 200108 // 报价结束
)

const (
	SIGN_PRICE_PREFIX         = "fjdlsjfdlsQ@#^&*"
	SIGN_DECLARE_PRICE_PREFIX = "fjdlsjfdlsQ@#^&*32131"
)

const (
	SERVICE_TYPE_PICKUP  = 1  // 接机
	SERVICE_TYPE_DROPOFF = 2  // 送机
	SERVICE_TYPE_DAILY   = 3  // 包车
	SERVICE_TYPE_PTP     = 4  // 单次接送
	SERVICE_TYPE_PUPDOFF = 12 // 接送机
)

var SERVICE_TYPE_NAME = map[int]string{SERVICE_TYPE_PICKUP: "接机", SERVICE_TYPE_DROPOFF: "送机", SERVICE_TYPE_DAILY: "包车", SERVICE_TYPE_PTP: "单次接送"}

const (
	// TOUR_TYPE_HALFDAY   = 101
	TOUR_TYPE_NOCAR     = 0
	TOUR_TYPE_DAILY     = 101
	TOUR_TYPE_SURRONDS  = 201
	TOUR_TYPE_PERIPHERY = 301
)

const (
	TOUR_ORDER_TYPE_DAILY         = 2000 // 全日包车
	TOUR_ORDER_TYPE_PICKUP        = 1000 // 仅接机
	TOUR_ORDER_TYPE_DROPOFF       = 1001 // 仅送机
	TOUR_ORDER_TYPE_DAILYPICKUP   = 2010 // 包车+接机
	TOUR_ORDER_TYPE_DAILYDROPPOFF = 2030 // 包车+送机
	TOUR_ORDER_TYPE_NOSERVICE     = 5000 // 无车服务
)

const (
	GUIDE_DECLARE_OK          = "0"  // 我要接单
	GUIDE_DECLARE_LETMESEESEE = "-1" // 我再想象
	GUIDE_DECLARE_REFUSE      = "-2" // 拒绝接单
)

const (
	SMS_ID_BINDGUIDEFAIL    = 510001
	SMS_ID_NOGUIDE          = 510002
	SMS_ID_BINDGUIDESUCCESS = 510003
)

const (
	CHANNEL_MIS             = 100
	CHANNEL_CAPP            = 1001 // 官方手机端APP
	CHANNEL_PC              = 2001 // 官方PC端
	CHANNEL_TMALL           = 3001 // 天猫
	CHANNEL_CTRIP           = 4001 // 携程API
	CHANNEL_CTRIP_LOCALPLAY = 4002 // 携程当地玩乐
	CHANNEL_FLIGGY          = 5001 // 飞猪
	// 6001
	CHANNEL_JD          = 7001  // 京东
	CHANNEL_WECHAT_MINI = 8001  // 微信小程序
	CHANNEL_WECHAT      = 8002  // 微信
	CHANNEL_OFFLINE     = 9001  // 线下订单
	CHANNEL_ACTIVITY    = 10001 // 活动
	CHANNEL_MAFENGWO    = 11001 // sku 马蜂窝
	CHANNEL_MEITUAN     = 12001 // sku 美团
	CHANNEL_KLOOK       = 13001 // sku 客路
	CHANNEL_TRAVEL      = 14001 // 水帘洞
)

var CHANNEL_LIST = map[int]string{CHANNEL_MIS: "MIS", CHANNEL_CAPP: "CAPP", CHANNEL_PC: "PC", CHANNEL_CTRIP: "CTRIP", CHANNEL_TMALL: "TMALL",
	CHANNEL_CTRIP_LOCALPLAY: "CTRIP_LOCALPLAY", CHANNEL_FLIGGY: "FLIGGY", CHANNEL_JD: "JD", CHANNEL_WECHAT: "WECHAT", CHANNEL_WECHAT_MINI: "WECHAT_MINI",
	CHANNEL_OFFLINE: "OFFLINE", CHANNEL_ACTIVITY: "ACTIVITY", CHANNEL_MAFENGWO: "MAFENGWO",
	CHANNEL_MEITUAN: "MEITUAN", CHANNEL_KLOOK: "KLOOK", CHANNEL_TRAVEL: "TRAVEL"}

// 携程配置
const (
	CHANNEL_CTRIP_RESPONSE_OK             = "OK"             // 成功
	CHANNEL_CTRIP_RESPONSE_NO_SERVICE     = "NO_SERVICE"     // 该城市无服务
	CHANNEL_CTRIP_RESPONSE_NO_VEHICLE     = "NO_VEHICLE"     // 车型不提供服务
	CHANNEL_CTRIP_RESPONSE_NO_TIMESERVICE = "NO_TIMESERVICE" // 用车时间接 近或者超过业务预定时间，无法 预定
	CHANNEL_CTRIP_RESPONSE_ERROR          = "ERROR"          // 其他错误(正常情况下 不应该出现。)
	CHANNEL_CTRIP_CODE_CS005              = "CS005"
	CHANNEL_CTRIP_CODE_CS006              = "CS006"
	CHANNEL_CTRIP_CODE_CS007              = "CS007" // 儿童座椅(首个免费)
	CHANNEL_CTRIP_CODE_CS008              = "CS008"
	CHANNEL_CTRIP_CODE_PC002              = "PC002" // 举牌接机服务(收费)
	CHANNEL_CTRIP_CODE_PC003              = "PC003"
	CHANNEL_CTRIP_CODE_OCHS001            = "OCHS001"

	CTRIP_PATTERNTYPE_PICKUP  = 1
	CTRIP_PATTERNTYPE_DROPOFF = 2

	CTRIP_DAY_DRIVER_LANGUAGE_DEFAULT       = 0 // 当地司机中 文导游
	CTRIP_DAY_DRIVER_LANGUAGE_LOCAL         = 1 // 当地司机
	CTRIP_DAY_DRIVER_LANGUAGE_CHINESE       = 2 // 中文司机
	CTRIP_DAY_DRIVER_LANGUAGE_LOCAL_CHINESE = 3 // 当地司机中文导游
)

// redis key
const (
	REDIS_KEY_PRICE_PUPDOFF_CHANNEL_SNAPSHOT           = "REDIS_KEY_PRICE_PUPDOFF_CHANNEL_SNAPSHOT"
	REDIS_KEY_PRICE_PUPDOFF_CHANNEL_CTRIP_DAY_SNAPSHOT = "REDIS_KEY_PRICE_PUPDOFF_CHANNEL_CTRIP_DAY_SNAPSHOT"
)

const (
	STATUS_DELETE = -1
	STATUS_INIT   = 0
	STATUS_OK     = 1
)
