package sqlkit

import (
	"github.com/Masterminds/squirrel"
)

type InsertBuilder struct {
	inner     squirrel.InsertBuilder
	modelMeta ModelMeta
}

func (b InsertBuilder) Sql() (string, []interface{}) {
	return b.inner.MustSql()
}
func (b InsertBuilder) ToSql() (string, []interface{}, error) {
	return b.inner.ToSql()
}

// SQL methods

// Prefix 在 sql 前写入语句
func (b InsertBuilder) Prefix(sql string, args ...interface{}) InsertBuilder {
	b.inner = b.inner.Prefix(sql, args...)
	return b
}
func (b InsertBuilder) PrefixExpr(expr InsertBuilder) InsertBuilder {
	b.inner = b.inner.PrefixExpr(expr)
	return b
}
func (b InsertBuilder) Suffix(sql string, args ...interface{}) InsertBuilder {
	b.inner = b.inner.Suffix(sql, args...)
	return b
}
func (b InsertBuilder) SuffixExpr(expr InsertBuilder) InsertBuilder {
	b.inner = b.inner.SuffixExpr(expr)
	return b
}

func (b InsertBuilder) Options(options ...string) InsertBuilder {
	b.inner = b.inner.Options(options...)
	return b
}

// Columns adds insert columns to the query.
func (b InsertBuilder) Columns(columns ...string) InsertBuilder {
	b.inner = b.inner.Columns(columns...)
	return b
}

// Values adds a single row's values to the query.
func (b InsertBuilder) Values(values ...interface{}) InsertBuilder {
	b.inner = b.inner.Values(values...)
	return b
}

func (b InsertBuilder) Select(sb SelectBuilder) InsertBuilder {
	b.inner = b.inner.Values(sb.inner)
	return b
}

//func (b InsertBuilder) SetMap(clauses map[string]interface{}) InsertBuilder {
//	b.inner = b.inner.SetMap(clauses)
//	return b
//}
