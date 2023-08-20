package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mizuki1412/go-core-kit/class/exception"
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

// Builder 结构化语句
func (dao Dao[T]) Builder() SQLBuilder {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	sb := SQLBuilder{modelMeta: dao.modelMeta}
	switch dao.dataSource.Driver {
	case Postgres:
		sb.inner = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	default:
		// todo 未处理oracle
		sb.inner = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
	return sb
}

func (dao Dao[T]) getModelMeta() ModelMeta {
	return dao.modelMeta
}

// SetDataSource 设置数据源，同时 init modelMeta
func (dao Dao[T]) SetDataSource(ds *DataSource) Dao[T] {
	dao.dataSource = ds
	dao.modelMeta.dateSource = ds
	dao.modelMeta.init(dao.meta)
	return dao
}
func (dao Dao[T]) DataSource() *DataSource {
	return dao.dataSource
}

func (dao Dao[T]) SetResultType(rt byte) Dao[T] {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	dao.ResultType = rt
	return dao
}

func (dao Dao[T]) Query(sql string, args ...any) *sqlx.Rows {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	return dao.dataSource.Query(sql, args...)
}

func (dao Dao[T]) Exec(sql string, args ...any) {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	dao.dataSource.Exec(sql, args...)
}

// ScanList 取值封装list
func (dao Dao[T]) ScanList(sql string, args ...any) []*T {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
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

func (dao Dao[T]) ScanOne(sql string, args ...any) *T {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
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

func (dao Dao[T]) ScanOneMap(sql string, args []any) map[string]any {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	rows := dao.Query(sql, args...)
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

func (dao Dao[T]) ScanListMap(sql string, args ...any) []map[string]any {
	if dao.dataSource == nil {
		dao.dataSource = DefaultDataSource()
	}
	rows := dao.Query(sql, args...)
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
