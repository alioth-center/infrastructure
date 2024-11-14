package network

import (
	"github.com/alioth-center/infrastructure/utils/values"
	"net"
)

const (
	maxPort = 1 << 16
)

var privateIpBlocks = []*net.IPNet{
	// IPv4 private address blocks
	{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
	{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
	{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	// IPv6 private address blocks
	{IP: net.ParseIP("fc00::"), Mask: net.CIDRMask(7, 128)},
	{IP: net.ParseIP("fe80::"), Mask: net.CIDRMask(10, 128)},
}

// IsPrivateIP determines whether a given IP address is a private IP address.
// It parses the IP address from a string and checks it against a list of
// predefined private IP blocks.
//
// Parameters:
//
//	ip (string): The IP address in string format to be checked.
//
// Returns:
//
//	isPrivate (bool): Returns true if the IP address is a private IP address,
//	                  false otherwise.
func IsPrivateIP(ip string) (isPrivate bool) {
	// Parse string ip to net.IP
	addr := net.ParseIP(ip)
	if addr == nil {
		// Invalid IP string
		return false
	}

	// Check if IP is in private IP blocks
	for _, block := range privateIpBlocks {
		if block.Contains(addr) {
			return true
		}
	}

	// If no private IP block contains the IP, return false
	return false
}

// IsValidIP determines whether a given IP address is a valid IP address.
// It parses the IP address from a string and checks if the parsed IP is not nil.
//
// Parameters:
//
//	ip (string): The IP address in string format to be checked.
//
// Returns:
//
//	isValid (bool): Returns true if the IP address is a valid IP address,
//	                false otherwise.
func IsValidIP(ip string) (isValid bool) {
	// Parse string ip to net.IP, if err is nil, it's a valid IP
	return net.ParseIP(ip) != nil
}

// IsValidIPOrCIDR determines whether a given string is a valid IP address or CIDR.
// It first parses the string as an IP address, if it fails, it then parses the string as a CIDR.
// If both parsing operations fail, the string is not a valid IP address or CIDR.
//
// Parameters:
//
//	ipOrCIDR (string): The IP address or CIDR in string format to be checked.
//
// Returns:
//
//	isValid (bool): Returns true if the string is a valid IP address or CIDR,
//	                false otherwise.
func IsValidIPOrCIDR(ipOrCIDR string) (isValid bool) {
	// Parse string ip to net.IP, if error is nil, it's a valid IP
	if net.ParseIP(ipOrCIDR) != nil {
		return true
	}

	// Parse string ip to net.IPNet, if error is nil, it's a valid CIDR
	_, _, err := net.ParseCIDR(ipOrCIDR)
	return err == nil
}

// IPInCIDR determines whether a given IP address is in a given CIDR.
// It parses the IP address and CIDR from strings and checks if the IP address is in the CIDR.
//
// Parameters:
//
//	ip (string): The IP address in string format to be checked.
//	cidr (string): The CIDR in string format to be checked.
//
// Returns:
//
//	inCIDR (bool): Returns true if the IP address is in the CIDR, false otherwise.
func IPInCIDR(ip, cidr string) (inCIDR bool) {
	// Parse IP and CIDR, if error is not nil, return false
	addr := net.ParseIP(ip)
	_, ipNet, err := net.ParseCIDR(cidr)
	if addr == nil || err != nil {
		return false
	}

	// Check if IP is in CIDR
	return ipNet.Contains(addr)
}

func IsValidHostPort(hostPort string) (isValid bool) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return false
	}

	portValue := values.StringToInt(port, -1)
	return IsValidIP(host) && portValue < maxPort && portValue > 0
}
