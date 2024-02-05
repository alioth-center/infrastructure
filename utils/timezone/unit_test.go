package timezone

import (
	"testing"
)

func TestTimeZone(t *testing.T) {
	t.Run("SetLocalTimeZone", func(t *testing.T) {
		_ = SetLocal("Asia/ShanghaiNotExist")
		_ = SetLocal("Asia/Shanghai")
		_ = SetLocal("Asia/Shanghai")
		_ = SetLocal("Asia/Tokyo")
	})

	t.Run("GetFixedTimeZone", func(t *testing.T) {
		_ = GetFixedTimeZone("Asia/Shanghai")
		_ = GetFixedTimeZone("Asia/ShanghaiNotExist")
	})

	t.Run("NowIn", func(t *testing.T) {
		NowInLocalTime()
		NowInZeroTime()
		NowInTimezone(LocationUTC)
		NowInTimezone("NOT_EXIST")
	})

	t.Run("NowInUnix", func(t *testing.T) {
		NowInLocalTimeUnix()
		NowInZeroTimeUnix()
		NowInTimezoneUnix(LocationUTC)
		NowInTimezoneUnix("NOT_EXIST")
	})

	t.Run("NowInUnixMilli", func(t *testing.T) {
		NowInLocalTimeUnixMilli()
		NowInZeroTimeUnixMilli()
		NowInTimezoneUnixMilli(LocationUTC)
		NowInTimezoneUnixMilli("NOT_EXIST")
	})

	t.Run("NowInUnixMicro", func(t *testing.T) {
		NowInLocalTimeUnixMicro()
		NowInZeroTimeUnixMicro()
		NowInTimezoneUnixMicro(LocationUTC)
		NowInTimezoneUnixMicro("NOT_EXIST")
	})

	t.Run("NowInUnixNano", func(t *testing.T) {
		NowInLocalTimeUnixNano()
		NowInZeroTimeUnixNano()
		NowInTimezoneUnixNano(LocationUTC)
		NowInTimezoneUnixNano("NOT_EXIST")
	})

	t.Run("UnixIn", func(t *testing.T) {
		UnixTimestampInLocal(0)
		UnixTimestampInZero(0)
		UnixTimestampInTimezone(0, LocationUTC)
		UnixTimestampInTimezone(0, "NOT_EXIST")
	})

	t.Run("UnixMilliIn", func(t *testing.T) {
		UnixMilliTimestampInLocal(0)
		UnixMilliTimestampInZero(0)
		UnixMilliTimestampInTimezone(0, LocationUTC)
		UnixMilliTimestampInTimezone(0, "NOT_EXIST")
	})

	t.Run("UnixMicroIn", func(t *testing.T) {
		UnixMicroTimestampInLocal(0)
		UnixMicroTimestampInZero(0)
		UnixMicroTimestampInTimezone(0, LocationUTC)
		UnixMicroTimestampInTimezone(0, "NOT_EXIST")
	})

	t.Run("UnixNanoIn", func(t *testing.T) {
		UnixNanoTimestampInLocal(0)
		UnixNanoTimestampInZero(0)
		UnixNanoTimestampInTimezone(0, LocationUTC)
		UnixNanoTimestampInTimezone(0, "NOT_EXIST")
	})

	t.Run("FixedTimestampIn", func(t *testing.T) {
		FixedTimestampUnixInZero(0, LocationUTC)
		FixedTimestampUnixInZero(0, "NOT_EXIST")
		FixedTimestampUnixMilliInZero(0, LocationUTC)
		FixedTimestampUnixMilliInZero(0, "NOT_EXIST")
		FixedTimestampUnixMicroInZero(0, LocationUTC)
		FixedTimestampUnixMicroInZero(0, "NOT_EXIST")
		FixedTimestampUnixNanoInZero(0, LocationUTC)
		FixedTimestampUnixNanoInZero(0, "NOT_EXIST")
	})
}
