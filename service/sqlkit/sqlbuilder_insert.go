package sqlkit

import (
	"github.com/Masterminds/squirrel"
	"github.com/mizuki1412/go-core-kit/class/exception"
)

type InsertBuilder struct {
	inner     squirrel.InsertBuilder
	modelMeta ModelMeta
	driver    string
}

func (b InsertBuilder) Sql() (string, []interface{}) {
	sql, args, err := b.ToSql()
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return sql, args
}
func (b InsertBuilder) ToSql() (string, []interface{}, error) {
	b.inner = b.inner.PlaceholderFormat(placeholder(b.driver))
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
	b.inner = b.inner.Columns(b.modelMeta.escapeNames(columns)...)
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
