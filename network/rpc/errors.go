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
