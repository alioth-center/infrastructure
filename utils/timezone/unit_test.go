package timezone

import (
	"testing"
	"time"
)

func TestTimeZone(t *testing.T) {
	_ = SetLocal(LocationShanghai)
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now()
	t.Log("now:", now)
	utc := now.In(time.UTC).Unix()
	t.Log("utc:", time.Unix(utc, 0))
	tky := now.In(loc).Unix()
	t.Log("tokyo:", time.Unix(tky, 0))
	fixednow := TimestampInLocal(utc)
	t.Log("fixed now:", fixednow)
	convertnow := FixedTimestampInLocal(utc, LocationTokyo)
	t.Log("convert now:", convertnow)
}
