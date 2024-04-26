package standard

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
