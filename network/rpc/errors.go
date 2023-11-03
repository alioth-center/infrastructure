package rpc

type GetRPCClientIPFailedError struct{}

func (e GetRPCClientIPFailedError) Error() string {
	return "get rpc client ip failed"
}

func NewGetRPCClientIPFailedError() GetRPCClientIPFailedError {
	return GetRPCClientIPFailedError{}
}

type UnsupportedNetworkError struct {
	Network string
}

func (e UnsupportedNetworkError) Error() string {
	return "unsupported network: " + e.Network
}

func NewUnsupportedNetworkError(network string) UnsupportedNetworkError {
	return UnsupportedNetworkError{Network: network}
}

type InvalidIPAddressError struct {
	IPAddress string
}

func (e InvalidIPAddressError) Error() string {
	return "invalid ip address: " + e.IPAddress
}

func NewInvalidIPAddressError(ipAddress string) InvalidIPAddressError {
	return InvalidIPAddressError{IPAddress: ipAddress}
}

type ServerAlreadyServingError struct {
	Address string
}

func (e ServerAlreadyServingError) Error() string {
	return "rpc server is already serving at " + e.Address
}

func NewServerAlreadyServingError(address string) ServerAlreadyServingError {
	return ServerAlreadyServingError{Address: address}
}

type CheckRequestFailedError struct {
	Reason string
}

func (e CheckRequestFailedError) Error() string {
	return "check request failed: " + e.Reason
}

func NewCheckRequestFailedError(reason string) CheckRequestFailedError {
	return CheckRequestFailedError{Reason: reason}
}

type CheckResponseFailedError struct {
	Reason string
}

func (e CheckResponseFailedError) Error() string {
	return "check response failed: " + e.Reason
}

func NewCheckResponseFailedError(reason string) CheckResponseFailedError {
	return CheckResponseFailedError{Reason: reason}
}
