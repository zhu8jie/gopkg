package xip

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type IpArea struct {
	StartLong int
	EndLong   int
	AreaID    int
	Longitude float64
	Latitude  float64
}

var ipAreaTable []*IpArea

//GetAreaIDByIPLong 二分查找
func GetAreaByIPLong(ipLong int) *IpArea {
	if len(ipAreaTable) == 0 {
		fmt.Println("yanw_test")
		return nil
	}
	//查找start之后的
	l := len(ipAreaTable)
	i, j, m := 0, l-1, (l-1)/2
	for i < m {
		//判断是否在此区间
		//如果ip大于等于中start，小于等于中end，在中
		//如果ip小于中start，在左
		//如果ip大于中end，在右
		if ipLong < ipAreaTable[m].StartLong {
			//左
			j = m
			m = (i + j) / 2
		} else if ipLong > ipAreaTable[m].EndLong {
			//右
			i = m
			m = (i + j) / 2
		} else {
			//中，命中
			return ipAreaTable[m]
		}
	}
	return nil
}

func init() {
	fi, err := os.Open("./ip.labe")
	if err != nil {
		fmt.Printf("init file open error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		data := strings.Split(string(a), ",")
		if len(data) == 5 {
			startLong, _ := strconv.Atoi(data[0])
			endLong, _ := strconv.Atoi(data[1])
			areaID, _ := strconv.Atoi(data[2])
			longitude, _ := strconv.ParseFloat(data[3], 64)
			latitude, _ := strconv.ParseFloat(data[4], 64)

			ipAreaTable = append(ipAreaTable, &IpArea{
				StartLong: startLong,
				EndLong:   endLong,
				AreaID:    areaID,
				Longitude: longitude,
				Latitude:  latitude,
			})
		}
	}
}

func Ip2Long(ip string) int {
	ipSegs := strings.Split(ip, ".")
	var ipInt int = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

var chineseProvinceCode = []int{
	137101101100100,
	137101102100100,
	137101103100100,
	137101104100100,
	137101105100100,
	137102101100100,
	137102102100100,
	137102103100100,
	137103101100100,
	137103102100100,
	137103103100100,
	137103104100100,
	137103105100100,
	137103106100100,
	137103107100100,
	137104101100100,
	137104102100100,
	137104103100100,
	137104104100100,
	137104105100100,
	137104106100100,
	137105101100100,
	137105102100100,
	137105103100100,
	137105104100100,
	137105105100100,
	137106101100100,
	137106102100100,
	137106103100100,
	137106104100100,
	137106105100100,
	137107101100100,
	137107102100100,
	137107103100100,
}
