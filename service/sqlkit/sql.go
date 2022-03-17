package sqlkit

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
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
	if !configkit.Exist(configkey.DBDriver) || !configkit.Exist(configkey.DBHost) || !configkit.Exist(configkey.DBPort) || !configkit.Exist(configkey.DBPwd) || !configkit.Exist(configkey.DBUser) || !configkit.Exist(configkey.DBName) {
		panic("sqlkit: database config error")
	}
	if db == nil {
		driver = configkit.GetString(configkey.DBDriver, "")
		var param string
		if driver == "postgres" {
			param = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configkit.GetString(configkey.DBHost, ""), configkit.GetInt(configkey.DBPort, 0), configkit.GetString(configkey.DBUser, ""), configkit.GetString(configkey.DBPwd, ""), configkit.GetString(configkey.DBName, ""))
		} else if driver == "mysql" {
			param = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", configkit.GetString(configkey.DBUser, ""), configkit.GetString(configkey.DBPwd, ""), configkit.GetString(configkey.DBHost, ""), configkit.GetString(configkey.DBPort, ""), configkit.GetString(configkey.DBName, ""))
		}
		db = sqlx.MustConnect(driver, param)
		lt := cast.ToInt(configkit.GetInt(configkey.DBMaxLife, 20))
		db.SetConnMaxLifetime(time.Duration(lt) * time.Minute)
		db.SetMaxOpenConns(configkit.GetInt(configkey.DBMaxOpen, 25))
		db.SetMaxIdleConns(configkit.GetInt(configkey.DBMaxIdle, 5))
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

// NewTX 在context中调用，如果要单独调用，注意commit和rollback的处理
func NewTX(schema string) *Dao {
	return &Dao{
		Schema: schema,
		TX:     connector().MustBegin(),
		DB:     connector(),
	}
}

// NewHelper 用于处理子类Dao New中的通用逻辑
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

// GetTable 根据类获取tablename，并判断schema
func (dao *Dao) GetTable(dest any) string {
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

// GetTableD 直接获取schema+tableName
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

// Commit transaction
func (dao *Dao) Commit() {
	if dao.TX != nil {
		err := dao.TX.Commit()
		if err != nil {
			logkit.Error(exception.New(err.Error()))
			//panic(exception.New(err.Error(), 2))
		}
	}
}
func (dao *Dao) Rollback() {
	if dao.TX != nil {
		err := dao.TX.Rollback()
		if err != nil {
			logkit.Error(exception.New(err.Error()))
		}
	}
}

func (dao *Dao) Query(sql string, args ...any) *sqlx.Rows {
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
func (dao *Dao) Exec(sql string, args ...any) {
	if dao.TX != nil {
		dao.TX.MustExec(sql, args...)
	} else {
		dao.DB.MustExec(sql, args...)
	}
}

// dest a struct
// todo select会引起no-struct错误（Scan()导致）；structScan 对any报错
func (dao *Dao) QueryStruct(destType func(rs *sqlx.Rows) (any, error), sql string, args []any, err error) []any {
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println("sqlkit:", sql, args)
	//err = DB().Select(dest, sql, args...)
	var list []any
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

func (dao *Dao) QueryMap(sql string, args []any, err error) []map[string]any {
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	var list []map[string]any
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := map[string]any{}
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
//func (dao *Dao) QueryById(dest any, selects ...string) {
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
func (dao *Dao) Insert(dest any) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Insert(getTable(rt, dao.Schema))
	var pks []string
	var columns []string
	var vals []any
	for i := 0; i < rt.NumField(); i++ {
		// 自增的排除
		if t, ok := rt.Field(i).Tag.Lookup("autoincrement"); ok && t == "true" {
			name := rt.Field(i).Tag.Get("db")
			if name == "" {
				panic(exception.New("field "+rt.Field(i).Name+" no db tag", 2))
			}
			pks = append(pks, name)
			continue
		}
		db, ok := rt.Field(i).Tag.Lookup("db")
		var val any
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
	//log.Println(sql, args)
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
func (dao *Dao) Update(dest any) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Update(getTable(rt, dao.Schema))
	for i := 0; i < rt.NumField(); i++ {
		// db字段
		dbKey, ok := rt.Field(i).Tag.Lookup("db")
		pk := rt.Field(i).Tag.Get("pk")
		var val any
		// 判断field是否指针
		if rt.Field(i).Type.Kind() == reflect.Ptr && rv.Field(i).Elem().IsValid() {
			val = rv.Field(i).Elem().Interface()
		} else if rt.Field(i).Type.Kind() != reflect.Ptr {
			val = rv.Field(i).Interface()
		}
		method := rv.Field(i).MethodByName("IsValid")
		// method.Call([]reflect.Value{})[0].Interface() != nil
		// 非主键 或 IsValid函数不存在 或 IsValid==true
		if ok && pk != "true" && val != nil && (!method.IsValid() || (method.IsValid() && method.Call([]reflect.Value{})[0].Bool())) {
			// 针对class.MapString 采用merge方式
			if rt.Field(i).Type.String() == "class.MapString" {
				builder = builder.Set(dbKey, squirrel.Expr("coalesce("+dbKey+",'{}'::jsonb) || ?", val))
			} else {
				builder = builder.Set(dbKey, val)
			}
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
	dao.Exec(sql, args...)
}

/// dest should be elem
func (dao *Dao) Delete(dest any) {
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

func (dao *Dao) DeleteOff(dest any) {
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

// 多或单主键
func getPKs(rt reflect.Type, rv reflect.Value) map[string]any {
	pks := map[string]any{}
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
