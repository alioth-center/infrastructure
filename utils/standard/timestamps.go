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
	// TimestampFormatRFC822 RFC 822 时间格式，更新于RFC 1123以允许4位年份
	TimestampFormatRFC822 TimestampFormat = "Mon, 02 Jan 06 15:04:05 MST"

	// TimestampFormatRFC850 RFC 850 时间格式，后来被RFC 1036取代
	TimestampFormatRFC850 TimestampFormat = "Monday, 02-Jan-06 15:04:05 MST"

	// TimestampFormatRFC1036 RFC 1036 时间格式，更新了RFC 850的时间格式
	TimestampFormatRFC1036 TimestampFormat = "Mon, 02 Jan 06 15:04:05 MST"

	// TimestampFormatRFC1123 RFC 1123 时间格式，基于RFC 822，广泛用于HTTP和其他互联网协议
	TimestampFormatRFC1123 TimestampFormat = "Mon, 02 Jan 2006 15:04:05 MST"

	// TimestampFormatRFC2822 RFC 2822 时间格式，用于电子邮件，是RFC 822的后继
	TimestampFormatRFC2822 TimestampFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

	// TimestampFormatRFC3339 RFC 3339 时间格式，基于ISO 8601，用于互联网时间戳
	TimestampFormatRFC3339 TimestampFormat = "2006-01-02T15:04:05Z07:00"

	// TimestampFormatRFC5322 RFC 5322 时间格式，更新了RFC 2822，用于电子邮件
	TimestampFormatRFC5322 TimestampFormat = "Mon, 02 Jan 2006 15:04:05 -0700"
)
