package configkey

// RestServerPort rest server相关配置
const RestServerPort = "rest.port"

// RestServerBase base path
const RestServerBase = "rest.base"

// RestRequestBodySize 单位MB
const RestRequestBodySize = "rest.requestBodySize"

// RestPPROF 是否开启rest server 的pprof接口： /debug/pprof
const RestPPROF = "rest.pprof"

const SessionExpire = "rest.sessionExpire"

// SessionSecure 上传cookie时是否需要https，关系到浏览器的跨域策略和具体是否用https部署服务
const SessionSecure = "rest.sessionSecure"
