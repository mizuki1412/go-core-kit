package sqlkit

import (
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

// NewDataSource 创建一个数据源
func NewDataSource(param DataSourceParam) *DataSource {
	db := getDB(param)
	ds := &DataSource{
		Driver: param.Driver,
		DBPool: db,
	}
	if param.Driver == "postgres" {
		ds.Schema = "public"
	}
	return ds
}

func getDB(param DataSourceParam) *sqlx.DB {
	db := sqlx.MustConnect(GetDataSourceName(param))
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
	if driver == "postgres" {
		ds.Schema = "public"
	}
	return ds
}

// 获取 schema 修饰tableName
func (ds *DataSource) getDecoSchema() string {
	if ds.Schema != "" {
		return ds.Schema + "."
	} else if ds.Driver == "postgres" {
		return "public."
	}
	return ""
}

func (ds *DataSource) setSchema(schema string) *DataSource {
	ds.Schema = schema
	return ds
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
