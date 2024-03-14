package influx2kit

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/spf13/cast"
	"sync"
)

var once sync.Once
var client influxdb2.Client

func DefaultClient() influxdb2.Client {
	once.Do(func() {
		option := influxdb2.DefaultOptions()
		option.HTTPOptions().SetHTTPRequestTimeout(cast.ToUint(configkit.GetInt(configkey.InfluxReqTimeout)))
		client = influxdb2.NewClientWithOptions(configkit.GetString(configkey.InfluxURL), configkit.GetString(configkey.InfluxToken), option)
	})
	return client
}
