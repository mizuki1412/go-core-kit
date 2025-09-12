package sqlkit

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"github.com/mizuki1412/go-core-kit/v2/service/logkit"
)

type DataSource struct {
	// 用于postgres, 或oracle/dm
	Schema string
	// 数据源（事务时使用）
	TX *sqlx.Tx
	// 指定数据源（原始数据源连接池）
	DBPool *sqlx.DB
	// 连接时的driver
	Driver string
}

// ColumnSchema 表结构字段
type ColumnSchema struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Comment  string `json:"comment"`
}

var defaultDB *sqlx.DB

// DataSourceParam 创建数据源的参数
type DataSourceParam struct {
	Driver  string `json:"driver"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	User    string `json:"username"`
	Pwd     string `json:"pwd"`
	Name    string `json:"db"`
	MaxOpen int    `json:"maxOpen"`
	MaxIdle int    `json:"maxIdle"`
	MaxLife int    `json:"maxLife"`
}

func getDataSourceName(p DataSourceParam) (string, string) {
	var param string
	switch p.Driver {
	case sqlconst.Postgres, sqlconst.Kingbase:
		param = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Pwd, p.Name)
		if p.Host == "" || p.Port == "" {
			panic(exception.New("sqlkit: database config error"))
		}
	case sqlconst.Mysql:
		param = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s", p.User, p.Pwd, p.Host, p.Port, p.Name, "Asia%2FShanghai")
		if p.Host == "" || p.Port == "" {
			panic(exception.New("sqlkit: database config error"))
		}
	case sqlconst.SqlServer:
		param = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s", p.Host, p.User, p.Pwd, p.Port, p.Name)
		if p.Host == "" || p.Port == "" {
			panic(exception.New("sqlkit: database config error"))
		}
	case sqlconst.Sqlite3:
		param = p.Name
		if p.Name == "" {
			panic(exception.New("sqlkit: dbName error"))
		}
	case sqlconst.DM:
		// https://eco.dameng.com/document/dm/zh-cn/pm/go-rogramming-guide.html
		param = fmt.Sprintf("dm://%s:%s@%s:%s", p.User, p.Pwd, p.Host, p.Port)
		if p.Host == "" || p.Port == "" {
			panic(exception.New("sqlkit: database config error"))
		}
	case sqlconst.TaosSql:
		param = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", p.User, p.Pwd, p.Host, p.Port, p.Name)
		if p.Host == "" || p.Port == "" {
			panic(exception.New("sqlkit: database config error"))
		}
	case sqlconst.TaosWS:
		param = fmt.Sprintf("%s:%s@ws(%s:%s)/%s", p.User, p.Pwd, p.Host, p.Port, p.Name)
		if p.Host == "" || p.Port == "" {
			panic(exception.New("sqlkit: database config error"))
		}
	default:
		panic(exception.New("driver not supported"))
	}
	return p.Driver, param
}

// NewDataSource 创建一个数据源
func NewDataSource(param DataSourceParam) *DataSource {
	db := getDB(param)
	ds := &DataSource{
		Driver: param.Driver,
		DBPool: db,
	}
	if sqlconst.IsPostgresType(param.Driver) {
		ds.Schema = "public"
	} else if sqlconst.IsSingleDBSchema(param.Driver) {
		ds.Schema = param.Name
	}
	return ds
}

func getDB(param DataSourceParam) *sqlx.DB {
	var db *sqlx.DB
	if param.Driver == sqlconst.Sqlite3 {
		db = sqlx.MustOpen(getDataSourceName(param))
	} else {
		db = sqlx.MustConnect(getDataSourceName(param))
		if param.MaxLife > 0 {
			db.SetConnMaxLifetime(time.Duration(param.MaxLife) * time.Minute)
		}
		if param.MaxOpen > 0 {
			db.SetMaxOpenConns(param.MaxOpen)
		}
		if param.MaxIdle > 0 {
			db.SetMaxIdleConns(param.MaxIdle)
		}
	}
	err := db.Ping()
	if err != nil {
		logkit.Info("db: " + param.Name + " ping success")
	} else {
		logkit.Error("db: " + param.Name + " ping failed")
	}
	return db
}

var once sync.Once

func DefaultDataSource() *DataSource {
	once.Do(func() {
		defaultDB = getDB(DataSourceParam{
			Driver:  configkit.GetString(configkey.DBDriver),
			Host:    configkit.GetString(configkey.DBHost),
			Port:    configkit.GetString(configkey.DBPort),
			User:    configkit.GetString(configkey.DBUser),
			Pwd:     configkit.GetString(configkey.DBPwd),
			Name:    configkit.GetString(configkey.DBName),
			MaxOpen: configkit.GetInt(configkey.DBMaxOpen),
			MaxIdle: configkit.GetInt(configkey.DBMaxIdle),
			MaxLife: configkit.GetInt(configkey.DBMaxLife),
		})
	})
	driver := configkit.GetString(configkey.DBDriver)
	ds := &DataSource{
		Driver: driver,
		DBPool: defaultDB,
	}
	if sqlconst.IsPostgresType(driver) {
		ds.Schema = "public"
	} else {
		ds.Schema = configkit.GetString(configkey.DBName)
	}
	return ds
}

// DecoTableName 获取 schema 修饰的转义的tableName
func (ds *DataSource) DecoTableName(tableName string) string {
	s := ""
	if sqlconst.IsPostgresType(ds.Driver) {
		if ds.Schema != "" {
			s = ds.EscapeName(ds.Schema) + "."
		} else {
			s = "public."
		}
	} else if ds.Schema != "" && sqlconst.IsSingleDBSchema(ds.Driver) {
		s = ds.EscapeName(ds.Schema) + "."
	}
	return s + ds.EscapeName(tableName)
}

// EscapeName 表名列名的转义符添加
func (ds *DataSource) EscapeName(name string) string {
	switch ds.Driver {
	case sqlconst.Mysql:
		return "`" + name + "`"
	case sqlconst.TaosWS, sqlconst.TaosSql:
		// key=tabname时，特殊处理，不加``，如果加了就代表是自定义表字段
		if name == "tbname" {
			return name
		}
		return "`" + name + "`"
	case sqlconst.DM, sqlconst.Oracle:
		// 注意大写了
		return "\"" + strings.ToUpper(name) + "\""
	default:
		return "\"" + name + "\""
	}
}

func (ds *DataSource) Commit() {
	if ds.TX == nil {
		return
	}
	err := ds.TX.Commit()
	if err != nil {
		panic(exception.New(err.Error()))
	}
}

func (ds *DataSource) Rollback() {
	if ds.TX == nil {
		return
	}
	err := ds.TX.Rollback()
	if err != nil {
		panic(exception.New(err.Error()))
	}
}

func (ds *DataSource) BeginTX() *sqlx.Tx {
	return ds.DBPool.MustBegin()
}

func (ds *DataSource) Query(sql string, args []any) *sqlx.Rows {
	var rows *sqlx.Rows
	var err error
	if ds.TX != nil {
		rows, err = ds.TX.Queryx(sql, args...)
	} else {
		rows, err = ds.DBPool.Queryx(sql, args...)
	}
	if err != nil {
		panic(exception.New(err.Error()+" ["+sql+"]", 2))
	}
	return rows
}

func (ds *DataSource) Exec(sql string, args []any) sql.Result {
	if ds.TX != nil {
		return ds.TX.MustExec(sql, args...)
	} else {
		return ds.DBPool.MustExec(sql, args...)
	}
}
