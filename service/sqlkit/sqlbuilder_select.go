package sqlkit

import (
	"github.com/Masterminds/squirrel"
)

type SelectBuilder struct {
	inner     squirrel.SelectBuilder
	modelMeta ModelMeta
	fromAs    string
}

func (b SelectBuilder) Sql() (string, []interface{}) {
	if b.fromAs == "" {
		b.inner = b.inner.From(b.modelMeta.GetTableName())
	} else {
		b.inner = b.inner.From(b.modelMeta.GetTableName(b.fromAs))
	}
	return b.inner.MustSql()
}
func (b SelectBuilder) ToSql() (string, []interface{}, error) {
	if b.fromAs == "" {
		b.inner = b.inner.From(b.modelMeta.GetTableName())
	} else {
		b.inner = b.inner.From(b.modelMeta.GetTableName(b.fromAs))
	}
	return b.inner.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (b SelectBuilder) Prefix(sql string, args ...interface{}) SelectBuilder {
	b.inner = b.inner.Prefix(sql, args...)
	return b
}
func (b SelectBuilder) PrefixExpr(expr SelectBuilder) SelectBuilder {
	b.inner = b.inner.PrefixExpr(expr)
	return b
}
func (b SelectBuilder) Suffix(sql string, args ...interface{}) SelectBuilder {
	b.inner = b.inner.Suffix(sql, args...)
	return b
}
func (b SelectBuilder) SuffixExpr(expr SelectBuilder) SelectBuilder {
	b.inner = b.inner.SuffixExpr(expr)
	return b
}

// Columns select 中额外增加 column
func (b SelectBuilder) Columns(cs ...string) SelectBuilder {
	b.inner = b.inner.Columns(cs...)
	return b
}

func (b SelectBuilder) FromAs(alias string) SelectBuilder {
	b.fromAs = alias
	return b
}
func (b SelectBuilder) FromSubQuery(sub SelectBuilder, alias string) SelectBuilder {
	b.inner = b.inner.FromSelect(sub.inner, alias)
	return b
}

func (b SelectBuilder) Distinct() SelectBuilder {
	b.inner = b.inner.Distinct()
	return b
}

// Options adds select option to the query
func (b SelectBuilder) Options(options ...string) SelectBuilder {
	b.inner = b.inner.Options(options...)
	return b
}

func (b SelectBuilder) Join(dm DaoModelMeta, as string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.Join(dm.modelMeta().GetTableName(as), rest...)
	return b
}
func (b SelectBuilder) LeftJoin(dm DaoModelMeta, as string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.LeftJoin(dm.modelMeta().GetTableName(as), rest...)
	return b
}
func (b SelectBuilder) RightJoin(dm DaoModelMeta, as string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.RightJoin(dm.modelMeta().GetTableName(as), rest...)
	return b
}
func (b SelectBuilder) InnerJoin(dm DaoModelMeta, as string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.InnerJoin(dm.modelMeta().GetTableName(as), rest...)
	return b
}
func (b SelectBuilder) CrossJoin(dm DaoModelMeta, as string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.CrossJoin(dm.modelMeta().GetTableName(as), rest...)
	return b
}

func (b SelectBuilder) JoinRaw(join string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.Join(join, rest...)
	return b
}
func (b SelectBuilder) LeftJoinRaw(join string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.LeftJoin(join, rest...)
	return b
}
func (b SelectBuilder) RightJoinRaw(join string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.RightJoin(join, rest...)
	return b
}
func (b SelectBuilder) InnerJoinRaw(join string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.InnerJoin(join, rest...)
	return b
}
func (b SelectBuilder) CrossJoinRaw(join string, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.CrossJoin(join, rest...)
	return b
}

func (b SelectBuilder) Where(pred interface{}, args ...interface{}) SelectBuilder {
	b.inner = b.inner.Where(pred, args...)
	return b
}
func (b SelectBuilder) Having(pred interface{}, rest ...interface{}) SelectBuilder {
	b.inner = b.inner.Where(pred, rest...)
	return b
}

func (b SelectBuilder) GroupBy(groupBys ...string) SelectBuilder {
	b.inner = b.inner.GroupBy(b.modelMeta.escapeNames(groupBys)...)
	return b
}
func (b SelectBuilder) GroupByRow(groupBys ...string) SelectBuilder {
	b.inner = b.inner.GroupBy(groupBys...)
	return b
}

func (b SelectBuilder) OrderBy(field string) SelectBuilder {
	b.inner = b.inner.OrderBy(b.modelMeta.escapeName(field))
	return b
}
func (b SelectBuilder) OrderByDesc(field string) SelectBuilder {
	b.inner = b.inner.OrderBy(b.modelMeta.escapeName(field) + " DESC")
	return b
}

func (b SelectBuilder) Limit(limit uint64) SelectBuilder {
	b.inner = b.inner.Limit(limit)
	return b
}

func (b SelectBuilder) Offset(offset uint64) SelectBuilder {
	b.inner = b.inner.Offset(offset)
	return b
}
