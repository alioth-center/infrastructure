package timezone

import (
	"github.com/alioth-center/infrastructure/errors"
	"time"
)

var (
	isSetLocal = false

	Zero  = LocationLondon // 零时区
	Local = Zero           // 默认时区，零时区
)

// SetLocal set timezone as Local timezone
// example:
//
//	err := timezone.SetLocal(timezone.LocationBeijing)
func SetLocal(timezone TimeZone) (err error) {
	if Local == timezone {
		// if Local is already set to timezone, return nil
		return nil
	}

	if _, exist := timezoneData[timezone]; !exist {
		// if timezone is not exist, return error
		return errors.NewInvalidTimezoneError(string(timezone))
	}

	if !isSetLocal {
		// if Local is not set, set it to timezone
		Local = timezone
		isSetLocal = true
		return nil
	}

	return errors.NewLocalTimezoneAlreadySetError()
}

type TimeZone string

const (
	LocationBeijing      TimeZone = "Asia/Beijing"
	LocationShanghai     TimeZone = "Asia/Shanghai"
	LocationTokyo        TimeZone = "Asia/Tokyo"
	LocationLondon       TimeZone = "Europe/London"
	LocationParis        TimeZone = "Europe/Paris"
	LocationNewYork      TimeZone = "America/New_York"
	LocationLosAnge      TimeZone = "America/Los_Angeles"
	LocationMoscow       TimeZone = "Europe/Moscow"
	LocationMonaco       TimeZone = "Europe/Monaco"
	LocationKarachi      TimeZone = "Asia/Karachi"
	LocationSydney       TimeZone = "Australia/Sydney"
	LocationSeoul        TimeZone = "Asia/Seoul"
	LocationNewDelhi     TimeZone = "Asia/New_Delhi"
	LocationJohannesburg TimeZone = "Africa/Johannesburg"
	LocationMadrid       TimeZone = "Europe/Madrid"
	LocationDubai        TimeZone = "Asia/Dubai"
	LocationBerlin       TimeZone = "Europe/Berlin"
	LocationRome         TimeZone = "Europe/Rome"
	LocationHonolulu     TimeZone = "Pacific/Honolulu"
	LocationPrague       TimeZone = "Europe/Prague"
	LocationBucharest    TimeZone = "Europe/Bucharest"
	LocationWarsaw       TimeZone = "Europe/Warsaw"
	LocationAthens       TimeZone = "Europe/Athens"
	LocationHelsinki     TimeZone = "Europe/Helsinki"
	LocationStockholm    TimeZone = "Europe/Stockholm"
	LocationCopenhagen   TimeZone = "Europe/Copenhagen"
	LocationOslo         TimeZone = "Europe/Oslo"
	LocationVienna       TimeZone = "Europe/Vienna"
	LocationBrussels     TimeZone = "Europe/Brussels"
	LocationAmsterdam    TimeZone = "Europe/Amsterdam"
	LocationDublin       TimeZone = "Europe/Dublin"
	LocationLisbon       TimeZone = "Europe/Lisbon"
	LocationBern         TimeZone = "Europe/Bern"
	LocationZurich       TimeZone = "Europe/Zurich"
	LocationKeiv         TimeZone = "Europe/Keiv"

	LocationUTC           TimeZone = "UTC"
	LocationUTCEast1      TimeZone = "UTC+1"
	LocationUTCEast2      TimeZone = "UTC+2"
	LocationUTCEast3      TimeZone = "UTC+3"
	LocationUTCEast3Dot5  TimeZone = "UTC+3.5"
	LocationUTCEast4      TimeZone = "UTC+4"
	LocationUTCEast4Dot5  TimeZone = "UTC+4.5"
	LocationUTCEast5      TimeZone = "UTC+5"
	LocationUTCEast5Dot5  TimeZone = "UTC+5.5"
	LocationUTCEast5Dot75 TimeZone = "UTC+5.75"
	LocationUTCEast6      TimeZone = "UTC+6"
	LocationUTCEast7      TimeZone = "UTC+7"
	LocationUTCEast8      TimeZone = "UTC+8"
	LocationUTCEast9      TimeZone = "UTC+9"
	LocationUTCEast9Dot5  TimeZone = "UTC+9.5"
	LocationUTCEast10     TimeZone = "UTC+10"
	LocationUTCEast10Dot5 TimeZone = "UTC+10.5"
	LocationUTCEast11     TimeZone = "UTC+11"
	LocationUTCEast12     TimeZone = "UTC+12"
	LocationUTCEast13     TimeZone = "UTC+13"
	LocationUTCEast14     TimeZone = "UTC+14"
	LocationUTCWest1      TimeZone = "UTC-1"
	LocationUTCWest2      TimeZone = "UTC-2"
	LocationUTCWest3      TimeZone = "UTC-3"
	LocationUTCWest4      TimeZone = "UTC-4"
	LocationUTCWest5      TimeZone = "UTC-5"
	LocationUTCWest6      TimeZone = "UTC-6"
	LocationUTCWest7      TimeZone = "UTC-7"
	LocationUTCWest8      TimeZone = "UTC-8"
	LocationUTCWest9      TimeZone = "UTC-9"
	LocationUTCWest10     TimeZone = "UTC-10"
	LocationUTCWest11     TimeZone = "UTC-11"
	LocationUTCWest12     TimeZone = "UTC-12"
)

var (
	timezoneData = map[TimeZone]*time.Location{
		LocationBeijing:      time.FixedZone("Asia/Beijing", 8*60*60),
		LocationShanghai:     time.FixedZone("Asia/Shanghai", 8*60*60),
		LocationTokyo:        time.FixedZone("Asia/Tokyo", 9*60*60),
		LocationLondon:       time.FixedZone("Europe/London", 0),
		LocationParis:        time.FixedZone("Europe/Paris", 1*60*60),
		LocationNewYork:      time.FixedZone("America/New_York", -5*60*60),
		LocationLosAnge:      time.FixedZone("America/Los_Angeles", -8*60*60),
		LocationMoscow:       time.FixedZone("Europe/Moscow", 3*60*60),
		LocationMonaco:       time.FixedZone("Europe/Monaco", 1*60*60),
		LocationKarachi:      time.FixedZone("Asia/Karachi", 5*60*60),
		LocationSydney:       time.FixedZone("Australia/Sydney", 10*60*60),
		LocationSeoul:        time.FixedZone("Asia/Seoul", 9*60*60),
		LocationNewDelhi:     time.FixedZone("Asia/New_Delhi", 5*60*60+30*60),
		LocationJohannesburg: time.FixedZone("Africa/Johannesburg", 2*60*60),
		LocationMadrid:       time.FixedZone("Europe/Madrid", 1*60*60),
		LocationDubai:        time.FixedZone("Asia/Dubai", 4*60*60),
		LocationBerlin:       time.FixedZone("Europe/Berlin", 1*60*60),
		LocationRome:         time.FixedZone("Europe/Rome", 1*60*60),
		LocationHonolulu:     time.FixedZone("Pacific/Honolulu", -10*60*60),
		LocationPrague:       time.FixedZone("Europe/Prague", 1*60*60),
		LocationBucharest:    time.FixedZone("Europe/Bucharest", 2*60*60),
		LocationWarsaw:       time.FixedZone("Europe/Warsaw", 1*60*60),
		LocationAthens:       time.FixedZone("Europe/Athens", 2*60*60),
		LocationHelsinki:     time.FixedZone("Europe/Helsinki", 2*60*60),
		LocationStockholm:    time.FixedZone("Europe/Stockholm", 1*60*60),
		LocationCopenhagen:   time.FixedZone("Europe/Copenhagen", 1*60*60),
		LocationOslo:         time.FixedZone("Europe/Oslo", 1*60*60),
		LocationVienna:       time.FixedZone("Europe/Vienna", 1*60*60),
		LocationBrussels:     time.FixedZone("Europe/Brussels", 1*60*60),
		LocationAmsterdam:    time.FixedZone("Europe/Amsterdam", 1*60*60),
		LocationDublin:       time.FixedZone("Europe/Dublin", 0),
		LocationLisbon:       time.FixedZone("Europe/Lisbon", 0),
		LocationBern:         time.FixedZone("Europe/Bern", 1*60*60),
		LocationZurich:       time.FixedZone("Europe/Zurich", 1*60*60),
		LocationKeiv:         time.FixedZone("Europe/Keiv", 2*60*60),

		LocationUTC:           time.UTC,
		LocationUTCEast1:      time.FixedZone("UTC+1", 1*60*60),
		LocationUTCEast2:      time.FixedZone("UTC+2", 2*60*60),
		LocationUTCEast3:      time.FixedZone("UTC+3", 3*60*60),
		LocationUTCEast3Dot5:  time.FixedZone("UTC+3.5", 3*60*60+30*60),
		LocationUTCEast4:      time.FixedZone("UTC+4", 4*60*60),
		LocationUTCEast4Dot5:  time.FixedZone("UTC+4.5", 4*60*60+30*60),
		LocationUTCEast5:      time.FixedZone("UTC+5", 5*60*60),
		LocationUTCEast5Dot5:  time.FixedZone("UTC+5.5", 5*60*60+30*60),
		LocationUTCEast5Dot75: time.FixedZone("UTC+5.75", 5*60*60+45*60),
		LocationUTCEast6:      time.FixedZone("UTC+6", 6*60*60),
		LocationUTCEast7:      time.FixedZone("UTC+7", 7*60*60),
		LocationUTCEast8:      time.FixedZone("UTC+8", 8*60*60),
		LocationUTCEast9:      time.FixedZone("UTC+9", 9*60*60),
		LocationUTCEast9Dot5:  time.FixedZone("UTC+9.5", 9*60*60+30*60),
		LocationUTCEast10:     time.FixedZone("UTC+10", 10*60*60),
		LocationUTCEast10Dot5: time.FixedZone("UTC+10.5", 10*60*60+30*60),
		LocationUTCEast11:     time.FixedZone("UTC+11", 11*60*60),
		LocationUTCEast12:     time.FixedZone("UTC+12", 12*60*60),
		LocationUTCEast13:     time.FixedZone("UTC+13", 13*60*60),
		LocationUTCEast14:     time.FixedZone("UTC+14", 14*60*60),
		LocationUTCWest1:      time.FixedZone("UTC-1", -1*60*60),
		LocationUTCWest2:      time.FixedZone("UTC-2", -2*60*60),
		LocationUTCWest3:      time.FixedZone("UTC-3", -3*60*60),
		LocationUTCWest4:      time.FixedZone("UTC-4", -4*60*60),
		LocationUTCWest5:      time.FixedZone("UTC-5", -5*60*60),
		LocationUTCWest6:      time.FixedZone("UTC-6", -6*60*60),
		LocationUTCWest7:      time.FixedZone("UTC-7", -7*60*60),
		LocationUTCWest8:      time.FixedZone("UTC-8", -8*60*60),
		LocationUTCWest9:      time.FixedZone("UTC-9", -9*60*60),
		LocationUTCWest10:     time.FixedZone("UTC-10", -10*60*60),
		LocationUTCWest11:     time.FixedZone("UTC-11", -11*60*60),
		LocationUTCWest12:     time.FixedZone("UTC-12", -12*60*60),
	}
)

// GetFixedTimeZone returns fixed timezone, if timezone not exist, use Zero timezone
func GetFixedTimeZone(zone TimeZone) *time.Location {
	if _, exist := timezoneData[zone]; !exist {
		zone = LocationUTC
	}

	return timezoneData[zone]
}

// NowInLocalTime return current time in Local timezone, if Local not set, use Zero timezone
func NowInLocalTime() time.Time {
	return time.Now().In(timezoneData[Local])
}

// NowInZeroTime return current time in Zero timezone
func NowInZeroTime() time.Time {
	return time.Now().In(timezoneData[Zero])
}

// NowInTimezone return current time in specified timezone, if timezone not exist, use Zero timezone
func NowInTimezone(location TimeZone) time.Time {
	if _, exist := timezoneData[location]; !exist {
		return NowInZeroTime()
	}

	return time.Now().In(timezoneData[location])
}

// NowInLocalTimeUnix return current timestamp with unix format in Local timezone, if Local not set, use Zero timezone
func NowInLocalTimeUnix() int64 {
	return NowInLocalTime().Unix()
}

// NowInZeroTimeUnix return current timestamp with unix format in Zero timezone
func NowInZeroTimeUnix() int64 {
	return NowInZeroTime().Unix()
}

// NowInTimezoneUnix return current timestamp with unix format in specified timezone, if timezone not exist, use Zero timezone
func NowInTimezoneUnix(location TimeZone) int64 {
	return NowInTimezone(location).Unix()
}

// NowInLocalTimeUnixMilli return current timestamp with unix millisecond format in Local timezone, if Local not set, use Zero timezone
func NowInLocalTimeUnixMilli() int64 {
	return NowInLocalTime().UnixMilli()
}

// NowInZeroTimeUnixMilli return current timestamp with unix millisecond format in Zero timezone
func NowInZeroTimeUnixMilli() int64 {
	return NowInZeroTime().UnixMilli()
}

// NowInTimezoneUnixMilli return current timestamp with unix millisecond format in specified timezone, if timezone not exist, use Zero timezone
func NowInTimezoneUnixMilli(location TimeZone) int64 {
	return NowInTimezone(location).UnixMilli()
}

// NowInLocalTimeUnixMicro return current timestamp with unix microsecond format in Local timezone, if Local not set, use Zero timezone
func NowInLocalTimeUnixMicro() int64 {
	return NowInLocalTime().UnixMicro()
}

// NowInZeroTimeUnixMicro return current timestamp with unix microsecond format in Zero timezone
func NowInZeroTimeUnixMicro() int64 {
	return NowInZeroTime().UnixMicro()
}

// NowInTimezoneUnixMicro return current timestamp with unix microsecond format in specified timezone, if timezone not exist, use Zero timezone
func NowInTimezoneUnixMicro(location TimeZone) int64 {
	return NowInTimezone(location).UnixMicro()
}

// NowInLocalTimeUnixNano return current timestamp with unix nanosecond format in Local timezone, if Local not set, use Zero timezone
func NowInLocalTimeUnixNano() int64 {
	return NowInLocalTime().UnixNano()
}

// NowInZeroTimeUnixNano return current timestamp with unix nanosecond format in Zero timezone
func NowInZeroTimeUnixNano() int64 {
	return NowInZeroTime().UnixNano()
}

// NowInTimezoneUnixNano return current timestamp with unix nanosecond format in specified timezone, if timezone not exist, use Zero timezone
func NowInTimezoneUnixNano(location TimeZone) int64 {
	return NowInTimezone(location).UnixNano()
}

// UnixTimestampInLocal return time in Local timezone with unix format
func UnixTimestampInLocal(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(timezoneData[Local])
}

// UnixTimestampInZero return time in Zero timezone with unix format
func UnixTimestampInZero(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(timezoneData[Zero])
}

// UnixTimestampInTimezone return time in specified timezone with unix format, if timezone not exist, use Zero timezone
func UnixTimestampInTimezone(timestamp int64, location TimeZone) time.Time {
	if _, exist := timezoneData[location]; !exist {
		return time.Unix(timestamp, 0).In(timezoneData[Zero])
	}

	return time.Unix(timestamp, 0).In(timezoneData[location])
}

// UnixMilliTimestampInLocal return time in Local timezone with unix millisecond format
func UnixMilliTimestampInLocal(timestamp int64) time.Time {
	return time.UnixMilli(timestamp).In(timezoneData[Local])
}

// UnixMilliTimestampInZero return time in Zero timezone with unix millisecond format
func UnixMilliTimestampInZero(timestamp int64) time.Time {
	return time.UnixMilli(timestamp).In(timezoneData[Zero])
}

// UnixMilliTimestampInTimezone return time in specified timezone with unix millisecond format, if timezone not exist, use Zero timezone
func UnixMilliTimestampInTimezone(timestamp int64, location TimeZone) time.Time {
	if _, exist := timezoneData[location]; !exist {
		return time.UnixMilli(timestamp).In(timezoneData[Zero])
	}

	return time.UnixMilli(timestamp).In(timezoneData[location])
}

// UnixMicroTimestampInLocal return time in Local timezone with unix microsecond format
func UnixMicroTimestampInLocal(timestamp int64) time.Time {
	return time.UnixMicro(timestamp).In(timezoneData[Local])
}

// UnixMicroTimestampInZero return time in Zero timezone with unix microsecond format
func UnixMicroTimestampInZero(timestamp int64) time.Time {
	return time.UnixMicro(timestamp).In(timezoneData[Zero])
}

// UnixMicroTimestampInTimezone return time in specified timezone with unix microsecond format, if timezone not exist, use Zero timezone
func UnixMicroTimestampInTimezone(timestamp int64, location TimeZone) time.Time {
	if _, exist := timezoneData[location]; !exist {
		return time.UnixMicro(timestamp).In(timezoneData[Zero])
	}

	return time.UnixMicro(timestamp).In(timezoneData[location])
}

// UnixNanoTimestampInLocal return time in Local timezone with unix nanosecond format
func UnixNanoTimestampInLocal(timestamp int64) time.Time {
	return time.Unix(timestamp/1e9, (timestamp%1e9)*1e3).In(timezoneData[Local])
}

// UnixNanoTimestampInZero return time in Zero timezone with unix nanosecond format
func UnixNanoTimestampInZero(timestamp int64) time.Time {
	return time.Unix(timestamp/1e9, (timestamp%1e9)*1e3).In(timezoneData[Zero])
}

// UnixNanoTimestampInTimezone return time in specified timezone with unix nanosecond format, if timezone not exist, use Zero timezone
func UnixNanoTimestampInTimezone(timestamp int64, location TimeZone) time.Time {
	if _, exist := timezoneData[location]; !exist {
		return time.Unix(timestamp/1e9, (timestamp%1e9)*1e3).In(timezoneData[Zero])
	}

	return time.Unix(timestamp/1e9, (timestamp%1e9)*1e3).In(timezoneData[location])
}

// FixedTimestampUnixInZero return time in Zero timezone with unix format timestamp from specified timezone
func FixedTimestampUnixInZero(timestamp int64, fixedLocation TimeZone) time.Time {
	if _, exist := timezoneData[fixedLocation]; !exist {
		return time.Unix(timestamp, 0).In(timezoneData[Zero])
	}

	return time.Unix(timestamp, 0).In(timezoneData[fixedLocation]).In(timezoneData[Zero])
}

// FixedTimestampUnixMilliInZero return time in Zero timezone with unix millisecond format timestamp from specified timezone
func FixedTimestampUnixMilliInZero(timestamp int64, fixedLocation TimeZone) int64 {
	if _, exist := timezoneData[fixedLocation]; !exist {
		return time.UnixMilli(timestamp).In(timezoneData[Zero]).UnixMilli()
	}

	return time.UnixMilli(timestamp).In(timezoneData[fixedLocation]).In(timezoneData[Zero]).UnixMilli()
}

// FixedTimestampUnixMicroInZero return time in Zero timezone with unix microsecond format timestamp from specified timezone
func FixedTimestampUnixMicroInZero(timestamp int64, fixedLocation TimeZone) int64 {
	if _, exist := timezoneData[fixedLocation]; !exist {
		return time.UnixMicro(timestamp).In(timezoneData[Zero]).UnixMicro()
	}

	return time.UnixMicro(timestamp).In(timezoneData[fixedLocation]).In(timezoneData[Zero]).UnixMicro()
}

// FixedTimestampUnixNanoInZero return time in Zero timezone with unix nanosecond format timestamp from specified timezone
func FixedTimestampUnixNanoInZero(timestamp int64, fixedLocation TimeZone) int64 {
	if _, exist := timezoneData[fixedLocation]; !exist {
		return time.Unix(timestamp/1e9, (timestamp%1e9)*1e3).In(timezoneData[Zero]).UnixNano()
	}

	return time.Unix(timestamp/1e9, (timestamp%1e9)*1e3).In(timezoneData[fixedLocation]).In(timezoneData[Zero]).UnixNano()
}
