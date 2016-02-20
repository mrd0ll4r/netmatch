package netmatch

import "net"

var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

// Key converts a net.IP to a [16]byte.
// If ip is an IPv4 address, it will be prefixed with the net.v4InV6Prefix
// to get a valid IPv6 address.
// The resulting array can then be used for the Trie.
func Key(ip net.IP) [16]byte {
	var array [16]byte

	if len(ip) == net.IPv4len {
		copy(array[:], v4InV6Prefix)
		copy(array[12:], ip)
	} else {
		copy(array[:], ip)
	}
	return array
}

// Subnet parses a subnet in CIDR notation and returns everything necessary
// to add that subnet to the Trie.
func Subnet(subnet string) ([16]byte, int, error) {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return [16]byte{}, 0, err
	}

	key := Key(ipnet.IP)
	size, _ := ipnet.Mask.Size()
	if len(ipnet.IP) == net.IPv4len {
		size = size + 12*8
	}

	return key, size, nil
}
