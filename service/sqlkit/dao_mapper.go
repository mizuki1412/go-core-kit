package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
)

func (dao Dao[T]) Insert(dest *T) {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	builder := dao.Builder().Insert()
	var columns []string
	var vals []any
	rv := reflect.ValueOf(dest)
	for _, e := range dao.modelMeta.allInsertKeys {
		var val = e.val(rv)
		if val == nil {
			continue
		}
		columns = append(columns, e.Key)
		vals = append(vals, val)
	}
	if len(columns) == 0 {
		panic(exception.New("no fields", 2))
	}
	builder = builder.Columns(columns...).Values(vals...)
	builder = builder.Suffix("returning *")
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

func (dao Dao[T]) Update(dest any) {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	builder := dao.Builder().Update()
	rv := reflect.ValueOf(dest)
	for _, e := range dao.modelMeta.allUpdateKeys {
		var val = e.val(rv)
		if val == nil {
			continue
		}
		// 针对class.MapString 采用merge方式 todo mysql
		if e.RStruct.Type.String() == "class.MapString" && dao.dataSource.Driver == Postgres {
			builder = builder.Set(e.Key, squirrel.Expr("coalesce("+e.Key+",'{}'::jsonb) || ?", val))
		} else {
			builder = builder.Set(e.Key, val)
		}
	}
	for _, e := range dao.modelMeta.allPKs {
		v := e.val(rv)
		if v == nil {
			panic(exception.New("pk val is nil"))
		}
		builder = builder.Where(e.Key+"=?", v)
	}
	sql, args, err := builder.ToSql()
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	dao.Exec(sql, args...)
}

// Delete dest should be elem
func (dao Dao[T]) Delete(dest any) {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	rv := reflect.ValueOf(dest)
	var sql string
	var args []interface{}
	var err error
	if dao.modelMeta.logicDelKey != "" {
		builder := dao.Builder().Update()
		ldv := LogicDelVal[0]
		if len(dao.LogicDelVal) > 0 {
			ldv = dao.LogicDelVal[0]
		}
		builder.Set(dao.modelMeta.logicDelKey, ldv)
		for _, e := range dao.modelMeta.allPKs {
			v := e.val(rv)
			if v == nil {
				panic(exception.New("pk val is nil"))
			}
			builder = builder.Where(e.Key+"=?", v)
		}
		sql, args, err = builder.ToSql()
	} else {
		builder := dao.Builder().Delete()
		for _, e := range dao.modelMeta.allPKs {
			v := e.val(rv)
			if v == nil {
				panic(exception.New("pk val is nil"))
			}
			builder = builder.Where(e.Key+"=?", v)
		}
		sql, args, err = builder.ToSql()
	}
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	//log.Println(sql, args)
	dao.Exec(sql, args...)
}

func (dao Dao[T]) SelectOneById(id ...any) *T {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	builder := dao.Builder().Select()
	if len(id) != len(dao.modelMeta.allPKs) {
		panic(exception.New("主键数量不匹配"))
	}
	for i := 0; i < len(dao.modelMeta.allPKs); i++ {
		builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
	}
	sql, args := builder.Sql()
	return dao.ScanOne(sql, args)
}
