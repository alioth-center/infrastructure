package network

import (
	"net"
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
