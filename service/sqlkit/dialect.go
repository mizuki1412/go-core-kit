package sqlkit

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

func GetDataSourceName(p DataSourceParam) (string, string) {
	if p.Driver == "" || p.Host == "" || p.Port == "" {
		panic(exception.New("sqlkit: database config error"))
	}
	var param string
	switch p.Driver {
	case "postgres":
		param = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Pwd, p.Name)
	case "mysql":
		param = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s", p.User, p.Pwd, p.Host, p.Port, p.Name, "Asia%2FShanghai")
	case "mssql":
		param = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s", p.Host, p.User, p.Pwd, p.Port, p.Name)
	default:
		panic(exception.New("driver not supported"))
	}
	return p.Driver, param
}

func getStatementBuilderType(ds *DataSource) squirrel.StatementBuilderType {
	switch ds.Driver {
	case "postgres":
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	default:
		// todo 未处理oracle
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
}
