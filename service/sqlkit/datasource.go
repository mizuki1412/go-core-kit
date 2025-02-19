package sqlkit

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/v2/class/const/sqlconst"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
	"sync"
	"time"
)

type DataSource struct {
	// 用于postgres
	Schema string
	// 数据源（事务时使用）
	TX *sqlx.Tx
	// 指定数据源（原始数据源连接池）
	DBPool *sqlx.DB
	// 连接时的driver
	Driver string
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
	case sqlconst.Postgres:
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
	if param.Driver == sqlconst.Postgres {
		ds.Schema = "public"
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
	if driver == sqlconst.Postgres {
		ds.Schema = "public"
	}
	return ds
}

// DecoTableName 获取 schema 修饰的转义的tableName
func (ds *DataSource) DecoTableName(tableName string) string {
	s := ""
	if ds.Driver == sqlconst.Postgres {
		if ds.Schema != "" {
			s = ds.Schema + "."
		} else {
			s = "public."
		}
	}
	return s + ds.EscapeName(tableName)
}

// EscapeName 表名列名的转义符添加
func (ds *DataSource) EscapeName(name string) string {
	switch ds.Driver {
	case sqlconst.Mysql:
		return "`" + name + "`"
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
		panic(exception.New(err.Error(), 2))
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
