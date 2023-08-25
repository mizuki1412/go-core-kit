package sqlkit

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

type SelectBuilder struct {
	inner     squirrel.SelectBuilder
	modelMeta ModelMeta
	logicDel  []any
	driver    string
	fromAs    string
}

func (b SelectBuilder) Sql() (string, []any) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return sql, args
}

// 默认占位符的，一般用于子查询
func (b SelectBuilder) sqlOriginPlaceholder() (string, []any) {
	if b.fromAs == "" {
		b.inner = b.inner.From(b.modelMeta.getTable())
	} else {
		b.inner = b.inner.From(b.modelMeta.getTable(b.fromAs))
	}
	b.inner = b.inner.PlaceholderFormat(squirrel.Question)
	return b.inner.MustSql()
}
func (b SelectBuilder) ToSql() (string, []any, error) {
	if b.fromAs == "" {
		b.inner = b.inner.From(b.modelMeta.getTable())
	} else {
		b.inner = b.inner.From(b.modelMeta.getTable(b.fromAs))
	}
	b.inner = b.inner.PlaceholderFormat(placeholder(b.driver))
	return b.inner.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (b SelectBuilder) Prefix(sql string, args ...any) SelectBuilder {
	b.inner = b.inner.Prefix(sql, args...)
	return b
}
func (b SelectBuilder) PrefixExpr(expr SelectBuilder) SelectBuilder {
	b.inner = b.inner.PrefixExpr(expr)
	return b
}
func (b SelectBuilder) Suffix(sql string, args ...any) SelectBuilder {
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

func (b SelectBuilder) Join(dm DaoModelMeta, as string, rest ...any) SelectBuilder {
	b.inner = b.inner.Join(dm.getModelMeta().getTable(as), rest...)
	return b
}
func (b SelectBuilder) LeftJoin(dm DaoModelMeta, as string, rest ...any) SelectBuilder {
	b.inner = b.inner.LeftJoin(dm.getModelMeta().getTable(as), rest...)
	return b
}
func (b SelectBuilder) RightJoin(dm DaoModelMeta, as string, rest ...any) SelectBuilder {
	b.inner = b.inner.RightJoin(dm.getModelMeta().getTable(as), rest...)
	return b
}
func (b SelectBuilder) InnerJoin(dm DaoModelMeta, as string, rest ...any) SelectBuilder {
	b.inner = b.inner.InnerJoin(dm.getModelMeta().getTable(as), rest...)
	return b
}
func (b SelectBuilder) CrossJoin(dm DaoModelMeta, as string, rest ...any) SelectBuilder {
	b.inner = b.inner.CrossJoin(dm.getModelMeta().getTable(as), rest...)
	return b
}

func (b SelectBuilder) JoinRaw(join string, rest ...any) SelectBuilder {
	b.inner = b.inner.Join(join, rest...)
	return b
}
func (b SelectBuilder) LeftJoinRaw(join string, rest ...any) SelectBuilder {
	b.inner = b.inner.LeftJoin(join, rest...)
	return b
}
func (b SelectBuilder) RightJoinRaw(join string, rest ...any) SelectBuilder {
	b.inner = b.inner.RightJoin(join, rest...)
	return b
}
func (b SelectBuilder) InnerJoinRaw(join string, rest ...any) SelectBuilder {
	b.inner = b.inner.InnerJoin(join, rest...)
	return b
}
func (b SelectBuilder) CrossJoinRaw(join string, rest ...any) SelectBuilder {
	b.inner = b.inner.CrossJoin(join, rest...)
	return b
}

func (b SelectBuilder) Where(pred any, args ...any) SelectBuilder {
	b.inner = b.inner.Where(pred, args...)
	return b
}
func (b SelectBuilder) Having(pred any, rest ...any) SelectBuilder {
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

// custom

func (b SelectBuilder) WhereNLogicDel() SelectBuilder {
	if b.modelMeta.logicDelKey.Key != "" {
		return b.Where(b.modelMeta.logicDelKey.Key+"<>?", b.logicDel[0])
	}
	return b
}

// 生成sql中: sth in (select unnest(Array[?,?,?])) []any
// 注意使用时 args...
func (b SelectBuilder) whereUnnest(arr any, key, flag string) SelectBuilder {
	switch b.driver {
	case Postgres:
		s, v := pgArray(arr)
		return b.Where(fmt.Sprintf("%s %s (select unnest(%s))", b.modelMeta.escapeName(key), flag, s), v...)
	default:
		panic(exception.New("whereUnnest not supported"))
	}
}
func (b SelectBuilder) WhereUnnestIn(key string, arr any) SelectBuilder {
	return b.whereUnnest(arr, key, "IN")
}
func (b SelectBuilder) WhereUnnestNotIn(key string, arr any) SelectBuilder {
	return b.whereUnnest(arr, key, "NOT IN")
}

func (b SelectBuilder) WhereIn(key string, sub SelectBuilder) SelectBuilder {
	sql, args := sub.sqlOriginPlaceholder()
	return b.Where(squirrel.Expr(key+" IN ("+sql+")", args...))
}
