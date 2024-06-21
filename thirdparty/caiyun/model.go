package caiyun

import "github.com/alioth-center/infrastructure/utils/values"

type BaseResponse struct {
	Status     string      `json:"status"`      // 返回状态
	ApiVersion string      `json:"api_version"` // API版本
	ApiStatus  string      `json:"api_status"`  // API状态
	Language   string      `json:"lang"`        // 语言
	Unit       string      `json:"unit"`        // 单位
	TzShift    int         `json:"tzshift"`     // 时区偏移
	Timezone   string      `json:"timezone"`    // 时区
	ServerTime int64       `json:"server_time"` // 服务器时间
	Location   Location    `json:"location"`    // 位置
	Result     BaseResults `json:"result"`      // 结果
}

type BaseResults struct {
	AlertResult      AlertResult    `json:"alert,omitempty"`             // 预警信息
	RealtimeResult   RealtimeResult `json:"realtime,omitempty"`          // 实时天气信息
	Primary          int            `json:"primary,omitempty"`           // 主要天气信息
	ForecastKeypoint string         `json:"forecast_keypoint,omitempty"` // 预报关键点
}

// Location is the struct for the location, which contains the latitude and longitude
type Location [2]float64

func (l Location) Longitude() float64 {
	return l[0]
}

func (l Location) Latitude() float64 {
	return l[1]
}

// NewLocation creates a new location with the given latitude and longitude
func NewLocation(longitude, latitude float64) Location {
	return Location{longitude, latitude}
}

// LocationSetSourceItem is the struct for the location set source item
type LocationSetSourceItem struct {
	Adcode           string `csv:"adcode"`            // 行政区划代码
	FormattedAddress string `csv:"formatted_address"` // 格式化地址
	Longitude        string `csv:"lng"`               // 经度
	Latitude         string `csv:"lat"`               // 纬度
}

// LocationSetItem is the struct for the location set item
type LocationSetItem struct {
	Adcode           string  `json:"adcode"`            // 行政区划代码
	FormattedAddress string  `json:"formatted_address"` // 格式化地址
	Longitude        float64 `json:"lng"`               // 经度
	Latitude         float64 `json:"lat"`               // 纬度
}

// NewLocationSetFromSourceItem creates a new location set item from the source item
func NewLocationSetFromSourceItem(sourceItem LocationSetSourceItem) LocationSetItem {
	return LocationSetItem{
		Adcode:           sourceItem.Adcode,
		FormattedAddress: sourceItem.FormattedAddress,
		Longitude:        values.StringToFloat64(sourceItem.Longitude, 0.0),
		Latitude:         values.StringToFloat64(sourceItem.Latitude, 0.0),
	}
}

// LifeIndexType is the type for the life index
//
// Reference: [life-index]
//
// [life-index]: https://docs.caiyunapp.com/weather-api/v2/v2.6/tables/lifeindex.html
type LifeIndexType string

// LifeIndexType values
const (
	LifeIndexTypeUltraviolet LifeIndexType = "ultraviolet" // 紫外线指数
	LifeIndexTypeDressing    LifeIndexType = "dressing"    // 穿衣指数
	LifeIndexTypeComfort     LifeIndexType = "comfort"     // 舒适度指数
	LifeIndexTypeCarWashing  LifeIndexType = "car_washing" // 洗车指数
	LifeIndexTypeColdRisk    LifeIndexType = "cold_risk"   // 感冒指数
)

// LifeIndexDetails is the struct for the life index details
type LifeIndexDetails struct {
	Index       float64 `json:"index"` // 生活指数
	Description string  `json:"desc"`  // 生活指数描述
}

// SkyconType is the type for the skycon
//
// Reference: [skycon]
//
// [skycon]: https://docs.caiyunapp.com/weather-api/v2/v2.6/tables/skycon.html
type SkyconType string

// SkyconType values
const (
	SkyconTypeClearDay          SkyconType = "CLEAR_DAY"           // 晴（白天）
	SkyconTypeClearNight        SkyconType = "CLEAR_NIGHT"         // 晴（夜间）
	SkyconTypePartlyCloudyDay   SkyconType = "PARTLY_CLOUDY_DAY"   // 多云（白天）
	SkyconTypePartlyCloudyNight SkyconType = "PARTLY_CLOUDY_NIGHT" // 多云（夜间）
	SkyconTypeCloudy            SkyconType = "CLOUDY"              // 阴
	SkyconTypeLightHaze         SkyconType = "LIGHT_HAZE"          // 轻度雾霾
	SkyconTypeModerateHaze      SkyconType = "MODERATE_HAZE"       // 中度雾霾
	SkyconTypeHeavyHaze         SkyconType = "HEAVY_HAZE"          // 重度雾霾
	SkyconTypeLightRain         SkyconType = "LIGHT_RAIN"          // 小雨
	SkyconTypeModerateRain      SkyconType = "MODERATE_RAIN"       // 中雨
	SkyconTypeHeavyRain         SkyconType = "HEAVY_RAIN"          // 大雨
	SkyconTypeStorm             SkyconType = "STORM_RAIN"          // 暴雨
	SkyconTypeFog               SkyconType = "FOG"                 // 雾
	SkyconTypeLightSnow         SkyconType = "LIGHT_SNOW"          // 小雪
	SkyconTypeModerateSnow      SkyconType = "MODERATE_SNOW"       // 中雪
	SkyconTypeHeavySnow         SkyconType = "HEAVY_SNOW"          // 大雪
	SkyconTypeStormSnow         SkyconType = "STORM_SNOW"          // 暴雪
	SkyconTypeDust              SkyconType = "DUST"                // 浮尘
	SkyconTypeSand              SkyconType = "SAND"                // 沙尘
	SkyconTypeWind              SkyconType = "WIND"                // 大风
)

// RealtimeResult is the struct for the response of the realtime weather API
//
// Reference: [realtime-weather], [skycon] and [life-index]
//
// [realtime-weather]: https://docs.caiyunapp.com/weather-api/v2/v2.6/1-realtime.html
// [skycon]: https://docs.caiyunapp.com/weather-api/v2/v2.6/tables/skycon.html
// [life-index]: https://docs.caiyunapp.com/weather-api/v2/v2.6/tables/lifeindex.html
type RealtimeResult struct {
	Status              string                             `json:"status"`               // 返回状态
	Temperature         float64                            `json:"temperature"`          // 地表两米气温
	Humidity            float64                            `json:"humidity"`             // 地表两米相对湿度百分比
	Cloudrate           float64                            `json:"cloudrate"`            // 云量百分比
	Skycon              SkyconType                         `json:"skycon"`               // 天气现象
	Visibility          float64                            `json:"visibility"`           // 能见度
	Dswrf               float64                            `json:"dswrf"`                // 短波辐射
	Wind                WindResult                         `json:"wind"`                 // 风向信息
	Pressure            float64                            `json:"pressure"`             // 地面大气压强
	ApparentTemperature float64                            `json:"apparent_temperature"` // 体感温度
	Precipitation       PrecipitationResult                `json:"precipitation"`        // 降水信息
	AirQuality          AirQualityResult                   `json:"air_quality"`          // 空气质量
	LifeIndex           map[LifeIndexType]LifeIndexDetails `json:"life_index"`           // 生活指数
}

// WindResult is the struct for the wind information
type WindResult struct {
	Speed     float64 `json:"speed"`     // 地表两米风速
	Direction float64 `json:"direction"` // 地表两米风向
}

// PrecipitationResult is the struct for the precipitation information
type PrecipitationResult struct {
	Local   LocalPrecipitation   `json:"local"`   // 本地降水信息
	Nearest NearestPrecipitation `json:"nearest"` // 最近降水信息
}

// LocalPrecipitation is the struct for the local precipitation information
type LocalPrecipitation struct {
	Status     string  `json:"status"`     // 返回状态
	Datasource string  `json:"datasource"` // 数据源
	Intensity  float64 `json:"intensity"`  // 降水强度
}

// NearestPrecipitation is the struct for the nearest precipitation information
type NearestPrecipitation struct {
	Status    string  `json:"status"`    // 返回状态
	Distance  float64 `json:"distance"`  // 最近降水带距离
	Intensity float64 `json:"intensity"` // 最近降水带降水强度
}

// AirQualityResult is the struct for the air quality information
type AirQualityResult struct {
	Pm25    float64        `json:"pm25"`        // PM2.5
	Pm10    float64        `json:"pm10"`        // PM10
	O3      float64        `json:"o3"`          // 臭氧
	So2     float64        `json:"so2"`         // 二氧化硫
	No2     float64        `json:"no2"`         // 二氧化氮
	Co      float64        `json:"co"`          // 一氧化碳
	Aqi     AqiSet         `json:"aqi"`         // 空气质量指数
	AqiDesc AqiDescription `json:"description"` // 空气质量描述
}

// AqiSet is the struct for the air quality index set
type AqiSet struct {
	Chn int `json:"chn"` // 中国标准
	Usa int `json:"usa"` // 美国标准
}

// AqiDescription is the struct for the air quality index description
type AqiDescription struct {
	Chn string `json:"chn"` // 中国标准
	Usa string `json:"usa"` // 美国标准
}

// AlertResult is the struct for the response of the alert API
//
// Reference: [alert]
//
// [alert]: https://docs.caiyunapp.com/weather-api/v2/v2.6/5-alert.html
type AlertResult struct {
	Status  string         `json:"status"`  // 返回状态
	Content []AlertContent `json:"content"` // 预警内容
	Adcodes []Adcode       `json:"adcodes"` // 行政区划代码
}

// AlertContent is the struct for the alert content
type AlertContent struct {
	Province      string   `json:"province"`       // 省份
	Status        string   `json:"status"`         // 预警状态
	Code          string   `json:"code"`           // 预警代码
	Description   string   `json:"description"`    // 预警描述
	RegionID      string   `json:"regionId"`       // 区域ID
	Country       string   `json:"country"`        // 县区
	Pubtimestamp  int64    `json:"pubtimestamp"`   // 发布时间
	Latlon        Location `json:"latlon"`         // 经纬度
	City          string   `json:"city"`           // 城市
	AlertID       string   `json:"alertId"`        // 预警ID
	Title         string   `json:"title"`          // 预警标题
	Adcode        string   `json:"adcode"`         // 行政区划代码
	Source        string   `json:"source"`         // 预警来源
	Location      string   `json:"location"`       // 预警地点
	RequestStatus string   `json:"request_status"` // 请求状态
}

// Adcode is the struct for the adcode
type Adcode struct {
	Adcode string `json:"adcode"` // 行政区划代码
	Name   string `json:"name"`   // 地区名称
}
