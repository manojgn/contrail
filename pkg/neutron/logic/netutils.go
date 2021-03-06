package logic

import (
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	ipV4 = 4
	ipV6 = 6
)

func getIPVersionFromCIDR(cidr string) (int8, error) {
	prefix := strings.Split(cidr, "/")
	return getIPVersion(prefix[0])
}

func getIPVersion(ipAddress string) (int8, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return -1, errors.Errorf("Invalid ip address: %s", ipAddress)
	}

	if v := ip.To4(); v != nil {
		return ipV4, nil
	}

	return ipV6, nil
}

func getIPPrefixAndPrefixLen(cidr string) (prefixIP string, prefixNetIP string, prefixLen int64, err error) {
	ip, netIP, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", 0, err
	}
	prefix := strings.Split(netIP.String(), "/")

	prefixIP = ip.String()
	prefixNetIP = prefix[0]
	prefixLen, err = strconv.ParseInt(prefix[1], 10, 64)
	if err != nil {
		return "", "", 0, err
	}

	return prefixIP, prefixNetIP, prefixLen, nil
}
