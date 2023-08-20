package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

type SQLBuilder struct {
	inner0    squirrel.StatementBuilderType
	modelMeta ModelMeta
	// 逻辑删除标记
	logicDel []any
	driver   string
}

// Select 默认取modelmeta中的columns，并装饰引号；fields中不装饰，因为可能存在表达式
func (b SQLBuilder) _select(fields ...string) SelectBuilder {
	return SelectBuilder{
		inner:     b.inner0.Select(fields...),
		modelMeta: b.modelMeta,
		logicDel:  b.logicDel,
		driver:    b.driver,
	}
}

func (b SQLBuilder) Select(fields ...string) SelectBuilder {
	if len(fields) == 0 {
		return b.SelectWithout()
	} else {
		return b._select(fields...)
	}
}

// SelectWithout 在modelmeta columns中去掉指定的字段
func (b SQLBuilder) SelectWithout(fields ...string) SelectBuilder {
	return b.SelectPrefix("", fields...)
}

// SelectPrefix 在modelmeta的字段前增加prefix
func (b SQLBuilder) SelectPrefix(prefix string, without ...string) SelectBuilder {
	if b.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return b._select(b.modelMeta.getSelectColumnsWithPrefix(prefix, without...)...)
}

func (b SQLBuilder) Update() UpdateBuilder {
	if b.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return UpdateBuilder{
		inner:     b.inner0.Update(b.modelMeta.getTable()),
		modelMeta: b.modelMeta,
		logicDel:  b.logicDel,
	}
}

func (b SQLBuilder) Delete() DeleteBuilder {
	if b.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return DeleteBuilder{
		inner:     b.inner0.Delete(b.modelMeta.getTable()),
		modelMeta: b.modelMeta,
	}
}

func (b SQLBuilder) Insert() InsertBuilder {
	if b.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return InsertBuilder{
		inner:     b.inner0.Insert(b.modelMeta.getTable()),
		modelMeta: b.modelMeta,
	}
}

func (b SQLBuilder) Replace() InsertBuilder {
	if b.modelMeta.tableName == "" {
		panic(exception.New("sqlbuilder modelmeta null"))
	}
	return InsertBuilder{
		inner:     b.inner0.Replace(b.modelMeta.getTable()),
		modelMeta: b.modelMeta,
	}
}
