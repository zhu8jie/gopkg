package xip

import (
	"fmt"
	"testing"
)

func TestGetAreaByIPLong(t *testing.T) {
	Init("./ip.labe")

	fmt.Println(GetChineseProvinceCode())

	ip := "116.176.191.122"
	ipLong := Ip2Long(ip)
	fmt.Println(ipLong)

	ipArea := GetAreaByIPLong(ipLong)
	fmt.Println(ipArea == nil)
	if ipArea != nil {
		fmt.Println(ipArea.Longitude, ipArea.Latitude, ipArea.AreaID)
	}

}
