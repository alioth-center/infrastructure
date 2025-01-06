package standard

import "time"

// TimestampFormat timestamp format, contains enums in [RFC 822], [RFC 850], [RFC 1036], [RFC 1123], [RFC 2822], [RFC 3339], [RFC 5322]
//
// [RFC 822]: https://datatracker.ietf.org/doc/html/rfc822
// [RFC 850]: https://datatracker.ietf.org/doc/html/rfc850
// [RFC 1036]: https://datatracker.ietf.org/doc/html/rfc1036
// [RFC 1123]: https://datatracker.ietf.org/doc/html/rfc1123
// [RFC 2822]: https://datatracker.ietf.org/doc/html/rfc2822
// [RFC 3339]: https://datatracker.ietf.org/doc/html/rfc3339
// [RFC 5322]: https://datatracker.ietf.org/doc/html/rfc5322
type TimestampFormat = string

const (
	// TimestampFormatRFC822 RFC 822 time format, see [RFC 822]
	//
	// example: Mon, 02 Jan 06 15:04:05 MST
	//
	// [RFC 822]: https://datatracker.ietf.org/doc/html/rfc822
	TimestampFormatRFC822 TimestampFormat = "Mon, 02 Jan 06 15:04:05 MST"

	// TimestampFormatRFC850 RFC 850 time format, see [RFC 850]
	//
	// example: Monday, 02-Jan-06 15:04:05 MST
	//
	// [RFC 850]: https://datatracker.ietf.org/doc/html/rfc850
	TimestampFormatRFC850 TimestampFormat = "Monday, 02-Jan-06 15:04:05 MST"

	// TimestampFormatRFC1036 RFC 1036 time format, see [RFC 1036]
	//
	// example: Mon, 02 Jan 06 15:04:05 MST
	//
	// [RFC 1036]: https://datatracker.ietf.org/doc/html/rfc1036
	TimestampFormatRFC1036 TimestampFormat = "Mon, 02 Jan 06 15:04:05 MST"

	// TimestampFormatRFC1123 RFC 1123 time format, see [RFC 1123]
	//
	// example: Mon, 02 Jan 2006 15:04:05 MST
	//
	// [RFC 1123]: https://datatracker.ietf.org/doc/html/rfc1123
	TimestampFormatRFC1123 TimestampFormat = "Mon, 02 Jan 2006 15:04:05 MST"

	// TimestampFormatRFC2822 RFC 2822 time format, see [RFC 2822]
	//
	// example: Mon, 02 Jan 2006 15:04:05 -0700
	//
	// [RFC 2822]: https://datatracker.ietf.org/doc/html/rfc2822
	TimestampFormatRFC2822 TimestampFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

	// TimestampFormatRFC3339 RFC 3339 time format, see [RFC 3339]
	//
	// example: 2006-01-02T15:04:05Z07:00
	//
	// [RFC 3339]: https://datatracker.ietf.org/doc/html/rfc3339
	TimestampFormatRFC3339 TimestampFormat = "2006-01-02T15:04:05Z07:00"

	// TimestampFormatRFC5322 RFC 5322 time format, see [RFC 5322]
	//
	// example: Mon, 02 Jan 2006 15:04:05 -0700
	//
	// [RFC 5322]: https://datatracker.ietf.org/doc/html/rfc5322
	TimestampFormatRFC5322 TimestampFormat = "Mon, 02 Jan 2006 15:04:05 -0700"
)

const (
	// Minute represents a minute
	Minute = time.Minute

	// FiveMinute represents five minutes
	FiveMinute = time.Minute * 5

	// TenMinute represents ten minutes
	TenMinute = time.Minute * 10

	// FifteenMinute represents fifteen minutes
	FifteenMinute = time.Minute * 15

	// HalfHour represents half an hour
	HalfHour = time.Minute * 30

	// Hour represents an hour
	Hour = time.Hour

	// HalfDay represents half a day
	HalfDay = time.Hour * 12

	// Day represents a day
	Day = time.Hour * 24

	// ThreeDay represents three days
	ThreeDay = time.Hour * 24 * 3

	// Week represents a week
	Week = time.Hour * 24 * 7

	// TwoWeek represents two weeks
	TwoWeek = time.Hour * 24 * 14

	// Month represents a month
	Month = time.Hour * 24 * 30

	// Season represents a season, 3 months, 90 days
	Season = time.Hour * 24 * 30 * 3

	// HalfYear represents half a year, 6 months, 180 days
	HalfYear = time.Hour * 24 * 30 * 6

	// Year represents a year, 12 months, 365 days
	Year = time.Hour * 24 * 365
)
