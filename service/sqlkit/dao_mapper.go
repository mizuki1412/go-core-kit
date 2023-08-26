package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"reflect"
)

func (dao Dao[T]) InsertObj(dest *T) {
	builder := dao.Insert()
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
	builder.ReturnOne(dest)
}

func (dao Dao[T]) UpdateObj(dest *T) int64 {
	builder := dao.Update()
	rv := reflect.ValueOf(dest).Elem()
	for _, e := range dao.modelMeta.allUpdateKeys {
		var val = e.val(rv)
		if val == nil {
			continue
		}
		// 针对class.MapString 采用merge方式 todo mysql
		if (e.RStruct.Type.String() == "class.MapString" || e.RStruct.Type.String() == "class.MapStringSync") && dao.dataSource.Driver == Postgres {
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
	return builder.Exec()
}

func (dao Dao[T]) DeleteById(id ...any) int64 {
	if len(id) != len(dao.modelMeta.allPKs) {
		panic(exception.New("主键数量不匹配"))
	}
	if dao.modelMeta.logicDelKey.Key != "" {
		builder := dao.Update()
		builder = builder.Set(dao.modelMeta.logicDelKey.OriKey, builder.LogicDelVal[0])
		for i := 0; i < len(dao.modelMeta.allPKs); i++ {
			builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
		}
		return builder.Exec()
	} else {
		builder := dao.Delete()
		for i := 0; i < len(dao.modelMeta.allPKs); i++ {
			builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
		}
		return builder.Exec()
	}
}

func (dao Dao[T]) SelectOneById(id ...any) *T {
	builder := dao.Select()
	if len(id) != len(dao.modelMeta.allPKs) {
		panic(exception.New("主键数量不匹配"))
	}
	for i := 0; i < len(dao.modelMeta.allPKs); i++ {
		builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
	}
	return builder.One()
}

func (dao Dao[T]) SelectOneWithDelById(id ...any) *T {
	builder := dao.Select()
	if len(id) != len(dao.modelMeta.allPKs) {
		panic(exception.New("主键数量不匹配"))
	}
	for i := 0; i < len(dao.modelMeta.allPKs); i++ {
		builder = builder.Where(dao.modelMeta.allPKs[i].Key+"=?", id[i])
	}
	return builder.IgnoreLogicDel().One()
}
