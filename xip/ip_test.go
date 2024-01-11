package xip

import (
	"fmt"
	"testing"
)

func TestGetAreaByIPLong(t *testing.T) {
	ip := "111.206.120.199"
	ipLong := Ip2Long(ip)
	fmt.Println(ipLong)

	ipArea := GetAreaByIPLong(ipLong)
	fmt.Println(ipArea == nil)
	fmt.Println(ipArea.Longitude, ipArea.Latitude)

}
