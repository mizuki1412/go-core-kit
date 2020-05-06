package sqlkit

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"mizuki/project/core-kit/class/exception"
	"mizuki/project/core-kit/service/configkit"
	"reflect"
)

/**
数据库连接和CRUD的通用接口
rows转struct和map
在项目的main中须要导入driver
//todo middleware中加入事务
*/

var driver string
var db *sqlx.DB

const SchemaDefault = "public"

func Driver() string {
	if db == nil {
		DB()
	}
	return driver
}
func DB() *sqlx.DB {
	if !configkit.Exist(ConfigKeyDBDriver) || !configkit.Exist(ConfigKeyDBHost) || !configkit.Exist(ConfigKeyDBPort) || !configkit.Exist(ConfigKeyDBPwd) || !configkit.Exist(ConfigKeyDBUser) || !configkit.Exist(ConfigKeyDBName) {
		panic("sqlkit: database config error")
	}
	if db == nil {
		driver = configkit.GetString(ConfigKeyDBDriver, "")
		db = sqlx.MustConnect(driver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configkit.GetString(ConfigKeyDBHost, ""), configkit.GetInt(ConfigKeyDBPort, 0), configkit.GetString(ConfigKeyDBUser, ""), configkit.GetString(ConfigKeyDBPwd, ""), configkit.GetString(ConfigKeyDBName, "")))
		//db.SetConnMaxLifetime()
	}
	return db
}

func GetTable(dest interface{}, schema ...string) string {
	rt := reflect.TypeOf(dest).Elem()
	return getTable(rt, schema...)
}

func getTable(rt reflect.Type, schema ...string) string {
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
	schema0 := SchemaDefault
	if schema != nil && len(schema) > 0 {
		schema0 = schema[0]
	}
	return schema0 + "." + tname
}

// dest a struct
// todo select会引起no-struct错误（Scan()导致）；structScan 对interface{}报错
func QueryStruct(destType func(rs *sqlx.Rows) (interface{}, error), sql string, args []interface{}, err error) []interface{} {
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	log.Println("sqlkit:", sql, args)
	//err = DB().Select(dest, sql, args...)
	var list []interface{}
	rows, err := DB().Queryx(sql, args...)
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	for rows.Next() {
		m, err := destType(rows)
		if err != nil {
			panic(exception.New(err.Error(), 2))
		}
		list = append(list, m)
	}
	return list
}

func QueryMap(sql string, args []interface{}, err error) []map[string]interface{} {
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	log.Println("sqlkit: ", sql)
	var list []map[string]interface{}
	rows, _ := DB().Queryx(sql, args...)
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
func QueryById(dest interface{}, schema ...string) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	pks := getPKs(rt, rv)
	builder := Builder().Select("*").From(getTable(rt, schema...))
	for k, v := range pks {
		builder = builder.Where(k+"=?", v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	rows, _ := DB().Queryx(sql, args...)
	for rows.Next() {
		err := rows.StructScan(dest)
		if err != nil {
			panic(exception.New(err.Error(), 2))
		}
		break
	}
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

func Insert(dest interface{}, schema ...string) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Insert(getTable(rt, schema...))
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
		method := rv.Field(i).MethodByName("Value")
		if ok && val != nil && (!method.IsValid() || (method.IsValid() && method.Call([]reflect.Value{})[0].Interface() != nil)) {
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
	rows, err := DB().Queryx(sql, args...)
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	for rows.Next() {
		// return 赋值
		err := rows.StructScan(dest)
		if err != nil {
			panic(exception.New(err.Error(), 2))
		}
		break
	}
}

func Update(dest interface{}, schema ...string) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Update(getTable(rt, schema...))
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
	DB().MustExec(sql, args...)
}

func Delete(dest interface{}, schema ...string) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Delete(getTable(rt, schema...))
	pks := getPKs(rt, rv)
	for k, v := range pks {
		builder = builder.Where(k+"=?", v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println(sql, args)
	DB().MustExec(sql, args...)
}

func DeleteOff(dest interface{}, schema ...string) {
	rt := reflect.TypeOf(dest).Elem()
	rv := reflect.ValueOf(dest).Elem()
	builder := Builder().Update(getTable(rt, schema...)).Set("off", true)
	pks := getPKs(rt, rv)
	for k, v := range pks {
		builder = builder.Where(k+"=?", v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println(sql, args)
	DB().MustExec(sql, args...)
}

/**
struct tag 包括 db:name, pk:true, tablename:name(只在一个tag), autoincrement:"true"
*/
