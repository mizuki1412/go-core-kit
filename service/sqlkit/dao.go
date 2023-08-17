package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
)

type Dao[T any] struct {
	// 返回级联的类型
	ResultType byte
	// 级联实现的函数
	Cascade func(*T)
	// 数据源
	DataSource *DataSource
}

// Builder 结构化语句
func (dao *Dao[T]) Builder() squirrel.StatementBuilderType {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
	return getStatementBuilderType(dao.DataSource)
}

func (dao *Dao[T]) SetResultType(rt byte) *Dao[T] {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
	dao.ResultType = rt
	return dao
}

func (dao *Dao[T]) SetSchema(schema string) *Dao[T] {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
	dao.DataSource.Schema = schema
	return dao
}

// GetTable 根据类获取tablename，并判断schema
func (dao *Dao[T]) GetTable(dest any) string {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
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
	return GetSchemaTable(schema, tname)
}

func (dao *Dao[T]) Query(sql string, args ...any) *sqlx.Rows {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
	var rows *sqlx.Rows
	var err error
	if dao.TX != nil {
		rows, err = dao.TX.Queryx(sql, args...)
	} else {
		rows, err = dao.Connector().Queryx(sql, args...)
	}
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	return rows
}
func (dao *Dao[T]) Exec(sql string, args ...any) {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
	if dao.TX != nil {
		dao.TX.MustExec(sql, args...)
	} else {
		dao.Connector().MustExec(sql, args...)
	}
}

// dest a struct
// todo select会引起no-struct错误（Scan()导致）；structScan 对any报错
func (dao *Dao[T]) QueryStruct(destType func(rs *sqlx.Rows) (any, error), sql string, args []any, err error) []any {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
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

func (dao *Dao[T]) QueryMap(sql string, args []any, err error) []map[string]any {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
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
//func (dao *Dao[T]) QueryById(dest any, selects ...string) {
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

// Insert dest should be elem
func (dao *Dao[T]) Insert(dest any) {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
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

// Update dest should be elem
func (dao *Dao[T]) Update(dest any) {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
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

// Delete dest should be elem
func (dao *Dao[T]) Delete(dest any) {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
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

func (dao *Dao[T]) DeleteOff(dest any) {
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

// ScanList 取值封装list
func (dao *Dao[T]) ScanList(sql string, args []any) []*T {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
	rows := dao.Query(sql, args...)
	list := make([]*T, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := new(T)
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	if dao.Cascade != nil {
		for i := range list {
			dao.Cascade(list[i])
		}
	}
	return list
}

func (dao *Dao[T]) ScanOne(sql string, args []any) *T {
	if dao.DataSource == nil {
		dao.DataSource = DefaultDataSource()
	}
	rows := dao.Query(sql, args...)
	defer rows.Close()
	for rows.Next() {
		m := new(T)
		err := rows.StructScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		if dao.Cascade != nil {
			dao.Cascade(m)
		}
		return m
	}
	return nil
}
