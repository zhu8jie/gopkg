package xip

import (
	"fmt"
	"testing"
)

func TestGetAreaByIPLong(t *testing.T) {
	ip := "120.229.5.193"
	ipLong := Ip2Long(ip)
	fmt.Println(ipLong)

	ipArea := GetAreaByIPLong(ipLong)
	fmt.Println(ipArea == nil)
	fmt.Println(ipArea.Longitude, ipArea.Latitude, ipArea.AreaID)

}
