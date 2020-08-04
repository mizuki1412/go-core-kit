package locationkit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/httpkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/tidwall/gjson"
	"strings"
)

type Location struct {
	Lon          class.Decimal `json:"lon" description:"经度"`
	Lat          class.Decimal `json:"lat" description:"纬度"`
	ProvinceName string        `json:"provinceName"`
	CityName     string        `json:"cityName"`
	AreaName     string        `json:"areaName"`
	ProvinceCode string        `json:"provinceCode"`
	CityCode     string        `json:"cityCode"`
	AreaCode     string        `json:"areaCode"`
	Address      string        `json:"address"`
}

// 只返回lat/lon
func Geo(cityCode, address string) (loc *Location) {
	defer func() {
		if err := recover(); err != nil {
			loc = nil
		}
	}()
	// 应用amap
	if len(cityCode) == 4 {
		cityCode += "00"
	}
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=%s&address=%s&city=%s", configkit.GetStringD(ConfigKeyAmapKey), address, cityCode)
	ret, _ := httpkit.Request(httpkit.Req{
		Method: httpkit.MethodGet,
		Url:    url,
	})
	if gjson.Get(ret, "status").String() == "1" && len(gjson.Get(ret, "geocodes").Array()) > 0 {
		l := gjson.Get(ret, "geocodes").Array()[0].Get("location").String()
		if l != "" {
			arrs := stringkit.Split(l, ",")
			if len(arrs) == 2 {
				location := &Location{}
				location.Lon.Set(arrs[0])
				location.Lat.Set(arrs[1])
				return location
			}
		} else {
			panic(exception.New("geo location result null"))
		}
	} else {
		panic(exception.New(gjson.Get(ret, "info").String()))
	}
	panic(exception.New("geo 失败"))
}

// 不返回lat/lon
func ReGeo(lon, lat class.Decimal) (loc *Location) {
	defer func() {
		if err := recover(); err != nil {
			loc = nil
		}
	}()
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%s,%s", configkit.GetStringD(ConfigKeyAmapKey), lon.Decimal.String(), lat.Decimal.String())
	ret, _ := httpkit.Request(httpkit.Req{
		Method: httpkit.MethodGet,
		Url:    url,
	})
	res1 := gjson.Get(ret, "regeocode")
	if gjson.Get(ret, "status").String() == "1" && res1.Exists() {
		location := &Location{}
		address := res1.Get("formatted_address").String()
		res2 := res1.Get("addressComponent")
		if res2.Exists() {
			location.ProvinceName = res2.Get("province").String()
			location.CityName = res2.Get("city").String()
			location.AreaName = res2.Get("district").String()
			location.Address = strings.ReplaceAll(address, location.ProvinceName+location.CityName+location.AreaName, "")
			adcode := res2.Get("adcode").String()
			if adcode != "" {
				location.ProvinceCode = adcode[0:2]
				location.CityCode = adcode[0:4]
				location.AreaCode = adcode
				return location
			}
		}
	}
	panic(exception.New("regeo 失败"))
}

// todo weather

var header = map[string]string{
	"Accept":     "*/*",
	"DNT":        "1",
	"Origin":     "http://cellid.cn",
	"Referer":    "http://cellid.cn/",
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36",
}

// 基站转， 十进制的lac和ci
func LacCiTransfer(lac, ci int32) (loc *Location) {
	defer func() {
		if err := recover(); err != nil {
			loc = nil
		}
	}()
	ret, _ := httpkit.Request(httpkit.Req{
		Method: httpkit.MethodGet,
		Url:    "http://cellid.cn/",
		Header: header,
	})
	arr := strings.Split(ret, "<input type=\"hidden\" id=\"flag\" name=\"flag\" value=\"")
	if len(arr) < 2 {
		panic(exception.New("lac ci transfer error 1"))
	}
	flag := ""
	for i := 0; i < len(arr[1]); i++ {
		c := arr[1][i : i+1]
		if c != "\"" {
			flag += c
		} else {
			break
		}
	}
	ret, _ = httpkit.Request(httpkit.Req{
		Method: httpkit.MethodPost,
		Url:    fmt.Sprintf("http://cellid.cn/cidInfo.php?lac=%d&cell_id=%d&hex=false&flag=%s", lac, ci, flag),
		Header: header,
	})
	// cidMap(30.30xxxx,120.xxxx,
	arr = strings.Split(ret, ",")
	if len(arr) < 2 || !strings.Contains(arr[0], "(") {
		panic(exception.New("lac ci transfer error 1"))
	}
	location := &Location{}
	location.Lon.Set(arr[1])
	location.Lat.Set(strings.Split(arr[0], "(")[1])
	return location
}
