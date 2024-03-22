package values

import "unsafe"

// UnsafeStringToBytes converts a string to a byte slice without memory allocation. before you use this function, you should know what you are doing.
//   - see also: [k8s usage]
//
// [k8s usage]: https://github.com/kubernetes/apiserver/blob/2a8bc69060e4f3b030c957c5172c0957b4fcd80e/pkg/authentication/token/cache/cached_token_authenticator.go#L277-L286
func UnsafeStringToBytes(s string) []byte {
	if len(s) == 0 {
		return []byte{}
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeBytesToString converts a byte slice to a string without copying the data. before you use this function, you should know what you are doing.
//   - see also: [k8s usage]
//
// [k8s usage]: https://github.com/kubernetes/apiserver/blob/2a8bc69060e4f3b030c957c5172c0957b4fcd80e/pkg/authentication/token/cache/cached_token_authenticator.go#L288-L297
func UnsafeBytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
