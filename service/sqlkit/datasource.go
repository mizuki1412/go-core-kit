package sqlkit

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/cli/configkey"
	"github.com/mizuki1412/go-core-kit/service/configkit"
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
	Driver  string
	Host    string
	Port    string
	User    string
	Pwd     string
	Name    string // dbName
	MaxOpen int
	MaxIdle int
	MaxLife int
}

const (
	Postgres  = "postgres"
	Mysql     = "mysql"
	SqlServer = "mssql"
	Oracle    = "oracle"
	KingBase  = "kingbase"
)

func getDataSourceName(p DataSourceParam) (string, string) {
	if p.Driver == "" || p.Host == "" || p.Port == "" {
		panic(exception.New("sqlkit: database config error"))
	}
	var param string
	switch p.Driver {
	case Postgres:
		param = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Pwd, p.Name)
	case Mysql:
		param = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=%s", p.User, p.Pwd, p.Host, p.Port, p.Name, "Asia%2FShanghai")
	case SqlServer:
		param = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s", p.Host, p.User, p.Pwd, p.Port, p.Name)
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
	if param.Driver == Postgres {
		ds.Schema = "public"
	}
	return ds
}

func getDB(param DataSourceParam) *sqlx.DB {
	db := sqlx.MustConnect(getDataSourceName(param))
	if param.MaxLife > 0 {
		db.SetConnMaxLifetime(time.Duration(param.MaxLife) * time.Minute)
	}
	if param.MaxOpen > 0 {
		db.SetMaxOpenConns(param.MaxOpen)
	}
	if param.MaxIdle > 0 {
		db.SetMaxIdleConns(param.MaxIdle)
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
	if driver == Postgres {
		ds.Schema = "public"
	}
	return ds
}

// 获取 schema 修饰的转义的tableName
func (ds *DataSource) decoTableName(tableName string) string {
	s := ""
	if ds.Schema != "" {
		s = ds.Schema + "."
	} else if ds.Driver == Postgres {
		s = "public."
	}
	return s + ds.escapeName(tableName)
}

// 表名列名的转义符添加
func (ds *DataSource) escapeName(name string) string {
	switch ds.Driver {
	case Mysql:
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

func (ds *DataSource) Query(sql string, args ...any) *sqlx.Rows {
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

func (ds *DataSource) Exec(sql string, args ...any) {
	if ds.TX != nil {
		ds.TX.MustExec(sql, args...)
	} else {
		ds.DBPool.MustExec(sql, args...)
	}
}
