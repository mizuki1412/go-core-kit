package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

// 默认取modelmeta中的columns，并装饰引号；fields中不装饰，因为可能存在表达式
func (dao Dao[T]) _select(fields ...string) SelectDao[T] {
	ldv := LogicDelVal
	if len(dao.LogicDelVal) > 0 {
		ldv = dao.LogicDelVal
	}
	d := SelectDao[T]{
		builder: squirrel.Select(fields...),
	}
	d.meta = dao.meta
	d.dataSource = dao.dataSource
	d.modelMeta = dao.modelMeta
	d.LogicDelVal = ldv
	d.ResultType = dao.ResultType
	d.Cascade = dao.Cascade
	return d
}

func (dao Dao[T]) Select(fields ...string) SelectDao[T] {
	if len(fields) == 0 {
		return dao.SelectEx()
	} else {
		return dao._select(fields...)
	}
}

// SelectEx 在modelmeta columns中去掉指定的字段
func (dao Dao[T]) SelectEx(fields ...string) SelectDao[T] {
	return dao.SelectPrefix("", fields...)
}

// SelectPrefix 在modelmeta的字段前增加prefix
func (dao Dao[T]) SelectPrefix(prefix string, without ...string) SelectDao[T] {
	if dao.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return dao._select(dao.modelMeta.getSelectColumnsWithPrefix(prefix, without...)...)
}

func (dao Dao[T]) Update() UpdateDao[T] {
	if dao.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	ldv := LogicDelVal
	if len(dao.LogicDelVal) > 0 {
		ldv = dao.LogicDelVal
	}
	d := UpdateDao[T]{
		builder: squirrel.Update(dao.modelMeta.getTable()),
	}
	d.meta = dao.meta
	d.dataSource = dao.dataSource
	d.modelMeta = dao.modelMeta
	d.LogicDelVal = ldv
	d.ResultType = dao.ResultType
	d.Cascade = dao.Cascade
	return d
}

func (dao Dao[T]) Delete() DeleteDao[T] {
	if dao.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	d := DeleteDao[T]{
		builder: squirrel.Delete(dao.modelMeta.getTable()),
	}
	ldv := LogicDelVal
	if len(dao.LogicDelVal) > 0 {
		ldv = dao.LogicDelVal
	}
	d.meta = dao.meta
	d.dataSource = dao.dataSource
	d.modelMeta = dao.modelMeta
	d.LogicDelVal = ldv
	d.ResultType = dao.ResultType
	d.Cascade = dao.Cascade
	return d
}

func (dao Dao[T]) Insert() InsertDao[T] {
	if dao.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	d := InsertDao[T]{
		builder: squirrel.Insert(dao.modelMeta.getTable()),
	}
	ldv := LogicDelVal
	if len(dao.LogicDelVal) > 0 {
		ldv = dao.LogicDelVal
	}
	d.meta = dao.meta
	d.dataSource = dao.dataSource
	d.modelMeta = dao.modelMeta
	d.LogicDelVal = ldv
	d.ResultType = dao.ResultType
	d.Cascade = dao.Cascade
	return d
}

func (dao Dao[T]) Replace() InsertDao[T] {
	if dao.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	d := InsertDao[T]{
		builder: squirrel.Replace(dao.modelMeta.getTable()),
	}
	ldv := LogicDelVal
	if len(dao.LogicDelVal) > 0 {
		ldv = dao.LogicDelVal
	}
	d.meta = dao.meta
	d.dataSource = dao.dataSource
	d.modelMeta = dao.modelMeta
	d.LogicDelVal = ldv
	d.ResultType = dao.ResultType
	d.Cascade = dao.Cascade
	return d
}
