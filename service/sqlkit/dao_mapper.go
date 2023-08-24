package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
)

func (dao Dao[T]) Insert(dest *T) {
	builder := dao.Builder().Insert()
	var columns []string
	var vals []any
	rv := reflect.ValueOf(dest).Elem()
	for _, e := range dao.modelMeta.allInsertKeys {
		var val = e.val(rv)
		if val == nil {
			continue
		}
		columns = append(columns, e.OriKey)
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
	rows := dao.Query(sql, args)
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

func (dao Dao[T]) Update(dest *T) {
	builder := dao.Builder().Update()
	rv := reflect.ValueOf(dest).Elem()
	for _, e := range dao.modelMeta.allUpdateKeys {
		var val = e.val(rv)
		if val == nil {
			continue
		}
		// 针对class.MapString 采用merge方式 todo mysql
		if e.RStruct.Type.String() == "class.MapString" && dao.dataSource.Driver == Postgres {
			builder = builder.Set(e.OriKey, squirrel.Expr("coalesce("+e.OriKey+",'{}'::jsonb) || ?", val))
		} else {
			builder = builder.Set(e.OriKey, val)
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
	dao.Exec(sql, args)
}

func (dao Dao[T]) DeleteById(id ...any) {
	var sql string
	var args []interface{}
	var err error
	if len(id) != len(dao.modelMeta.allPKs) {
		panic(exception.New("主键数量不匹配"))
	}
	if dao.modelMeta.logicDelKey.Key != "" {
		builder := dao.Builder().Update()
		builder = builder.Set(dao.modelMeta.logicDelKey.OriKey, builder.logicDel[0])
		for i := 0; i < len(dao.modelMeta.allPKs); i++ {
			builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
		}
		sql, args, err = builder.ToSql()
	} else {
		builder := dao.Builder().Delete()
		for i := 0; i < len(dao.modelMeta.allPKs); i++ {
			builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
		}
		sql, args, err = builder.ToSql()
	}
	if err != nil {
		panic(exception.New(err.Error(), 2))
	}
	dao.Exec(sql, args)
}

func (dao Dao[T]) SelectOneById(id ...any) *T {
	builder := dao.Builder().Select()
	if len(id) != len(dao.modelMeta.allPKs) {
		panic(exception.New("主键数量不匹配"))
	}
	for i := 0; i < len(dao.modelMeta.allPKs); i++ {
		builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
	}
	builder = builder.WhereNLogicDel()
	sql, args := builder.Sql()
	return dao.ScanOne(sql, args)
}

func (dao Dao[T]) SelectOneWithDelById(id ...any) *T {
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
