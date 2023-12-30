package tag

// ParamIn openapi的参数路径in参数
var ParamIn = Tag{Name: "in"}

// Schema openapi schema配置
var Schema = Tag{Name: "schema"}

// SchemaIgnore openapi不序列化
const SchemaIgnore = "ignore"

// RetData 专用于RestRet，自定义时表示data所在的区域
var RetData = Tag{Name: "data", Value: "true"}
