package utils

import (
	"github.com/alioth-center/infrastructure/errors"
	"time"
)

var (
	isSetLocal = false

	Zero  TimeZone = LocationLondon // 零时区
	Local TimeZone = Zero           // 默认时区，零时区
)

func SetLocal(timezone TimeZone) (err error) {
	if !isSetLocal {
		Local = timezone
		isSetLocal = true
		return nil
	} else if Local == timezone {
		return nil
	} else {
		return errors.NewLocalTimezoneAlreadySetError()
	}
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
	}
)

func LocalTimeNow() time.Time {
	return time.Now().In(timezoneData[Local])
}

func ZeroTimeNow() time.Time {
	return time.Now().In(timezoneData[Zero])
}

func TimeInLocalUnix(t time.Time) int64 {
	return t.In(timezoneData[Local]).Unix()
}

func TimeInLocationUnix(t time.Time, location TimeZone) int64 {
	return t.In(timezoneData[location]).Unix()
}

func TimeInZeroUnix(t time.Time) int64 {
	return t.In(timezoneData[Zero]).Unix()
}

func TimestampInLocal(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(timezoneData[Local])
}

func TimestampInZero(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(timezoneData[Zero])
}

func TimeStampInLocation(timestamp int64, location TimeZone) time.Time {
	return time.Unix(timestamp, 0).In(timezoneData[location])
}

func FixedTimestampInLocal(timestamp int64, location TimeZone) time.Time {
	return TimeStampInLocation(timestamp, location).In(timezoneData[Local])
}
