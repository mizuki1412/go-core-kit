package sqlkit

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"log"
	"reflect"
	"time"
)

/**
数据库连接和CRUD的通用接口
在项目的main中须要导入driver
*/

var driver string
var db *sqlx.DB

const SchemaDefault = "public"

func driverName() string {
	if db == nil {
		connector()
	}
	return driver
}
func connector() *sqlx.DB {
	if !configkit.Exist(ConfigKeyDBDriver) || !configkit.Exist(ConfigKeyDBHost) || !configkit.Exist(ConfigKeyDBPort) || !configkit.Exist(ConfigKeyDBPwd) || !configkit.Exist(ConfigKeyDBUser) || !configkit.Exist(ConfigKeyDBName) {
		panic("sqlkit: database config error")
	}
	if db == nil {
		driver = configkit.GetString(ConfigKeyDBDriver, "")
		db = sqlx.MustConnect(driver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configkit.GetString(ConfigKeyDBHost, ""), configkit.GetInt(ConfigKeyDBPort, 0), configkit.GetString(ConfigKeyDBUser, ""), configkit.GetString(ConfigKeyDBPwd, ""), configkit.GetString(ConfigKeyDBName, "")))
		db.SetConnMaxLifetime(time.Minute * 60)
		// todo 需要调优
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
	}
	return db
}

type Dao struct {
	ResultType byte
	Schema     string
	TX         *sqlx.Tx
	DB         *sqlx.DB
}

func New(schema string) *Dao {
	return &Dao{
		Schema: schema,
		DB:     connector(),
	}
}

// 在context中调用，如果要单独调用，注意commit和rollback的处理
func NewTX(schema string) *Dao {
	return &Dao{
		Schema: schema,
		TX:     connector().MustBegin(),
		DB:     connector(),
	}
}

// 用于处理子类Dao New中的通用逻辑
func (dao *Dao) NewHelper(schema string, tx ...*Dao) {
	dao.Schema = schema
	if tx != nil && len(tx) > 0 {
		dao.DB = tx[0].DB
		dao.TX = tx[0].TX
	} else {
		dao.DB = New(schema).DB
	}
}

func (dao *Dao) SetResultType(rt byte) *Dao {
	dao.ResultType = rt
	return dao
}

func (dao *Dao) SetSchema(schema string) *Dao {
	dao.Schema = schema
	return dao
}

/// 根据类获取tablename，并判断schema
func (dao *Dao) GetTable(dest interface{}) string {
	rt := reflect.TypeOf(dest).Elem()
	return getTable(rt, dao.Schema)
}
func getTable(rt reflect.Type, schema string) string {
	var tname string
	for i := 0; i < rt.NumField(); i++ {
		if t, ok := rt.Field(i).Tag.Lookup("tablename"); ok {
			tname = t
			break
		}
	}
	if tname == "" {
		panic(exception.New("tablename未设置", 2))
	}
	var schema0 string
	if schema != "" {
		schema0 = schema
	} else if driver == "postgres" {
		schema0 = SchemaDefault
	} else {
		schema0 = ""
	}
	if schema0 == "" {
		return tname
	} else {
		return schema0 + "." + tname
	}
}

/// 直接获取schema+tableName
func (dao *Dao) GetTableD(tName string) string {
	var schema0 string
	if dao.Schema != "" {
		schema0 = dao.Schema
	} else if driver == "postgres" {
		schema0 = SchemaDefault
	} else {
		schema0 = ""
	}
	if schema0 == "" {
		return tName
	} else {
		return schema0 + "." + tName
	}
}

// transaction
func (dao *Dao) Commit() {
	if dao.TX != nil {
		err := dao.TX.Commit()
		if err != nil {
			logkit.Error(err.Error())
			//panic(exception.New(err.Error(), 2))
		}
	}
}
func (dao *Dao) Rollback() {
	if dao.TX != nil {
		err := dao.TX.Rollback()
		if err != nil {
			logkit.Error(err.Error())
		}
	}
}

func (dao *Dao) Query(sql string, args ...interface{}) *sqlx.Rows {
	var rows *sqlx.Rows
	var err error
	if dao.TX != nil {
		rows, err = dao.TX.Queryx(sql, args...)
	} else {
		rows, err = dao.DB.Queryx(sql, args...)
	}
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	return rows
}
func (dao *Dao) Exec(sql string, args ...interface{}) {
	if dao.TX != nil {
		dao.TX.MustExec(sql, args...)
	} else {
		dao.DB.MustExec(sql, args...)
	}
}

// dest a struct
// todo select会引起no-struct错误（Scan()导致）；structScan 对interface{}报错
func (dao *Dao) QueryStruct(destType func(rs *sqlx.Rows) (interface{}, error), sql string, args []interface{}, err error) []interface{} {
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println("sqlkit:", sql, args)
	//err = DB().Select(dest, sql, args...)
	var list []interface{}
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m, err := destType(rows)
		if err != nil {
			panic(exception.New(err.Error(), 2))
		}
		list = append(list, m)
	}
	return list
}

func (dao *Dao) QueryMap(sql string, args []interface{}, err error) []map[string]interface{} {
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	var list []map[string]interface{}
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := map[string]interface{}{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error(), 2))
		}
		list = append(list, m)
	}
	return list
}

// dest a pointer
// todo 此方法不能表示nil
//func (dao *Dao) QueryById(dest interface{}, selects ...string) {
//	rt := reflect.TypeOf(dest).Elem()
//	rv := reflect.ValueOf(dest).Elem()
//	pks := getPKs(rt, rv)
//	builder0 := Builder()
//	var builder squirrel.SelectBuilder
//	if len(selects) > 0 {
//		builder = builder0.Select(selects...)
//	} else {
//		builder = builder0.Select("*")
//	}
//	builder = builder.From(getTable(rt, dao.Schema))
//	for k, v := range pks {
//		builder = builder.Where(k+"=?", v)
//	}
//	sql, args, err := builder.ToSql()
//	if err != nil {
//		panic(exception.New(err.Error(), 2))
//	}
//	rows := dao.Query(sql, args...)
//  defer rows.Close()
//	for rows.Next() {
//		err := rows.StructScan(dest)
//		if err != nil {
//			panic(exception.New(err.Error(), 2))
//		}
//		break
//	}
//}

/// dest should be elem
func (dao *Dao) Insert(dest interface{}) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Insert(getTable(rt, dao.Schema))
	var pks []string
	var columns []string
	var vals []interface{}
	for i := 0; i < rt.NumField(); i++ {
		// 自增
		if t, ok := rt.Field(i).Tag.Lookup("autoincrement"); ok && t == "true" {
			name := rt.Field(i).Tag.Get("db")
			if name == "" {
				panic(exception.New("field "+rt.Field(i).Name+" no db tag", 2))
			}
			pks = append(pks, name)
			continue
		}
		db, ok := rt.Field(i).Tag.Lookup("db")
		var val interface{}
		// 判断field是否指针
		if rt.Field(i).Type.Kind() == reflect.Ptr && rv.Field(i).Elem().IsValid() {
			val = rv.Field(i).Elem().Interface()
		} else if rt.Field(i).Type.Kind() != reflect.Ptr {
			val = rv.Field(i).Interface()
		}
		// eg: MapString, 根据类的Value返回判断此field是否可以insert
		method := rv.Field(i).MethodByName("Value")
		if ok && val != nil && (!method.IsValid() || method.Call([]reflect.Value{})[0].Interface() != nil) {
			columns = append(columns, db)
			vals = append(vals, val)
		}
	}
	if len(columns) == 0 {
		panic(exception.New("no fields", 2))
	}
	builder = builder.Columns(columns...).Values(vals...)
	// 暂支持一个return
	if len(pks) > 0 {
		builder = builder.Suffix("returning " + pks[0])
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	log.Println(sql, args)
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		// return 赋值
		err := rows.StructScan(dest)
		if err != nil {
			panic(exception.New(err.Error(), 2))
		}
		break
	}
}

/// dest should be elem
func (dao *Dao) Update(dest interface{}) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Update(getTable(rt, dao.Schema))
	for i := 0; i < rt.NumField(); i++ {
		db, ok := rt.Field(i).Tag.Lookup("db")
		pk := rt.Field(i).Tag.Get("pk")
		var val interface{}
		// 判断field是否指针
		if rt.Field(i).Type.Kind() == reflect.Ptr && rv.Field(i).Elem().IsValid() {
			val = rv.Field(i).Elem().Interface()
		} else if rt.Field(i).Type.Kind() != reflect.Ptr {
			val = rv.Field(i).Interface()
		}
		method := rv.Field(i).MethodByName("Value")
		if ok && pk != "true" && val != nil && method.IsValid() && method.Call([]reflect.Value{})[0].Interface() != nil {
			builder = builder.Set(db, val)
		}
	}
	pks := getPKs(rt, rv)
	for k, v := range pks {
		builder = builder.Where(k+"=?", v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println(sql, args)
	dao.Exec(sql, args...)
}

/// dest should be elem
func (dao *Dao) Delete(dest interface{}) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Delete(getTable(rt, dao.Schema))
	pks := getPKs(rt, rv)
	for k, v := range pks {
		builder = builder.Where(k+"=?", v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println(sql, args)
	dao.Exec(sql, args...)
}

func (dao *Dao) DeleteOff(dest interface{}) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Update(getTable(rt, dao.Schema)).Set("off", true)
	pks := getPKs(rt, rv)
	for k, v := range pks {
		builder = builder.Where(k+"=?", v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println(sql, args)
	dao.Exec(sql, args...)
}

func getPKs(rt reflect.Type, rv reflect.Value) map[string]interface{} {
	pks := map[string]interface{}{}
	for i := 0; i < rt.NumField(); i++ {
		if t, ok := rt.Field(i).Tag.Lookup("pk"); ok {
			if t == "true" {
				name := rt.Field(i).Tag.Get("db")
				if name == "" {
					panic(exception.New("field "+rt.Field(i).Name+" no db tag", 2))
				}
				pks[name] = rv.Field(i).Interface()
			}
		}
	}
	if len(pks) == 0 {
		panic(exception.New("未设置pk", 2))
	}
	return pks
}
