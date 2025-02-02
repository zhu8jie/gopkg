package xip

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/zhu8jie/gopkg/xutils"
)

type IpArea struct {
	StartLong int64
	EndLong   int64
	AreaID    int64
	Longitude float64
	Latitude  float64
}

var ipAreaTable []*IpArea
var areaIds []int64

// GetAreaIDByIPLong 二分查找
func GetAreaByIPLong(ipLong int64) *IpArea {
	if len(ipAreaTable) == 0 {
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

func Init(ipFile string) error {
	fi, err := os.Open(ipFile)
	if err != nil {
		// fmt.Printf("init file open error: %s\n", err)
		return err
	}
	defer fi.Close()

	areaIdMap := make(map[int64]struct{})

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		data := strings.Split(string(a), ",")
		if len(data) == 5 {
			startLong := xutils.StrToInt64(data[0])
			endLong := xutils.StrToInt64(data[1])
			areaID := xutils.StrToInt64(data[2])
			longitude := xutils.StrToFloat64(data[3])
			latitude := xutils.StrToFloat64(data[4])

			ipAreaTable = append(ipAreaTable, &IpArea{
				StartLong: startLong,
				EndLong:   endLong,
				AreaID:    areaID,
				Longitude: longitude,
				Latitude:  latitude,
			})
			if _, exist := areaIdMap[areaID]; !exist {
				areaIds = append(areaIds, areaID)
				areaIdMap[areaID] = struct{}{}
			}

		}
	}
	return nil
}

func Ip2Long(ip string) int64 {
	ipSegs := strings.Split(ip, ".")
	var ipInt int64 = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt := xutils.StrToInt64(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

var chineseProvinceCode = []int64{
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

func GetChineseProvinceCode() []int64 {
	return areaIds
}
