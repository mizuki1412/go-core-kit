package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
)

// LogicDelVal 全局逻辑删除的 value
var LogicDelVal = []any{true, false}

type Dao[T any] struct {
	meta T
	// 逻辑删除的字段，可替代全局的LogicDelVal
	LogicDelVal []any
	// 返回级联的类型
	ResultType byte
	// 级联实现的函数
	Cascade func(*T)
	// 数据源
	dataSource *DataSource
	// 目标表结构
	modelMeta ModelMeta
}

type DaoModelMeta interface {
	getModelMeta() ModelMeta
}

// New 必须从初始化函数生成 dao
func New[T any](ds ...*DataSource) Dao[T] {
	dao := Dao[T]{}
	if len(ds) > 0 {
		dao.dataSource = ds[0]
	} else {
		dao.dataSource = DefaultDataSource()
	}
	dao.modelMeta.dateSource = dao.dataSource
	dao.modelMeta = dao.modelMeta.init(dao.meta)
	return dao
}

//func (dao *Dao[T]) Init(ds ...*DataSource) {
//	if len(ds) > 0 {
//		dao.dataSource = ds[0]
//	} else {
//		dao.dataSource = DefaultDataSource()
//	}
//	dao.modelMeta.dateSource = dao.dataSource
//	dao.modelMeta.init(dao.meta)
//}

// Builder 结构化语句
func (dao Dao[T]) Builder() SQLBuilder {
	sb := SQLBuilder{modelMeta: dao.modelMeta}
	ldv := LogicDelVal
	if len(dao.LogicDelVal) > 0 {
		ldv = dao.LogicDelVal
	}
	sb.logicDel = ldv
	sb.driver = dao.dataSource.Driver
	switch dao.dataSource.Driver {
	case Postgres:
		sb.inner0 = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	default:
		// todo 未处理oracle
		sb.inner0 = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
	return sb
}

func (dao Dao[T]) getModelMeta() ModelMeta {
	return dao.modelMeta
}

func (dao Dao[T]) DataSource() *DataSource {
	return dao.dataSource
}

func (dao Dao[T]) Query(sql string, args []any) *sqlx.Rows {
	logkit.DebugConcat(sql, " | args:", jsonkit.ToString(args))
	return dao.dataSource.Query(sql, args)
}

func (dao Dao[T]) Exec(sql string, args []any) {
	logkit.DebugConcat(sql, " | args:", jsonkit.ToString(args))
	dao.dataSource.Exec(sql, args)
}

// ScanList 取值封装list
func (dao Dao[T]) ScanList(sql string, args []any) []*T {
	rows := dao.Query(sql, args)
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

func (dao Dao[T]) ScanOne(sql string, args []any) *T {
	rows := dao.Query(sql, args)
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

func (dao Dao[T]) ScanOneMap(sql string, args []any) map[string]any {
	rows := dao.Query(sql, args)
	defer rows.Close()
	for rows.Next() {
		m := map[string]any{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		return m
	}
	return nil
}

func (dao Dao[T]) ScanListMap(sql string, args []any) []map[string]any {
	rows := dao.Query(sql, args)
	list := make([]map[string]any, 0, 5)
	defer rows.Close()
	for rows.Next() {
		m := map[string]any{}
		err := rows.MapScan(m)
		if err != nil {
			panic(exception.New(err.Error()))
		}
		list = append(list, m)
	}
	return list
}

func (dao Dao[T]) SelectColumns(excludes ...string) []string {
	return dao.modelMeta.getSelectColumns(excludes...)
}
func (dao Dao[T]) SelectColumnsWithP(prefix string, excludes ...string) []string {
	return dao.modelMeta.getSelectColumnsWithPrefix(prefix, excludes...)
}
func (dao Dao[T]) Table(alias ...string) string {
	return dao.modelMeta.getTable(alias...)
}

func (dao Dao[T]) EscapeNames(name ...string) []string {
	return dao.modelMeta.escapeNames(name)
}
