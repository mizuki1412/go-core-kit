package locationkit

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/v2/class"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/library/httpkit"
	"github.com/mizuki1412/go-core-kit/v2/library/stringkit"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

type Location struct {
	Lon          class.Decimal `json:"lon" comment:"经度" precision:"8"`
	Lat          class.Decimal `json:"lat" comment:"纬度" precision:"8"`
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
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=%s&address=%s&city=%s", configkit.GetString(configkey.AmapKey), address, cityCode)
	ret, _ := httpkit.Request(httpkit.Req{
		Method: http.MethodGet,
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
	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%s,%s", configkit.GetString(configkey.AmapKey), lon.Decimal.String(), lat.Decimal.String())
	ret, _ := httpkit.Request(httpkit.Req{
		Method: http.MethodGet,
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

// params：city code
func Weather(city string) []map[string]any {
	ret, _ := httpkit.Request(httpkit.Req{
		Method: http.MethodGet,
		Url:    "https://restapi.amap.com/v3/weather/weatherInfo?key=" + configkit.GetString(configkey.AmapKey) + "&city=" + city + "&extensions=all",
	})
	rs := gjson.Parse(ret).Get("forecasts").Array()
	cast := rs[0].Get("casts").Array()
	data := make([]map[string]any, 0, len(cast))
	for _, v := range cast {
		tmp := v.Map()
		m := map[string]any{}
		for k, v := range tmp {
			m[k] = v.Value()
		}
		data = append(data, m)
	}
	return data
}

var header = map[string]string{
	"Accept":     "*/*",
	"DNT":        "1",
	"Origin":     "http://cellid.cn",
	"Referer":    "http://cellid.cn/",
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36",
}

// 基站转， 十进制的lac和ci
// http://www.cellocation.com/api/
func LacCiTransfer(mnc, lac, ci int32) (loc *Location) {
	defer func() {
		if err := recover(); err != nil {
			loc = nil
		}
	}()
	ret, _ := httpkit.Request(httpkit.Req{
		Method: http.MethodGet,
		Url:    fmt.Sprintf("http://api.cellocation.com:81/cell/?mcc=460&mnc=%d&lac=%d&ci=%d&output=json", mnc, lac, ci),
		Header: header,
	})
	if ret != "" && gjson.Get(ret, "errcode").Int() == 0 {
		location := &Location{}
		location.Lon.Set(gjson.Get(ret, "lon").Float())
		location.Lat.Set(gjson.Get(ret, "lat").Float())
		return location
	}
	return nil
}

//func LacCiTransfer(lac, ci int32) (loc *Location) {
//	defer func() {
//		if err := recover(); err != nil {
//			loc = nil
//		}
//	}()
//	ret, _ := httpkit.Request(httpkit.Req{
//		Method: http.MethodGet,
//		Url:    "http://cellid.cn/",
//		Header: header,
//	})
//	arr := strings.Split(ret, "<input type=\"hidden\" id=\"flag\" name=\"flag\" value=\"")
//	if len(arr) < 2 {
//		panic(exception.New("lac ci transfer error 1"))
//	}
//	flag := ""
//	for i := 0; i < len(arr[1]); i++ {
//		c := arr[1][i : i+1]
//		if c != "\"" {
//			flag += c
//		} else {
//			break
//		}
//	}
//	ret, _ = httpkit.Request(httpkit.Req{
//		Method: http.MethodPost,
//		Url:    fmt.Sprintf("http://cellid.cn/cidInfo.php?lac=%d&cell_id=%d&hex=false&flag=%s", lac, ci, flag),
//		Header: header,
//	})
//	// cidMap(30.30xxxx,120.xxxx,
//	arr = strings.Split(ret, ",")
//	if len(arr) < 2 || !strings.Contains(arr[0], "(") {
//		panic(exception.New("lac ci transfer error 1"))
//	}
//	location := &Location{}
//	location.Lon.Set(arr[1])
//	location.Lat.Set(strings.Split(arr[0], "(")[1])
//	return location
//}
